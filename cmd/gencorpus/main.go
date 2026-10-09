// Command gencorpus generates the synthetic TelArg corpus (seed 10433) and
// writes it, with its persisted train/test split, to corpus.db.
package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	out := flag.String("out", "corpus.db", "path of the corpus.db to create (must not exist)")
	report := flag.String("report", "docs/generator_validation_report.md", "path of the validation report")
	flag.Parse()

	res, err := run(Options{Config: DefaultConfig(), OutPath: *out, ReportPath: *report})
	if err != nil {
		fmt.Fprintln(os.Stderr, "gencorpus:", err)
		os.Exit(1)
	}
	fmt.Printf("gencorpus: wrote %s (%d rows, sha256 %s)\n", *out, res.Persisted.Rows, res.Persisted.SHA256)
	fmt.Printf("gencorpus: validation report at %s\n", *report)
}
