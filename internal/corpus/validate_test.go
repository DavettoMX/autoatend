package corpus

import (
	"reflect"
	"strings"
	"testing"
	"time"
)

func validationConfig() Config {
	c := DefaultConfig()
	c.Total = 5000
	return c
}

func generated(t *testing.T) ([]Comment, Config) {
	t.Helper()
	cfg := validationConfig()
	cs, err := Generate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	return cs, cfg
}

func checkByName(t *testing.T, v Validation, name string) Check {
	t.Helper()
	for _, c := range v.Checks {
		if c.Name == name {
			return c
		}
	}
	t.Fatalf("battery has no check named %q (have %v)", name, checkNames(v))
	return Check{}
}

func checkNames(v Validation) []string {
	var n []string
	for _, c := range v.Checks {
		n = append(n, c.Name)
	}
	return n
}

func TestValidationPassesOnGeneratedCorpus(t *testing.T) {
	cs, cfg := generated(t)
	v := Validate(cs, cfg)
	for _, c := range v.Checks {
		if !c.Pass {
			t.Errorf("check %q failed: %s", c.Name, c.Detail)
		}
	}
	if !v.OK() {
		t.Fatal("OK() = false on a freshly generated corpus")
	}
	want := []string{"category counts", "no empty texts", "no unfilled slots", "dates", "noise rate", "lexical overlap", "sample review"}
	if !reflect.DeepEqual(checkNames(v), want) {
		t.Errorf("check names = %v, want %v", checkNames(v), want)
	}
}

func TestValidationDetectsWrongCategoryCount(t *testing.T) {
	cs, cfg := generated(t)
	cs = cs[:len(cs)-1]
	v := Validate(cs, cfg)
	if checkByName(t, v, "category counts").Pass || v.OK() {
		t.Fatal("dropping a row must fail the category-count check")
	}
}

func TestValidationDetectsEmptyText(t *testing.T) {
	cs, cfg := generated(t)
	cs[10].Texto = "   \t"
	v := Validate(cs, cfg)
	if checkByName(t, v, "no empty texts").Pass || v.OK() {
		t.Fatal("whitespace-only text must fail")
	}
}

func TestValidationDetectsUnfilledSlot(t *testing.T) {
	cs, cfg := generated(t)
	cs[10].Texto = "Me cobraron {MONTO} por nada."
	v := Validate(cs, cfg)
	if checkByName(t, v, "no unfilled slots").Pass || v.OK() {
		t.Fatal("leftover slot must fail")
	}
}

func TestValidationDetectsMalformedOrUnorderedDates(t *testing.T) {
	cs, cfg := generated(t)
	cs[3].Fecha = "2024/08/30"
	if v := Validate(cs, cfg); checkByName(t, v, "dates").Pass {
		t.Fatal("malformed date must fail")
	}
	cs, cfg = generated(t)
	cs[0], cs[len(cs)-1] = cs[len(cs)-1], cs[0] // breaks chronology and per-day counts
	if v := Validate(cs, cfg); checkByName(t, v, "dates").Pass {
		t.Fatal("non-chronological dates must fail")
	}
}

func TestValidationDetectsMiscalibratedNoise(t *testing.T) {
	cs, cfg := generated(t)
	for i := range cs {
		cs[i].Noisy = false
	}
	v := Validate(cs, cfg)
	if checkByName(t, v, "noise rate").Pass {
		t.Fatal("0% realized noise vs 30% configured must fail")
	}
	if v.RealizedNoise != 0 {
		t.Errorf("RealizedNoise = %v, want 0", v.RealizedNoise)
	}
}

