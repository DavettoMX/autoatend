package build

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/testutil"
)

func newCorpus(t *testing.T) (*sql.DB, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "corpus.db")
	d, err := CreateCorpus(path)
	if err != nil {
		t.Fatalf("CreateCorpus: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d, path
}

func newResultados(t *testing.T) (*sql.DB, string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "resultados.db")
	d, err := CreateResultados(path)
	if err != nil {
		t.Fatalf("CreateResultados: %v", err)
	}
	t.Cleanup(func() { d.Close() })
	return d, path
}

func seedCategorias(t *testing.T, d *sql.DB) {
	t.Helper()
	if err := InsertCategorias(d); err != nil {
		t.Fatalf("InsertCategorias: %v", err)
	}
}

func TestCorpusAcceptsExactlyTheFiveCategories(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)

	var n int
	if err := d.QueryRow(`SELECT COUNT(*) FROM categorias`).Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 5 {
		t.Fatalf("categorias rows = %d, want 5", n)
	}
	for i, c := range db.Categories {
		if c.ID != i+1 {
			t.Errorf("Categories[%d].ID = %d, want %d", i, c.ID, i+1)
		}
	}
}

func TestCorpusRejectsSixthCategory(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	testutil.ExpectErr(t, d, "CHECK constraint failed", `INSERT INTO categorias(nombre) VALUES ('Sugerencias')`)
}

func TestCorpusRejectsDuplicateCategoryName(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	testutil.ExpectErr(t, d, "UNIQUE constraint failed", `INSERT INTO categorias(nombre) VALUES ('Quejas')`)
}

func TestCorpusForeignKeyRestrict(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	testutil.MustExec(t, d, `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', 'hola', 1)`)

	// ON DELETE RESTRICT: a referenced category cannot be deleted.
	testutil.ExpectErr(t, d, "FOREIGN KEY constraint failed", `DELETE FROM categorias WHERE id = 1`)
	// Dangling reference rejected on insert.
	testutil.ExpectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', 'x', 99)`)
	// An unreferenced category can still be deleted.
	testutil.MustExec(t, d, `DELETE FROM categorias WHERE id = 5`)
}

