package ml_test

import (
	"encoding/json"
	"math"
	"os"
	"testing"

	"autoatend/internal/ml"
)

// Artifacts written by `make train` and versioned; make verify only reads them.
const (
	modeloPath      = "../../data/modelo.json"
	alternativoPath = "../../data/modelo_alternativo.json"
	goldPath        = "../../data/gold5.json"
	paridadPath     = "../../testdata/paridad_1000.json"
)

type paridadScores struct {
	Prob        [5]float64 `json:"prob"`
	ProbSklearn [5]float64 `json:"prob_sklearn"`
	Clase       int        `json:"clase"`
}

type paridadCase struct {
	ID          int           `json:"id"`
	Texto       string        `json:"texto"`
	Vector      [][2]float64  `json:"vector"`
	Principal   paridadScores `json:"principal"`
	Alternativo paridadScores `json:"alternativo"`
}

type paridadFile struct {
	Semilla int           `json:"semilla"`
	Casos   []paridadCase `json:"casos"`
}

func readJSON(t *testing.T, path string, v any) {
	t.Helper()
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("%s: %v (run make train)", path, err)
	}
	if err := json.Unmarshal(raw, v); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
}

func loadParidad(t *testing.T) paridadFile {
	t.Helper()
	var f paridadFile
	readJSON(t, paridadPath, &f)
	return f
}

func mustLoad(t *testing.T, path string) *ml.Model {
	t.Helper()
	m, err := ml.Load(path)
	if err != nil {
		t.Fatal(err)
	}
	return m
}

// TestCPA06_ParidadPythonGo: the Go vectorizer must reproduce the Python TF-IDF
// vector of 1,000 test comments with the same support and bitwise-equal values
// (L1 = 0), and the same predicted class. The probabilities go through math.Exp
// in Go and libm in Python, so they are compared with a tolerance and their
// maximum difference is reported.
func TestCPA06_ParidadPythonGo(t *testing.T) {
	m := mustLoad(t, modeloPath)
	f := loadParidad(t)
	if len(f.Casos) != 1000 {
		t.Fatalf("fixture has %d cases, want 1000", len(f.Casos))
	}
	var l1, maxDP float64
	classOK, supportDiff := 0, 0
	for _, c := range f.Casos {
		x := m.Vectorize(c.Texto)
		if len(x) != len(c.Vector) {
			supportDiff++
			t.Errorf("id %d: Go support has %d terms, Python %d", c.ID, len(x), len(c.Vector))
			continue
		}
		for i, term := range x {
			if float64(term.Index) != c.Vector[i][0] {
				supportDiff++
				t.Errorf("id %d: index %d vs %v", c.ID, term.Index, c.Vector[i][0])
				break
			}
			if term.Value != c.Vector[i][1] {
				t.Errorf("id %d: value at index %d differs: Go %v, Python %v", c.ID, term.Index, term.Value, c.Vector[i][1])
			}
			l1 += math.Abs(term.Value - c.Vector[i][1])
		}
		p := m.PredictSparse(x)
		if p.CategoriaID == c.Principal.Clase {
			classOK++
		} else {
			t.Errorf("id %d: class Go %d, Python %d", c.ID, p.CategoriaID, c.Principal.Clase)
		}
		for k := range p.Confianzas {
			maxDP = math.Max(maxDP, math.Abs(p.Confianzas[k]-c.Principal.Prob[k]))
		}
	}
	t.Logf("PARIDAD n=%d L1=%g support_diff=%d class_agree=%d max_dprob=%.3e", len(f.Casos), l1, supportDiff, classOK, maxDP)
	if l1 != 0 || supportDiff != 0 {
		t.Errorf("vector parity broken: L1 = %g, support differences = %d", l1, supportDiff)
	}
	if maxDP > 1e-12 {
		t.Errorf("max |dprob| = %g, want <= 1e-12", maxDP)
	}
}

// TestModeloLineal_EquivalenciaSklearn: the linear scorer plus softmax of the
// exported weights reproduces predict_proba of scikit-learn for both models, so
// Naive Bayes (w = feature_log_prob_, b = class_log_prior_) and logistic
// regression share one inference path.
func TestModeloLineal_EquivalenciaSklearn(t *testing.T) {
	f := loadParidad(t)
	for _, tc := range []struct {
		name, path string
		pick       func(paridadCase) paridadScores
	}{
		{"principal", modeloPath, func(c paridadCase) paridadScores { return c.Principal }},
		{"alternativo", alternativoPath, func(c paridadCase) paridadScores { return c.Alternativo }},
	} {
		m := mustLoad(t, tc.path)
		var maxD float64
		for _, c := range f.Casos {
			s := tc.pick(c)
			p := m.Predict(c.Texto)
			for k := range p.Confianzas {
				maxD = math.Max(maxD, math.Abs(p.Confianzas[k]-s.ProbSklearn[k]))
			}
			if p.CategoriaID != s.Clase {
				t.Errorf("%s id %d: class Go %d, sklearn %d", tc.name, c.ID, p.CategoriaID, s.Clase)
			}
		}
		t.Logf("%s (%s/%s): max |p_go - p_sklearn| = %.3e", tc.name, m.Familia, m.Estrategia, maxD)
		if maxD > 1e-9 {
			t.Errorf("%s: linear scorer differs from sklearn predict_proba by %g", tc.name, maxD)
		}
	}
}

