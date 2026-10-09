package main

import (
	"math"
	"math/rand"
	"strings"
	"testing"
	"unicode"
)

func only(set func(*NoiseConfig)) NoiseConfig {
	c := NoiseConfig{Rate: 1}
	set(&c)
	return c
}

func TestNoiseRateZeroLeavesTextUntouched(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := DefaultNoise
	cfg.Rate = 0
	in := "Me cobraron de más en la factura."
	for i := 0; i < 100; i++ {
		out, noisy := injectNoise(in, cfg, rng)
		if out != in || noisy {
			t.Fatalf("Rate=0 altered text: %q noisy=%v", out, noisy)
		}
	}
}

func TestNoiseRateOneAlwaysChangesText(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := DefaultNoise
	cfg.Rate = 1
	for _, in := range []string{"Hola.", "No tengo señal en mi casa.", "Gracias por todo.", "Atención rápida y amable."} {
		for i := 0; i < 200; i++ {
			out, noisy := injectNoise(in, cfg, rng)
			if !noisy || out == in {
				t.Fatalf("Rate=1 left %q unchanged (noisy=%v)", in, noisy)
			}
		}
	}
}

func TestNoiseDropsTildesButKeepsEnye(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := only(func(c *NoiseConfig) { c.TildeDrop = 1 })
	out, _ := injectNoise("Facturación móvil está aquí, niño Ángel", cfg, rng)
	if out != "Facturacion movil esta aqui, niño Angel" {
		t.Fatalf("got %q", out)
	}
}

func TestNoiseAbbreviatesWholeWordsOnly(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := only(func(c *NoiseConfig) { c.Abbrev = 1 })
	out, _ := injectNoise("quiero que me llamen por teléfono sobre mi factura, no por el portal, porque sí", cfg, rng)
	want := "quiero q me llamen x tel sobre mi fact, no x el portal, porque sí"
	if out != want {
		t.Fatalf("got  %q\nwant %q", out, want)
	}
}

func TestNoiseCaseJitterOnlyChangesCase(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := only(func(c *NoiseConfig) { c.CaseJitter = 1 })
	in := "Hola Mundo Cruel"
	seen := map[string]bool{}
	for i := 0; i < 300; i++ {
		out, _ := injectNoise(in, cfg, rng)
		if out == in || !strings.EqualFold(out, in) {
			t.Fatalf("case jitter produced %q from %q", out, in)
		}
		seen[out] = true
	}
	if !seen["HOLA MUNDO CRUEL"] || !seen["hola mundo cruel"] {
		t.Errorf("expected both ALLCAPS and all-lowercase modes among outputs, got %d distinct", len(seen))
	}
}

func TestNoiseTyposStayLocal(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := only(func(c *NoiseConfig) { c.Typo = 1 })
	in := "necesito revisar urgente mi servicio"
	words := len(strings.Fields(in))
	for i := 0; i < 200; i++ {
		out, _ := injectNoise(in, cfg, rng)
		if out == in {
			t.Fatalf("Typo=1 left text unchanged")
		}
		if math.Abs(float64(len(out)-len(in))) > float64(words) {
			t.Fatalf("typos changed length too much: %q -> %q", in, out)
		}
		if len(strings.Fields(out)) != words {
			t.Fatalf("typos changed word count: %q", out)
		}
	}
}

func lettersOnly(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func TestNoiseIrregularPunctuationKeepsLetters(t *testing.T) {
	rng := rand.New(rand.NewSource(1))
	cfg := only(func(c *NoiseConfig) { c.Punct = 1 })
	in := "No funciona mi internet."
	for i := 0; i < 100; i++ {
		out, _ := injectNoise(in, cfg, rng)
		if out == in {
			t.Fatal("Punct=1 left text unchanged")
		}
		if lettersOnly(out) != lettersOnly(in) {
			t.Fatalf("punctuation noise altered letters: %q", out)
		}
	}
}

func TestNoiseIsDeterministic(t *testing.T) {
	run := func() []string {
		rng := rand.New(rand.NewSource(42))
		var outs []string
		for i := 0; i < 200; i++ {
			o, _ := injectNoise("Me cobraron de más por teléfono, qué problema con la factura.", DefaultNoise, rng)
			outs = append(outs, o)
		}
		return outs
	}
	a, b := run(), run()
	for i := range a {
		if a[i] != b[i] {
			t.Fatalf("draw %d differs: %q vs %q", i, a[i], b[i])
		}
	}
}

func TestNoiseRealizedRateMatchesConfiguredRate(t *testing.T) {
	rng := rand.New(rand.NewSource(10433))
	n, noisyCount := 20000, 0
	in := "No tengo señal desde ayer en mi colonia y quiero que me expliquen por qué."
	for i := 0; i < n; i++ {
		out, noisy := injectNoise(in, DefaultNoise, rng)
		if noisy {
			noisyCount++
			if out == in {
				t.Fatal("flagged noisy but text unchanged")
			}
		}
	}
	got := float64(noisyCount) / float64(n)
	if math.Abs(got-DefaultNoise.Rate) > 0.015 {
		t.Fatalf("realized noise rate %.4f, configured %.2f (tolerance 0.015)", got, DefaultNoise.Rate)
	}
}
