package lexer

import (
	"fmt"
	"math"
	"sort"
)

// NumCategories is the number of case labels (db.Categories).
const NumCategories = 5

// VariantCount is a surface form with its frequency.
type VariantCount struct {
	Form string
	N    int
}

// Entry is one lexeme of the symbol table: all the surface variants that share
// a stem, with their frequencies globally and per category.
type Entry struct {
	Stem     string         // grouping key (Stem of every variant)
	Variants map[string]int // surface form -> occurrences
	PerCat   [NumCategories]int
	Global   int
}

// Canonical is the name of the lexeme: its most frequent variant, ties broken
// by the smaller string, so it does not depend on map iteration order.
func (e *Entry) Canonical() string {
	best, bestN := "", -1
	for form, n := range e.Variants {
		if n > bestN || (n == bestN && form < best) {
			best, bestN = form, n
		}
	}
	return best
}

// SortedVariants lists the variants ordered by form (byte order).
func (e *Entry) SortedVariants() []VariantCount {
	out := make([]VariantCount, 0, len(e.Variants))
	for form, n := range e.Variants {
		out = append(out, VariantCount{form, n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Form < out[j].Form })
	return out
}

// Table is the symbol table of the sweep (TE §6.2.2): a hash table from stem to
// lexeme, so insertion and lookup are O(1) on average.
type Table struct {
	byStem map[string]*Entry
}

// NewTable returns an empty symbol table.
func NewTable() *Table { return &Table{byStem: make(map[string]*Entry)} }

// Insert records one occurrence of the surface form in category cat (1..5) and
// returns the lexeme it belongs to.
func (t *Table) Insert(form string, cat int) (*Entry, error) {
	if cat < 1 || cat > NumCategories {
		return nil, fmt.Errorf("lexer: categoria_id %d out of range 1..%d", cat, NumCategories)
	}
	key := Stem(form)
	e := t.byStem[key]
	if e == nil {
		e = &Entry{Stem: key, Variants: make(map[string]int)}
		t.byStem[key] = e
	}
	// TE §6.3.4: no counter may silently wrap around; the insert is rejected
	// before any of them changes.
	if e.Global == math.MaxInt || e.PerCat[cat-1] == math.MaxInt || e.Variants[form] == math.MaxInt {
		return nil, fmt.Errorf("lexer: frequency counter of %q at its maximum (%d)", form, math.MaxInt)
	}
	e.Variants[form]++
	e.PerCat[cat-1]++
	e.Global++
	return e, nil
}

// Lookup returns the lexeme of a surface form (the canonical form or any other
// variant of it).
func (t *Table) Lookup(form string) (*Entry, bool) {
	e, ok := t.byStem[Stem(form)]
	return e, ok
}

// Entries lists every lexeme ordered by canonical form, so that whatever is
// written from it is deterministic.
func (t *Table) Entries() []*Entry {
	type named struct {
		canon string
		e     *Entry
	}
	all := make([]named, 0, len(t.byStem))
	for _, e := range t.byStem {
		all = append(all, named{e.Canonical(), e})
	}
	sort.Slice(all, func(i, j int) bool { return all[i].canon < all[j].canon })
	out := make([]*Entry, len(all))
	for i, n := range all {
		out[i] = n.e
	}
	return out
}