func TestModelo_RealLoadsAndFamilies(t *testing.T) {
	a := mustLoad(t, modeloPath)
	b := mustLoad(t, alternativoPath)
	if a.Familia == b.Familia {
		t.Errorf("both models are %q; the alternative must be the best of the other family", a.Familia)
	}
	for _, m := range []*ml.Model{a, b} {
		if m.Familia != "nbm" && m.Familia != "rls" {
			t.Errorf("familia = %q", m.Familia)
		}
		if len(m.Vocabulario) == 0 || len(m.Vocabulario) != len(m.IDF) {
			t.Errorf("vocabulary %d, idf %d", len(m.Vocabulario), len(m.IDF))
		}
		if m.Normalizacion != "l2" {
			t.Errorf("normalizacion = %q", m.Normalizacion)
		}
	}
}

type goldItem struct {
	ID          int     `json:"id"`
	Texto       string  `json:"texto"`
	CategoriaID int     `json:"categoria_id"`
	Categoria   string  `json:"categoria"`
	Confianza   float64 `json:"confianza_python"`
}

// TestGold5_ClasificadosEnGo: the five demo comments are classified correctly,
// with confidence >= 0.80, by the Go inference of modelo.json.
func TestGold5_ClasificadosEnGo(t *testing.T) {
	m := mustLoad(t, modeloPath)
	var gold []goldItem
	readJSON(t, goldPath, &gold)
	if len(gold) != 5 {
		t.Fatalf("gold5 has %d items", len(gold))
	}
	seen := map[int]bool{}
	for _, g := range gold {
		seen[g.CategoriaID] = true
		p := m.Predict(g.Texto)
		if p.CategoriaID != g.CategoriaID {
			t.Errorf("gold %d: Go predicts %d (%s), expected %d (%s)", g.ID, p.CategoriaID, p.Categoria, g.CategoriaID, g.Categoria)
		}
		if c := p.Confianzas[g.CategoriaID-1]; c < 0.80 {
			t.Errorf("gold %d: confidence %.3f < 0.80", g.ID, c)
		}
		if len(p.LexemasClave) == 0 {
			t.Errorf("gold %d: no key lexemes", g.ID)
		}
	}
	if len(seen) != 5 {
		t.Errorf("gold5 covers %d categories, want 5", len(seen))
	}
}

// TestCPM01_MacroF1 and TestCPM02_RecallReconocimientos report the test-set
// metrics embedded in modelo.json against the TE thresholds (0.85 and 0.70).
// Per PLAN.md §1.3 thresholds are measured and reported, not asserted: a
// failing threshold is logged, never tuned away.
func TestCPM01_MacroF1(t *testing.T) {
	m := mustLoad(t, modeloPath)
	for name, e := range map[string]ml.Eval{"test": m.Metricas.Test, "test_sin_duplicados": m.Metricas.TestSinDuplicados, "test_sin_ambiguos": m.Metricas.TestSinAmbiguos} {
		if e.N == 0 || e.MacroF1 <= 0 || e.MacroF1 > 1 {
			t.Errorf("%s: metrics missing or out of range: %+v", name, e)
		}
		verdict := "CUMPLE"
		if e.MacroF1 < 0.85 {
			verdict = "NO CUMPLE"
		}
		t.Logf("CPM01 %s: n=%d Macro-F1=%.4f (umbral 0.85): %s", name, e.N, e.MacroF1, verdict)
	}
}

func TestCPM02_RecallReconocimientos(t *testing.T) {
	m := mustLoad(t, modeloPath)
	for name, e := range map[string]ml.Eval{"test": m.Metricas.Test, "test_sin_duplicados": m.Metricas.TestSinDuplicados, "test_sin_ambiguos": m.Metricas.TestSinAmbiguos} {
		if e.N == 0 || e.RecallReconocimientos <= 0 || e.RecallReconocimientos > 1 {
			t.Errorf("%s: metrics missing or out of range: %+v", name, e)
		}
		verdict := "CUMPLE"
		if e.RecallReconocimientos < 0.70 {
			verdict = "NO CUMPLE"
		}
		t.Logf("CPM02 %s: n=%d recall Reconocimientos=%.4f (umbral 0.70): %s", name, e.N, e.RecallReconocimientos, verdict)
	}
}
