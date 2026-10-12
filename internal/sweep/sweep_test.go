package sweep_test

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
	"autoatend/internal/sweep"
	"autoatend/internal/testutil"
)

// teI2Query is the frequency-conservation query of TE §6.3.3, verbatim.
const teI2Query = `
	SELECT l.forma_canonica,
	       SUM(CASE WHEN f.categoria_id IS NULL THEN f.valor ELSE 0 END) AS global,
	       SUM(CASE WHEN f.categoria_id IS NOT NULL THEN f.valor ELSE 0 END) AS suma_categorias
	FROM lexemas l
	JOIN frecuencias f ON f.lexema_id = l.id
	GROUP BY l.id
	HAVING global != suma_categorias`

type row struct {
	fecha string
	texto string
	cat   int
}

// fixtureRows is a 20-comment corpus whose lexemes were counted by hand.
// With N = 20 the cut-off is the date of the 17th row (db.SplitCutoff), so
// rows 1..16 are train and rows 17..20 are test.
var fixtureRows = []row{
	{"2025-01-01", "La factura llegó mal", 1},
	{"2025-01-02", "Facturas de más y la factura de abril", 1},
	{"2025-01-03", "Me cobraron de más, cobro indebido", 2},
	{"2025-01-04", "Cobros indebidos $599.00", 2},
	{"2025-01-05", "factura factura factura", 3},
	{"2025-01-06", "Mal servicio, el servicio es malo", 4},
	{"2025-01-07", "Gracias por el servicio", 5},
	{"2025-01-08", "Sin señal en mi casa", 1},
	{"2025-01-09", "Sin señal en mi oficina", 1},
	{"2025-01-10", "No hay internet", 1},
	{"2025-01-11", "De la y el", 4},
	{"2025-01-12", "Pésima atención al cliente", 4},
	{"2025-01-13", "Gracias por la atención", 5},
	{"2025-01-14", "Excelente atención, gracias", 5},
	{"2025-01-15", "Mi recibo no llegó", 3},
	{"2025-01-16", "El recibo llegó tarde", 3},
	{"2025-01-17", "Sin señal otra vez", 1},
	{"2025-01-18", "Error en la facturación", 3},
	{"2025-01-19", "cobro cobro", 2},
	{"2025-01-20", "Gracias", 5},
}

const (
	fixtureCut     = "2025-01-17"
	fixtureTokens  = 56 // kept words, counted by hand
	fixtureLexemas = 27 // distinct stems among them, counted by hand
	fixtureZeroDF  = 3  // otra, vez, error: only in test rows
)

func writeCorpus(t *testing.T, path string, rows []row) {
	t.Helper()
	d, err := build.CreateCorpus(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if err := build.InsertCategorias(d); err != nil {
		t.Fatal(err)
	}
	for _, r := range rows {
		testutil.MustExec(t, d, `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, ?, ?)`, r.fecha, r.texto, r.cat)
	}
	if err := build.Finalize(d); err != nil {
		t.Fatal(err)
	}
}

func fixtureCorpus(t *testing.T) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "corpus.db")
	writeCorpus(t, p, fixtureRows)
	return p
}

func runFixture(t *testing.T) (sweep.Result, *sql.DB, string, string) {
	t.Helper()
	dir := t.TempDir()
	out := filepath.Join(dir, "resultados.db")
	df := filepath.Join(dir, "df_train.json")
	res, err := sweep.Run(sweep.Options{CorpusPath: fixtureCorpus(t), OutPath: out, DFPath: df, Seed: 10433})
	if err != nil {
		t.Fatal(err)
	}
	d, err := db.OpenReadOnly(out)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return res, d, out, df
}

func scalarInt(t *testing.T, d *sql.DB, q string, args ...any) int {
	t.Helper()
	var n int
	if err := d.QueryRow(q, args...).Scan(&n); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	return n
}

func TestRun_FixtureCounts(t *testing.T) {
	res, d, _, _ := runFixture(t)
	testutil.Scalar(t, d, fixtureLexemas, `SELECT COUNT(*) FROM lexemas`)
	testutil.Scalar(t, d, 1, `SELECT COUNT(*) FROM corridas`)
	testutil.Scalar(t, d, fixtureTokens, `SELECT SUM(valor) FROM frecuencias WHERE categoria_id IS NULL`)
	testutil.Scalar(t, d, fixtureTokens, `SELECT SUM(valor) FROM frecuencias WHERE categoria_id IS NOT NULL`)
	testutil.Scalar(t, d, fixtureTokens, `SELECT SUM(frecuencia) FROM variantes`)
	if res.Lexemas != fixtureLexemas || res.Tokens != fixtureTokens {
		t.Errorf("Result lexemas/tokens = %d/%d, want %d/%d", res.Lexemas, res.Tokens, fixtureLexemas, fixtureTokens)
	}
	if res.Comentarios != len(fixtureRows) || res.TrainDocs != 16 || res.Cutoff != fixtureCut {
		t.Errorf("Result comentarios/train/cutoff = %d/%d/%s", res.Comentarios, res.TrainDocs, res.Cutoff)
	}
	if res.ZeroDF != fixtureZeroDF {
		t.Errorf("ZeroDF = %d, want %d", res.ZeroDF, fixtureZeroDF)
	}
}

