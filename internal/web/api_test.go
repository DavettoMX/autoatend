package web_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sort"
	"strings"
	"testing"
	"time"

	"autoatend/internal/db"
	"autoatend/internal/web"
)

type clasificarResp struct {
	Categoria    string             `json:"categoria"`
	Confianzas   map[string]float64 `json:"confianzas"`
	LexemasClave []struct {
		Lexema           string  `json:"lexema"`
		ExpresionRegular string  `json:"expresion_regular"`
		Contribucion     float64 `json:"contribucion"`
	} `json:"lexemas_clave"`
	Vacio bool `json:"vacio"`
}

type errorResp struct {
	Error string `json:"error"`
}

func categoryNames() map[string]bool {
	m := map[string]bool{}
	for _, c := range db.Categories {
		m[c.Nombre] = true
	}
	return m
}

// openOracle opens the extracted resultados.db read-only, independently of the code under test.
func openOracle(t *testing.T, dir string) *sql.DB {
	t.Helper()
	d, err := db.OpenReadOnly(filepath.Join(dir, "resultados.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

// CPA01: valid comment -> HTTP 200, categoria among the 5 case labels, 5 confidences summing to 1 (TE §6.3.7).
func TestCPA01_ClasificarValido(t *testing.T) {
	a := bootApp(t)
	h := a.Handler()
	oracle := openOracle(t, a.Dir)
	names := categoryNames()

	for _, g := range loadGold(t) {
		rec := postClasificar(t, h, g.Texto)
		if rec.Code != http.StatusOK {
			t.Fatalf("gold %d: status %d: %s", g.ID, rec.Code, rec.Body)
		}
		if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
			t.Errorf("Content-Type = %q", ct)
		}
		var r clasificarResp
		decode(t, rec, &r)
		if !names[r.Categoria] {
			t.Errorf("gold %d: categoria %q is not one of the 5 case labels", g.ID, r.Categoria)
		}
		if r.Categoria != g.Categoria {
			t.Errorf("gold %d: categoria %q, want %q", g.ID, r.Categoria, g.Categoria)
		}
		if len(r.Confianzas) != 5 {
			t.Fatalf("gold %d: %d confidences, want 5: %v", g.ID, len(r.Confianzas), r.Confianzas)
		}
		var sum float64
		for n, p := range r.Confianzas {
			if !names[n] {
				t.Errorf("confianzas has unknown key %q", n)
			}
			if p < 0 || p > 1 {
				t.Errorf("confianzas[%q] = %v outside [0,1]", n, p)
			}
			sum += p
		}
		if math.Abs(sum-1) > 1e-9 {
			t.Errorf("gold %d: confidences sum to %v, want 1", g.ID, sum)
		}
		if r.Vacio {
			t.Errorf("gold %d: vacio = true for a comment with catalogue words", g.ID)
		}
		if len(r.LexemasClave) == 0 || len(r.LexemasClave) > 5 {
			t.Errorf("gold %d: %d lexemas_clave, want 1..5", g.ID, len(r.LexemasClave))
		}
		// Every key lexeme carries the regular expression stored in the catalogue of resultados.db.
		for _, lc := range r.LexemasClave {
			var er string
			err := oracle.QueryRow(`SELECT expresion_regular FROM lexemas WHERE forma_canonica = ?`, lc.Lexema).Scan(&er)
			if err != nil {
				t.Fatalf("lexeme %q is not in the catalogue: %v", lc.Lexema, err)
			}
			if lc.ExpresionRegular != er || er == "" {
				t.Errorf("lexeme %q: expresion_regular %q, catalogue has %q", lc.Lexema, lc.ExpresionRegular, er)
			}
			if lc.Contribucion == 0 {
				t.Errorf("lexeme %q: contribucion is 0", lc.Lexema)
			}
		}
	}
}

func TestCPA01_TextoSinVocabularioMarcaVacio(t *testing.T) {
	h := bootApp(t).Handler()
	rec := postClasificar(t, h, "zzzz qqqq xxxx")
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d: %s", rec.Code, rec.Body)
	}
	var raw map[string]json.RawMessage
	decode(t, rec, &raw)
	if string(raw["vacio"]) != "true" {
		t.Errorf("vacio = %s, want true", raw["vacio"])
	}
	if string(raw["lexemas_clave"]) != "[]" {
		t.Errorf("lexemas_clave = %s, want []", raw["lexemas_clave"])
	}
}

// CPA02: empty text -> HTTP 400 with a descriptive message and no panic.
func TestCPA02_TextoVacio(t *testing.T) {
	h := bootApp(t).Handler()
	cases := map[string]string{
		"vacío":              `{"texto": ""}`,
		"solo espacios":      `{"texto": "  \n\t "}`,
		"sin el campo":       `{}`,
		"JSON mal formado":   `{"texto": `,
		"no es un objeto":    `[1,2]`,
		"texto no es string": `{"texto": 5}`,
		"cuerpo vacío":       ``,
	}
	for name, body := range cases {
		rec := do(t, h, http.MethodPost, "/api/clasificar", strings.NewReader(body))
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, want 400", name, rec.Code)
			continue
		}
		var e errorResp
		decode(t, rec, &e)
		if len(e.Error) < 10 {
			t.Errorf("%s: error message %q is not descriptive", name, e.Error)
		}
	}
	// The spec example mentions the field name: the message for an empty text says what to send.
	rec := postClasificar(t, h, "")
	if !strings.Contains(rec.Body.String(), "texto") {
		t.Errorf("empty-text message does not mention the field: %s", rec.Body)
	}
}

