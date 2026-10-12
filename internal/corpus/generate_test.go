package corpus

import (
	"autoatend/internal/corpus/templates"
	"math/rand"
	"reflect"
	"strings"
	"testing"
)

func TestDefaultConfigIsTheCaseConfig(t *testing.T) {
	c := DefaultConfig()
	if c.Seed != 10433 {
		t.Errorf("Seed = %d, want 10433", c.Seed)
	}
	if c.Total != 73000 {
		t.Errorf("Total = %d, want 73000", c.Total)
	}
	if c.PerDay != 100 {
		t.Errorf("PerDay = %d, want 100", c.PerDay)
	}
}

func TestTargetCountsFollow30_20_20_20_10(t *testing.T) {
	cases := map[int][5]int{
		73000: {21900, 14600, 14600, 14600, 7300},
		1000:  {300, 200, 200, 200, 100},
		100:   {30, 20, 20, 20, 10},
	}
	for total, want := range cases {
		if got := targetCounts(total); got != want {
			t.Errorf("targetCounts(%d) = %v, want %v", total, got, want)
		}
	}
	// Non-multiples of 100: the remainder goes to the largest category and the sum is preserved.
	got := targetCounts(1003)
	sum := 0
	for _, n := range got {
		sum += n
	}
	if sum != 1003 {
		t.Errorf("targetCounts(1003) sums to %d", sum)
	}
}

func TestTemplateCountPerCategoryIsBetween30And60(t *testing.T) {
	if len(templates.Templates) != 5 {
		t.Fatalf("templates has %d categories, want 5", len(templates.Templates))
	}
	for i, ts := range templates.Templates {
		if len(ts) < 30 || len(ts) > 60 {
			t.Errorf("category %d has %d templates, want 30..60", i+1, len(ts))
		}
		seen := map[string]bool{}
		for _, tpl := range ts {
			if seen[tpl] {
				t.Errorf("category %d: duplicate template %q", i+1, tpl)
			}
			seen[tpl] = true
		}
	}
}

func TestEveryTemplateFillsCleanly(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	for ci, ts := range templates.Templates {
		for _, tpl := range ts {
			for k := 0; k < 20; k++ {
				out, err := newSlotFiller(rng).fill(tpl)
				if err != nil {
					t.Fatalf("category %d template %q: %v", ci+1, tpl, err)
				}
				if strings.TrimSpace(out) == "" || strings.ContainsAny(out, "{}") {
					t.Fatalf("category %d template %q produced %q", ci+1, tpl, out)
				}
			}
		}
	}
}

func TestFillRejectsUnknownAndUnterminatedSlots(t *testing.T) {
	f := newSlotFiller(rand.New(rand.NewSource(1)))
	if _, err := f.fill("hola {NOPE}"); err == nil {
		t.Error("unknown slot: expected error")
	}
	if _, err := f.fill("hola {PLAN"); err == nil {
		t.Error("unterminated slot: expected error")
	}
}

func TestSecondAmountSlotDiffersFromFirst(t *testing.T) {
	rng := rand.New(rand.NewSource(2))
	for k := 0; k < 500; k++ {
		f := newSlotFiller(rng)
		if _, err := f.fill("{MONTO} y {MONTO_B}"); err != nil {
			t.Fatal(err)
		}
		if f.amount == f.prevAmount {
			t.Fatalf("iteration %d: both amount slots drew %d", k, f.amount)
		}
	}
}

func smallConfig(seed int64) Config {
	c := DefaultConfig()
	c.Seed = seed
	c.Total = 2000
	return c
}

func TestGenerateIsDeterministic(t *testing.T) {
	a, err := Generate(smallConfig(10433))
	if err != nil {
		t.Fatal(err)
	}
	b, err := Generate(smallConfig(10433))
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(a, b) {
		t.Fatal("same seed produced different comment sets")
	}
	c, err := Generate(smallConfig(10434))
	if err != nil {
		t.Fatal(err)
	}
	if reflect.DeepEqual(a, c) {
		t.Fatal("different seeds produced identical comment sets")
	}
}