func TestRun_FixtureLexemeFactura(t *testing.T) {
	_, d, _, _ := runFixture(t)
	var id int
	var er string
	if err := d.QueryRow(`SELECT id, expresion_regular FROM lexemas WHERE forma_canonica = 'factura'`).Scan(&id, &er); err != nil {
		t.Fatal(err)
	}
	if er != `factur(a|ación|as)` {
		t.Errorf("ER = %q, want factur(a|ación|as)", er)
	}
	testutil.Scalar(t, d, 7, `SELECT valor FROM frecuencias WHERE lexema_id = `+itoa(id)+` AND categoria_id IS NULL`)
	testutil.Scalar(t, d, 3, `SELECT valor FROM frecuencias WHERE lexema_id = `+itoa(id)+` AND categoria_id = 1`)
	testutil.Scalar(t, d, 4, `SELECT valor FROM frecuencias WHERE lexema_id = `+itoa(id)+` AND categoria_id = 3`)
	// Per-category rows exist only for valor > 0.
	testutil.Scalar(t, d, 3, `SELECT COUNT(*) FROM frecuencias WHERE lexema_id = `+itoa(id))
	rows, err := d.Query(`SELECT forma, frecuencia FROM variantes WHERE lexema_id = ? ORDER BY forma`, id)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	got := map[string]int{}
	for rows.Next() {
		var f string
		var n int
		if err := rows.Scan(&f, &n); err != nil {
			t.Fatal(err)
		}
		got[f] = n
	}
	want := map[string]int{"factura": 5, "facturas": 1, "facturación": 1}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("variantes = %v, want %v", got, want)
	}
}

func itoa(n int) string { b, _ := json.Marshal(n); return string(b) }

func TestRun_FixtureOtherLexemes(t *testing.T) {
	_, d, _, _ := runFixture(t)
	type want struct {
		canon  string
		er     string
		global int
	}
	for _, w := range []want{
		{"cobro", `cobr(aron|o|os)`, 5},
		{"servicio", `servicio`, 3},
	} {
		var er string
		var g int
		err := d.QueryRow(`SELECT l.expresion_regular, f.valor FROM lexemas l JOIN frecuencias f ON f.lexema_id = l.id
			WHERE l.forma_canonica = ? AND f.categoria_id IS NULL`, w.canon).Scan(&er, &g)
		if err != nil {
			t.Fatalf("%s: %v", w.canon, err)
		}
		if er != w.er || g != w.global {
			t.Errorf("%s: ER %q global %d, want %q %d", w.canon, er, g, w.er, w.global)
		}
	}
	// mal/malo are one lexeme (stem "mal"): the variant "mal" is the root, so the ER backs off one rune.
	var er string
	if err := d.QueryRow(`SELECT expresion_regular FROM lexemas WHERE forma_canonica = 'mal'`).Scan(&er); err != nil {
		t.Fatal(err)
	}
	if er != `ma(l|lo)` {
		t.Errorf("mal ER = %q, want ma(l|lo)", er)
	}
	// Digits are recognized by the DFA but never catalogued.
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM lexemas WHERE forma_canonica GLOB '*[0-9]*'`)
	// Stopwords never appear as lexemes or variants.
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM variantes WHERE forma IN ('de','la','y','el','en','por','sin','al')`)
}

