package corpus

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
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
	} {
		if !strings.Contains(string(rep), want) {
			t.Errorf("report lacks %q", want)
		}
	}
}

func TestReportIsDeterministic(t *testing.T) {
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
