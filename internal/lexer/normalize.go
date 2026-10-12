// Package lexer is the lexical front end of AutoAtend (TE §6.2.2): the
// normalizer, the 4-state DFA tokenizer, the closed-class stoplist, the light
// Spanish stemmer that groups surface variants into lexemes, the symbol table
// and the regular-expression generator. Analyze is the single lexical
// definition shared by the sweep, the simulator and the runtime inference.
// Nothing here uses the stdlib regexp package or an NLP library.
package lexer

import "unicode"

// Alphabet is the TE §6.2.1 alphabet (lowercase letters of Spanish), documented
// for the run parameters. Digits and the space separator complete the input
// alphabet of the tokenizer.
const Alphabet = "a-záéíóúüñ"

// composeKey is a base letter followed by a combining mark.
type composeKey struct{ base, mark rune }

const (
	markAcute     = '\u0301'
	markDiaeresis = '\u0308'
	markTilde     = '\u0303'
)

// compose is the restricted NFC of the normalizer: only the 14 pairs that
// build the accented letters of the alphabet, in lowercase and uppercase.
var compose = map[composeKey]rune{
	{'a', markAcute}: 'á', {'e', markAcute}: 'é', {'i', markAcute}: 'í', {'o', markAcute}: 'ó', {'u', markAcute}: 'ú',
	{'u', markDiaeresis}: 'ü', {'n', markTilde}: 'ñ',
	{'A', markAcute}: 'Á', {'E', markAcute}: 'É', {'I', markAcute}: 'Í', {'O', markAcute}: 'Ó', {'U', markAcute}: 'Ú',
	{'U', markDiaeresis}: 'Ü', {'N', markTilde}: 'Ñ',
}

// isLetter reports whether r is a (lowercase) letter of the alphabet.
func isLetter(r rune) bool {
	switch {
	case r >= 'a' && r <= 'z':
		return true
	case r == 'á', r == 'é', r == 'í', r == 'ó', r == 'ú', r == 'ü', r == 'ñ':
		return true
	}
	return false
}

func isDigit(r rune) bool { return r >= '0' && r <= '9' }

// Normalize lowercases s, composes the combining marks of the alphabet
// (base + U+0301/U+0308/U+0303 into the precomposed letter) and replaces every
// character outside LETRA, DIGITO and SEPARADOR (punctuation, symbols, emoji,
// letters of other alphabets, stray combining marks) with one space. Every
// whitespace character becomes a space. It runs in O(|s|) and its output
// contains only the alphabet letters, ASCII digits and ' ', so it is idempotent.
func Normalize(s string) string {
	out := make([]rune, 0, len(s))
	// prev is the raw rune that produced the last output rune, and composable
	// tells whether that output rune is still open for a combining mark.
	var prev rune
	composable := false
	for _, r := range s {
		if r == markAcute || r == markDiaeresis || r == markTilde {
			if composable {
				if c, ok := compose[composeKey{prev, r}]; ok {
					out[len(out)-1] = unicode.ToLower(c)
					composable = false
					continue
				}
			}
			out = append(out, ' ')
			composable = false
			continue
		}
		lr := unicode.ToLower(r)
		switch {
		case isLetter(lr) || isDigit(lr):
			out = append(out, lr)
			prev, composable = r, true
		default:
			out = append(out, ' ')
			composable = false
		}
	}
	return string(out)
}
