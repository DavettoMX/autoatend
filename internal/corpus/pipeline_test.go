package corpus

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"autoatend/internal/corpus/templates"
	"autoatend/internal/db"
)

func smallOptions(t *testing.T) Options {
	t.Helper()
	cfg := DefaultConfig()
	cfg.Total = 2000
	dir := t.TempDir()
	return Options{Config: cfg, OutPath: filepath.Join(dir, "corpus.db"), ReportPath: filepath.Join(dir, "report.md")}
}

func TestPersistIsAllOrNothing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "corpus.db")
	cs := []Comment{
		{Fecha: "2024-08-30", Texto: "bien", CategoriaID: 1},
		{Fecha: "2024-08-30", Texto: "mal", CategoriaID: 9}, // violates the FK
	}
	if err := Persist(path, cs); err == nil {
		t.Fatal("Persist with an invalid categoria_id: expected error")
	}
	if _, err := os.Stat(path); err == nil {
		t.Fatal("failed Persist left a partial corpus.db behind")
	}
}

func TestRunRefusesToPersistWhenValidationFails(t *testing.T) {
	opts := smallOptions(t)
	opts.generate = func(cfg Config) ([]Comment, error) {
		cs, err := Generate(cfg)
		if err == nil {
			cs[5].Texto = ""
		}
		return cs, err
	}
	_, err := Run(opts)
	if err == nil {
		t.Fatal("run: expected validation error")
	}
	var ve *ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("error %T is not *ValidationError: %v", err, err)
	}
	if _, statErr := os.Stat(opts.OutPath); statErr == nil {
		t.Fatal("corpus.db was written although validation failed")
	}
}

func TestRunProducesVerifiedCorpusAndReport(t *testing.T) {
	pendingTemplates(t)
	opts := smallOptions(t)
	res, err := Run(opts)
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	p := res.Persisted
	if p.Rows != 2000 || p.PerCategory != [5]int{600, 400, 400, 400, 200} || p.Empty != 0 {
		t.Fatalf("persisted stats = %+v", p)
	}
	for i := 0; i < 5; i++ {
		if p.Train[i]+p.Test[i] != p.PerCategory[i] || p.Train[i] == 0 || p.Test[i] == 0 {
			t.Errorf("category %d split %d/%d of %d", i+1, p.Train[i], p.Test[i], p.PerCategory[i])
		}
		if !(p.MaxTrainDate[i] < p.Cutoff && p.Cutoff <= p.MinTestDate[i]) {
			t.Errorf("category %d: train ends %s, test starts %s, cut-off %s", i+1, p.MaxTrainDate[i], p.MinTestDate[i], p.Cutoff)
		}
	}
	if len(p.SHA256) != 64 {
		t.Errorf("sha256 = %q", p.SHA256)
	}

	rep, err := os.ReadFile(opts.ReportPath)
	if err != nil {
		t.Fatalf("report not written: %v", err)
	}
	for _, want := range []string{
		"# Generator validation report", "## Class distribution", "## Validation battery",
		"## Lexical overlap", "## Noise calibration", "## Persisted corpus verification",
		"## Sample review", "10433", "Falta de Servicio", "Reconocimientos", p.SHA256,
		"## Generator parameters",
		fmt.Sprintf("| Templates per category | %d / %d / %d / %d / %d |", len(templates.Templates[0]), len(templates.Templates[1]),
			len(templates.Templates[2]), len(templates.Templates[3]), len(templates.Templates[4])),
		fmt.Sprintf("| Complements per category | %d / %d / %d / %d / %d |", len(templates.Complements[0]), len(templates.Complements[1]),
			len(templates.Complements[2]), len(templates.Complements[3]), len(templates.Complements[4])),
		fmt.Sprintf("| Shared openers / closers | %d / %d |", len(templates.SharedOpeners), len(templates.SharedClosers)),
		"| Opener / complement / closer probability | 0.25 / 0.40 / 0.30 |",
		fmt.Sprintf("| Amount pool / surface formats | %d / 6 |", len(templates.AmountPool)),
		fmt.Sprintf("| Folio prefixes / digits | %d / 5–7 |", len(templates.FolioPrefixes)),
		fmt.Sprintf("Test comments with an identical text in train: %d of %d", p.TestInTrain, p.Rows-sum5(p.Train)),
	} {
		if !strings.Contains(string(rep), want) {
			t.Errorf("report lacks %q", want)
		}
	}
}

func TestReportIsDeterministic(t *testing.T) {
	pendingTemplates(t)
	render := func() ([]byte, string) {
		opts := smallOptions(t)
		res, err := Run(opts)
		if err != nil {
			t.Fatal(err)
		}
		var b bytes.Buffer
		if err := writeReport(&b, opts.Config, res); err != nil {
			t.Fatal(err)
		}
		return b.Bytes(), res.Persisted.SHA256
	}
	a, shaA := render()
	b, shaB := render()
	if !bytes.Equal(a, b) {
		t.Fatal("report differs between identical runs")
	}
	if shaA != shaB {
		t.Fatalf("corpus.db SHA-256 differs between identical runs: %s vs %s", shaA, shaB)
	}
}

// The train/test duplicate count is checked against an oracle built in Go from
// the persisted rows, not against the SQL used by verifyCorpus.
func TestPersistedTestDuplicatesInTrain(t *testing.T) {
	pendingTemplates(t)
	opts := smallOptions(t)
	res, err := Run(opts)
	if err != nil {
		t.Fatal(err)
	}
	d, err := db.OpenReadOnly(opts.OutPath)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	rows, err := d.Query(`SELECT fecha, texto FROM comentarios`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	train := map[string]bool{}
	var test []string
	for rows.Next() {
		var fecha, texto string
		if err := rows.Scan(&fecha, &texto); err != nil {
			t.Fatal(err)
		}
		if fecha < res.Persisted.Cutoff {
			train[texto] = true
		} else {
			test = append(test, texto)
		}
	}
	want := 0
	for _, texto := range test {
		if train[texto] {
			want++
		}
	}
	if want == 0 {
		t.Fatal("fixture has no test/train duplicates; the check would be vacuous")
	}
	if res.Persisted.TestInTrain != want {
		t.Fatalf("TestInTrain = %d, want %d", res.Persisted.TestInTrain, want)
	}
}