// synthetic builds a corpus that satisfies every check except possibly the lexical one.
func synthetic(textFor func(cat, i int) string) ([]Comment, Config) {
	cfg := DefaultConfig()
	cfg.Total = 1000
	cfg.PerDay = 10
	start, _ := time.Parse(dateLayout, cfg.Start)
	var cs []Comment
	for cat, cnt := range targetCounts(cfg.Total) {
		for j := 0; j < cnt; j++ {
			cs = append(cs, Comment{CategoriaID: cat + 1, Texto: textFor(cat+1, j)})
		}
	}
	for i := range cs {
		cs[i].Fecha = start.AddDate(0, 0, i/cfg.PerDay).Format(dateLayout)
		cs[i].Noisy = i%100 < 30
	}
	return cs, cfg
}

func words(prefix string, n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		b.WriteString(prefix)
		b.WriteString(strings.Repeat(string(rune('a'+i%26)), 1+i/26))
		b.WriteByte(' ')
	}
	return b.String()
}

func TestValidationRejectsIdenticalVocabularyAcrossCategories(t *testing.T) {
	shared := words("comun", 60)
	cs, cfg := synthetic(func(cat, i int) string { return shared + "texto" })
	v := Validate(cs, cfg)
	if c := checkByName(t, v, "lexical overlap"); c.Pass {
		t.Fatalf("identical vocabularies must exceed the overlap upper bound: %s", c.Detail)
	}
}

func TestValidationRejectsDisjointVocabularies(t *testing.T) {
	cs, cfg := synthetic(func(cat, i int) string {
		return words(strings.Repeat(string(rune('a'+cat)), 3), 60)
	})
	v := Validate(cs, cfg)
	if c := checkByName(t, v, "lexical overlap"); c.Pass {
		t.Fatalf("disjoint vocabularies (trivially separable) must fall below the overlap lower bound: %s", c.Detail)
	}
}

func TestValidationSampleIs50PerCategoryAndDeterministic(t *testing.T) {
	cs, cfg := generated(t)
	a, b := Validate(cs, cfg), Validate(cs, cfg)
	for i := 0; i < 5; i++ {
		if len(a.Samples[i]) != 50 {
			t.Errorf("category %d sample size = %d, want 50", i+1, len(a.Samples[i]))
		}
	}
	if !reflect.DeepEqual(a.Samples, b.Samples) {
		t.Error("sample selection is not deterministic")
	}
	cfg2 := cfg
	cfg2.Seed++
	if c := Validate(cs, cfg2); reflect.DeepEqual(a.Samples, c.Samples) {
		t.Error("sample selection ignores the seed")
	}
}

func TestValidationSampleReviewFlagsDegenerateTexts(t *testing.T) {
	cs, cfg := generated(t)
	for i := range cs {
		cs[i].Texto = "ok"
	}
	v := Validate(cs, cfg)
	if checkByName(t, v, "sample review").Pass {
		t.Fatal("two-letter texts must fail the sample review")
	}
}

func TestValidationReportsDuplicateRateAndOverlapMatrix(t *testing.T) {
	cs, cfg := generated(t)
	v := Validate(cs, cfg)
	if v.DuplicateRate < 0 || v.DuplicateRate >= 1 {
		t.Fatalf("DuplicateRate = %v", v.DuplicateRate)
	}
	if v.Length.Min < 1 || v.Length.Max <= v.Length.Min || v.Length.Mean < float64(v.Length.Min) || v.Length.Mean > float64(v.Length.Max) {
		t.Errorf("Length = %+v, want heterogeneous lengths with min < mean < max", v.Length)
	}
	for i := 0; i < 5; i++ {
		if v.Overlap[i][i] != 1 {
			t.Errorf("Overlap[%d][%d] = %v, want 1", i, i, v.Overlap[i][i])
		}
		for j := 0; j < 5; j++ {
			if v.Overlap[i][j] != v.Overlap[j][i] || v.Overlap[i][j] < 0 || v.Overlap[i][j] > 1 {
				t.Errorf("Overlap[%d][%d] = %v (symmetric: %v)", i, j, v.Overlap[i][j], v.Overlap[j][i])
			}
		}
	}
}
