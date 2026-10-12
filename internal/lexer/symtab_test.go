package lexer

import (
	"math/rand"
	"reflect"
	"sort"
	"testing"
)

func mustInsert(t *testing.T, tab *Table, form string, cat int) {
	t.Helper()
	if _, err := tab.Insert(form, cat); err != nil {
		t.Fatalf("Insert(%q, %d): %v", form, cat, err)
	}
}

// CPU09: the same surface form twice in category 1.
func TestCPU09_SameFormTwiceInCategory1(t *testing.T) {
	tab := NewTable()
	mustInsert(t, tab, "factura", 1)
	mustInsert(t, tab, "factura", 1)
	e, ok := tab.Lookup("factura")
	if !ok {
		t.Fatal("factura not in the table")
	}
	if e.Global != 2 {
		t.Errorf("Global = %d, want 2", e.Global)
	}
	if e.PerCat != [5]int{2, 0, 0, 0, 0} {
		t.Errorf("PerCat = %v, want [2 0 0 0 0]", e.PerCat)
	}
	if len(e.Variants) != 1 || e.Variants["factura"] != 2 {
		t.Errorf("Variants = %v, want {factura: 2}", e.Variants)
	}
}

// CPU10: two surface forms of one lexeme in different categories.
func TestCPU10_VariantsShareALexeme(t *testing.T) {
	tab := NewTable()
	mustInsert(t, tab, "factura", 1)
	mustInsert(t, tab, "facturas", 2)
	if got := len(tab.Entries()); got != 1 {
		t.Fatalf("lexemes = %d, want 1", got)
	}
	e, _ := tab.Lookup("factura")
	if e.Global != 2 {
		t.Errorf("Global = %d, want 2", e.Global)
	}
	if len(e.Variants) != 2 {
		t.Errorf("variants = %d, want 2", len(e.Variants))
	}
	if e.PerCat != [5]int{1, 1, 0, 0, 0} {
		t.Errorf("PerCat = %v, want [1 1 0 0 0]", e.PerCat)
	}
	if e2, ok := tab.Lookup("facturas"); !ok || e2 != e {
		t.Error("Lookup of the other variant must return the same entry")
	}
}

func TestInsert_RejectsCategoriesOutsideOneToFive(t *testing.T) {
	tab := NewTable()
	for _, cat := range []int{-1, 0, 6, 100} {
		if _, err := tab.Insert("factura", cat); err == nil {
			t.Errorf("Insert with categoria %d: expected error", cat)
		}
	}
	if got := len(tab.Entries()); got != 0 {
		t.Errorf("a rejected insert created %d entries", got)
	}
	for cat := 1; cat <= 5; cat++ {
		if _, err := tab.Insert("factura", cat); err != nil {
			t.Errorf("Insert with categoria %d: %v", cat, err)
		}
	}
}

func TestCanonical_MostFrequentVariantTieBreaksLexicographically(t *testing.T) {
	tab := NewTable()
	mustInsert(t, tab, "facturas", 1)
	mustInsert(t, tab, "facturas", 1)
	mustInsert(t, tab, "factura", 1)
	mustInsert(t, tab, "factura", 2)
	e, _ := tab.Lookup("factura")
	// factura 2, facturas 2: the tie goes to the smaller string.
	if got := e.Canonical(); got != "factura" {
		t.Errorf("Canonical = %q, want factura (tie)", got)
	}
	mustInsert(t, tab, "facturas", 3)
	if got := e.Canonical(); got != "facturas" {
		t.Errorf("Canonical = %q, want facturas (most frequent)", got)
	}
}

func TestSortedVariants(t *testing.T) {
	tab := NewTable()
	for _, f := range []string{"facturas", "factura", "facturación", "factura"} {
		mustInsert(t, tab, f, 1)
	}
	e, _ := tab.Lookup("factura")
	want := []VariantCount{{"factura", 2}, {"facturación", 1}, {"facturas", 1}}
	if got := e.SortedVariants(); !reflect.DeepEqual(got, want) {
		t.Errorf("SortedVariants = %v, want %v", got, want)
	}
}

func TestEntries_SortedByCanonical(t *testing.T) {
	tab := NewTable()
	for _, f := range []string{"zona", "cobro", "internet", "señal", "factura", "cobros", "cobros"} {
		mustInsert(t, tab, f, 3)
	}
	var got []string
	for _, e := range tab.Entries() {
		got = append(got, e.Canonical())
	}
	if !sort.StringsAreSorted(got) {
		t.Errorf("Entries not sorted by canonical: %v", got)
	}
	want := []string{"cobros", "factura", "internet", "señal", "zona"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("canonicals = %v, want %v", got, want)
	}
}

func TestLookup_Missing(t *testing.T) {
	tab := NewTable()
	if _, ok := tab.Lookup("factura"); ok {
		t.Error("Lookup on an empty table found something")
	}
}

// Global equals the sum over categories after N random inserts, and the sum of
// the variant frequencies equals Global (checked against a plain counter that
// does not use the Table).
func TestInsert_GlobalIsSumOfCategories(t *testing.T) {
	rng := rand.New(rand.NewSource(7))
	forms := []string{"factura", "facturas", "cobro", "cobros", "cobraron", "señal", "internet", "mes", "meses"}
	tab := NewTable()
	oracle := map[string][5]int{} // by surface form
	const n = 2000
	for i := 0; i < n; i++ {
		f := forms[rng.Intn(len(forms))]
		cat := rng.Intn(5) + 1
		mustInsert(t, tab, f, cat)
		c := oracle[f]
		c[cat-1]++
		oracle[f] = c
	}
	total := 0
	for _, e := range tab.Entries() {
		sumCat, sumVar := 0, 0
		for _, v := range e.PerCat {
			sumCat += v
		}
		for form, freq := range e.Variants {
			sumVar += freq
			var s int
			for _, v := range oracle[form] {
				s += v
			}
			if s != freq {
				t.Errorf("variant %q frequency %d, oracle %d", form, freq, s)
			}
		}
		if e.Global != sumCat || e.Global != sumVar {
			t.Errorf("%s: Global %d, sum categories %d, sum variants %d", e.Canonical(), e.Global, sumCat, sumVar)
		}
		total += e.Global
	}
	if total != n {
		t.Errorf("total = %d, want %d", total, n)
	}
}
