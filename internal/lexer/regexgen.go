package lexer

import (
	"sort"
	"strings"
)

// metachars are the operators of the TE §6.2.3 grammar (| * + ? ( ) [ ]), the
// escape itself and '-', which the generator escapes so that a variant can
// never be read as syntax.
const metachars = `\|*+?()[]-`

// escape prefixes every metacharacter of s with '\'.
func escape(s string) string {
	if !strings.ContainsAny(s, metachars) {
		return s
	}
	var b strings.Builder
	for _, r := range s {
		if strings.ContainsRune(metachars, r) {
			b.WriteByte('\\')
		}
		b.WriteRune(r)
	}
	return b.String()
}

// GenerateER builds the regular expression of a lexeme from its observed
// surface variants (TE §6.2.2, reading B: the ER captures the morphological
// variants of the lexeme). The variants are sorted bytewise and deduplicated;
// the longest common prefix (in runes) is factored out and the suffixes become
// an alternation: factur(a|ación|ado|ando|as). The TE grammar has no empty
// alternative, so when a variant is itself the common prefix the prefix is
// shortened, one rune at a time, until no suffix is empty. With a single
// variant the result is that form, escaped; with an empty prefix it is a bare
// alternation (v1|v2). Empty variants are ignored; no variants give "".
func GenerateER(variants []string) string {
	vs := make([]string, 0, len(variants))
	for _, v := range variants {
		if v != "" {
			vs = append(vs, v)
		}
	}
	sort.Strings(vs)
	out := vs[:0]
	for i, v := range vs {
		if i == 0 || v != vs[i-1] {
			out = append(out, v)
		}
	}
	vs = out
	switch len(vs) {
	case 0:
		return ""
	case 1:
		return escape(vs[0])
	}

	// With the variants sorted, the common prefix of all is the common prefix
	// of the first and the last.
	prefix := []rune(commonPrefix(vs[0], vs[len(vs)-1]))
	// A suffix is empty iff the variant is the prefix itself; only the first
	// variant can be that short, since the list is sorted and deduplicated.
	for len(prefix) > 0 && len(prefix) == len([]rune(vs[0])) {
		prefix = prefix[:len(prefix)-1]
	}
	p := string(prefix)

	var b strings.Builder
	b.WriteString(escape(p))
	b.WriteByte('(')
	for i, v := range vs {
		if i > 0 {
			b.WriteByte('|')
		}
		b.WriteString(escape(v[len(p):]))
	}
	b.WriteByte(')')
	return b.String()
}

// commonPrefix returns the longest common prefix of a and b, cut at a rune boundary.
func commonPrefix(a, b string) string {
	ra, rb := []rune(a), []rune(b)
	n := 0
	for n < len(ra) && n < len(rb) && ra[n] == rb[n] {
		n++
	}
	return string(ra[:n])
}
