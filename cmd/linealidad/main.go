// Command linealidad measures invariant I-4 (TE §6.3.3 and §6.3.8): the sweep time over the first
// N comments of data/corpus.db for N = 1k, 10k, 30k, 50k and 73k, a least-squares line, R2 and
// T(73k). The measurement lives in internal/sweep; this command only parses flags and writes
// docs/reports/i4_linealidad.md.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"autoatend/internal/sweep"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("linealidad", flag.ContinueOnError)
	fs.SetOutput(stderr)
	corpus := fs.String("corpus", "data/corpus.db", "path of the corpus.db to read (never modified)")
	report := fs.String("report", "docs/reports/i4_linealidad.md", "path of the Markdown report to write")
	sizesFlag := fs.String("sizes", "1000,10000,30000,50000,73000", "comma-separated subcorpus sizes")
	runs := fs.Int("runs", 3, "runs per size (the median is kept)")
	equipo := fs.String("equipo", fmt.Sprintf("%s, %d CPU, %s", runtime.GOARCH, runtime.NumCPU(), runtime.Version()), "description of the machine for the report")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	var sizes []int
	for _, s := range strings.Split(*sizesFlag, ",") {
		n, err := strconv.Atoi(strings.TrimSpace(s))
		if err != nil {
			fmt.Fprintf(stderr, "linealidad: invalid size %q\n", s)
			return 2
		}
		sizes = append(sizes, n)
	}

	work, err := os.MkdirTemp("", "linealidad-")
	if err != nil {
		fmt.Fprintln(stderr, "linealidad:", err)
		return 1
	}
	defer os.RemoveAll(work)

	l, err := sweep.MeasureLinearity(*corpus, sizes, *runs, work)
	if err != nil {
		fmt.Fprintln(stderr, "linealidad:", err)
		return 1
	}
	md := l.Markdown(*equipo, time.Now().Format("2006-01-02"))
	if err := os.WriteFile(*report, []byte(md), 0o644); err != nil {
		fmt.Fprintln(stderr, "linealidad:", err)
		return 1
	}
	for _, p := range l.Points {
		fmt.Fprintf(stdout, "N=%6d  T=%.3f s\n", p.N, p.Median.Seconds())
	}
	fmt.Fprintf(stdout, "linealidad: T(N) = %.4f + %.3e*N s; R2 = %.5f; wrote %s\n", l.Intercept, l.Slope, l.R2, *report)
	return 0
}
