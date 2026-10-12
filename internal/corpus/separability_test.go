package corpus

import (
	"fmt"
	"math"
	"strings"
	"testing"
)

// Pre-registered separability criterion K1–K3 (PLAN.md §3, 2026-10-12).

// stopwordPrefix is made of stoplist v1 entries only.
const stopwordPrefix = "de la el en que y a los las por con del al un una unos unas lo sin hasta desde " +
	"durante mediante ante tras contra hacia pero sino porque pues o u e ni "

func TestK1_OverlapIsMeasuredAfterStoplist(t *testing.T) {
	// Categories share only stopwords; without the stoplist the shared function words
	// alone would put the Jaccard inside [0.15, 0.80].
	cs, cfg := synthetic(func(cat, i int) string {
		return stopwordPrefix + words(strings.Repeat(string(rune('k'+cat)), 3), 20)
	})
	v := Validate(cs, cfg)
	if c := checkByName(t, v, "lexical overlap"); c.Pass {
		t.Fatalf("categories that share only stopwords must fail K1: %s", c.Detail)
	}
	if v.Overlap[0][1] != 0 {
		t.Errorf("Overlap[0][1] = %v, want 0 once stopwords are removed", v.Overlap[0][1])
	}
}

// collisionCorpus gives share[k] of category k+1 the same base sentence, shared by
// every category; the rest get a base sentence of their own category.
func collisionCorpus(share [5]float64) ([]Comment, Config) {
	cs, cfg := synthetic(func(cat, i int) string { return "" })
	seen := [5]int{}
	counts := targetCounts(cfg.Total)
	for i := range cs {
		k := cs[i].CategoriaID - 1
		own := "propio " + strings.Repeat(string(rune('a'+k)), 3) + " " + words("w", 1+seen[k]%7)
		if float64(seen[k]) < share[k]*float64(counts[k]) {
			cs[i].Base = "Sigo esperando una respuesta clara"
		} else {
			cs[i].Base = own
		}
		cs[i].Texto = cs[i].Base + " " + words(strings.Repeat(string(rune('k'+k)), 2), 12)
		seen[k]++
	}
	return cs, cfg
}

func TestK2_CollisionRateWithinPreRegisteredRange(t *testing.T) {
	share := [5]float64{0.06, 0.06, 0.06, 0.06, 0.04}
	cs, cfg := collisionCorpus(share)
	v := Validate(cs, cfg)
	for k := 0; k < 5; k++ {
		if math.Abs(v.Collision[k]-share[k]) > 1e-12 {
			t.Errorf("Collision[%d] = %v, want %v", k, v.Collision[k], share[k])
		}
	}
	// Oracle: 0.3*0.06 + 0.2*0.06*3 + 0.1*0.04 = 0.058.
	if math.Abs(v.CollisionGlobal-0.058) > 1e-12 {
		t.Errorf("CollisionGlobal = %v, want 0.058", v.CollisionGlobal)
	}
	if c := checkByName(t, v, "base collision"); !c.Pass {
		t.Fatalf("K2 must pass inside the range: %s", c.Detail)
	}
}

func TestK2_CollisionRateOutsideRangeFails(t *testing.T) {
	cases := map[string][5]float64{
		"no collisions":              {0, 0, 0, 0, 0},
		"global above 0.10":          {0.12, 0.12, 0.12, 0.12, 0.04},
		"one category below 0.03":    {0.07, 0.07, 0.02, 0.07, 0.04},
		"Reconocimientos above 0.05": {0.06, 0.06, 0.06, 0.06, 0.06},
	}
	for name, share := range cases {
		cs, cfg := collisionCorpus(share)
		if c := checkByName(t, Validate(cs, cfg), "base collision"); c.Pass {
			t.Errorf("%s: K2 must fail: %s", name, c.Detail)
		}
	}
}

func TestK2_IdenticalBaseInsideOneCategoryIsNotACollision(t *testing.T) {
	cs, cfg := synthetic(func(cat, i int) string { return "" })
	for i := range cs {
		cs[i].Base = "mismo texto " + strings.Repeat(string(rune('a'+cs[i].CategoriaID)), 3)
		cs[i].Texto = cs[i].Base
	}
	if v := Validate(cs, cfg); v.CollisionGlobal != 0 {
		t.Errorf("CollisionGlobal = %v, want 0", v.CollisionGlobal)
	}
}

func TestK3_ProbeRejectsTrivialAndIndistinguishableCorpora(t *testing.T) {
	disjoint, cfg := synthetic(func(cat, i int) string {
		return words(strings.Repeat(string(rune('a'+cat)), 3), 8+i%5)
	})
	v := Validate(disjoint, cfg)
	if v.ProbeF1 != 1 {
		t.Errorf("disjoint vocabularies: ProbeF1 = %v, want 1", v.ProbeF1)
	}
	if c := checkByName(t, v, "probe classifier"); c.Pass {
		t.Errorf("trivially separable corpus must fail K3: %s", c.Detail)
	}
	same, cfg := synthetic(func(cat, i int) string { return words("comun", 8+i%5) })
	v = Validate(same, cfg)
	if v.ProbeF1 >= probeMin {
		t.Errorf("identical vocabularies: ProbeF1 = %v, want < %v", v.ProbeF1, probeMin)
	}
	if c := checkByName(t, v, "probe classifier"); c.Pass {
		t.Errorf("indistinguishable corpus must fail K3: %s", c.Detail)
	}
}

func TestK3_ProbePassesWithSomeAmbiguousComments(t *testing.T) {
	// One comment in ten is the same generic text in every category.
	cs, cfg := synthetic(func(cat, i int) string {
		if i%10 == 0 {
			return "sigo esperando respuesta nadie resuelve"
		}
		return words(strings.Repeat(string(rune('a'+cat)), 3), 6+i%5)
	})
	v := Validate(cs, cfg)
	if c := checkByName(t, v, "probe classifier"); !c.Pass {
		t.Fatalf("K3 must pass: %s", c.Detail)
	}
	if v2 := Validate(cs, cfg); v2.ProbeF1 != v.ProbeF1 {
		t.Errorf("probe is not deterministic: %v vs %v", v.ProbeF1, v2.ProbeF1)
	}
}

func TestK3_MacroF1Oracle(t *testing.T) {
	// Hand-computed: class 1 P=2/3 R=1 F1=0.8; class 2 P=1 R=0.5 F1=2/3; classes 3-5 absent
	// from gold and pred are skipped; Macro-F1 = (0.8 + 2/3) / 2.
	gold := []int{1, 1, 2, 2}
	pred := []int{1, 1, 1, 2}
	if got, want := macroF1(gold, pred), (0.8+2.0/3)/2; math.Abs(got-want) > 1e-12 {
		t.Errorf("macroF1 = %v, want %v", got, want)
	}
}

// TestPreRegistration_CurrentTemplatesFailK1K2K3 is the evidence that the criterion
// precedes the template changes: the corpus of b810bad2… fails all three.
func TestPreRegistration_CurrentTemplatesFailK1K2K3(t *testing.T) {
	if testing.Short() {
		t.Skip("generates the full corpus")
	}
	cfg := DefaultConfig()
	cs, err := Generate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	v := Validate(cs, cfg)
	for _, name := range []string{"lexical overlap", "base collision", "probe classifier"} {
		c := checkByName(t, v, name)
		t.Log(fmt.Sprintf("%s: %s", name, c.Detail))
		if c.Pass {
			t.Errorf("%s passes on the current templates; the pre-registration expects a failure", name)
		}
	}
}