// CPA03: SQL injection in the text -> 200 or 400, no database error and resultados.db is unchanged.
func TestCPA03_InyeccionSQL(t *testing.T) {
	a := bootApp(t)
	h := a.Handler()
	res := filepath.Join(a.Dir, "resultados.db")
	corpus := filepath.Join(a.Dir, "corpus.db")
	before, beforeCorpus := sha(t, res), sha(t, corpus)

	payloads := []string{
		`'; DROP TABLE lexemas; --`,
		`" OR 1=1 --`,
		`cobraron'); DELETE FROM frecuencias; --`,
		`1; ATTACH DATABASE '/tmp/x.db' AS x; --`,
		`Robert'); UPDATE lexemas SET expresion_regular='' WHERE 1=1; --`,
	}
	for _, p := range payloads {
		rec := postClasificar(t, h, p)
		if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest {
			t.Errorf("payload %q: status %d, want 200 or 400: %s", p, rec.Code, rec.Body)
		}
		// The same payloads through the query parameters of the catalogue endpoint.
		rec = do(t, h, http.MethodGet, "/api/lexemas?q="+httptestEscape(p), nil)
		if rec.Code != http.StatusOK && rec.Code != http.StatusBadRequest {
			t.Errorf("q=%q: status %d: %s", p, rec.Code, rec.Body)
		}
	}
	if got := sha(t, res); got != before {
		t.Errorf("resultados.db changed: %s -> %s", before, got)
	}
	if got := sha(t, corpus); got != beforeCorpus {
		t.Errorf("corpus.db changed: %s -> %s", beforeCorpus, got)
	}
	// The catalogue is intact and the service still answers.
	var n int
	if err := openOracle(t, a.Dir).QueryRow(`SELECT COUNT(*) FROM lexemas`).Scan(&n); err != nil || n == 0 {
		t.Fatalf("lexemas after injection: %d, %v", n, err)
	}
	if rec := postClasificar(t, h, "me cobraron de más"); rec.Code != http.StatusOK {
		t.Errorf("service broken after injection: %d", rec.Code)
	}
}

func httptestEscape(s string) string {
	r := strings.NewReplacer(" ", "%20", "'", "%27", `"`, "%22", ";", "%3B", "#", "%23", "&", "%26", "+", "%2B")
	return r.Replace(s)
}

