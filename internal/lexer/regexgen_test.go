package lexer

import (
	"math/rand"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGenerateER_TEExample(t *testing.T) {
	got := GenerateER([]string{"factura", "facturas", "facturación", "facturado", "facturando"})
	want := `factur(a|ación|ado|ando|as)`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestCPP01_GenerateERCobro(t *testing.T) {
	// Bytewise order puts "ado" before "aron" ('d' < 'r').
	got := GenerateER([]string{"cobro", "cobros", "cobraron", "cobrado"})
	want := `cobr(ado|aron|o|os)`
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGenerateER_PrefixBacksOffWhenASuffixWouldBeEmpty(t *testing.T) {
	cases := []struct {
		variants []string
		want     string
	}{
		{[]string{"cobro", "cobros"}, `cobr(o|os)`},
		{[]string{"cobros", "cobro"}, `cobr(o|os)`},
		{[]string{"señal", "señales"}, `seña(l|les)`},
		{[]string{"mes", "meses"}, `me(s|ses)`},
		{[]string{"a", "ab"}, `(a|ab)`},
	}
	for _, c := range cases {
		if got := GenerateER(c.variants); got != c.want {
			t.Errorf("GenerateER(%v) = %q, want %q", c.variants, got, c.want)
		}
	}
}

func TestGenerateER_SingleVariant(t *testing.T) {
	if got := GenerateER([]string{"internet"}); got != "internet" {
		t.Errorf("got %q, want internet", got)
	}
	if got := GenerateER([]string{"señal", "señal"}); got != "señal" {
		t.Errorf("duplicates must collapse: got %q", got)
	}
}

func TestGenerateER_EmptyPrefix(t *testing.T) {
	if got, want := GenerateER([]string{"cd", "ab"}), `(ab|cd)`; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGenerateER_NoVariants(t *testing.T) {
	if got := GenerateER(nil); got != "" {
		t.Errorf("got %q, want empty", got)
	}
	if got := GenerateER([]string{""}); got != "" {
		t.Errorf("got %q, want empty", got)
	}
}

func TestGenerateER_EscapesMetacharacters(t *testing.T) {
	cases := []struct {
		variants []string
		want     string
	}{
		{[]string{`a|b`}, `a\|b`},
		{[]string{`a*b+c?`}, `a\*b\+c\?`},
		{[]string{`(a)[b]`}, `\(a\)\[b\]`},
		{[]string{`a\b`}, `a\\b`},
		{[]string{`a-b`}, `a\-b`},
		{[]string{`a|b`, `a|c`}, `a\|(b|c)`},
		{[]string{`x`, `x.`}, `(x|x.)`},
	}
	for _, c := range cases {
		if got := GenerateER(c.variants); got != c.want {
			t.Errorf("GenerateER(%q) = %q, want %q", c.variants, got, c.want)
		}
	}
}

func TestGenerateER_DoesNotModifyItsInput(t *testing.T) {
	in := []string{"cobros", "cobro", "cobros"}
	cp := append([]string(nil), in...)
	GenerateER(in)
	if !reflect.DeepEqual(in, cp) {
		t.Errorf("input modified: %v", in)
	}
}

// expandER is an oracle for the only shape GenerateER emits,
// prefix(s1|s2|...) or (v1|v2|...) or a bare escaped word: it returns the set
// of strings the expression accepts, and fails on empty alternatives.
func expandER(t *testing.T, er string) []string {
	t.Helper()
	var prefix strings.Builder
	rs := []rune(er)
	i := 0
	for i < len(rs) && rs[i] != '(' {
		if rs[i] == '\\' {
			i++
		}
		prefix.WriteRune(rs[i])
		i++
	}
	if i == len(rs) {
		return []string{prefix.String()}
	}
	if rs[len(rs)-1] != ')' {
		t.Fatalf("%q: unterminated group", er)
	}
	var alts []string
	var cur strings.Builder
	for j := i + 1; j < len(rs)-1; j++ {
		switch rs[j] {
		case '\\':
			j++
			cur.WriteRune(rs[j])
		case '|':
			alts = append(alts, cur.String())
			cur.Reset()
		default:
			cur.WriteRune(rs[j])
		}
	}
	alts = append(alts, cur.String())
	var out []string
	for _, a := range alts {
		if a == "" {
			t.Fatalf("%q has an empty alternative", er)
		}
		out = append(out, prefix.String()+a)
	}
	return out
}

func TestGenerateER_ExpansionEqualsVariantSet(t *testing.T) {
	rng := rand.New(rand.NewSource(11))
	letters := []rune("abcñá")
	for iter := 0; iter < 500; iter++ {
		n := rng.Intn(6) + 1
		set := map[string]bool{}
		var variants []string
		for k := 0; k < n; k++ {
			l := rng.Intn(5) + 1
			r := make([]rune, l)
			for x := range r {
				r[x] = letters[rng.Intn(len(letters))]
			}
			v := string(r)
			variants = append(variants, v)
			set[v] = true
		}
		er := GenerateER(variants)
		got := expandER(t, er)
		var want []string
		for v := range set {
			want = append(want, v)
		}
		sort.Strings(got)
		sort.Strings(want)
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("variants %v -> %q expands to %v, want %v", variants, er, got, want)
		}
		if !utf8Valid(er) {
			t.Fatalf("%q is not valid UTF-8 (prefix cut inside a rune)", er)
		}
	}
}

func utf8Valid(s string) bool { return strings.ToValidUTF8(s, "�") == s }
