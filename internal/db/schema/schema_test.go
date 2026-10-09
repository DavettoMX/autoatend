package schema

import (
	"strings"
	"testing"
)

func TestStatementCounts(t *testing.T) {
	if n := len(Corpus()); n != 6 {
		t.Errorf("Corpus() has %d statements, want 6", n)
	}
	if n := len(Resultados()); n != 8 {
		t.Errorf("Resultados() has %d statements, want 8", n)
	}
}

func TestEveryStatementStartsWithCreate(t *testing.T) {
	for name, stmts := range map[string][]string{"corpus": Corpus(), "resultados": Resultados()} {
		for i, s := range stmts {
			if !strings.HasPrefix(s, "CREATE ") {
				t.Errorf("%s statement %d does not start with CREATE: %q", name, i, s)
			}
			if strings.HasSuffix(s, ";") || s != strings.TrimSpace(s) {
				t.Errorf("%s statement %d is not trimmed or still terminated: %q", name, i, s)
			}
		}
	}
}