// A 1 MB body is rejected with 413 (http.MaxBytesReader), quickly (TE §6.3.7: 413 or < 50 ms).
func TestCuerpoDe1MBSeRechaza413(t *testing.T) {
	h := bootApp(t).Handler()
	big := strings.Repeat("a", 1<<20)
	start := time.Now()
	rec := postClasificar(t, h, big)
	elapsed := time.Since(start)
	if rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status %d, want 413", rec.Code)
	}
	var e errorResp
	decode(t, rec, &e)
	if e.Error == "" {
		t.Error("413 without a message")
	}
	t.Logf("1 MB body rejected in %v", elapsed)
	if !testing.Short() && elapsed > 50*time.Millisecond {
		t.Errorf("1 MB body took %v, want < 50 ms", elapsed)
	}
	// A body just under the limit is processed.
	ok := postClasificar(t, h, strings.Repeat("hola ", 1000))
	if ok.Code != http.StatusOK {
		t.Errorf("5 KB body: status %d", ok.Code)
	}
}

func TestMetodoYRutaDesconocidos(t *testing.T) {
	h := bootApp(t).Handler()
	for _, c := range []struct {
		method, path string
		want         int
	}{
		{http.MethodGet, "/api/clasificar", http.StatusMethodNotAllowed},
		{http.MethodPost, "/api/health", http.StatusMethodNotAllowed},
		{http.MethodDelete, "/api/lexemas", http.StatusMethodNotAllowed},
		{http.MethodPost, "/", http.StatusMethodNotAllowed},
		{http.MethodGet, "/api/nada", http.StatusNotFound},
		{http.MethodGet, "/otra", http.StatusNotFound},
	} {
		rec := do(t, h, c.method, c.path, nil)
		if rec.Code != c.want {
			t.Errorf("%s %s: status %d, want %d", c.method, c.path, rec.Code, c.want)
			continue
		}
		var e errorResp
		decode(t, rec, &e)
		if e.Error == "" {
			t.Errorf("%s %s: no JSON error message", c.method, c.path)
		}
	}
}

// /api/categorias: exactly the 5 case categories with their ids.
func TestApiCategorias(t *testing.T) {
	h := bootApp(t).Handler()
	rec := do(t, h, http.MethodGet, "/api/categorias", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var got []struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
	}
	decode(t, rec, &got)
	if len(got) != 5 {
		t.Fatalf("%d categories, want 5: %v", len(got), got)
	}
	for i, c := range db.Categories {
		if got[i].ID != c.ID || got[i].Nombre != c.Nombre {
			t.Errorf("categoria %d = %+v, want {%d %s}", i, got[i], c.ID, c.Nombre)
		}
	}
}

