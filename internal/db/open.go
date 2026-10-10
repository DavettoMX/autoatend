package db

import (
	"database/sql"
	"net/url"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

// FileURI builds a modernc.org/sqlite connection string (a SQLite file: URI) for
// the file at path, with the given PRAGMAs and any extra query parameters.
func FileURI(path string, pragmas []string, extra url.Values) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	return fileURI(abs, pragmas, extra), nil
}

// fileURI renders the URI of an absolute path. Windows paths (C:\...) become
// file:///C:/...: without the leading slash SQLite would parse "C:..." as the URI
// authority and reject it.
func fileURI(abs string, pragmas []string, extra url.Values) string {
	abs = filepath.ToSlash(abs)
	if !strings.HasPrefix(abs, "/") {
		abs = "/" + abs
	}
	q := url.Values{}
	for k, v := range extra {
		q[k] = v
	}
	for _, p := range pragmas {
		q.Add("_pragma", p)
	}
	u := url.URL{Scheme: "file", Path: abs, RawQuery: q.Encode()}
	return u.String()
}

// OpenReadOnly opens an existing database file for reading only
// (mode=ro plus PRAGMA query_only=ON). It fails if the file does not exist.
func OpenReadOnly(path string) (*sql.DB, error) {
	name, err := FileURI(path, []string{"query_only(1)", "foreign_keys(1)"}, url.Values{"mode": {"ro"}})
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
