package db_test

import (
	"os"
	"path/filepath"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
	"autoatend/internal/testutil"
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

// TE §5.1.4: the final artifact opens corpus.db with mode=ro&immutable=1. A WAL
// database opened with mode=ro alone needs to create its -shm sidecar, which
// fails in a read-only location; immutable=1 reads the file as a fixed image.
func TestOpenReadOnlyWALInReadOnlyDir(t *testing.T) {
	src := newCorpusFile(t)
	raw, err := os.ReadFile(src)
	if err != nil {
		t.Fatal(err)
	}
	dir := filepath.Join(t.TempDir(), "ro")
	if err := os.Mkdir(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "corpus.db")
	if err := os.WriteFile(path, raw, 0o444); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(dir, 0o555); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { os.Chmod(dir, 0o755) })

	ro, err := db.OpenReadOnly(path)
	if err != nil {
		t.Fatalf("OpenReadOnly in a read-only directory: %v", err)
	}
	defer ro.Close()
	testutil.Scalar(t, ro, 5, `SELECT COUNT(*) FROM categorias`)

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("read-only open left %d files in %s, want only corpus.db", len(entries), dir)
	}
}
