package lexer

import (
	"reflect"
	"testing"
)

func tableOf(t *testing.T, forms ...string) *Table {
	t.Helper()
	tab := NewTable()
	for _, f := range forms {
		mustInsert(t, tab, f, 1)
	}
	return tab
}

func TestPrefixIndex_GroupsByFirstThreeRunes(t *testing.T) {
	tab := tableOf(t, "factura", "facturas", "factible", "falla", "cobro", "cobros", "cobertura", "señal", "señales", "se")
	idx := BuildPrefixIndex(tab.Entries())

	cases := []struct {
		token string
		want  []string
	}{
		{"factura", []string{"factible", "factura"}},
		{"facil", []string{"factible", "factura"}}, // same first 3 runes "fac"
		{"falla", []string{"falla"}},
		{"cobro", []string{"cobertura", "cobro"}},
		{"señalar", []string{"señal"}},
		{"se", []string{"se"}},
		{"xyz", nil},
		{"", nil},
	}
	for _, c := range cases {
		if got := idx.Candidates(c.token); !reflect.DeepEqual(got, c.want) {
			t.Errorf("Candidates(%q) = %v, want %v", c.token, got, c.want)
		}
	}
}

func TestPrefixIndex_ShortTokensKeyOnTheWholeToken(t *testing.T) {
	// A canonical shorter than 3 runes is its own key.
	tab := tableOf(t, "ok", "okey")
	idx := BuildPrefixIndex(tab.Entries())
	if got, want := idx.Candidates("ok"), []string{"ok"}; !reflect.DeepEqual(got, want) {
		t.Errorf("Candidates(ok) = %v, want %v", got, want)
	}
	if got := idx.Candidates("o"); got != nil {
		t.Errorf("Candidates(o) = %v, want nil", got)
	}
}

func TestPrefixIndex_CandidatesAreSortedAndIsolated(t *testing.T) {
	tab := tableOf(t, "carta", "cargo", "cargos", "cara", "casa")
	idx := BuildPrefixIndex(tab.Entries())
	got := idx.Candidates("car")
	want := []string{"cara", "cargo", "carta"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	got[0] = "mutated"
	if again := idx.Candidates("car"); again[0] != "cara" {
		t.Error("Candidates must return a copy")
	}
}

func TestPrefixIndex_EveryVariantFindsItsLexeme(t *testing.T) {
	tab := tableOf(t, "factura", "facturas", "facturación", "facturado", "cobro", "cobros", "cobraron", "mes", "meses", "señal", "señales")
	idx := BuildPrefixIndex(tab.Entries())
	for _, e := range tab.Entries() {
		canon := e.Canonical()
		for form := range e.Variants {
			if !contains(idx.Candidates(form), canon) {
				t.Errorf("Candidates(%q) = %v lacks the lexeme %q", form, idx.Candidates(form), canon)
			}
		}
	}
}
