package schema

import (
	"strings"
	"testing"
)

func TestStatementCounts(t *testing.T) {
	if n := len(Corpus()); n != 5 {
		t.Errorf("Corpus() has %d statements, want 5", n)
	}
	if n := len(Resultados()); n != 7 {
		t.Errorf("Resultados() has %d statements, want 7", n)
	}
}

func TestEveryStatementStartsWithCreate(t *testing.T) {
	for name, stmts := range map[string][]string{"corpus": Corpus(), "resultados": Resultados()} {
		for i, s := range stmts {
			// Anexo I writes every statement as CREATE TABLE/INDEX IF NOT EXISTS.
			if !strings.HasPrefix(s, "CREATE TABLE IF NOT EXISTS ") && !strings.HasPrefix(s, "CREATE INDEX IF NOT EXISTS ") {
				t.Errorf("%s statement %d is not CREATE TABLE/INDEX IF NOT EXISTS: %q", name, i, s)
			}
			if strings.HasSuffix(s, ";") || s != strings.TrimSpace(s) {
				t.Errorf("%s statement %d is not trimmed or still terminated: %q", name, i, s)
			}
		}
	}
}
