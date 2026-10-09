package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"path/filepath"
)

// CheckCategoriaIntegrity verifies the logical cross-file foreign key
// frecuencias.categoria_id -> corpus.categorias(id) by attaching corpusPath
// read-only to the resultados connection and returning every dangling
// categoria_id, sorted ascending. NULL (global) rows are ignored.
//
// Generation pipeline only (ADR-007): it needs both files on disk, which is not
// the case at runtime, where the databases are in-memory deserialized buffers.
func CheckCategoriaIntegrity(resultados *sql.DB, corpusPath string) ([]int64, error) {
	abs, err := filepath.Abs(corpusPath)
	if err != nil {
		return nil, err
	}
	uri := (&url.URL{Scheme: "file", Path: abs, RawQuery: url.Values{"mode": {"ro"}}.Encode()}).String()

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
