package db_test

import (
	"os"
	"path/filepath"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
)

// newCorpusFile creates a finished corpus.db fixture with the five categories.
func newCorpusFile(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "corpus.db")
	d, err := build.CreateCorpus(path)
	if err != nil {
		t.Fatalf("CreateCorpus: %v", err)
	}
	if err := build.InsertCategorias(d); err != nil {
		t.Fatalf("InsertCategorias: %v", err)
	}
	if err := build.Finalize(d); err != nil {
		t.Fatal(err)
	}
	d.Close()
	return path
}

func TestOpenReadOnlyRejectsWrites(t *testing.T) {
	path := newCorpusFile(t)

	ro, err := db.OpenReadOnly(path)
	if err != nil {
		t.Fatalf("OpenReadOnly: %v", err)
	}
	defer ro.Close()

	var n int
	if err := ro.QueryRow(`SELECT COUNT(*) FROM categorias`).Scan(&n); err != nil || n != 5 {
		t.Fatalf("read via read-only opener: n=%d err=%v", n, err)
	}
	if _, err := ro.Exec(`DELETE FROM categorias`); err == nil {
		t.Fatal("write through read-only opener: expected error")
	}
	var qo int
	if err := ro.QueryRow(`PRAGMA query_only`).Scan(&qo); err != nil || qo != 1 {
		t.Fatalf("PRAGMA query_only = %d (err %v), want 1", qo, err)
	}
}

func TestOpenReadOnlyMissingFile(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.db")
	if _, err := db.OpenReadOnly(missing); err == nil {
		t.Fatal("OpenReadOnly on a missing file: expected error")
	}
	if _, err := os.Stat(missing); err == nil {
		t.Fatal("OpenReadOnly must not create the missing file")
	}
}
