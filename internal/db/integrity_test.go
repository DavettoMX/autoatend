package db_test

import (
	"path/filepath"
	"reflect"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
	"autoatend/internal/testutil"
)

// newResultadosFile creates a finished resultados.db with one global row and the given
// categoria_id values, and returns its path.
func newResultadosFile(t *testing.T, categorias ...int) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "resultados.db")
	d, err := build.CreateResultados(path)
	if err != nil {
		t.Fatalf("CreateResultados: %v", err)
	}
	testutil.MustExec(t, d, `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (1, '{}', 1)`)
	testutil.MustExec(t, d, `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('a', 'a', 1)`)
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 9)`)
	for _, c := range categorias {
		testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, ?, 1)`, c)
	}
	if err := build.Finalize(d); err != nil {
		t.Fatal(err)
	}
	d.Close()
	return path
}

// The runtime opens both files read-only and checks the cross-file foreign key itself,
// without importing internal/db/build (TE §6.1.3).
func TestCheckCategoriaIntegrityRuntimeReadOnly(t *testing.T) {
	corpus := newCorpusFile(t)

	clean, err := db.OpenReadOnly(newResultadosFile(t, 1, 2, 3, 4, 5))
	if err != nil {
		t.Fatal(err)
	}
	defer clean.Close()
	orphans, err := db.CheckCategoriaIntegrity(clean, corpus)
	if err != nil || len(orphans) != 0 {
		t.Fatalf("clean: orphans = %v, err = %v; want none", orphans, err)
	}

	dirty, err := db.OpenReadOnly(newResultadosFile(t, 1, 9, 7))
	if err != nil {
		t.Fatal(err)
	}
	defer dirty.Close()
	orphans, err = db.CheckCategoriaIntegrity(dirty, corpus)
	if err != nil {
		t.Fatal(err)
	}
	if want := []int64{7, 9}; !reflect.DeepEqual(orphans, want) {
		t.Fatalf("orphans = %v, want %v", orphans, want)
	}
	// The check must leave the connection pool without attached databases.
	testutil.Scalar(t, dirty, 1, `SELECT COUNT(*) FROM pragma_database_list`)
}

func TestCheckCategoriaIntegrityMissingCorpus(t *testing.T) {
	res, err := db.OpenReadOnly(newResultadosFile(t, 1))
	if err != nil {
		t.Fatal(err)
	}
	defer res.Close()
	if _, err := db.CheckCategoriaIntegrity(res, filepath.Join(t.TempDir(), "absent.db")); err == nil {
		t.Fatal("expected an error for a missing corpus file")
	}
}