func TestCorpusFechaGlob(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	testutil.MustExec(t, d, `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-02-29', 'ok', 1)`)
	for _, bad := range []string{"20240229", "2024/02/29", "2024-2-29", "24-02-29", "2024-02-29T10:00"} {
		testutil.ExpectErr(t, d, "CHECK constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, 'x', 1)`, bad)
	}
}

func TestCorpusNotNullColumns(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	testutil.ExpectErr(t, d, "NOT NULL constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', NULL, 1)`)
	testutil.ExpectErr(t, d, "NOT NULL constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (NULL, 'x', 1)`)
	testutil.ExpectErr(t, d, "NOT NULL constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', 'x', NULL)`)
}

func objectNames(t *testing.T, d *sql.DB, typ string) map[string]bool {
	t.Helper()
	rows, err := d.Query(`SELECT name FROM sqlite_master WHERE type = ? AND name NOT LIKE 'sqlite_%'`, typ)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	out := map[string]bool{}
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		out[n] = true
	}
	return out
}

func TestCorpusIndexes(t *testing.T) {
	d, _ := newCorpus(t)
	got := objectNames(t, d, "index")
	for _, want := range []string{"idx_comentarios_categoria", "idx_comentarios_fecha", "idx_comentarios_cat_fecha"} {
		if !got[want] {
			t.Errorf("missing index %s (have %v)", want, got)
		}
	}
}

func TestForeignKeysEnabled(t *testing.T) {
	for name, mk := range map[string]func(*testing.T) (*sql.DB, string){"corpus": newCorpus, "resultados": newResultados} {
		d, _ := mk(t)
		var on int
		if err := d.QueryRow(`PRAGMA foreign_keys`).Scan(&on); err != nil {
			t.Fatal(err)
		}
		if on != 1 {
			t.Errorf("%s: PRAGMA foreign_keys = %d, want 1", name, on)
		}
	}
}

// Anexo I: corpus.db uses WAL and resultados.db the rollback journal (DELETE).
// After Finalize and Close each file must be a single self-contained image
// (no -wal/-shm/-journal sidecar), so it can be embedded as-is.
func TestJournalModeMatchesAnexoI(t *testing.T) {
	cases := []struct {
		name    string
		mk      func(*testing.T) (*sql.DB, string)
		mode    string
		version byte // SQLite header bytes 18/19: 1 = rollback journal, 2 = WAL
	}{
		{"corpus", newCorpus, "wal", 2},
		{"resultados", newResultados, "delete", 1},
	}
	for _, c := range cases {
		d, path := c.mk(t)
		var mode string
		if err := d.QueryRow(`PRAGMA journal_mode`).Scan(&mode); err != nil {
			t.Fatal(err)
		}
		if mode != c.mode {
			t.Errorf("%s: journal_mode = %q, want %s", c.name, mode, c.mode)
		}
		if err := Finalize(d); err != nil {
			t.Fatalf("%s: Finalize: %v", c.name, err)
		}
		d.Close()

		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		if raw[18] != c.version || raw[19] != c.version {
			t.Errorf("%s: header write/read version = %d/%d, want %d/%d", c.name, raw[18], raw[19], c.version, c.version)
		}
		for _, suffix := range []string{"-wal", "-shm", "-journal"} {
			if _, err := os.Stat(path + suffix); err == nil {
				t.Errorf("%s: unexpected sidecar file %s", c.name, path+suffix)
			}
		}
	}
}

func TestCreateRefusesExistingFile(t *testing.T) {
	_, path := newCorpus(t)
	if _, err := CreateCorpus(path); err == nil {
		t.Fatal("CreateCorpus over an existing file: expected error")
	}
}

func TestResultadosTablesAndIndexes(t *testing.T) {
	d, _ := newResultados(t)
	tables := objectNames(t, d, "table")
	for _, want := range []string{"corridas", "lexemas", "frecuencias", "variantes"} {
		if !tables[want] {
			t.Errorf("missing table %s", want)
		}
	}
	indexes := objectNames(t, d, "index")
	for _, want := range []string{"idx_frecuencias_lexema", "idx_frecuencias_corrida", "idx_variantes_lexema", "uq_frecuencias_global"} {
		if !indexes[want] {
			t.Errorf("missing index %s", want)
		}
	}
}

// Determinism (I-3): corridas must carry no wall-clock column.
func TestCorridasHasNoTimestampColumn(t *testing.T) {
	d, _ := newResultados(t)
	rows, err := d.Query(`SELECT name FROM pragma_table_info('corridas') ORDER BY cid`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var cols []string
	for rows.Next() {
		var n string
		if err := rows.Scan(&n); err != nil {
			t.Fatal(err)
		}
		cols = append(cols, n)
	}
	want := "id,semilla,parametros,total_lexemas"
	if strings.Join(cols, ",") != want {
		t.Errorf("corridas columns = %v, want %s", cols, want)
	}
}

func seedResultados(t *testing.T, d *sql.DB) {
	t.Helper()
	testutil.MustExec(t, d, `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (10433, '{}', 2)`)
	testutil.MustExec(t, d, `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('factura', 'factur(a|as)', 1)`)
	testutil.MustExec(t, d, `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('cobro', 'cobro', 1)`)
}

// ADR-006: SQLite treats NULLs as distinct in a PRIMARY KEY, so the composite PK alone
// does not stop duplicate global rows; the partial unique index does.
func TestFrecuenciasGlobalRowIsUnique(t *testing.T) {
	d, _ := newResultados(t)
	seedResultados(t, d)
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 10)`)
	testutil.ExpectErr(t, d, "UNIQUE constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 3)`)

	// Different lexeme or different corrida: another global row is fine.
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (2, 1, NULL, 4)`)
	testutil.MustExec(t, d, `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (10433, '{}', 1)`)
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 2, NULL, 7)`)
}

func TestFrecuenciasPerCategoryRowsKeepCompositePK(t *testing.T) {
	d, _ := newResultados(t)
	seedResultados(t, d)
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 1, 6)`)
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 2, 4)`)
	testutil.ExpectErr(t, d, "UNIQUE constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 1, 9)`)
	// A global row coexists with per-category rows of the same (lexema, corrida).
	testutil.MustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 10)`)
}

func TestResultadosCheckAndForeignKeyConstraints(t *testing.T) {
	d, _ := newResultados(t)
	testutil.ExpectErr(t, d, "CHECK constraint failed", `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (10433, '{}', 0)`)
	seedResultados(t, d)
	testutil.ExpectErr(t, d, "CHECK constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 1, 0)`)
	testutil.ExpectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (99, 1, 1, 1)`)
	testutil.ExpectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('x', 'x', 42)`)
	testutil.ExpectErr(t, d, "UNIQUE constraint failed", `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('factura', 'f', 1)`)

	testutil.MustExec(t, d, `INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (1, 'factura', 3)`)
	testutil.ExpectErr(t, d, "UNIQUE constraint failed", `INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (1, 'factura', 1)`)
	testutil.ExpectErr(t, d, "CHECK constraint failed", `INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (1, 'facturas', 0)`)
}

