// Command gencorpus generates the synthetic TelArg corpus (seed 10433) and
// writes it, with its persisted train/test split, to data/corpus.db, and the
// generator validation report to docs/reports. The pipeline lives in
// internal/corpus; this command only parses flags.
package main

import (
	"flag"
	"fmt"
	"os"

	"autoatend/internal/corpus"
)

func main() {
	out := flag.String("out", "data/corpus.db", "path of the corpus.db to create (must not exist)")
	report := flag.String("report", "docs/reports/generator_validation_report.md", "path of the validation report")
	ambiguous := flag.String("ambiguos", "data/corpus_ambiguos.json", "path of the ids of the ambiguous (K2 collision) subset")
	seed := flag.Int64("seed", 10433, "random seed")
	total := flag.Int("total", 73000, "approximate total comments")
	noiseRate := flag.Float64("noise-rate", 0.30, "lexical noise rate")
	flag.Parse()

	cfg := corpus.DefaultConfig()
	cfg.Seed = *seed
	cfg.Total = *total
	cfg.Noise.Rate = *noiseRate

	res, err := corpus.Run(corpus.Options{Config: cfg, OutPath: *out, ReportPath: *report, AmbiguousPath: *ambiguous})
	if err != nil {
		fmt.Fprintln(os.Stderr, "gencorpus:", err)
		os.Exit(1)
	}
	fmt.Printf("gencorpus: wrote %s (%d rows, sha256 %s)\n", *out, res.Persisted.Rows, res.Persisted.SHA256)
	fmt.Printf("gencorpus: validation report at %s\n", *report)
}
