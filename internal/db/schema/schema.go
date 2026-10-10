// Package schema holds the DDL of the two SQLite files of the case: corpus.db
// (baseline Anexo I, part A) and resultados.db (part B). The statements live in
// corpus.sql and resultados.sql, written as in Anexo I (CREATE ... IF NOT EXISTS),
// and are embedded at build time.
//
// The Anexo I PRAGMAs are not part of the scripts: journal_mode (WAL for
// corpus.db, DELETE for resultados.db) and foreign_keys are set through the
// connection in internal/db/build.
//
// frecuencias.categoria_id references corpus.db's categorias(id), but that
// cross-file check is not expressible in DDL: it is verified with ATTACH
// DATABASE (build.CheckCategoriaIntegrity), as Anexo I states.
package schema

import (
	_ "embed"
	"strings"
)

//go:embed corpus.sql
var corpusSQL string

//go:embed resultados.sql
var resultadosSQL string

// Corpus returns the statements that create corpus.db, in application order.
// Each statement is returned exactly as written, trimmed of surrounding whitespace.
func Corpus() []string { return split(corpusSQL) }

// Resultados returns the statements that create resultados.db, in application order.
// Each statement is returned exactly as written, trimmed of surrounding whitespace.
func Resultados() []string { return split(resultadosSQL) }

// split breaks a script on ';' and drops empty pieces. SQLite stores the statement
// text verbatim in sqlite_master, so the trimmed text must match the source exactly.
func split(script string) []string {
	var out []string
	for _, piece := range strings.Split(script, ";") {
		if s := strings.TrimSpace(piece); s != "" {
			out = append(out, s)
		}
	}
	return out
}
