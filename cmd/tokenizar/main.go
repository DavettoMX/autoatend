// Command tokenizar maps comments to catalogue lexemes for the Python trainer
// (CPA06, TE §6.3.6): JSON lines {"id":..,"texto":..} on stdin, JSON lines
// {"id":..,"lexemas":[..]} on stdout. The logic lives in internal/ml; this
// command only parses flags.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"autoatend/internal/ml"
)

func main() { os.Exit(run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr)) }

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("tokenizar", flag.ContinueOnError)
	fs.SetOutput(stderr)
	cat := fs.String("catalogo", "data/resultados.db", "path of the resultados.db whose lexemas table is the lexicon")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	lex, err := ml.LoadCatalogue(*cat)
	if err != nil {
		fmt.Fprintln(stderr, "tokenizar:", err)
		return 1
	}
	if err := ml.RunTokenizar(stdin, stdout, lex); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	return 0
}
