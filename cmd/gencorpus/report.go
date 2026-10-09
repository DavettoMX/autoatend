package main

import (
	"fmt"
	"io"
	"strings"

	"autoatend/internal/db"
)

// writeReport renders the generator validation report. The output is a pure
// function of cfg and res: no timestamps, no map iteration, so it is reproducible.
func writeReport(w io.Writer, cfg Config, res Result) error {
	v, p := res.Validation, res.Persisted
	var b strings.Builder
	pf := func(format string, args ...any) { fmt.Fprintf(&b, format, args...) }

	pf("# Generator validation report — corpus.db\n\n")
	pf("Produced by `cmd/gencorpus`. Seed **%d**, %d comments, %d per day from %s (%s to %s in the file).\n",
		cfg.Seed, cfg.Total, cfg.PerDay, cfg.Start, p.FirstDate, p.LastDate)
	pf("The corpus is persisted only if every check of the battery passes (verdict: **%s**).\n\n", verdict(v.OK()))

	pf("## Class distribution\n\n| Category | Target | Generated | Persisted | Share |\n|---|---:|---:|---:|---:|\n")
	for i, c := range db.Categories {
		pf("| %s | %d | %d | %d | %.2f%% |\n", c.Nombre, v.Targets[i], v.Counts[i], p.PerCategory[i],
			100*float64(p.PerCategory[i])/float64(p.Rows))
	}
	pf("| **Total** | %d | %d | %d | 100.00%% |\n\n", cfg.Total, sum5(v.Counts), p.Rows)

	pf("## Validation battery\n\n| Check | Result | Detail |\n|---|---|---|\n")
	for _, c := range v.Checks {
		pf("| %s | %s | %s |\n", c.Name, verdict(c.Pass), c.Detail)
	}
	pf("\nCorpus statistics: %.2f%% duplicate texts; text length %d / %.1f / %d runes (min / mean / max).\n\n",
		100*v.DuplicateRate, v.Length.Min, v.Length.Mean, v.Length.Max)

	pf("## Lexical overlap\n\nJaccard similarity of the %d most frequent tokens per category (accepted range for distinct pairs: [%.2f, %.2f]).\n\n",
		overlapTopK, overlapMin, overlapMax)
	pf("| |")
	for _, c := range db.Categories {
		pf(" %s |", c.Nombre)
	}
	pf("\n|---|---:|---:|---:|---:|---:|\n")
	for i, c := range db.Categories {
		pf("| %s |", c.Nombre)
		for j := 0; j < 5; j++ {
			pf(" %.3f |", v.Overlap[i][j])
		}
		pf("\n")
	}

	n := cfg.Noise
	pf("\n## Noise calibration\n\n| Parameter | Value |\n|---|---:|\n")
	pf("| Comment selection rate | %.2f |\n| Tilde drop (per accented letter) | %.2f |\n| Abbreviation q/x/tel/fact (per word) | %.2f |\n", n.Rate, n.TildeDrop, n.Abbrev)
	pf("| Keyboard typo (per word of 4+ letters) | %.2f |\n| Case jitter (per comment) | %.2f |\n| Irregular punctuation (per comment) | %.2f |\n", n.Typo, n.CaseJitter, n.Punct)
	pf("| **Realized selection rate** | **%.4f** (tolerance ±%.3f) |\n\n", v.RealizedNoise, noiseTolerance)

	pf("## Persisted corpus verification\n\nMeasured by re-opening the finished file read-only.\n\n")
	pf("| Category | Rows | Train | Test | max(train.fecha) | min(test.fecha) |\n|---|---:|---:|---:|---|---|\n")
	for i, c := range db.Categories {
		pf("| %s | %d | %d | %d | %s | %s |\n", c.Nombre, p.PerCategory[i], p.Train[i], p.Test[i], p.MaxTrainDate[i], p.MinTestDate[i])
	}
	pf("\nEmpty texts: %d. Split rule: per category, ordered by (fecha, id), oldest %d%% train, newest remainder test.\n", p.Empty, db.TrainPercent)
	pf("\ncorpus.db SHA-256: `%s`\n\n", p.SHA256)

	pf("## Sample review\n\n%d comments per category, drawn with a dedicated rng stream (seed+1) for manual inspection.\n", sampleSize)
	for i, c := range db.Categories {
		pf("\n### %s\n\n", c.Nombre)
		for k, t := range v.Samples[i] {
			pf("%d. %s\n", k+1, strings.ReplaceAll(t, "\n", " "))
		}
	}

	_, err := io.WriteString(w, b.String())
	return err
}

func verdict(ok bool) string {
	if ok {
		return "PASS"
	}
	return "FAIL"
}

func sum5(a [5]int) int {
	s := 0
	for _, n := range a {
		s += n
	}
	return s
}
