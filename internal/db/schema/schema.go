// Package schema holds the DDL of the two SQLite files of the case: corpus.db
// (baseline Anexo I, part A) and resultados.db (part B). The statements live in
// corpus.sql and resultados.sql and are embedded at build time.
//
// Deviation (ADR-007): both files use journal_mode=DELETE, not WAL. It is set
// through the connection PRAGMAs in internal/db/build, not here.
//
// Deviation (ADR-006): resultados.db declares uq_frecuencias_global, which restores
// uniqueness of the global (categoria_id IS NULL) row. The composite PRIMARY KEY
// of frecuencias cannot enforce it because SQLite treats NULLs as distinct in a
// PRIMARY KEY.
//
// Phase 1 addition: corpus.db also carries the particiones table, the
// materialized deterministic train/test assignment (see build.ApplySplit).
//
// frecuencias.categoria_id references corpus.db's categorias(id), but that
// cross-file check is not expressible in DDL: it lives in build.CheckCategoriaIntegrity
// (generation pipeline only, ADR-007).
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
