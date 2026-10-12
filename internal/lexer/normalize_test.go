package lexer

import (
	"strings"
	"testing"
	"unicode/utf8"
)

func TestNormalize_ComposesCombiningMarks(t *testing.T) {
	cases := []struct{ in, want string }{
		{"a\u0301", "\u00e1"}, {"e\u0301", "\u00e9"}, {"i\u0301", "\u00ed"}, {"o\u0301", "\u00f3"}, {"u\u0301", "\u00fa"},
		{"u\u0308", "\u00fc"}, {"n\u0303", "\u00f1"},
		{"A\u0301", "\u00e1"}, {"E\u0301", "\u00e9"}, {"I\u0301", "\u00ed"}, {"O\u0301", "\u00f3"}, {"U\u0301", "\u00fa"},
		{"U\u0308", "\u00fc"}, {"N\u0303", "\u00f1"},
		{"factura\u0301", "factur\u00e1"},
		{"Pen\u0303a", "pe\u00f1a"},
		{"Ca\u0301lculo", "c\u00e1lculo"},
		{"\u00c1\u00c9", "\u00e1\u00e9"},
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_LowercasesAlphabet(t *testing.T) {
	cases := []struct{ in, want string }{
		{"Ñandú ÁÉÍÓÚ Ü", "ñandú áéíóú ü"},
		{"FACTURA", "factura"},
		{"Ñ", "ñ"},
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_OutsideAlphabetBecomesSeparator(t *testing.T) {
	cases := []struct{ in, want string }{
		{"¿Qué?", " qué "},
		{"$599.00", " 599 00"},
		{"a-b", "a b"},
		{"façade", "fa ade"},
		{"hola\tmundo\nfin", "hola mundo fin"},
		{"", ""},
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_EmojiAndURLAreSeparators(t *testing.T) {
	got := strings.Fields(Normalize("Genial \U0001F600\U0001F44D ver https://telarg.mx/ayuda?id=7"))
	want := []string{"genial", "ver", "https", "telarg", "mx", "ayuda", "id", "7"}
	if strings.Join(got, "|") != strings.Join(want, "|") {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestNormalize_StrayMarkIsASeparator(t *testing.T) {
	cases := []struct{ in, want string }{
		{"b\u0301", "b "},            // b + acute is outside the alphabet
		{"?\u0301", "  "},            // mark after a replaced character
		{"\u0301a", " a"},            // mark at the start
		{"a\u0301\u0301", "\u00e1 "}, // the second mark has nothing left to compose with
		{"a\u0308", "a "},            // a + diaeresis is not in the alphabet
	}
	for _, c := range cases {
		if got := Normalize(c.in); got != c.want {
			t.Errorf("Normalize(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestNormalize_OutputAlphabetAndIdempotence(t *testing.T) {
	inputs := []string{"¿Qué tal?", "Ñ̃ Á", "x\U0001F600y", "ÁÉ 123 $", "ǅ İ ß", "ä", "\x80\xff bad utf8"}
	for _, in := range inputs {
		out := Normalize(in)
		for _, r := range out {
			if !isLetter(r) && !(r >= '0' && r <= '9') && r != ' ' {
				t.Errorf("Normalize(%q) = %q contains %q outside the alphabet", in, out, r)
			}
		}
		if again := Normalize(out); again != out {
			t.Errorf("not idempotent on %q: %q then %q", in, out, again)
		}
		if utf8.RuneCountInString(out) > utf8.RuneCountInString(in) {
			t.Errorf("Normalize(%q) grew: %q", in, out)
		}
	}
}
