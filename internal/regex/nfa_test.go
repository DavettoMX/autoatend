package regex

import (
	"strings"
	"testing"
	"time"
)

// mustCompile compiles src and fails the test on error.
func mustCompile(t *testing.T, src string) *NFA {
	t.Helper()
	m, err := Compile(src)
	if err != nil {
		t.Fatalf("Compile(%q): %v", src, err)
	}
	return m
}

func contains(set []int, s int) bool {
	for _, x := range set {
		if x == s {
			return true
		}
	}
	return false
}

// CPU04 (TE §6.2.4): NodoCaracter('a') gives a fragment with exactly two
// states and one labelled transition.
func TestCPU04_ThompsonChar(t *testing.T) {
	m := Thompson(&Char{'a'})
	if m.NumStates() != 2 || m.NumLabelled() != 1 || m.NumEpsilon() != 0 {
		t.Fatalf("char: %d states, %d labelled, %d epsilon; want 2, 1, 0",
			m.NumStates(), m.NumLabelled(), m.NumEpsilon())
	}
	if m.Start() == m.Accept() {
		t.Error("start and accept must be different states")
	}
	if !m.Match("a") || m.Match("") || m.Match("b") || m.Match("aa") {
		t.Error("char a must accept exactly \"a\"")
	}
}

// CPU05 (TE §6.2.4): one test per composition rule, with the state and
// transition counts the construction dictates.
func TestCPU05_ThompsonComposition(t *testing.T) {
	cases := []struct {
		src                         string
		states, epsilon, transition int
	}{
		{"a", 2, 0, 1},
		{"(a)", 2, 0, 1},    // grouping adds nothing
		{"a|b", 6, 4, 2},    // new start and accept, 2 ε out and 2 ε in
		{"ab", 4, 1, 2},     // one ε from accept(N1) to start(N2)
		{"abc", 6, 2, 3},    // chained concatenation
		{"a*", 4, 4, 1},     // start→N, start→accept, N.accept→N.start, N.accept→accept
		{"a+", 4, 3, 1},     // no start→accept
		{"a?", 4, 3, 1},     // no N.accept→N.start
		{"[a-cx]", 2, 0, 2}, // one labelled transition per range
		{"[a]", 2, 0, 1},
		{"a|b|c", 10, 8, 3}, // alternation is built pairwise: 6 + 2 + 2
		{"factur(a|as)", 20, 11, 9},
	}
	for _, c := range cases {
		m := mustCompile(t, c.src)
		if m.NumStates() != c.states || m.NumEpsilon() != c.epsilon || m.NumLabelled() != c.transition {
			t.Errorf("%q: %d states, %d epsilon, %d labelled; want %d, %d, %d",
				c.src, m.NumStates(), m.NumEpsilon(), m.NumLabelled(), c.states, c.epsilon, c.transition)
		}
	}
}

// The construction is bounded by 2·|r| states, |r| being the length of the
// expression in runes (TE §6.2.4).
func TestCPU05_ThompsonStateBound(t *testing.T) {
	for _, src := range []string{
		"a", "a|b", "ab", "a*", "(a)", "(a|b)*c", "[a-z]+", "factur(a|as|ación|ado|ando)",
		"((a|b)(c|d))*", "a?b?c?d?", nested(50, 500), `a\|b\*c`,
	} {
		m := mustCompile(t, src)
		if limit := 2 * len([]rune(src)); m.NumStates() > limit {
			t.Errorf("%q: %d states exceeds 2·|r| = %d", src, m.NumStates(), limit)
		}
	}
}

