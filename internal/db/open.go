package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// buildPragmas are the fixed PRAGMAs applied to every build-time connection.
// page_size and encoding only take effect on an empty database, so they must come first.
var buildPragmas = []string{
	"page_size(4096)",
	"encoding('UTF-8')",
	"journal_mode(DELETE)",
	"foreign_keys(1)",
}

func dsn(path string, pragmas []string, extra url.Values) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	q := url.Values{}
	for k, v := range extra {
		q[k] = v
	}
	for _, p := range pragmas {
		q.Add("_pragma", p)
	}
	u := url.URL{Scheme: "file", Path: abs, RawQuery: q.Encode()}
	return u.String(), nil
}

// create builds a brand-new database file and applies ddl atomically.
// The pool is pinned to one connection so per-connection PRAGMAs always hold.
func create(path string, ddl []string) (*sql.DB, error) {
	if _, err := os.Stat(path); err == nil {
		return nil, fmt.Errorf("db: refusing to overwrite existing file %s", path)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	name, err := dsn(path, buildPragmas, nil)
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
func CreateCorpus(path string) (*sql.DB, error) { return create(path, corpusDDL) }

// CreateResultados creates a new, empty resultados.db at path.
func CreateResultados(path string) (*sql.DB, error) { return create(path, resultadosDDL) }

// InsertCategorias inserts the five case labels with their fixed ids.
func InsertCategorias(d *sql.DB) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, c := range Categories {
		if _, err := tx.Exec(`INSERT INTO categorias(id, nombre) VALUES (?, ?)`, c.ID, c.Nombre); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Finalize forces journal_mode=DELETE (ADR-007) and then VACUUMs, leaving a
// compact, rollback-journal image that is safe to embed and deserialize later.
func Finalize(d *sql.DB) error {
	ctx := context.Background()
	c, err := d.Conn(ctx)
	if err != nil {
		return err
	}
	defer c.Close()
	var mode string
	if err := c.QueryRowContext(ctx, `PRAGMA journal_mode = DELETE`).Scan(&mode); err != nil {
		return err
	}
	if mode != "delete" {
		return fmt.Errorf("db: journal_mode = %q after setting DELETE", mode)
	}
	_, err = c.ExecContext(ctx, `VACUUM`)
	return err
}

// OpenReadOnly opens an existing database file for reading only
// (mode=ro plus PRAGMA query_only=ON). It fails if the file does not exist.
func OpenReadOnly(path string) (*sql.DB, error) {
	name, err := dsn(path, []string{"query_only(1)", "foreign_keys(1)"}, url.Values{"mode": {"ro"}})
	if err != nil {
		return nil, err
	}
	d, err := sql.Open("sqlite", name)
	if err != nil {
		return nil, err
	}
	if err := d.Ping(); err != nil {
		d.Close()
		return nil, err
	}
	return d, nil
}
