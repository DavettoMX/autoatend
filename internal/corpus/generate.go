package corpus

import (
	"autoatend/internal/corpus/templates"
	"fmt"
	"math/rand"
	"strings"
	"time"
	"unicode"
)

// Comment is one generated row of corpus.db. Noisy is generator metadata and is never persisted.
type Comment struct {
	Fecha       string
	Texto       string
	CategoriaID int
	Noisy       bool
}

// Config fully determines the generated corpus.
type Config struct {
	Seed   int64
	Total  int
	PerDay int
	Start  string // first day, YYYY-MM-DD
	Noise  NoiseConfig
}

// DefaultConfig is the case configuration: seed 10433, 73,000 comments,
// 100 per day over 730 days ending 2026-08-29.
func DefaultConfig() Config {
	return Config{Seed: 10433, Total: 73000, PerDay: 100, Start: "2024-08-30", Noise: DefaultNoise}
}

// categoryPercents is the class distribution, indexed by categoria_id-1.
var categoryPercents = [5]int{30, 20, 20, 20, 10}

// Probabilities of the optional composition parts.
const (
	pOpener     = 0.25
	pComplement = 0.40
	pCloser     = 0.30
	dateLayout  = "2006-01-02"
)

// targetCounts splits total by categoryPercents; any remainder goes to category 1.
func targetCounts(total int) [5]int {
	var out [5]int
	sum := 0
	for i, p := range categoryPercents {
		out[i] = total * p / 100
		sum += out[i]
	}
	out[0] += total - sum
	return out
}

// Generate produces the full corpus deterministically from cfg.Seed.
// The rng stream is consumed in a fixed order: category shuffle first, then
// one composition + noise pass per comment in chronological order.
func Generate(cfg Config) ([]Comment, error) {
	if cfg.Total <= 0 || cfg.PerDay <= 0 {
		return nil, fmt.Errorf("gencorpus: Total (%d) and PerDay (%d) must be positive", cfg.Total, cfg.PerDay)
	}
	start, err := time.Parse(dateLayout, cfg.Start)
	if err != nil {
		return nil, fmt.Errorf("gencorpus: start date: %w", err)
	}
	rng := rand.New(rand.NewSource(cfg.Seed))

	cats := make([]int, 0, cfg.Total)
	for i, n := range targetCounts(cfg.Total) {
		for j := 0; j < n; j++ {
			cats = append(cats, i+1)
		}
	}
	rng.Shuffle(len(cats), func(i, j int) { cats[i], cats[j] = cats[j], cats[i] })

	out := make([]Comment, len(cats))
	for i, cat := range cats {
		text, err := compose(cat, rng)
		if err != nil {
			return nil, err
		}
		text, noisy := injectNoise(text, cfg.Noise, rng)
		out[i] = Comment{
			Fecha:       start.AddDate(0, 0, i/cfg.PerDay).Format(dateLayout),
			Texto:       text,
			CategoriaID: cat,
			Noisy:       noisy,
		}
	}
	return out, nil
}

// compose builds one clean comment: optional opener, exactly one category base
// template, optional complement of the same category and optional closer.
// Opener, complement and closer never add a second incident: the comment talks
// about one problem.
func compose(cat int, rng *rand.Rand) (string, error) {
	ts := templates.Templates[cat-1]
	f := newSlotFiller(rng)
	var parts []string
	add := func(tpl string) error {
		s, err := f.fill(tpl)
		if err != nil {
			return err
		}
		parts = append(parts, capitalizeFirst(s))
		return nil
	}

	if rng.Float64() < pOpener {
		if err := add(pick(rng, templates.SharedOpeners, templates.Openers[cat-1])); err != nil {
			return "", err
		}
	}
	first := rng.Intn(len(ts))
	if err := add(ts[first]); err != nil {
		return "", err
	}
	if rng.Float64() < pComplement {
		cs := templates.Complements[cat-1]
		if err := add(cs[rng.Intn(len(cs))]); err != nil {
			return "", err
		}
	}
	if rng.Float64() < pCloser {
		if err := add(pick(rng, templates.SharedClosers, templates.Closers[cat-1])); err != nil {
			return "", err
		}
	}
	return strings.Join(parts, " "), nil
}

// pick draws uniformly from the concatenation of shared and own.
func pick(rng *rand.Rand, shared, own []string) string {
	i := rng.Intn(len(shared) + len(own))
	if i < len(shared) {
		return shared[i]
	}
	return own[i-len(shared)]
}

// capitalizeFirst uppercases the first letter, skipping leading punctuation such as "¿".
func capitalizeFirst(s string) string {
	rs := []rune(s)
	for i, r := range rs {
		if unicode.IsLetter(r) {
			rs[i] = unicode.ToUpper(r)
			break
		}
	}
	return string(rs)
}
