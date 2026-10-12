package lexer

import (
	"testing"
	"unicode/utf8"
)

func TestStem_Table(t *testing.T) {
	cases := []struct{ form, want string }{
		// factura family: one lexeme.
		{"factura", "factur"}, {"facturas", "factur"}, {"facturación", "factur"},
		{"facturaciones", "factur"}, {"facturado", "factur"}, {"facturando", "factur"},
		// cobro family.
		{"cobro", "cobr"}, {"cobros", "cobr"}, {"cobraron", "cobr"}, {"cobrado", "cobr"},
		{"cobrar", "cobr"}, {"cobrando", "cobr"},
		// servicio/servidor must not collapse.
		{"servicio", "servici"}, {"servicios", "servici"}, {"servidor", "servidor"},
		// the root keeps at least 3 runes.
		{"mes", "mes"}, {"meses", "mes"}, {"ido", "ido"}, {"as", "as"}, {"a", "a"}, {"", ""},
		// longest suffix first, only one suffix removed.
		{"aclaramiento", "aclar"}, {"cumplimiento", "cumpl"}, {"posiblemente", "posible"},
		{"mensajes", "mensaj"}, {"mensaje", "mensaj"},
		// no accent folding: the unaccented spelling is another surface form.
		{"facturacion", "facturacion"}, {"cobró", "cobró"},
		// words without a suffix stay as they are.
		{"internet", "internet"}, {"señal", "señal"}, {"plan", "plan"},
	}
	for _, c := range cases {
		if got := Stem(c.form); got != c.want {
			t.Errorf("Stem(%q) = %q, want %q", c.form, got, c.want)
		}
	}
}

func TestStem_GroupsFacturaFamily(t *testing.T) {
	want := Stem("factura")
	for _, f := range []string{"facturas", "facturación", "facturado", "facturando"} {
		if got := Stem(f); got != want {
			t.Errorf("Stem(%q) = %q, want %q", f, got, want)
		}
	}
}

func TestStem_Invariants(t *testing.T) {
	for _, f := range []string{"factura", "facturación", "mes", "ido", "cobraron", "servidor", "x", "añadidos", "señales"} {
		s := Stem(f)
		if len(s) > len(f) {
			t.Errorf("Stem(%q) = %q is longer than the form", f, s)
		}
		if s != f && utf8.RuneCountInString(s) < 3 {
			t.Errorf("Stem(%q) = %q has fewer than 3 runes", f, s)
		}
		if f[:len(s)] != s {
			t.Errorf("Stem(%q) = %q is not a prefix", f, s)
		}
	}
}

// The suffix list is ordered from longest to shortest, so the first match is
// the longest one.
func TestStemSuffixesAreLongestFirst(t *testing.T) {
	prev := 1 << 30
	for _, s := range suffixes {
		n := utf8.RuneCountInString(s)
		if n > prev {
			t.Errorf("suffix %q (%d runes) comes after a shorter one (%d)", s, n, prev)
		}
		prev = n
	}
	if got := len(suffixes); got != 29 {
		t.Errorf("suffix count = %d, want 29", got)
	}
}