// ε-closures per operator (TE §6.2.4): which states are reachable from the
// start state without consuming input.
func TestCPU05_EpsilonClosures(t *testing.T) {
	cases := []struct {
		src        string
		closure    int  // size of the ε-closure of the start state
		acceptsEps bool // whether that closure holds the accept state
	}{
		{"a", 1, false},
		{"ab", 1, false},
		{"a|b", 3, false}, // start and the two starts of the branches
		{"a*", 3, true},   // start, N.start and accept
		{"a?", 3, true},
		{"a+", 2, false}, // start and N.start, never the accept
		{"(a*)*", 5, true},
	}
	for _, c := range cases {
		m := mustCompile(t, c.src)
		cl := m.EpsilonClosure(m.Start())
		if len(cl) != c.closure {
			t.Errorf("%q: |ε-closure(start)| = %d (%v), want %d", c.src, len(cl), cl, c.closure)
		}
		if got := contains(cl, m.Accept()); got != c.acceptsEps {
			t.Errorf("%q: accept in ε-closure(start) = %v, want %v", c.src, got, c.acceptsEps)
		}
		if !contains(cl, m.Start()) {
			t.Errorf("%q: a state is always in its own ε-closure", c.src)
		}
	}
	// After reading 'a' in "ab" the closure holds the accept of 'a' and, via
	// the concatenation ε, the start of 'b'.
	m := mustCompile(t, "ab")
	if cl := m.EpsilonClosure(m.Start()); len(cl) != 1 {
		t.Fatalf("closure of start = %v", cl)
	}
	if got := m.step([]int{m.Start()}, 'a'); len(got) != 2 {
		t.Errorf("ε-closure(move({start}, 'a')) = %v, want 2 states", got)
	}
}

// CPU06 (TE §6.2.4): the NFA of factur(a|as) accepts "facturas".
func TestCPU06_AcceptsVariant(t *testing.T) {
	m := mustCompile(t, "factur(a|as)")
	for _, s := range []string{"facturas", "factura"} {
		if !m.Match(s) {
			t.Errorf("factur(a|as) must accept %q", s)
		}
	}
}

// CPU07: and rejects "facturó"; matching is on the whole string.
func TestCPU07_RejectsNonVariant(t *testing.T) {
	m := mustCompile(t, "factur(a|as)")
	for _, s := range []string{"facturó", "factur", "facturass", "xfacturas", "factura ", "", "FACTURAS"} {
		if m.Match(s) {
			t.Errorf("factur(a|as) must reject %q", s)
		}
	}
}

func TestMatchSemantics(t *testing.T) {
	cases := []struct {
		src    string
		accept []string
		reject []string
	}{
		{"a*", []string{"", "a", "aaaa"}, []string{"b", "ab"}},
		{"a+", []string{"a", "aaa"}, []string{"", "b"}},
		{"a?b", []string{"b", "ab"}, []string{"", "aab", "a"}},
		{"(a|b)*c", []string{"c", "abc", "bbac"}, []string{"", "ab", "cc"}},
		{"(ab)+", []string{"ab", "abab"}, []string{"", "a", "aba"}},
		{"[a-c]+x", []string{"ax", "cabx"}, []string{"x", "dx"}},
		{"[ñü]a", []string{"ña", "üa"}, []string{"na", "a"}},
		{"señ(al|ales)", []string{"señal", "señales"}, []string{"senal", "señale"}},
		{"cancelación|baja", []string{"cancelación", "baja"}, []string{"cancelacion", "bajas"}},
		{"7a", []string{"7a"}, []string{"a", "7"}},
		// Nested closures create ε-cycles; the simulation must still terminate.
		{"(a*)*", []string{"", "a", "aaa"}, []string{"b"}},
		{"(a?)+b", []string{"b", "aab"}, []string{"a"}},
		{"(a*|b*)*c", []string{"c", "abbac"}, []string{"d"}},
		// Escaped metacharacters are literals.
		{`a\|b`, []string{"a|b"}, []string{"a", "b"}},
		{`\(x\)\*`, []string{"(x)*"}, []string{"x"}},
		{`[\-a]+`, []string{"-", "a-a"}, []string{"b"}},
		{`a\\b`, []string{`a\b`}, []string{"ab"}},
	}
	for _, c := range cases {
		m := mustCompile(t, c.src)
		for _, s := range c.accept {
			if !m.Match(s) {
				t.Errorf("%q must accept %q", c.src, s)
			}
		}
		for _, s := range c.reject {
			if m.Match(s) {
				t.Errorf("%q must reject %q", c.src, s)
			}
		}
	}
}

func TestCompileError(t *testing.T) {
	if m, err := Compile("factur(a|as"); err == nil || m != nil {
		t.Errorf("Compile of an invalid ER = (%v, %v), want an error", m, err)
	}
}

