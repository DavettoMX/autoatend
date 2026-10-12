package build

import (
	"database/sql"

	"autoatend/internal/db"
)

// CheckCategoriaIntegrity verifies the logical cross-file foreign key
// frecuencias.categoria_id -> corpus.categorias(id) (TE §6.1.3). The read-only
// check itself lives in internal/db, because the runtime performs it at startup
// and must not import this generation-only package; this wrapper keeps the build
// pipeline's call site. See db.CheckCategoriaIntegrity.
func CheckCategoriaIntegrity(resultados *sql.DB, corpusPath string) ([]int64, error) {
	return db.CheckCategoriaIntegrity(resultados, corpusPath)
}
