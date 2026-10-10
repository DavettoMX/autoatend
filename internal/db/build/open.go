// Package build is the generation pipeline for corpus.db and resultados.db. It
// creates the files with their fixed build-time PRAGMAs (journal mode per
// Anexo I), applies the DDL from internal/db/schema, seeds the categories,
// VACUUMs the result and runs the cross-file integrity check (ATTACH).
// The runtime never imports this package.
package build

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"autoatend/internal/db"
	"autoatend/internal/db/schema"
)

// Journal modes of the two files, per baseline Anexo I.
const (
	journalCorpus     = "WAL"
	journalResultados = "DELETE"
)

// buildPragmas are the fixed PRAGMAs applied to every build-time connection.
// page_size and encoding only take effect on an empty database, so they must come first.
func buildPragmas(journal string) []string {
	return []string{
		"page_size(4096)",
		"encoding('UTF-8')",
		"journal_mode(" + journal + ")",
		"foreign_keys(1)",
	}
}

// create builds a brand-new database file in the given journal mode and applies ddl atomically,
// creating missing parent directories first. The pool is pinned to one
// connection so per-connection PRAGMAs always hold. If the DDL fails the
// partial file is removed; the parent directories are left in place.
func create(path string, ddl []string, journal string) (*sql.DB, error) {
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("db: refusing to overwrite existing file %s", path)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("db: create parent dir of %s: %w", path, err)
	}
	name, err := db.FileURI(path, buildPragmas(journal), nil)
	if err != nil {
		return nil, err
	}
	d, err := sql.Open("sqlite", name)
	if err != nil {
		return nil, err
	}
	d.SetMaxOpenConns(1)
	d.SetMaxIdleConns(1)
	d.SetConnMaxLifetime(0)
	if err := applyDDL(d, ddl); err != nil {
		d.Close()
		os.Remove(path)
		return nil, err
	}
	return d, nil
}

func applyDDL(d *sql.DB, ddl []string) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, s := range ddl {
		if _, err := tx.Exec(s); err != nil {
			return fmt.Errorf("db: ddl %q: %w", s, err)
		}
	}
	return tx.Commit()
}

// CreateCorpus creates a new, empty corpus.db at path.
func CreateCorpus(path string) (*sql.DB, error) { return create(path, schema.Corpus(), journalCorpus) }

// CreateResultados creates a new, empty resultados.db at path.
func CreateResultados(path string) (*sql.DB, error) {
	return create(path, schema.Resultados(), journalResultados)
}

// InsertCategorias inserts the five case labels with their fixed ids.
func InsertCategorias(d *sql.DB) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, c := range db.Categories {
		if _, err := tx.Exec(`INSERT INTO categorias(id, nombre) VALUES (?, ?)`, c.ID, c.Nombre); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Finalize VACUUMs the database and, in WAL mode, checkpoints the WAL back into
// the main file, leaving a compact, self-contained image that can be embedded.
// Only the two journal modes of Anexo I (delete, wal) are accepted.
func Finalize(d *sql.DB) error {
	ctx := context.Background()
	c, err := d.Conn(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	var mode string
	if err := c.QueryRowContext(ctx, `PRAGMA journal_mode`).Scan(&mode); err != nil {
		return err
	}
	if mode != "delete" && mode != "wal" {
		return fmt.Errorf("db: unsupported journal_mode %q (want delete or wal)", mode)
	}
	if _, err := c.ExecContext(ctx, `VACUUM`); err != nil {
		return err
	}
	if mode == "wal" {
		_, err = c.ExecContext(ctx, `PRAGMA wal_checkpoint(TRUNCATE)`)
	}
	return err
}
