package lexer

import "sort"

// prefixRunes is the length of the index key (TE §6.2.5).
const prefixRunes = 3

// PrefixIndex maps the first three runes of a word to the canonical forms of
// the lexemes that start with them (TE §6.2.5). The matcher uses it to
// simulate the AFN-ε of only a few lexemes instead of the whole catalogue.
// All the variants of a lexeme share the first three runes of its canonical
// form (a stem keeps at least three runes, see Stem), so a token and its
// lexeme always land under the same key.
type PrefixIndex struct {
	byPrefix map[string][]string
}

func prefixKey(s string) string {
	n := 0
	for i := range s {
		if n == prefixRunes {
			return s[:i]
		}
		n++
	}
	return s
}

// BuildPrefixIndex indexes the canonical forms of the lexemes.
func BuildPrefixIndex(entries []*Entry) *PrefixIndex {
	idx := &PrefixIndex{byPrefix: make(map[string][]string)}
	for _, e := range entries {
		c := e.Canonical()
		k := prefixKey(c)
		idx.byPrefix[k] = append(idx.byPrefix[k], c)
	}
	for _, list := range idx.byPrefix {
		sort.Strings(list)
	}
	return idx
}

// Candidates returns, sorted, the canonical forms whose first three runes are
// those of token. It returns a copy; nil when there is none.
func (p *PrefixIndex) Candidates(token string) []string {
	list := p.byPrefix[prefixKey(token)]
	if len(list) == 0 {
		return nil
	}
	return append([]string(nil), list...)
}