// CPU08 (TE §6.2.4): a* accepts a string of 10,000 'a'.
func TestCPU08_LongInput(t *testing.T) {
	m := mustCompile(t, "a*")
	if !m.Match(strings.Repeat("a", 10_000)) {
		t.Error("a* must accept 10,000 'a'")
	}
	if m.Match(strings.Repeat("a", 10_000) + "b") {
		t.Error("a* must reject 10,000 'a' followed by 'b'")
	}
}

// The classic exponential case for backtracking engines, (a?){n}a{n} against
// a^n, stays linear here because the simulation never backtracks.
func TestCPU08_NoBacktrackingBlowup(t *testing.T) {
	const n = 60
	m := mustCompile(t, strings.Repeat("a?", n)+strings.Repeat("a", n))
	start := time.Now()
	if !m.Match(strings.Repeat("a", n)) {
		t.Error("(a?){60}a{60} must accept a^60")
	}
	if m.Match(strings.Repeat("a", n-1)) {
		t.Error("(a?){60}a{60} must reject a^59")
	}
	if d := time.Since(start); d > time.Second {
		t.Errorf("took %v: the simulation is backtracking", d)
	}
}

// work is the deterministic cost of a match: the number of states the
// simulation visits plus the transitions it examines.
func work(m *NFA, s string) int {
	_, w := m.simulate(s)
	return w
}

// CPU08: O(n·m) time, checked first with a deterministic operation count and
// then with wall time at 10, 100, 1,000 and 10,000 characters.
func TestCPU08_LinearCost(t *testing.T) {
	m := mustCompile(t, "(x|y|z)*w")
	lengths := []int{10, 100, 1_000, 10_000}
	var ops []int
	for _, n := range lengths {
		s := strings.Repeat("xyz", n/3+1)[:n]
		w := work(m, s)
		ops = append(ops, w)
		if bound := 8 * n * m.NumStates(); w > bound {
			t.Errorf("n=%d: work %d exceeds 8·n·m = %d", n, w, bound)
		}
	}
	for i := 1; i < len(lengths); i++ {
		ratio := float64(ops[i]) / float64(ops[i-1])
		if ratio < 8 || ratio > 12 {
			t.Errorf("work grew %.2f× from n=%d to n=%d, want about 10×", ratio, lengths[i-1], lengths[i])
		}
	}
	t.Logf("work at n=%v: %v", lengths, ops)
}

func TestCPU08_LinearTime(t *testing.T) {
	if testing.Short() {
		t.Skip("timing probes skipped in -short mode")
	}
	m := mustCompile(t, "(x|y|z)*w")
	lengths := []int{10, 100, 1_000, 10_000}
	times := make([]time.Duration, len(lengths))
	for i, n := range lengths {
		s := strings.Repeat("xyz", n/3+1)[:n]
		reps := max(1, 400_000/n)
		best := time.Duration(1<<63 - 1)
		for trial := 0; trial < 5; trial++ {
			start := time.Now()
			for r := 0; r < reps; r++ {
				m.Match(s)
			}
			best = min(best, time.Since(start)/time.Duration(reps))
		}
		times[i] = best
	}
	t.Logf("time per match at n=%v (m=%d states): %v", lengths, m.NumStates(), times)
	// Linear growth is a factor of 10 per step; a quadratic algorithm gives
	// 100. The tolerance is generous because of timer noise and cache effects.
	for i := 2; i < len(lengths); i++ {
		if ratio := float64(times[i]) / float64(times[i-1]); ratio > 40 {
			t.Errorf("time grew %.1f× from n=%d to n=%d, want about 10× (linear)", ratio, lengths[i-1], lengths[i])
		}
	}
}

// The cost is also linear in the size of the NFA for a fixed input.
func TestCPU08_LinearInStates(t *testing.T) {
	small := mustCompile(t, "(a|b)*c")
	var alts []string
	for i := 0; i < 20; i++ {
		alts = append(alts, "a")
	}
	big := mustCompile(t, "("+strings.Join(alts, "|")+")*c")
	s := strings.Repeat("a", 1000)
	ws, wb := work(small, s), work(big, s)
	growth := float64(big.NumStates()) / float64(small.NumStates())
	if ratio := float64(wb) / float64(ws); ratio > 2*growth {
		t.Errorf("work grew %.1f× for an NFA %.1f× larger", ratio, growth)
	}
}
