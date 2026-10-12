package corpus

import (
	"fmt"
	"strings"
	"testing"

	"autoatend/internal/corpus/templates"
)

// kindRange returns the [lo, hi) index range of one template kind in Templates[cat].
func kindRange(cat int, kind string) (int, int) {
	b, c := templates.BaseCount[cat], templates.CrossCount[cat]
	switch kind {
	case "base":
		return 0, b
	case "cross":
		return b, b + c
	}
	return b + c, len(templates.Templates[cat])
}

func TestTemplateKinds_CountsWithinTE(t *testing.T) {
	for cat, ts := range templates.Templates {
		if n := len(ts); n < 30 || n > 60 {
			t.Errorf("category %d has %d templates, TE §6.1.4 requires 30–60", cat+1, n)
		}
		if c := templates.CrossCount[cat]; c < 5 || c > 6 {
			t.Errorf("category %d has %d cross-context templates, want 5–6", cat+1, c)
		}
		if templates.BaseCount[cat]+templates.CrossCount[cat] > len(ts) {
			t.Errorf("category %d: BaseCount+CrossCount exceeds the template list", cat+1)
		}
	}
}

func TestTemplateKinds_SharedTextsAreExactlyThePlannedPairs(t *testing.T) {
	owners := map[string][]int{}
	for cat, ts := range templates.Templates {
		for _, tpl := range ts {
			owners[tpl] = append(owners[tpl], cat+1)
		}
	}
	pairs := map[string]int{}
	for cat := range templates.Templates {
		for _, kind := range []string{"base", "cross", "shared"} {
			lo, hi := kindRange(cat, kind)
			for _, tpl := range templates.Templates[cat][lo:hi] {
				o := owners[tpl]
				if kind != "shared" && len(o) != 1 {
					t.Errorf("%s template of category %d also belongs to %v: %q", kind, cat+1, o, tpl)
				}
				if kind == "shared" {
					if len(o) != 2 {
						t.Errorf("shared template of category %d belongs to %v, want exactly 2: %q", cat+1, o, tpl)
						continue
					}
					if o[0] == cat+1 {
						pairs[fmt.Sprintf("%d-%d", o[0], o[1])]++
					}
				}
			}
		}
	}
	// PLAN: Falta↔Quejas 3, Cobros↔Facturación 3, Quejas↔Reconocimientos 2, Quejas↔Cobros 1, Quejas↔Facturación 1.
	want := map[string]int{"1-4": 3, "2-3": 3, "4-5": 2, "2-4": 1, "3-4": 1}
	if fmt.Sprint(pairs) != fmt.Sprint(want) {
		t.Errorf("shared pairs = %v, want %v", pairs, want)
	}
}

func TestTemplateKinds_RefFolioFrameInEveryCategory(t *testing.T) {
	for cat, ts := range templates.Templates {
		found := false
		for _, tpl := range ts {
			found = found || strings.Contains(tpl, "(ref. {FOLIO})")
		}
		if !found {
			t.Errorf("category %d has no \"(ref. {FOLIO})\" template", cat+1)
		}
	}
}
