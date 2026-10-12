package ml_test

import (
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/ml"
)

// miniModel is a hand-written model: 7 lexemes, exact binary fractions for the
// IDF, sparse weights per class. The expected values in the tests are computed
// by hand from these numbers.
const miniModel = `{
"version":1,
"familia":"rls",
"estrategia":"balanceo",
"hiperparametros":{"C":1.0},
"categorias":[{"id":1,"nombre":"Falta de Servicio"},{"id":2,"nombre":"Cobros Indebidos"},{"id":3,"nombre":"Problemas de Facturación"},{"id":4,"nombre":"Quejas"},{"id":5,"nombre":"Reconocimientos"}],
"tf":"conteo",
"idf_formula":"ln((1+N)/(1+df))+1",
"normalizacion":"l2",
"documentos_train":100,
"vocabulario":["factura","cobraron","servicio","excelente","mal","cliente","llamada"],
"idf":[1.5,2.0,1.25,3.0,1.0,2.0,1.0],
"pesos":[
 [[0,0.5],[2,1.0]],
 [[0,2.0],[1,1.5]],
 [[0,1.0]],
 [[0,1.0],[1,-2.0],[2,0.5],[3,3.0],[4,-4.0],[5,0.25],[6,3.0]],
 [[3,2.5]]
],
"sesgo":[0.1,-0.2,0.3,0.0,-0.5],
"metricas":{"test":{"n":10,"macro_f1":0.9,"recall_reconocimientos":0.8},"test_sin_duplicados":{"n":7,"macro_f1":0.8,"recall_reconocimientos":0.7}}
}`

func writeMini(t *testing.T, body string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), "modelo.json")
	if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func loadMini(t *testing.T) *ml.Model {
	t.Helper()
	m, err := ml.Load(writeMini(t, miniModel))
	if err != nil {
		t.Fatal(err)
	}
	return m
}

func TestModelo_Load(t *testing.T) {
	m := loadMini(t)
	if m.Familia != "rls" || m.Estrategia != "balanceo" || m.Version != 1 {
		t.Errorf("header = %q %q %d", m.Familia, m.Estrategia, m.Version)
	}
	if got := len(m.Vocabulario); got != 7 {
		t.Errorf("vocabulary size = %d, want 7", got)
	}
	cats := m.Categorias
	if len(cats) != 5 || cats[0].ID != 1 || cats[4].Nombre != "Reconocimientos" {
		t.Errorf("categories = %+v", cats)
	}
	if m.Metricas.Test.MacroF1 != 0.9 || m.Metricas.TestSinDuplicados.RecallReconocimientos != 0.7 {
		t.Errorf("metrics = %+v", m.Metricas)
	}
}

func TestModelo_LoadErrors(t *testing.T) {
	cases := []struct{ name, from, to, want string }{
		{"version", `"version":1`, `"version":2`, "version"},
		{"idf length", `"idf":[1.5,2.0,1.25,3.0,1.0,2.0,1.0]`, `"idf":[1.5]`, "idf"},
		{"weight index", `[[3,2.5]]`, `[[9,2.5]]`, "index"},
		{"weight order", `[[0,0.5],[2,1.0]]`, `[[2,1.0],[0,0.5]]`, "ascending"},
		{"bias length", `"sesgo":[0.1,-0.2,0.3,0.0,-0.5]`, `"sesgo":[0.1]`, "sesgo"},
		{"stem collision", `"factura","cobraron"`, `"factura","facturas"`, "stem"},
		{"normalization", `"normalizacion":"l2"`, `"normalizacion":"l1"`, "normalizacion"},
		{"category ids", `{"id":5,"nombre":"Reconocimientos"}`, `{"id":7,"nombre":"Reconocimientos"}`, "categor"},
	}
	for _, c := range cases {
		if !strings.Contains(miniModel, c.from) {
			t.Fatalf("%s: fixture lacks %q", c.name, c.from)
		}
		p := writeMini(t, strings.Replace(miniModel, c.from, c.to, 1))
		_, err := ml.Load(p)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: err = %v, want one containing %q", c.name, err, c.want)
		}
	}
	if _, err := ml.Load(filepath.Join(t.TempDir(), "missing.json")); err == nil {
		t.Error("Load of a missing file succeeded")
	}
	if _, err := ml.Load(writeMini(t, "{not json")); err == nil {
		t.Error("Load of malformed JSON succeeded")
	}
}