// /api/health: exactly {"ok","version","lexemas","puerto"}.
func TestApiHealth(t *testing.T) {
	a := bootWith(t, web.Options{Version: "v-test"})
	a.SetPort(8473)
	rec := do(t, a.Handler(), http.MethodGet, "/api/health", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	var m map[string]any
	decode(t, rec, &m)
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	if strings.Join(keys, ",") != "lexemas,ok,puerto,version" {
		t.Fatalf("fields = %v, want exactly ok, version, lexemas, puerto", keys)
	}
	if m["ok"] != true || m["version"] != "v-test" || m["puerto"] != float64(8473) {
		t.Errorf("health = %v", m)
	}
	var n int
	if err := openOracle(t, a.Dir).QueryRow(`SELECT COUNT(*) FROM lexemas`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if m["lexemas"] != float64(n) {
		t.Errorf("lexemas = %v, catalogue has %d", m["lexemas"], n)
	}
}

type lexemasResp struct {
	Total   int `json:"total"`
	Pagina  int `json:"pagina"`
	Tamanio int `json:"tamanio"`
	Lexemas []struct {
		Forma            string         `json:"forma_canonica"`
		ER               string         `json:"expresion_regular"`
		FrecuenciaGlobal int            `json:"frecuencia_global"`
		Frecuencias      map[string]int `json:"frecuencias"`
		Variantes        int            `json:"variantes"`
	} `json:"lexemas"`
}

func getLexemas(t *testing.T, h http.Handler, query string) lexemasResp {
	t.Helper()
	rec := do(t, h, http.MethodGet, "/api/lexemas"+query, nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /api/lexemas%s: status %d: %s", query, rec.Code, rec.Body)
	}
	var r lexemasResp
	decode(t, rec, &r)
	return r
}

// CPA04: /api/lexemas with no parameters -> 200, paginated list with regular expressions and frequencies.
func TestCPA04_LexemasPaginado(t *testing.T) {
	a := bootApp(t)
	h := a.Handler()
	oracle := openOracle(t, a.Dir)

	r := getLexemas(t, h, "")
	var total int
	if err := oracle.QueryRow(`SELECT COUNT(*) FROM lexemas`).Scan(&total); err != nil {
		t.Fatal(err)
	}
	if r.Total != total || r.Pagina != 1 || r.Tamanio != 100 || len(r.Lexemas) != 100 {
		t.Fatalf("total=%d pagina=%d tamanio=%d len=%d; want total=%d pagina=1 tamanio=100 len=100", r.Total, r.Pagina, r.Tamanio, len(r.Lexemas), total)
	}
	names := categoryNames()
	for _, l := range r.Lexemas {
		if l.Forma == "" || l.ER == "" || l.Variantes < 1 {
			t.Fatalf("incomplete lexeme: %+v", l)
		}
		if len(l.Frecuencias) != 5 {
			t.Fatalf("%s: %d category frequencies, want 5: %v", l.Forma, len(l.Frecuencias), l.Frecuencias)
		}
		sum := 0
		for n, v := range l.Frecuencias {
			if !names[n] {
				t.Errorf("%s: unknown category %q", l.Forma, n)
			}
			sum += v
		}
		// I-2 (TE §6.3.3): the global frequency is the sum over the categories.
		if l.FrecuenciaGlobal != sum {
			t.Errorf("%s: frecuencia_global %d != sum of categories %d", l.Forma, l.FrecuenciaGlobal, sum)
		}
	}
	// Spot-check one entry against the catalogue read directly.
	first := r.Lexemas[0]
	var er string
	var nv, glob int
	if err := oracle.QueryRow(`SELECT expresion_regular FROM lexemas WHERE forma_canonica = ?`, first.Forma).Scan(&er); err != nil {
		t.Fatal(err)
	}
	oracle.QueryRow(`SELECT COUNT(*) FROM variantes v JOIN lexemas l ON l.id = v.lexema_id WHERE l.forma_canonica = ?`, first.Forma).Scan(&nv)
	oracle.QueryRow(`SELECT f.valor FROM frecuencias f JOIN lexemas l ON l.id = f.lexema_id WHERE l.forma_canonica = ? AND f.categoria_id IS NULL`, first.Forma).Scan(&glob)
	if first.ER != er || first.Variantes != nv || first.FrecuenciaGlobal != glob {
		t.Errorf("first entry %+v differs from the catalogue (er=%q variantes=%d global=%d)", first, er, nv, glob)
	}
	for i := 1; i < len(r.Lexemas); i++ {
		if r.Lexemas[i-1].Forma >= r.Lexemas[i].Forma {
			t.Fatalf("default order is not ascending by forma_canonica at %d", i)
		}
	}
}

func TestLexemasPaginas(t *testing.T) {
	h := bootApp(t).Handler()
	p1 := getLexemas(t, h, "?pagina=1&tamanio=10")
	p2 := getLexemas(t, h, "?pagina=2&tamanio=10")
	if len(p1.Lexemas) != 10 || len(p2.Lexemas) != 10 || p2.Pagina != 2 {
		t.Fatalf("pages: %d, %d, pagina=%d", len(p1.Lexemas), len(p2.Lexemas), p2.Pagina)
	}
	if p1.Lexemas[9].Forma >= p2.Lexemas[0].Forma {
		t.Error("page 2 does not continue page 1")
	}
	// Past the last page: empty list, same total.
	far := getLexemas(t, h, "?pagina=100000&tamanio=100")
	if len(far.Lexemas) != 0 || far.Total != p1.Total {
		t.Errorf("past the end: %d lexemes, total %d", len(far.Lexemas), far.Total)
	}
	// The page size is capped.
	if big := getLexemas(t, h, "?tamanio=100000"); big.Tamanio != 500 || len(big.Lexemas) != 500 {
		t.Errorf("tamanio cap: tamanio=%d len=%d, want 500", big.Tamanio, len(big.Lexemas))
	}
	for _, bad := range []string{"?pagina=0", "?pagina=-1", "?pagina=x", "?tamanio=0", "?tamanio=abc", "?orden=alfabetico"} {
		rec := do(t, h, http.MethodGet, "/api/lexemas"+bad, nil)
		if rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, want 400", bad, rec.Code)
		}
	}
}

func TestLexemasFiltroPorPrefijo(t *testing.T) {
	a := bootApp(t)
	h := a.Handler()
	oracle := openOracle(t, a.Dir)

	r := getLexemas(t, h, "?q=cobr&tamanio=500")
	var want int
	if err := oracle.QueryRow(`SELECT COUNT(*) FROM lexemas WHERE substr(forma_canonica, 1, 4) = 'cobr'`).Scan(&want); err != nil {
		t.Fatal(err)
	}
	if want == 0 || r.Total != want {
		t.Fatalf("q=cobr: total %d, catalogue has %d", r.Total, want)
	}
	for _, l := range r.Lexemas {
		if !strings.HasPrefix(l.Forma, "cobr") {
			t.Errorf("%q does not start with cobr", l.Forma)
		}
	}
	// Case-insensitive prefix, and no match gives an empty page, not an error.
	if up := getLexemas(t, h, "?q=COBR&tamanio=500"); up.Total != want {
		t.Errorf("q=COBR: total %d, want %d", up.Total, want)
	}
	if none := getLexemas(t, h, "?q=qqqqzzzz"); none.Total != 0 || len(none.Lexemas) != 0 {
		t.Errorf("unmatched prefix: %+v", none)
	}
	// LIKE wildcards are plain characters.
	if pct := getLexemas(t, h, "?q=%25"); pct.Total != 0 {
		t.Errorf("q=%%: total %d, want 0", pct.Total)
	}
}

// orden=variantes puts the lexemes whose regular expression has the most alternatives first
// (demo step 4 shows alternation instead of mi, me, no).
func TestLexemasOrdenes(t *testing.T) {
	a := bootApp(t)
	h := a.Handler()
	oracle := openOracle(t, a.Dir)

	v := getLexemas(t, h, "?orden=variantes&tamanio=50")
	var max int
	if err := oracle.QueryRow(`SELECT MAX(c) FROM (SELECT COUNT(*) c FROM variantes GROUP BY lexema_id)`).Scan(&max); err != nil {
		t.Fatal(err)
	}
	if v.Lexemas[0].Variantes != max || max < 2 {
		t.Errorf("first by variantes has %d, catalogue maximum is %d", v.Lexemas[0].Variantes, max)
	}
	if !strings.Contains(v.Lexemas[0].ER, "|") {
		t.Errorf("most-variant ER has no alternation: %q", v.Lexemas[0].ER)
	}
	for i := 1; i < len(v.Lexemas); i++ {
		if v.Lexemas[i-1].Variantes < v.Lexemas[i].Variantes {
			t.Fatalf("orden=variantes not descending at %d", i)
		}
	}

	f := getLexemas(t, h, "?orden=frecuencia&tamanio=50")
	var maxFreq int
	if err := oracle.QueryRow(`SELECT MAX(valor) FROM frecuencias WHERE categoria_id IS NULL`).Scan(&maxFreq); err != nil {
		t.Fatal(err)
	}
	if f.Lexemas[0].FrecuenciaGlobal != maxFreq {
		t.Errorf("first by frecuencia has %d, maximum is %d", f.Lexemas[0].FrecuenciaGlobal, maxFreq)
	}
	for i := 1; i < len(f.Lexemas); i++ {
		if f.Lexemas[i-1].FrecuenciaGlobal < f.Lexemas[i].FrecuenciaGlobal {
			t.Fatalf("orden=frecuencia not descending at %d", i)
		}
	}
}

func TestLatencia_P99DeClasificar(t *testing.T) {
	a := bootApp(t)
	srv := httptest.NewServer(a.Handler())
	defer srv.Close()

	// A comment of mean length (by characters) from the corpus, read directly.
	var texto string
	err := openOracleCorpus(t, a.Dir).QueryRow(`
		SELECT texto FROM comentarios
		ORDER BY ABS(LENGTH(texto) - (SELECT AVG(LENGTH(texto)) FROM comentarios)), id LIMIT 1`).Scan(&texto)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(map[string]string{"texto": texto})

	const n = 100
	lat := make([]time.Duration, n)
	for i := range lat {
		start := time.Now()
		resp, err := http.Post(srv.URL+"/api/clasificar", "application/json", bytes.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		var r clasificarResp
		if err := json.NewDecoder(resp.Body).Decode(&r); err != nil || resp.StatusCode != 200 {
			t.Fatalf("call %d: status %d, %v", i, resp.StatusCode, err)
		}
		resp.Body.Close()
		lat[i] = time.Since(start)
	}
	sorted := append([]time.Duration(nil), lat...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	p99 := sorted[98] // nearest rank: ceil(0.99*100) = 99th value
	var sum time.Duration
	for _, d := range lat {
		sum += d
	}
	t.Logf("P99 /api/clasificar over %d consecutive calls (%d chars): %v (median %v, max %v, mean %v)",
		n, len([]rune(texto)), p99, sorted[49], sorted[n-1], sum/n)
	if !testing.Short() && p99 >= 50*time.Millisecond {
		t.Errorf("P99 = %v, want < 50 ms (TE §6.3.7)", p99)
	}
}

func openOracleCorpus(t *testing.T, dir string) *sql.DB {
	t.Helper()
	d, err := db.OpenReadOnly(filepath.Join(dir, "corpus.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	return d
}

// min_frecuencia hides the low-frequency lexemes (mostly typos of the noise injector) so the
// catalogue page shows real morphology; without it the contract of CPA04 is unchanged.
func TestLexemasFrecuenciaMinima(t *testing.T) {
	a := bootApp(t)
	h := a.Handler()
	oracle := openOracle(t, a.Dir)

	var want int
	if err := oracle.QueryRow(`SELECT COUNT(*) FROM frecuencias WHERE categoria_id IS NULL AND valor >= 20`).Scan(&want); err != nil {
		t.Fatal(err)
	}
	r := getLexemas(t, h, "?min_frecuencia=20&orden=variantes&tamanio=500")
	if want == 0 || r.Total != want {
		t.Fatalf("min_frecuencia=20: total %d, catalogue has %d", r.Total, want)
	}
	for _, l := range r.Lexemas {
		if l.FrecuenciaGlobal < 20 {
			t.Errorf("%q has global frequency %d < 20", l.Forma, l.FrecuenciaGlobal)
		}
	}
	all := getLexemas(t, h, "")
	if zero := getLexemas(t, h, "?min_frecuencia=0"); zero.Total != all.Total {
		t.Errorf("min_frecuencia=0: total %d, want %d", zero.Total, all.Total)
	}
	if both := getLexemas(t, h, "?min_frecuencia=20&q=cobr&tamanio=500"); both.Total > r.Total {
		t.Errorf("q and min_frecuencia together: total %d exceeds %d", both.Total, r.Total)
	}
	for _, bad := range []string{"?min_frecuencia=-1", "?min_frecuencia=abc"} {
		if rec := do(t, h, http.MethodGet, "/api/lexemas"+bad, nil); rec.Code != http.StatusBadRequest {
			t.Errorf("%s: status %d, want 400", bad, rec.Code)
		}
	}
}
