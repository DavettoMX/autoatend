package lexer

import (
	"math"
	"strings"
	"testing"
	"time"

	"autoatend/internal/db"
	"autoatend/internal/regex"
)

// Boundary values of TE §6.3.4 (análisis de valores límite) for the lexical
// pipeline. Each test is attached to the Anexo II case of the same component.

// Text length: 0, 1, 2, 999, 1000, 1001 and 10000 characters (CE-4, CE-6).
func TestCPE04_LongitudesLimite(t *testing.T) {
	word := "factura " // 8 runes, one catalogued token per repetition
	for _, n := range []int{0, 1, 2, 999, 1000, 1001, 10000} {
		text := strings.Repeat(word, n/len(word)+1)[:n]
		start := time.Now()
		got := Analyze(text)
		if d := time.Since(start); d > 500*time.Millisecond {
			t.Errorf("length %d took %v, want < 500ms", n, d)
		}
		// Oracle: every complete "factura" in the text is a token; a cut
		// trailing fragment ("f", "fa", ...) is a token of its own.
		want := strings.Fields(text)
		if len(got) != len(want) {
			t.Errorf("length %d: %d tokens, want %d", n, len(got), len(want))
		}
	}
}

// Number of tokens: 0, 1, 2 (CE-2, CE-3 a single token).
func TestCPE02_TokensLimite(t *testing.T) {
	cases := []struct {
		text string
		want []string
	}{
		{"", nil},
		{"servicio", []string{"servicio"}},
		{"sin servicio", []string{"servicio"}},
		{"servicio cortado", []string{"servicio", "cortado"}},
	}
	for _, c := range cases {
		got := Analyze(c.text)
		if strings.Join(got, " ") != strings.Join(c.want, " ") {
			t.Errorf("Analyze(%q) = %v, want %v", c.text, got, c.want)
		}
	}
}

// CE-3: a comment made of a single non-stoplist token gives a symbol table
// with exactly one entry of frequency 1.
func TestCPE01_TokenUnicoCE3(t *testing.T) {
	tab := NewTable()
	for _, form := range Analyze("Facturación") {
		mustInsert(t, tab, form, 3)
	}
	es := tab.Entries()
	if len(es) != 1 || es[0].Global != 1 || es[0].PerCat != [5]int{0, 0, 1, 0, 0} {
		t.Fatalf("entries = %d, want exactly one lexeme with frequency 1 in category 3", len(es))
	}
}

// Number of tokens: the maximum observed in the corpus is processed whole.
func TestCPE02_MaximoDeTokensDelCorpus(t *testing.T) {
	if testing.Short() {
		t.Skip("reads the full corpus")
	}
	d, err := db.OpenReadOnly("../../data/corpus.db")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	rows, err := d.Query(`SELECT texto FROM comentarios`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	maxTok, maxText := 0, ""
	for rows.Next() {
		var s string
		if err := rows.Scan(&s); err != nil {
			t.Fatal(err)
		}
		if n := len(Tokenize(Normalize(s))); n > maxTok {
			maxTok, maxText = n, s
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	if maxTok == 0 {
		t.Fatal("no tokens in the corpus")
	}
	toks := Tokenize(Normalize(maxText))
	if len(toks) != maxTok {
		t.Errorf("re-tokenizing the longest comment gave %d tokens, want %d", len(toks), maxTok)
	}
	t.Logf("maximum tokens in one comment: %d", maxTok)
}

// Lexeme frequency: 1, 2, MaxInt32 (Integer.MAX of the TE, 32 bits) and the
// limit of the counter type. Crossing 2^31 must be exact; reaching the maximum
// of int must be an error, never a silent wrap-around.
func TestCPU09_ContadoresSinDesbordamiento(t *testing.T) {
	tab := NewTable()
	e, err := tab.Insert("factura", 1)
	if err != nil {
		t.Fatal(err)
	}
	if e.Global != 1 {
		t.Fatalf("Global = %d after one insert, want 1", e.Global)
	}
	mustInsert(t, tab, "factura", 1)
	if e.Global != 2 {
		t.Fatalf("Global = %d after two inserts, want 2", e.Global)
	}

	// Integer.MAX/2 and Integer.MAX (32 bits): the int counters go past them exactly.
	for _, start := range []int{math.MaxInt32/2 - 1, math.MaxInt32 - 1} {
		e.Global, e.PerCat[0], e.Variants["factura"] = start, start, start
		mustInsert(t, tab, "factura", 1)
		if e.Global != start+1 || e.PerCat[0] != start+1 || e.Variants["factura"] != start+1 {
			t.Errorf("counter at %d did not increment exactly", start)
		}
	}

	// The maximum of the counter type: the insert is rejected and nothing changes.
	e.Global, e.PerCat[0], e.Variants["factura"] = math.MaxInt, math.MaxInt, math.MaxInt
	if _, err := tab.Insert("factura", 1); err == nil {
		t.Fatal("Insert at math.MaxInt returned no error (the counter would wrap around)")
	}
	if e.Global != math.MaxInt || e.PerCat[0] != math.MaxInt || e.Variants["factura"] != math.MaxInt {
		t.Error("a rejected insert modified the counters")
	}
}

// Number of variants: 1, 2, 10 and 50 per lexeme; the generated ER accepts
// exactly those variants (round trip with the project's own engine).
func TestCPP01_VariantesLimite(t *testing.T) {
	suffixes := []string{"a", "as", "ación", "ado", "ando", "aba", "aban", "aron", "e", "en"}
	for _, n := range []int{1, 2, 10, 50} {
		var vs []string
		for i := 0; len(vs) < n; i++ {
			v := "factur" + suffixes[i%len(suffixes)]
			if i >= len(suffixes) {
				v += strings.Repeat("s", i/len(suffixes))
			}
			vs = append(vs, v)
		}
		er := GenerateER(vs)
		m, err := regex.Compile(er)
		if err != nil {
			t.Fatalf("%d variants: ER %q does not compile: %v", n, er, err)
		}
		for _, v := range vs {
			if !m.Match(v) {
				t.Errorf("%d variants: ER %q rejects %q", n, er, v)
			}
		}
		if m.Match("factur") || m.Match("facturx") {
			t.Errorf("%d variants: ER %q accepts a non-variant", n, er)
		}
	}
}

// CE-7: the stoplist homographs of the TE (bajo, para, sobre) are retained and
// enter the symbol table.
func TestCPU12_HomografosCE7(t *testing.T) {
	got := Analyze("pago bajo protesta, para nada, sobre todo")
	tab := NewTable()
	for _, f := range got {
		mustInsert(t, tab, f, 4)
	}
	for _, h := range []string{"bajo", "para", "sobre"} {
		if _, ok := tab.Lookup(h); !ok {
			t.Errorf("homograph %q is not in the symbol table (got %v)", h, got)
		}
	}
}
