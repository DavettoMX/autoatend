// Command simulador sweeps data/corpus.db with the lexer and writes the lexical
// catalogue (lexemas, ER, frecuencias, variantes) to data/resultados.db, plus
// the document frequency of the training split to data/df_train.json. The
// pipeline lives in internal/sweep; this command only parses flags.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"autoatend/internal/sweep"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("simulador", flag.ContinueOnError)
	fs.SetOutput(stderr)
	corpus := fs.String("corpus", "data/corpus.db", "path of the corpus.db to read")
	out := fs.String("out", "data/resultados.db", "path of the resultados.db to create (must not exist)")
	df := fs.String("df", "data/df_train.json", "path of the df_train.json to create (must not exist; empty = skip)")
	seed := fs.Int64("seed", 10433, "seed recorded in corridas.semilla")
	if err := fs.Parse(args); err != nil {
		return 2
	}

	res, err := sweep.Run(sweep.Options{CorpusPath: *corpus, OutPath: *out, DFPath: *df, Seed: *seed})
	if err != nil {
		fmt.Fprintln(stderr, "simulador:", err)
		return 1
	}
	fmt.Fprintf(stdout, "simulador: wrote %s (%d lexemas, sha256 %s)\n", *out, res.Lexemas, res.SHA256)
	return 0
}