func TestCPP02_FixtureFrequencyConservation(t *testing.T) {
	_, d, _, _ := runFixture(t)
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM (`+teI2Query+`)`)
	// Σ variantes.frecuencia = global, per lexeme.
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM lexemas l
		WHERE (SELECT SUM(frecuencia) FROM variantes WHERE lexema_id = l.id) !=
		      (SELECT valor FROM frecuencias WHERE lexema_id = l.id AND categoria_id IS NULL)`)
	// No lexeme without variants or without frequencies.
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM lexemas WHERE id NOT IN (SELECT lexema_id FROM variantes)`)
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM lexemas WHERE id NOT IN (SELECT lexema_id FROM frecuencias)`)
	// Exactly one global row per lexeme.
	testutil.Scalar(t, d, fixtureLexemas, `SELECT COUNT(*) FROM frecuencias WHERE categoria_id IS NULL`)
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM pragma_foreign_key_check`)
}

func TestRun_CorridaRow(t *testing.T) {
	_, d, _, _ := runFixture(t)
	var id, semilla, total int
	var parametros string
	if err := d.QueryRow(`SELECT id, semilla, parametros, total_lexemas FROM corridas`).Scan(&id, &semilla, &parametros, &total); err != nil {
		t.Fatal(err)
	}
	if id != 1 || semilla != 10433 || total != fixtureLexemas {
		t.Errorf("corrida = id %d semilla %d total %d", id, semilla, total)
	}
	var p map[string]json.RawMessage
	if err := json.Unmarshal([]byte(parametros), &p); err != nil {
		t.Fatalf("parametros is not JSON: %v", err)
	}
	for _, k := range []string{"semilla", "comentarios", "ruido", "stoplist", "tokenizador", "stemmer", "corte_train"} {
		if _, ok := p[k]; !ok {
			t.Errorf("parametros lacks %q: %s", k, parametros)
		}
	}
	if string(p["comentarios"]) != "20" || string(p["stoplist"]) != `"v1"` || string(p["stemmer"]) != `"v1"` ||
		string(p["corte_train"]) != `"`+fixtureCut+`"` || string(p["semilla"]) != "10433" {
		t.Errorf("unexpected parametros: %s", parametros)
	}
	// Reproducible by construction: no timestamps, no absolute paths.
	for _, bad := range []string{"/tmp", "/home", ".db"} {
		if strings.Contains(parametros, bad) {
			t.Errorf("parametros contains %q: %s", bad, parametros)
		}
	}
}

func TestRun_DFTrainJSON(t *testing.T) {
	res, _, _, dfPath := runFixture(t)
	raw, err := os.ReadFile(dfPath)
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Corte           string `json:"corte"`
		DocumentosTrain int    `json:"documentos_train"`
		DF              []struct {
			Lexema string `json:"lexema"`
			DF     int    `json:"df"`
		} `json:"df"`
	}
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatalf("df_train.json is not JSON: %v", err)
	}
	if f.Corte != fixtureCut || f.DocumentosTrain != 16 {
		t.Errorf("corte/documentos_train = %s/%d", f.Corte, f.DocumentosTrain)
	}
	if len(f.DF) != fixtureLexemas-fixtureZeroDF || res.Lexemas-res.ZeroDF != len(f.DF) {
		t.Errorf("df entries = %d, want %d", len(f.DF), fixtureLexemas-fixtureZeroDF)
	}
	got := map[string]int{}
	var prev string
	for i, e := range f.DF {
		if e.DF <= 0 {
			t.Errorf("%s has df %d; only df > 0 is exported", e.Lexema, e.DF)
		}
		if i > 0 && e.Lexema <= prev {
			t.Errorf("df not sorted by lexeme: %q after %q", e.Lexema, prev)
		}
		prev = e.Lexema
		got[e.Lexema] = e.DF
	}
	// Document frequencies counted by hand over the 16 train rows (each stem
	// counted once per document).
	for lex, want := range map[string]int{"factura": 3, "cobro": 2, "servicio": 2, "gracias": 3, "señal": 2, "atención": 3, "recibo": 2, "mal": 2} {
		if got[lex] != want {
			t.Errorf("df[%s] = %d, want %d", lex, got[lex], want)
		}
	}
	for _, zero := range []string{"otra", "vez", "error"} {
		if _, ok := got[zero]; ok {
			t.Errorf("%s only occurs in test and must not be exported", zero)
		}
	}
}

func TestCPP03_FixtureIdempotence(t *testing.T) {
	corpus := fixtureCorpus(t)
	var shas []string
	var dfs [][]byte
	for i := 0; i < 2; i++ {
		dir := t.TempDir()
		out := filepath.Join(dir, "r.db")
		df := filepath.Join(dir, "df.json")
		res, err := sweep.Run(sweep.Options{CorpusPath: corpus, OutPath: out, DFPath: df, Seed: 10433})
		if err != nil {
			t.Fatal(err)
		}
		raw, err := os.ReadFile(out)
		if err != nil {
			t.Fatal(err)
		}
		sum := sha256.Sum256(raw)
		if hex.EncodeToString(sum[:]) != res.SHA256 {
			t.Errorf("Result.SHA256 %s differs from the file's %x", res.SHA256, sum)
		}
		shas = append(shas, res.SHA256)
		dfRaw, _ := os.ReadFile(df)
		dfs = append(dfs, dfRaw)
	}
	if shas[0] != shas[1] {
		t.Errorf("two runs produced different resultados.db: %s vs %s", shas[0], shas[1])
	}
	if !bytes.Equal(dfs[0], dfs[1]) {
		t.Error("two runs produced different df_train.json")
	}
}

func TestRun_CategoryIntegrity(t *testing.T) {
	corpus := fixtureCorpus(t)
	out := filepath.Join(t.TempDir(), "r.db")
	if _, err := sweep.Run(sweep.Options{CorpusPath: corpus, OutPath: out, Seed: 1}); err != nil {
		t.Fatal(err)
	}
	d, err := db.OpenReadOnly(out)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	// Direct SQL oracle: the categories used in frecuencias are those of the corpus.
	rows, err := d.Query(`SELECT DISTINCT categoria_id FROM frecuencias WHERE categoria_id IS NOT NULL ORDER BY 1`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var cats []int
	for rows.Next() {
		var c int
		if err := rows.Scan(&c); err != nil {
			t.Fatal(err)
		}
		cats = append(cats, c)
	}
	if !reflect.DeepEqual(cats, []int{1, 2, 3, 4, 5}) {
		t.Errorf("categorias = %v", cats)
	}
	orphans, err := build.CheckCategoriaIntegrity(d, corpus)
	if err != nil || len(orphans) != 0 {
		t.Errorf("orphans = %v, err = %v", orphans, err)
	}
}

func TestRun_IDsFollowCanonicalOrder(t *testing.T) {
	_, d, _, _ := runFixture(t)
	// lexemas.id is assigned in canonical order, so the table is deterministic.
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM lexemas a JOIN lexemas b ON b.id = a.id + 1 WHERE b.forma_canonica < a.forma_canonica`)
	testutil.Scalar(t, d, fixtureLexemas, `SELECT MAX(id) FROM lexemas`)
}

