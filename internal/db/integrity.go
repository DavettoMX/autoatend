package db

import (
	"context"
	"database/sql"
	"fmt"
)

// CheckCategoriaIntegrity verifies the logical cross-file foreign key
// frecuencias.categoria_id -> corpus.categorias(id) by attaching corpusPath
// read-only to the resultados connection and returning every dangling
// categoria_id, sorted ascending. NULL (global) rows are ignored.
// corpusPath must be a finished file (after Finalize): it is attached with
// immutable=1, so pending WAL content is not read. The check only reads, so the
// runtime runs it on its read-only connections at startup (TE §6.1.3) and
// internal/db/build delegates to it after generation.
func CheckCategoriaIntegrity(resultados *sql.DB, corpusPath string) ([]int64, error) {
	uri, err := FileURI(corpusPath, nil, ReadOnlyParams())
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	// ATTACH is per connection: pin one for attach, query and detach.
	c, err := resultados.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer c.Close()

	if _, err := c.ExecContext(ctx, `ATTACH DATABASE ? AS corpus`, uri); err != nil {
		return nil, fmt.Errorf("db: attach %s: %w", corpusPath, err)
	}
	defer c.ExecContext(ctx, `DETACH DATABASE corpus`)

	rows, err := c.QueryContext(ctx, `
		SELECT DISTINCT f.categoria_id
		FROM frecuencias f
		LEFT JOIN corpus.categorias c ON c.id = f.categoria_id
		WHERE f.categoria_id IS NOT NULL AND c.id IS NULL
		ORDER BY f.categoria_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var orphans []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		orphans = append(orphans, id)
	}
	return orphans, rows.Err()
}
