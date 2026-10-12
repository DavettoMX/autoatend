package corpus

import (
	"fmt"
	"io"
	"strings"

	"autoatend/internal/corpus/templates"
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

	pf("\n## Generator parameters\n\n| Parameter | Value |\n|---|---:|\n")
	ts := templates.Templates
	pf("| Templates per category | %d / %d / %d / %d / %d |\n", len(ts[0]), len(ts[1]), len(ts[2]), len(ts[3]), len(ts[4]))
	pf("| Complements per category | %d / %d / %d / %d / %d |\n", len(templates.Complements[0]), len(templates.Complements[1]), len(templates.Complements[2]), len(templates.Complements[3]), len(templates.Complements[4]))
	pf("| Shared openers / closers | %d / %d |\n", len(templates.SharedOpeners), len(templates.SharedClosers))
	pf("| Own openers per category | %d / %d / %d / %d / %d |\n", len(templates.Openers[0]), len(templates.Openers[1]), len(templates.Openers[2]), len(templates.Openers[3]), len(templates.Openers[4]))
	pf("| Own closers per category | %d / %d / %d / %d / %d |\n", len(templates.Closers[0]), len(templates.Closers[1]), len(templates.Closers[2]), len(templates.Closers[3]), len(templates.Closers[4]))
	pf("| Opener / complement / closer probability | %.2f / %.2f / %.2f |\n", pOpener, pComplement, pCloser)
	pf("| Amount pool / surface formats | %d / %d |\n", len(templates.AmountPool), amountStyles)
	pf("| Folio prefixes / digits | %d / %d–%d |\n", len(templates.FolioPrefixes), folioMinDigits, folioMaxDigits)

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
	pf("\nEmpty texts: %d. Split rule (global, temporal): the cut-off is the fecha at position %d%% of the corpus in chronological order; train is fecha < %s, test is fecha >= %s, so every training comment precedes every test comment.\n", p.Empty, db.TrainPercent, p.Cutoff, p.Cutoff)
	nTest := p.Rows - sum5(p.Train)
	pf("\nTest comments with an identical text in train: %d of %d (%.2f%%). A temporal split does not prevent exact repeats of a template draw.\n",
		p.TestInTrain, nTest, 100*float64(p.TestInTrain)/float64(nTest))
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
