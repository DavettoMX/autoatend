package lexer

import (
	"reflect"
	"testing"
)

func w(s string) Token { return Token{Text: s, Kind: Palabra} }
func n(s string) Token { return Token{Text: s, Kind: Numero} }

func TestTokenize_Words(t *testing.T) {
	cases := []struct {
		in   string
		want []Token
	}{
		{"", nil},
		{"   ", nil},
		{"hola", []Token{w("hola")}},
		{"  hola   mundo ", []Token{w("hola"), w("mundo")}},
		{"cobro mal", []Token{w("cobro"), w("mal")}},
		{"señal débil", []Token{w("señal"), w("débil")}},
		{"añü", []Token{w("añü")}},
	}
	for _, c := range cases {
		if got := Tokenize(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestTokenize_NumbersAreRecognized(t *testing.T) {
	cases := []struct {
		in   string
		want []Token
	}{
		{"599", []Token{n("599")}},
		{"plan 199 mxn", []Token{w("plan"), n("199"), w("mxn")}},
		// Normalize turns "1,499.00" into "1 499 00": three numbers.
		{"1 499 00", []Token{n("1"), n("499"), n("00")}},
	}
	for _, c := range cases {
		if got := Tokenize(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

func TestTokenize_NumberPriceFromNormalizedText(t *testing.T) {
	got := Tokenize(Normalize("$1,499.00"))
	want := []Token{n("1"), n("499"), n("00")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

// A change between letter and digit closes the current token and opens another
// one (documented in tokenizer.go): "abc123" is the word abc and the number 123.
func TestTokenize_LetterDigitBoundarySplits(t *testing.T) {
	cases := []struct {
		in   string
		want []Token
	}{
		{"abc123", []Token{w("abc"), n("123")}},
		{"123abc", []Token{n("123"), w("abc")}},
		{"tk35699a", []Token{w("tk"), n("35699"), w("a")}},
		{"a1b2", []Token{w("a"), n("1"), w("b"), n("2")}},
	}
	for _, c := range cases {
		if got := Tokenize(c.in); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Tokenize(%q) = %v, want %v", c.in, got, c.want)
		}
	}
}

// Any rune outside LETRA/DIGITO acts as a separator, so Tokenize is total even
// on text that was not normalized.
func TestTokenize_UnnormalizedInputSeparates(t *testing.T) {
	got := Tokenize("Hola, mundo!")
	want := []Token{w("ola"), w("mundo")}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestDFA_TransitionTable(t *testing.T) {
	type step struct {
		from  state
		class class
		to    state
		emit  bool
	}
	steps := []step{
		{stInicio, clLetra, stPalabra, false},
		{stInicio, clDigito, stNumero, false},
		{stInicio, clSep, stSeparador, false},
		{stPalabra, clLetra, stPalabra, false},
		{stPalabra, clDigito, stNumero, true},
		{stPalabra, clSep, stSeparador, true},
		{stNumero, clLetra, stPalabra, true},
		{stNumero, clDigito, stNumero, false},
		{stNumero, clSep, stSeparador, true},
		{stSeparador, clLetra, stPalabra, false},
		{stSeparador, clDigito, stNumero, false},
		{stSeparador, clSep, stSeparador, false},
	}
	for _, s := range steps {
		got := delta[s.from][s.class]
		if got.next != s.to || got.emit != s.emit {
			t.Errorf("delta[%v][%v] = %+v, want next=%v emit=%v", s.from, s.class, got, s.to, s.emit)
		}
	}
}

func TestClassify(t *testing.T) {
	cases := []struct {
		r    rune
		want class
	}{
		{'a', clLetra}, {'z', clLetra}, {'ñ', clLetra}, {'ü', clLetra}, {'á', clLetra},
		{'0', clDigito}, {'9', clDigito},
		{' ', clSep}, {'.', clSep}, {'A', clSep}, {'ç', clSep}, {'😀', clSep},
	}
	for _, c := range cases {
		if got := classify(c.r); got != c.want {
			t.Errorf("classify(%q) = %v, want %v", c.r, got, c.want)
		}
	}
}

func TestKindString(t *testing.T) {
	if Palabra.String() != "palabra" || Numero.String() != "numero" {
		t.Errorf("Kind strings: %q %q", Palabra, Numero)
	}
}