func TestRun_RefusesToOverwriteAndKeepsTheExistingFile(t *testing.T) {
	corpus := fixtureCorpus(t)
	out := filepath.Join(t.TempDir(), "r.db")
	if err := os.WriteFile(out, []byte("precious"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := sweep.Run(sweep.Options{CorpusPath: corpus, OutPath: out, Seed: 1}); err == nil {
		t.Fatal("expected an error when OutPath exists")
	}
	if b, _ := os.ReadFile(out); string(b) != "precious" {
		t.Errorf("existing file was modified or removed: %q", b)
	}
}

func TestRun_RefusesExistingDFBeforeDoingAnyWork(t *testing.T) {
	corpus := fixtureCorpus(t)
	dir := t.TempDir()
	out := filepath.Join(dir, "r.db")
	df := filepath.Join(dir, "df.json")
	if err := os.WriteFile(df, []byte("{}"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := sweep.Run(sweep.Options{CorpusPath: corpus, OutPath: out, DFPath: df, Seed: 1}); err == nil {
		t.Fatal("expected an error when DFPath exists")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("resultados.db was created although DFPath already exists")
	}
	if b, _ := os.ReadFile(df); string(b) != "{}" {
		t.Errorf("existing df file modified: %q", b)
	}
}

func TestRun_AllOrNothingOnFailure(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "r.db")
	df := filepath.Join(dir, "df.json")

	// Missing corpus.
	if _, err := sweep.Run(sweep.Options{CorpusPath: filepath.Join(dir, "nope.db"), OutPath: out, DFPath: df}); err == nil {
		t.Error("expected an error for a missing corpus")
	}
	// A corpus with no catalogable word has total_lexemas = 0, which the DDL rejects.
	empty := filepath.Join(dir, "empty.db")
	writeCorpus(t, empty, []row{{"2025-01-01", "de la y el", 1}, {"2025-01-02", "123", 2}})
	if _, err := sweep.Run(sweep.Options{CorpusPath: empty, OutPath: out, DFPath: df}); err == nil {
		t.Error("expected an error for a corpus without lexemes")
	}
	for _, p := range []string{out, df} {
		if _, err := os.Stat(p); !os.IsNotExist(err) {
			t.Errorf("%s left behind after a failed run", p)
		}
	}
}

func TestRun_WithoutDFPathWritesOnlyResultados(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "r.db")
	if _, err := sweep.Run(sweep.Options{CorpusPath: fixtureCorpus(t), OutPath: out, Seed: 1}); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(dir)
	if len(entries) != 1 {
		t.Errorf("dir has %d files, want only resultados.db", len(entries))
	}
}

// If the DF export fails after resultados.db was written, the run is all or nothing.
func TestRun_DFWriteFailureRemovesResultados(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "r.db")
	df := filepath.Join(dir, "missing-dir", "df.json")
	if _, err := sweep.Run(sweep.Options{CorpusPath: fixtureCorpus(t), OutPath: out, DFPath: df, Seed: 1}); err == nil {
		t.Fatal("expected an error when df_train.json cannot be created")
	}
	if _, err := os.Stat(out); !os.IsNotExist(err) {
		t.Error("resultados.db left behind after the DF export failed")
	}
}