func TestGenerateShape(t *testing.T) {
	cfg := smallConfig(10433)
	cs, err := Generate(cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(cs) != cfg.Total {
		t.Fatalf("len = %d, want %d", len(cs), cfg.Total)
	}
	want := targetCounts(cfg.Total)
	var got [5]int
	perDay := map[string]int{}
	prev := ""
	for i, c := range cs {
		if c.CategoriaID < 1 || c.CategoriaID > 5 {
			t.Fatalf("comment %d: categoria_id %d", i, c.CategoriaID)
		}
		got[c.CategoriaID-1]++
		if len(c.Fecha) != 10 || c.Fecha[4] != '-' || c.Fecha[7] != '-' {
			t.Fatalf("comment %d: malformed fecha %q", i, c.Fecha)
		}
		if c.Fecha < prev {
			t.Fatalf("comment %d: dates not chronological (%s after %s)", i, c.Fecha, prev)
		}
		prev = c.Fecha
		perDay[c.Fecha]++
	}
	if got != want {
		t.Errorf("category counts = %v, want %v", got, want)
	}
	if cs[0].Fecha != "2024-08-30" {
		t.Errorf("first fecha = %s, want 2024-08-30", cs[0].Fecha)
	}
	for d, n := range perDay {
		if n != cfg.PerDay {
			t.Errorf("day %s has %d comments, want %d", d, n, cfg.PerDay)
		}
	}
}

func TestPlanDerivedSlotsNeedAPlanFirst(t *testing.T) {
	for _, tpl := range []string{"cuesta {PRECIO}", "cobran {COBRO}", "de más {EXTRA}"} {
		if _, err := newSlotFiller(rand.New(rand.NewSource(1))).fill(tpl); err == nil {
			t.Errorf("%q: expected an error when no {PLAN} was drawn before", tpl)
		}
	}
}

func TestPlanDerivedSlotsAgreeWithThePlanDrawn(t *testing.T) {
	rng := rand.New(rand.NewSource(3))
	for k := 0; k < 2000; k++ {
		f := newSlotFiller(rng)
		out, err := f.fill("{PLAN}|{PRECIO}|{COBRO}|{EXTRA}")
		if err != nil {
			t.Fatal(err)
		}
		parts := strings.Split(out, "|")
		price, _ := firstAmount(parts[1])
		charge, _ := firstAmount(parts[2])
		extra, _ := firstAmount(parts[3])
		if want := planPrice[parts[0]]; price != want {
			t.Fatalf("%q: price %v, want %v", out, price, want)
		}
		if charge <= price || charge > 2*price || extra >= price || extra <= 0 {
			t.Fatalf("%q: charge %v and extra %v do not bracket the price %v", out, charge, extra, price)
		}
	}
}

func TestSmallAmountSlotStaysSmall(t *testing.T) {
	rng := rand.New(rand.NewSource(4))
	for k := 0; k < 1000; k++ {
		out, err := newSlotFiller(rng).fill("{MONTO_CHICO}")
		if err != nil {
			t.Fatal(err)
		}
		if v, ok := firstAmount(out); !ok || v >= 500 {
			t.Fatalf("small amount %q reaches 500", out)
		}
	}
	// Recharges and recurring add-on charges never reach the big amounts of the pool.
	for _, cat := range []int{1, 2} {
		for _, text := range cleanTexts(t, cat, 10000, int64(700+cat)) {
			if v, ok := amountAfter(text, "recarga de "); ok && v >= 500 {
				t.Fatalf("a recharge of %v is not plausible: %q", v, text)
			}
		}
	}
}

func TestGenerateRejectsInvalidConfig(t *testing.T) {
	c := DefaultConfig()
	c.Total = 0
	if _, err := Generate(c); err == nil {
		t.Error("Total=0: expected error")
	}
	c = DefaultConfig()
	c.PerDay = 0
	if _, err := Generate(c); err == nil {
		t.Error("PerDay=0: expected error")
	}
}
