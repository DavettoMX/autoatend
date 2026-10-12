package regex_test

import (
	"math/rand"
	"sort"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/lexer"
	"autoatend/internal/regex"
)

// resultadosPath is the versioned catalogue, relative to this package.
const resultadosPath = "../../data/resultados.db"

// The generator escapes the metacharacters of the grammar; the parser must
// read them back as literals (round trip lexer.GenerateER -> regex).
func TestCPP01_EscapesRoundTrip(t *testing.T) {
	variants := []string{`a|b`, `a*`, `a+b`, `a?`, `(x)`, `[y]`, `a-b`, `a\b`, `\|*+?()[]-`, `a`}
	er := lexer.GenerateER(variants)
	m, err := regex.Compile(er)
	if err != nil {
		t.Fatalf("Compile(%q): %v", er, err)
	}
	for _, v := range variants {
		if !m.Match(v) {
			t.Errorf("%q does not accept its variant %q", er, v)
		}
	}
	for _, s := range []string{"", "ab", "a|", "|b", "aa", `a\`, "x", "(x"} {
		if m.Match(s) {
			t.Errorf("%q accepts %q, which is not a variant", er, s)
		}
	}
}

// Randomized round trip: for random variant sets over an alphabet that mixes
// letters and every metacharacter, the ER of GenerateER accepts exactly the
// variants (it rejects other random strings that are not in the set).
func TestCPP01_GeneratorRoundTripRandom(t *testing.T) {
	alphabet := []rune("abñ|*+?()[]\\-")
	rng := rand.New(rand.NewSource(10433))
	word := func() string {
		rs := make([]rune, 1+rng.Intn(6))
		for i := range rs {
			rs[i] = alphabet[rng.Intn(len(alphabet))]
		}
		return string(rs)
	}
	for iter := 0; iter < 2000; iter++ {
		set := map[string]bool{}
		var variants []string
		for k := 1 + rng.Intn(5); k > 0; k-- {
			if w := word(); !set[w] {
				set[w] = true
				variants = append(variants, w)
			}
		}
		sort.Strings(variants)
		er := lexer.GenerateER(variants)
		m, err := regex.Compile(er)
		if err != nil {
			t.Fatalf("variants %q: ER %q does not parse: %v", variants, er, err)
		}
		for _, v := range variants {
			if !m.Match(v) {
				t.Fatalf("variants %q: ER %q rejects %q", variants, er, v)
			}
		}
		for probe := 0; probe < 20; probe++ {
			if w := word(); m.Match(w) != set[w] {
				t.Fatalf("variants %q: ER %q on %q = %v, want %v", variants, er, w, m.Match(w), set[w])
			}
		}
	}
}

type lexema struct {
	id       int
	canon    string
	er       string
	variants []string
}

// loadCatalogue reads every lexeme of data/resultados.db with its stored ER
// and its observed variants, straight from the database.
func loadCatalogue(t *testing.T) []lexema {
	t.Helper()
	d, err := db.OpenReadOnly(resultadosPath)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	rows, err := d.Query(`SELECT l.id, l.forma_canonica, l.expresion_regular, v.forma
		FROM lexemas l LEFT JOIN variantes v ON v.lexema_id = l.id
		ORDER BY l.id, v.forma`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var out []lexema
	for rows.Next() {
		var l lexema
		var form *string
		if err := rows.Scan(&l.id, &l.canon, &l.er, &form); err != nil {
			t.Fatal(err)
		}
		if n := len(out); n == 0 || out[n-1].id != l.id {
			out = append(out, l)
		}
		if form != nil {
			last := &out[len(out)-1]
			last.variants = append(last.variants, *form)
		}
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}
	return out
}

// Integration 4b (TE §6.3.5): every ER of the catalogue parses with the own
// parser, and its NFA respects the 2·|r| bound.
func TestIntegration4b_EveryERParses(t *testing.T) {
	if testing.Short() {
		t.Skip("catalogue integration skipped in -short mode")
	}
	cat := loadCatalogue(t)
	if len(cat) == 0 {
		t.Fatal("empty catalogue")
	}
	maxStates, sum := 0, 0
	for _, l := range cat {
		m, err := regex.Compile(l.er)
		if err != nil {
			t.Errorf("lexeme %q: ER %q does not parse: %v", l.canon, l.er, err)
			continue
		}
		if limit := 2 * len([]rune(l.er)); m.NumStates() > limit {
			t.Errorf("lexeme %q: %d states exceeds 2·|r| = %d", l.canon, m.NumStates(), limit)
		}
		maxStates = max(maxStates, m.NumStates())
		sum += m.NumStates()
	}
	t.Logf("%d ERs parsed; NFA states: max %d, mean %.2f", len(cat), maxStates, float64(sum)/float64(len(cat)))
}

// I-1 (TE §6.3.3): for every lexeme of data/resultados.db, its stored ER,
// compiled with the own engine, accepts every observed variant.
func TestCPP01_I1_RoundTripCatalogue(t *testing.T) {
	if testing.Short() {
		t.Skip("catalogue integration skipped in -short mode")
	}
	cat := loadCatalogue(t)
	variants, failures := 0, 0
	for _, l := range cat {
		m, err := regex.Compile(l.er)
		if err != nil {
			failures++
			t.Errorf("lexeme %q: ER %q does not compile: %v", l.canon, l.er, err)
			continue
		}
		if len(l.variants) == 0 {
			failures++
			t.Errorf("lexeme %q has no variants", l.canon)
		}
		for _, v := range l.variants {
			variants++
			if !m.Match(v) {
				failures++
				if failures <= 20 {
					t.Errorf("lexeme %q: ER %q rejects its variant %q", l.canon, l.er, v)
				}
			}
		}
	}
	t.Logf("I-1: %d lexemes, %d variants checked, %d failures", len(cat), variants, failures)
	if failures != 0 {
		t.Errorf("I-1 round trip: %d failures, want 0", failures)
	}
}
