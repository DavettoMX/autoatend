// Command casos renders docs/CASOS_PRUEBA.md (the Anexo II test-case matrix: case -> Go test ->
// measured result) from internal/assets/ui/casos_prueba.json, after checking the two-way
// traceability against the *_test.go files. The logic lives in internal/casos; this command only
// parses flags.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"autoatend/internal/casos"
)

func main() { os.Exit(run(os.Args[1:], os.Stdout, os.Stderr)) }

func run(args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("casos", flag.ContinueOnError)
	fs.SetOutput(stderr)
	in := fs.String("in", "internal/assets/ui/casos_prueba.json", "path of the matrix JSON")
	out := fs.String("out", "docs/CASOS_PRUEBA.md", "path of the Markdown to write")
	root := fs.String("root", ".", "module root, scanned for *_test.go files")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	raw, err := os.ReadFile(*in)
	if err != nil {
		fmt.Fprintln(stderr, "casos:", err)
		return 1
	}
	m, err := casos.Parse(raw)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	problems, err := casos.Trazabilidad(m, *root)
	if err != nil {
		fmt.Fprintln(stderr, "casos:", err)
		return 1
	}
	if len(problems) > 0 {
		for _, p := range problems {
			fmt.Fprintln(stderr, "casos: traceability:", p)
		}
		return 1
	}
	if err := os.WriteFile(*out, []byte(m.Markdown()), 0o644); err != nil {
		fmt.Fprintln(stderr, "casos:", err)
		return 1
	}
	anexo, ext, pasan := m.Resumen()
	fmt.Fprintf(stdout, "casos: wrote %s (%d Anexo II cases + %d extension, %d pass)\n", *out, anexo, ext, pasan)
	return 0
}
