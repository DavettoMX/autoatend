package lexer

import (
	"strings"
	"unicode/utf8"
)

// StemVersion identifies the suffix list, recorded in corridas.parametros.
const StemVersion = "v1"

// minRoot is the minimum number of runes a stem keeps after a suffix is removed.
const minRoot = 3

// suffixes is the light Spanish suffix list, ordered from longest to shortest
// (by runes) so that the first match is the longest one. There is no accent
// folding: "facturación" and "facturacion" are different surface forms.
var suffixes = []string{
	"aciones", "iciones", "amiento", "imiento",
	"ación", "ición", "mente", "ieron", "iendo",
	"aron", "ando", "ados", "adas", "idos", "idas",
	"ado", "ada", "ido", "ida",
	"ar", "er", "ir", "as", "es", "os",
	"a", "e", "o", "s",
}

// Stem returns the grouping key of a normalized surface form: the form without
// the first (longest) suffix that matches, provided the root keeps at least
// minRoot runes. If that suffix would leave a shorter root, the form is
// returned unchanged; no shorter suffix is tried. The stem is only a grouping
// criterion (TE §6.2.2 does not define one, ADR-003); the lexeme is named by
// its most frequent surface variant, not by the stem.
func Stem(form string) string {
	for _, suf := range suffixes {
		if !strings.HasSuffix(form, suf) {
			continue
		}
		root := form[:len(form)-len(suf)]
		if utf8.RuneCountInString(root) >= minRoot {
			return root
		}
		return form
	}
	return form
}
