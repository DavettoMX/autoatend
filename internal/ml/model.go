package ml

import (
	"encoding/json"
	"fmt"
	"os"

	"autoatend/internal/db"
)

// SupportedVersion is the only modelo.json format this package reads.
const SupportedVersion = 1

// NumClases is the number of categories (db.Categories).
const NumClases = 5

// Categoria is one label of the model.
type Categoria struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

// Eval is the evaluation of a model on a set of comments (TE §5.2.2). Only the
// fields the runtime and the tests read are decoded; modelo.json carries more.
type Eval struct {
	N                     int     `json:"n"`
	MacroF1               float64 `json:"macro_f1"`
	RecallReconocimientos float64 `json:"recall_reconocimientos"`
}

// Metricas are the metrics embedded in modelo.json.
type Metricas struct {
	Test              Eval `json:"test"`
	TestSinDuplicados Eval `json:"test_sin_duplicados"`
	TestSinAmbiguos   Eval `json:"test_sin_ambiguos"` // without the K2 collision subset (data/corpus_ambiguos.json)
}

// Model is a trained classifier: TF-IDF features and a linear scorer
// z_k = b_k + sum_i w_{k,i} x_i for both families (Naive Bayes uses
// w = feature_log_prob_, b = class_log_prior_; logistic regression uses
// coef_, intercept_), so there is one inference path.
type Model struct {
	Version         int
	Familia         string // "nbm" or "rls"
	Estrategia      string
	Hiperparametros json.RawMessage
	Categorias      []Categoria
	TF              string
	IDFFormula      string
	Normalizacion   string
	DocumentosTrain int
	Vocabulario     []string
	IDF             []float64
	Sesgo           [NumClases]float64
	Metricas        Metricas

	lex *Lexicon
	pos map[string]int // canonical -> vocabulary index
	w   [NumClases][]float64
}

// modelFile mirrors the JSON layout of modelo.json.
type modelFile struct {
	Version         int             `json:"version"`
	Familia         string          `json:"familia"`
	Estrategia      string          `json:"estrategia"`
	Hiperparametros json.RawMessage `json:"hiperparametros"`
	Categorias      []Categoria     `json:"categorias"`
	TF              string          `json:"tf"`
	IDFFormula      string          `json:"idf_formula"`
	Normalizacion   string          `json:"normalizacion"`
	DocumentosTrain int             `json:"documentos_train"`
	Vocabulario     []string        `json:"vocabulario"`
	IDF             []float64       `json:"idf"`
	Pesos           [][][2]float64  `json:"pesos"`
	Sesgo           []float64       `json:"sesgo"`
	Metricas        Metricas        `json:"metricas"`
}

// Load reads and validates a modelo.json. The IDF values are used exactly as
// stored: they were computed once, in Python, and are never recomputed here.
func Load(path string) (*Model, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("ml: %w", err)
	}
	var f modelFile
	if err := json.Unmarshal(raw, &f); err != nil {
		return nil, fmt.Errorf("ml: %s: %w", path, err)
	}
	return build(f, path)
}

func build(f modelFile, path string) (*Model, error) {
	fail := func(format string, a ...any) (*Model, error) {
		return nil, fmt.Errorf("ml: %s: "+format, append([]any{path}, a...)...)
	}
	if f.Version != SupportedVersion {
		return fail("unsupported version %d (want %d)", f.Version, SupportedVersion)
	}
	if f.Familia != "nbm" && f.Familia != "rls" {
		return fail("unknown familia %q", f.Familia)
	}
	if f.Normalizacion != "l2" {
		return fail("normalizacion %q is not supported (want l2)", f.Normalizacion)
	}
	if len(f.Categorias) != NumClases {
		return fail("%d categorias, want %d", len(f.Categorias), NumClases)
	}
	for i, c := range db.Categories {
		if f.Categorias[i].ID != c.ID || f.Categorias[i].Nombre != c.Nombre {
			return fail("categoria %d is %+v, db.Categories has {%d %s}", i, f.Categorias[i], c.ID, c.Nombre)
		}
	}
	if len(f.IDF) != len(f.Vocabulario) {
		return fail("idf has %d values for %d vocabulary entries", len(f.IDF), len(f.Vocabulario))
	}
	if len(f.Sesgo) != NumClases || len(f.Pesos) != NumClases {
		return fail("sesgo has %d and pesos %d entries, want %d each", len(f.Sesgo), len(f.Pesos), NumClases)
	}
	lex, err := NewLexicon(f.Vocabulario)
	if err != nil {
		return fail("vocabulary: %v", err)
	}
	m := &Model{
		Version: f.Version, Familia: f.Familia, Estrategia: f.Estrategia,
		Hiperparametros: f.Hiperparametros, Categorias: f.Categorias, TF: f.TF,
		IDFFormula: f.IDFFormula, Normalizacion: f.Normalizacion,
		DocumentosTrain: f.DocumentosTrain, Vocabulario: f.Vocabulario, IDF: f.IDF,
		Metricas: f.Metricas, lex: lex, pos: make(map[string]int, len(f.Vocabulario)),
	}
	copy(m.Sesgo[:], f.Sesgo)
	for i, v := range f.Vocabulario {
		m.pos[v] = i
	}
	for k, ws := range f.Pesos {
		m.w[k] = make([]float64, len(f.Vocabulario))
		prev := -1
		for _, p := range ws {
			idx := int(p[0])
			if float64(idx) != p[0] || idx < 0 || idx >= len(f.Vocabulario) {
				return fail("pesos[%d]: index %v out of range [0,%d)", k, p[0], len(f.Vocabulario))
			}
			if idx <= prev {
				return fail("pesos[%d]: indexes must be strictly ascending (%d after %d)", k, idx, prev)
			}
			prev = idx
			m.w[k][idx] = p[1]
		}
	}
	return m, nil
}