// closedCorpus returns a corpus.db handle that has already been closed, to
// exercise the error path of every operation that needs a live connection.
func closedCorpus(t *testing.T) *sql.DB {
	t.Helper()
	d, _ := newCorpus(t)
	if err := d.Close(); err != nil {
		t.Fatal(err)
	}
	return d
}

func TestCreateRemovesFileWhenDDLFails(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.db")
	_, err := create(path, []string{"CREATE TABLE a(x)", "NOT SQL"}, journalCorpus)
	if err == nil || !strings.Contains(err.Error(), "ddl") {
		t.Fatalf("create with invalid DDL: err = %v, want a ddl error", err)
	}
	for _, p := range []string{path, path + "-journal", path + "-wal", path + "-shm"} {
		if _, err := os.Stat(p); err == nil {
			t.Errorf("failed create left %s behind", p)
		}
	}
}

func TestCreateRejectsPathUnderAFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "not-a-dir")
	if err := os.WriteFile(file, []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(file, "corpus.db")
	if _, err := CreateCorpus(path); err == nil {
		t.Fatal("CreateCorpus under a regular file: expected error")
	}
	if _, err := os.Stat(path); err == nil {
		t.Fatalf("%s must not exist", path)
	}
}

func TestInsertCategoriasIsAtomicOnConflict(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	if err := InsertCategorias(d); err == nil || !strings.Contains(err.Error(), "UNIQUE constraint failed") {
		t.Fatalf("second InsertCategorias: err = %v, want UNIQUE violation", err)
	}
	testutil.Scalar(t, d, 5, `SELECT COUNT(*) FROM categorias`)
}

func TestFinalizeRejectsNonDeleteJournal(t *testing.T) {
	// An in-memory database cannot use the rollback journal: journal_mode stays "memory".
	d, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if err := Finalize(d); err == nil || !strings.Contains(err.Error(), "journal_mode") {
		t.Fatalf("Finalize on :memory:: err = %v, want journal_mode error", err)
	}
}

func TestOperationsFailOnClosedDB(t *testing.T) {
	_, corpusPath := newCorpus(t)
	ops := map[string]func(*sql.DB) error{
		"applyDDL":         func(d *sql.DB) error { return applyDDL(d, []string{"CREATE TABLE a(x)"}) },
		"InsertCategorias": InsertCategorias,
		"Finalize":         Finalize,
		"ApplySplit":       ApplySplit,
		"CheckCategoriaIntegrity": func(d *sql.DB) error {
			_, err := CheckCategoriaIntegrity(d, corpusPath)
			return err
		},
	}
	for name, op := range ops {
		t.Run(name, func(t *testing.T) {
			if err := op(closedCorpus(t)); err == nil {
				t.Fatalf("%s on a closed database: expected error", name)
			}
		})
	}
}

func TestCreateMakesParentDirectories(t *testing.T) {
	path := filepath.Join(t.TempDir(), "a", "b", "corpus.db")
	d, err := CreateCorpus(path)
	if err != nil {
		t.Fatalf("CreateCorpus in a missing directory: %v", err)
	}
	defer d.Close()
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("corpus.db not created: %v", err)
	}
	testutil.Scalar(t, d, 1, `SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'categorias'`)
}
