package build

import (
	"path/filepath"
	"reflect"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/testutil"
)

func TestCheckCategoriaIntegrityCleanAndOrphans(t *testing.T) {
	corpus, corpusPath := newCorpus(t)
	seedCategorias(t, corpus)
	if err := Finalize(corpus); err != nil {
		t.Fatal(err)
	}

	res, _ := newResultados(t)
	seedResultados(t, res)
	// Global row (NULL) plus all five valid categories.
	testutil.MustExec(t, res, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 15)`)
	for c := 1; c <= 5; c++ {
		testutil.MustExec(t, res, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, ?, 3)`, c)
	}

	orphans, err := CheckCategoriaIntegrity(res, corpusPath)
	if err != nil {
		t.Fatalf("CheckCategoriaIntegrity: %v", err)
	}
	if len(orphans) != 0 {
		t.Fatalf("orphans = %v, want none (NULL global rows must be ignored)", orphans)
	}

	// Introduce dangling categoria_id values (no FK across files, so SQLite accepts them).
	testutil.MustExec(t, res, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (2, 1, 9, 1)`)
	testutil.MustExec(t, res, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (2, 1, 7, 1)`)
	orphans, err = CheckCategoriaIntegrity(res, corpusPath)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{7, 9}; !reflect.DeepEqual(orphans, want) {
		t.Fatalf("orphans = %v, want %v", orphans, want)
	}
}

func TestCheckCategoriaIntegrityLeavesNothingAttached(t *testing.T) {
	corpus, corpusPath := newCorpus(t)
	seedCategorias(t, corpus)
	res, _ := newResultados(t)

	if _, err := CheckCategoriaIntegrity(res, corpusPath); err != nil {
		t.Fatal(err)
	}
	var n int
	if err := res.QueryRow(`SELECT COUNT(*) FROM pragma_database_list`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Fatalf("databases attached after check = %d, want 1 (main only)", n)
	}
}

func TestCheckCategoriaIntegrityMissingCorpusFile(t *testing.T) {
	res, _ := newResultados(t)
	missing := filepath.Join(t.TempDir(), "absent.db")
	if _, err := CheckCategoriaIntegrity(res, missing); err == nil {
		t.Fatal("expected error for missing corpus file")
	}
	if _, err := db.OpenReadOnly(missing); err == nil {
		t.Fatal("missing file must not have been created")
	}
}