func TestModelo_Categorias_MatchDB(t *testing.T) {
	m := loadMini(t)
	for i, c := range db.Categories {
		if m.Categorias[i].ID != c.ID || m.Categorias[i].Nombre != c.Name {
			t.Errorf("category %d = %+v, db has %+v", i, m.Categorias[i], c)
		}
	}
}

func TestVectorize_TfIdfL2(t *testing.T) {
	m := loadMini(t)
	// factura x2 (also as "facturación"), cobraron x1; "veces", "dos" and "la" are not in the vocabulary.
	x := m.Vectorize("Me cobraron la factura dos veces facturación")
	if len(x) != 2 || x[0].Index != 0 || x[1].Index != 1 {
		t.Fatalf("support = %+v, want indexes 0,1 ascending", x)
	}
	// tf*idf = 2*1.5 = 3 and 1*2 = 2; L2 norm sqrt(13).
	n := math.Sqrt(3*3 + 2*2)
	if x[0].Value != 3/n || x[1].Value != 2/n {
		t.Errorf("values = %v %v, want %v %v", x[0].Value, x[1].Value, 3/n, 2/n)
	}
	var s float64
	for _, tm := range x {
		s += tm.Value * tm.Value
	}
	if math.Abs(s-1) > 1e-15 {
		t.Errorf("squared norm = %v, want 1", s)
	}
}

func TestVectorize_Empty(t *testing.T) {
	m := loadMini(t)
	for _, txt := range []string{"", "   ", "el de la y", "xyzzy qwerty 123"} {
		if x := m.Vectorize(txt); len(x) != 0 {
			t.Errorf("Vectorize(%q) = %+v, want empty", txt, x)
		}
	}
}

func softmax(z [5]float64) [5]float64 {
	mx := z[0]
	for _, v := range z {
		mx = math.Max(mx, v)
	}
	var e [5]float64
	var s float64
	for i, v := range z {
		e[i] = math.Exp(v - mx)
		s += e[i]
	}
	for i := range e {
		e[i] /= s
	}
	return e
}

func TestPredict_LinearScoreAndSoftmax(t *testing.T) {
	m := loadMini(t)
	p := m.Predict("Me cobraron la factura dos veces facturación")
	n := math.Sqrt(13)
	x0, x1 := 3/n, 2/n
	z := [5]float64{
		0.1 + 0.5*x0,
		-0.2 + 2.0*x0 + 1.5*x1,
		0.3 + 1.0*x0,
		0.0 + 1.0*x0 + -2.0*x1,
		-0.5,
	}
	want := softmax(z)
	var sum float64
	for k := range want {
		if math.Abs(p.Confianzas[k]-want[k]) > 1e-14 {
			t.Errorf("confianza[%d] = %v, want %v", k, p.Confianzas[k], want[k])
		}
		sum += p.Confianzas[k]
	}
	if math.Abs(sum-1) > 1e-12 {
		t.Errorf("confidences sum to %v, want 1", sum)
	}
	if p.CategoriaID != 2 || p.Categoria != "Cobros Indebidos" {
		t.Errorf("prediction = %d %q, want 2 Cobros Indebidos", p.CategoriaID, p.Categoria)
	}
	if p.Vacio {
		t.Error("Vacio = true for a text with vocabulary lexemes")
	}
}

