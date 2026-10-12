package ml

import (
	"math"
	"sort"
)

// Term is one non-zero component of a sparse vector.
type Term struct {
	Index int
	Value float64
}

// Sparse is a sparse vector with strictly ascending indexes.
type Sparse []Term

// LexemaClave is a lexeme that pushed the winning class, with its signed
// contribution w_{k,i}·x_i to that class score.
type LexemaClave struct {
	Lexema       string  `json:"lexema"`
	Contribucion float64 `json:"contribucion"`
}

// MaxLexemasClave is the number of key lexemes reported (TE §6.4.1).
const MaxLexemasClave = 5

// Prediction is the classifier answer.
type Prediction struct {
	CategoriaID  int                `json:"categoria_id"`
	Categoria    string             `json:"categoria"`
	Confianzas   [NumClases]float64 `json:"confianzas"`
	LexemasClave []LexemaClave      `json:"lexemas_clave"`
	// Vacio is true when no token of the text is in the model vocabulary
	// (empty text, only stopwords or numbers, or only unknown words). The
	// vector is then zero, the scores are the biases and the answer is the
	// prior of the model, not evidence about the text.
	Vacio bool `json:"vacio"`
}

// Vectorize turns text into the L2-normalized TF-IDF vector of the model: raw
// term counts of the vocabulary lexemes times the stored IDF. The squared norm
// is summed in ascending index order with plain float64 operations (the
// explicit float64 conversions forbid fused multiply-add), so the Python
// trainer reproduces every value bit for bit.
func (m *Model) Vectorize(text string) Sparse {
	counts := make(map[int]int)
	for _, lx := range m.lex.Lexemes(text) {
		counts[m.pos[lx]]++
	}
	if len(counts) == 0 {
		return Sparse{}
	}
	x := make(Sparse, 0, len(counts))
	for i, c := range counts {
		x = append(x, Term{Index: i, Value: float64(float64(c) * m.IDF[i])})
	}
	sort.Slice(x, func(a, b int) bool { return x[a].Index < x[b].Index })
	var sumsq float64
	for _, t := range x {
		sumsq += float64(t.Value * t.Value)
	}
	norm := math.Sqrt(sumsq)
	for i := range x {
		x[i].Value /= norm
	}
	return x
}

// Predict classifies text.
func (m *Model) Predict(text string) Prediction { return m.PredictSparse(m.Vectorize(text)) }

// PredictSparse scores an already vectorized text: z_k = b_k + sum over the
// support of w_{k,i}·x_i (ascending index order), then softmax with the
// maximum subtracted. Ties between classes go to the smaller category id. Key
// lexemes are the (up to) five support terms with the largest |w·x| for the
// winning class, ties broken by the smaller index; terms whose weight in that
// class is zero are not listed.
func (m *Model) PredictSparse(x Sparse) Prediction {
	var z [NumClases]float64
	for k := range z {
		s := m.Sesgo[k]
		for _, t := range x {
			s += float64(m.w[k][t.Index] * t.Value)
		}
		z[k] = s
	}
	mx := z[0]
	for _, v := range z[1:] {
		mx = math.Max(mx, v)
	}
	var p Prediction
	var sum float64
	for k, v := range z {
		p.Confianzas[k] = math.Exp(v - mx)
		sum += p.Confianzas[k]
	}
	best := 0
	for k := range p.Confianzas {
		p.Confianzas[k] /= sum
		if p.Confianzas[k] > p.Confianzas[best] {
			best = k
		}
	}
	p.CategoriaID = m.Categorias[best].ID
	p.Categoria = m.Categorias[best].Nombre
	p.Vacio = len(x) == 0
	p.LexemasClave = m.keyLexemes(best, x)
	return p
}

func (m *Model) keyLexemes(k int, x Sparse) []LexemaClave {
	type cand struct {
		i int
		c float64
	}
	var cs []cand
	for _, t := range x {
		if c := float64(m.w[k][t.Index] * t.Value); c != 0 {
			cs = append(cs, cand{t.Index, c})
		}
	}
	// x is in ascending index order and the sort is stable, so ties keep the smaller index first.
	sort.SliceStable(cs, func(a, b int) bool { return math.Abs(cs[a].c) > math.Abs(cs[b].c) })
	if len(cs) > MaxLexemasClave {
		cs = cs[:MaxLexemasClave]
	}
	out := make([]LexemaClave, len(cs))
	for i, c := range cs {
		out[i] = LexemaClave{Lexema: m.Vocabulario[c.i], Contribucion: c.c}
	}
	return out
}
