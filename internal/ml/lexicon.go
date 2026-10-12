// Package ml is the runtime side of the classifier (consigna e, TE §6.4.2): it
// maps text to catalogue lexemes (Lexicon), vectorizes with the precomputed
// TF-IDF (Model.Vectorize) and scores with a linear model plus softmax
// (Model.Predict). Training happens offline in Python (ml/); the Python side
// calls cmd/tokenizar, which uses the same Lexicon, so both sides share one
// definition of a token (CPA06).
package ml

import (
	"fmt"

	"autoatend/internal/lexer"
)

// Lexicon maps surface forms to the canonical lexeme of a closed set. A form
// belongs to the lexeme whose canonical form has the same lexer.Stem; forms
// whose stem is not in the set are dropped.
type Lexicon struct {
	byStem map[string]string // lexer.Stem(canonical) -> canonical
}

// NewLexicon builds the lexicon of the given canonical forms. Two canonical
// forms with the same stem are an error: the catalogue groups by stem, so a
// well-formed set never has them.
func NewLexicon(canonicals []string) (*Lexicon, error) {
	l := &Lexicon{byStem: make(map[string]string, len(canonicals))}
	for _, c := range canonicals {
		s := lexer.Stem(c)
		if prev, dup := l.byStem[s]; dup {
			return nil, fmt.Errorf("ml: %q and %q share the stem %q", prev, c, s)
		}
		l.byStem[s] = c
	}
	return l, nil
}

// Size is the number of lexemes.
func (l *Lexicon) Size() int { return len(l.byStem) }

// Lexemes returns the canonical lexeme of every kept token of text, in order of
// appearance and with repetitions (the raw term count is the TF). Numbers and
// stopwords are dropped by lexer.Analyze; forms outside the lexicon are dropped
// here. The result is never nil.
func (l *Lexicon) Lexemes(text string) []string {
	out := []string{}
	for _, form := range lexer.Analyze(text) {
		if c, ok := l.byStem[lexer.Stem(form)]; ok {
			out = append(out, c)
		}
	}
	return out
}
