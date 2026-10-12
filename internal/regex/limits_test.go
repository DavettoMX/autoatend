package regex

import (
	"strings"
	"testing"
)

// ER length: 1, 100 and 500 characters (TE §6.3.4, análisis de valores límite):
// the parser and the Thompson construction handle long expressions, the NFA
// stays within 2·|r| states and accepts the string it describes.
func TestCPU03_LongitudesER(t *testing.T) {
	for _, n := range []int{1, 100, 500} {
		src := strings.Repeat("a", n)
		m, err := Compile(src)
		if err != nil {
			t.Fatalf("length %d: %v", n, err)
		}
		if m.NumStates() > 2*n {
			t.Errorf("length %d: %d states, want ≤ %d", n, m.NumStates(), 2*n)
		}
		if !m.Match(src) || m.Match(src+"a") || (n > 1 && m.Match(src[1:])) {
			t.Errorf("length %d: wrong acceptance", n)
		}
	}
	// A 100-character expression with operators: alternation, groups and closures.
	src := strings.Repeat("(ab|c)*d", 12) + "efgh"
	if len(src) != 100 {
		t.Fatalf("operator expression has %d characters, want 100", len(src))
	}
	m, err := Compile(src)
	if err != nil {
		t.Fatalf("100-character expression: %v", err)
	}
	if !m.Match(strings.Repeat("abcd", 12) + "efgh") {
		t.Error("100-character expression rejects a string in its language")
	}
}