func TestPredict_SoftmaxStableWithHugeScores(t *testing.T) {
	body := strings.Replace(miniModel, `"sesgo":[0.1,-0.2,0.3,0.0,-0.5]`, `"sesgo":[1000,1001,-1000,0,5]`, 1)
	m, err := ml.Load(writeMini(t, body))
	if err != nil {
		t.Fatal(err)
	}
	p := m.Predict("")
	var sum float64
	for _, c := range p.Confianzas {
		if math.IsNaN(c) || math.IsInf(c, 0) {
			t.Fatalf("non-finite confidence in %v", p.Confianzas)
		}
		sum += c
	}
	if math.Abs(sum-1) > 1e-12 || p.CategoriaID != 2 {
		t.Errorf("sum = %v class = %d", sum, p.CategoriaID)
	}
}

func TestPredict_EmptyTextIsPriorOnly(t *testing.T) {
	m := loadMini(t)
	want := softmax([5]float64{0.1, -0.2, 0.3, 0.0, -0.5})
	for _, txt := range []string{"", "el de la y", "xyzzy 42"} {
		p := m.Predict(txt)
		if !p.Vacio {
			t.Errorf("Predict(%q).Vacio = false", txt)
		}
		if p.Confianzas != want {
			t.Errorf("Predict(%q).Confianzas = %v, want %v", txt, p.Confianzas, want)
		}
		if p.CategoriaID != 3 {
			t.Errorf("Predict(%q) class = %d, want 3 (largest bias)", txt, p.CategoriaID)
		}
		if len(p.LexemasClave) != 0 {
			t.Errorf("Predict(%q).LexemasClave = %+v, want none", txt, p.LexemasClave)
		}
	}
}

func TestPredict_LexemasClaveTop5(t *testing.T) {
	m := loadMini(t)
	p := m.Predict("factura cobraron servicio excelente mal cliente llamada")
	if p.CategoriaID != 4 {
		t.Fatalf("class = %d, want 4", p.CategoriaID)
	}
	// Contributions of the winning class are proportional to w*idf =
	// [1.5, -4, 0.625, 9, -4, 0.5, 3]; by |.| descending, ties by index:
	// excelente(3), cobraron(1), mal(4), llamada(6), factura(0).
	wantNames := []string{"excelente", "cobraron", "mal", "llamada", "factura"}
	if len(p.LexemasClave) != 5 {
		t.Fatalf("got %d key lexemes, want 5", len(p.LexemasClave))
	}
	for i, lk := range p.LexemasClave {
		if lk.Lexema != wantNames[i] {
			t.Errorf("lexema[%d] = %q, want %q", i, lk.Lexema, wantNames[i])
		}
	}
	n := math.Sqrt(1.5*1.5 + 2*2 + 1.25*1.25 + 3*3 + 1*1 + 2*2 + 1*1)
	if got, want := p.LexemasClave[0].Contribucion, 3.0*(3.0/n); math.Abs(got-want) > 1e-15 {
		t.Errorf("contribution of excelente = %v, want %v", got, want)
	}
	if p.LexemasClave[1].Contribucion >= 0 {
		t.Errorf("cobraron has a negative weight, contribution = %v", p.LexemasClave[1].Contribucion)
	}
	abs := make([]float64, len(p.LexemasClave))
	for i, lk := range p.LexemasClave {
		abs[i] = math.Abs(lk.Contribucion)
	}
	if !sort.SliceIsSorted(abs, func(i, j int) bool { return abs[i] > abs[j] }) {
		t.Errorf("contributions not in |.| descending order: %v", abs)
	}
}

func TestPredict_LexemasClaveSkipsZeroWeights(t *testing.T) {
	m := loadMini(t)
	// Class 3 (Problemas de Facturación) only weighs "factura": "cliente" has weight 0 there.
	p := m.Predict("factura factura factura cliente")
	for _, lk := range p.LexemasClave {
		if lk.Contribucion == 0 {
			t.Errorf("zero contribution listed: %+v", lk)
		}
	}
}

func TestPredict_Deterministic(t *testing.T) {
	m := loadMini(t)
	a := m.Predict("Me cobraron la factura dos veces")
	b := m.Predict("Me cobraron la factura dos veces")
	if a.Confianzas != b.Confianzas || a.CategoriaID != b.CategoriaID {
		t.Error("two predictions of the same text differ")
	}
}
