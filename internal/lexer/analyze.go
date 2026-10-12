package lexer

// Analyze is the single lexical definition of AutoAtend (TE §6.2.2): Normalize,
// Tokenize, drop the numbers (recognized by the DFA but not catalogued), drop
// the forms of the stoplist, and return the surface forms in order of
// appearance. The sweep, the simulator and the runtime inference all call it,
// so the catalogue and the classifier can never disagree on what a word is.
// An empty text, or one made only of stopwords, yields an empty slice.
func Analyze(text string) []string {
	var forms []string
	for _, tok := range Tokenize(Normalize(text)) {
		if tok.Kind != Palabra {
			continue
		}
		if _, stop := Lookup(tok.Text); stop {
			continue
		}
		forms = append(forms, tok.Text)
	}
	return forms
}
