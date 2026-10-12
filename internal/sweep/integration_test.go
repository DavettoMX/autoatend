package sweep_test

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/lexer"
	"autoatend/internal/sweep"
	"autoatend/internal/testutil"
)

// corpusPath is the versioned 73,000-comment corpus, relative to this package.
const corpusPath = "../../data/corpus.db"

var (
	fullDir  string
	fullOnce sync.Once
	fullRes  sweep.Result
	fullErr  error
)

func TestMain(m *testing.M) {
	code := m.Run()
	if fullDir != "" {
		os.RemoveAll(fullDir)
	}
	os.Exit(code)
}

// fullRun sweeps the whole corpus once per test binary and shares the outcome.
func fullRun(t *testing.T) (sweep.Result, string) {
	t.Helper()
	if testing.Short() {
		t.Skip("full-corpus sweep skipped in -short mode")
	}
	fullOnce.Do(func() {
		fullDir, fullErr = os.MkdirTemp("", "sweep-full-")
		if fullErr != nil {
			return
		}
		fullRes, fullErr = sweep.Run(sweep.Options{
			CorpusPath: corpusPath,
			OutPath:    filepath.Join(fullDir, "resultados.db"),
			DFPath:     filepath.Join(fullDir, "df_train.json"),
			Seed:       10433,
		})
	})
	if fullErr != nil {
		t.Fatalf("full sweep: %v", fullErr)
	}
	return fullRes, filepath.Join(fullDir, "resultados.db")
}

func openFull(t *testing.T) (sweep.Result, *sql.DB) {
	t.Helper()
	res, path := fullRun(t)
	d, err := db.OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return res, d
}

// CPP02 over the 73,000-comment corpus: the TE §6.3.3 query returns no rows.
func TestCPP02_FullCorpusFrequencyConservation(t *testing.T) {
	res, h := openFull(t)
	t.Logf("sweep of %d comments: %d lexemes, %d tokens, %v", res.Comentarios, res.Lexemas, res.Tokens, res.Elapsed)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM (`+teI2Query+`)`)
	testutil.Scalar(t, h, 1, `SELECT COUNT(*) FROM corridas`)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM lexemas l
		WHERE (SELECT SUM(frecuencia) FROM variantes WHERE lexema_id = l.id) !=
		      (SELECT valor FROM frecuencias WHERE lexema_id = l.id AND categoria_id IS NULL)`)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM lexemas WHERE id NOT IN (SELECT lexema_id FROM variantes)`)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM lexemas WHERE id NOT IN (SELECT lexema_id FROM frecuencias)`)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM frecuencias WHERE valor <= 0`)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM pragma_foreign_key_check`)
	testutil.Scalar(t, h, res.Lexemas, `SELECT total_lexemas FROM corridas`)
	testutil.Scalar(t, h, res.Lexemas, `SELECT COUNT(*) FROM lexemas`)
}

func TestFullCorpus_NoDigitsNoStopwords(t *testing.T) {
	_, h := openFull(t)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM lexemas WHERE forma_canonica GLOB '*[0-9]*' OR forma_canonica = ''`)
	testutil.Scalar(t, h, 0, `SELECT COUNT(*) FROM variantes WHERE forma GLOB '*[0-9]*'`)
	rows, err := h.Query(`SELECT forma FROM variantes`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var f string
		if err := rows.Scan(&f); err != nil {
			t.Fatal(err)
		}
		if class, stop := lexer.Lookup(f); stop {
			t.Errorf("stopword %q (%s) was catalogued as a variant", f, class)
		}
	}
}

// The variant table is compared with an independent count: split every comment
// on non-alphabet characters, lowercase, drop stopwords, count each form.
func TestFullCorpus_VariantsMatchIndependentCount(t *testing.T) {
	res, h := openFull(t)
	src, err := db.OpenReadOnly(corpusPath)
	if err != nil {
		t.Fatal(err)
	}
	defer src.Close()
	rows, err := src.Query(`SELECT texto, categoria_id FROM comentarios`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	isAlpha := func(r rune) bool {
		return (r >= 'a' && r <= 'z') || strings.ContainsRune("áéíóúüñ", r)
	}
	byForm := map[string]int{}
	var perCat [5]int
	tokens := 0
	for rows.Next() {
		var text string
		var cat int
		if err := rows.Scan(&text, &cat); err != nil {
			t.Fatal(err)
		}
		for _, w := range strings.FieldsFunc(strings.ToLower(text), func(r rune) bool { return !isAlpha(r) }) {
			if _, stop := lexer.Lookup(w); stop {
				continue
			}
			byForm[w]++
			perCat[cat-1]++
			tokens++
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if res.Tokens != tokens {
		t.Errorf("Result.Tokens = %d, independent count %d", res.Tokens, tokens)
	}
	testutil.Scalar(t, h, tokens, `SELECT SUM(valor) FROM frecuencias WHERE categoria_id IS NULL`)
	for cat := 1; cat <= 5; cat++ {
		testutil.Scalar(t, h, perCat[cat-1], `SELECT COALESCE(SUM(valor), 0) FROM frecuencias WHERE categoria_id = `+itoa(cat))
	}
	vr, err := h.Query(`SELECT forma, frecuencia FROM variantes`)
	if err != nil {
		t.Fatal(err)
	}
	defer vr.Close()
	n := 0
	for vr.Next() {
		var f string
		var freq int
		if err := vr.Scan(&f, &freq); err != nil {
			t.Fatal(err)
		}
		if byForm[f] != freq {
			t.Errorf("variante %q: stored %d, independent count %d", f, freq, byForm[f])
		}
		n++
	}
	if n != len(byForm) {
		t.Errorf("variantes rows = %d, distinct forms counted independently = %d", n, len(byForm))
	}
}

// CPP03 over the full corpus: a second sweep yields a byte-identical resultados.db.
func TestCPP03_FullCorpusIdempotence(t *testing.T) {
	first, _ := fullRun(t)
	dir := t.TempDir()
	second, err := sweep.Run(sweep.Options{
		CorpusPath: corpusPath,
		OutPath:    filepath.Join(dir, "resultados.db"),
		DFPath:     filepath.Join(dir, "df_train.json"),
		Seed:       10433,
	})
	if err != nil {
		t.Fatal(err)
	}
	if first.SHA256 != second.SHA256 {
		t.Errorf("sha256 differs: %s vs %s", first.SHA256, second.SHA256)
	}
	a, _ := os.ReadFile(filepath.Join(fullDir, "df_train.json"))
	b, _ := os.ReadFile(filepath.Join(dir, "df_train.json"))
	if string(a) != string(b) {
		t.Error("df_train.json differs between runs")
	}
}
