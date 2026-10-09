package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func mustExec(t *testing.T, d *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := d.Exec(q, args...); err != nil {
		t.Fatalf("exec %q: %v", q, err)
	}
}

func expectErr(t *testing.T, d *sql.DB, wantSubstr, q string, args ...any) {
	t.Helper()
	_, err := d.Exec(q, args...)
	if err == nil {
		t.Fatalf("exec %q: expected error containing %q, got nil", q, wantSubstr)
	}
	if !strings.Contains(err.Error(), wantSubstr) {
		t.Fatalf("exec %q: error %q does not contain %q", q, err, wantSubstr)
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
	for i, c := range Categories {
		if c.ID != i+1 {
			t.Errorf("Categories[%d].ID = %d, want %d", i, c.ID, i+1)
		}
	}
}

func TestCorpusRejectsSixthCategory(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	expectErr(t, d, "CHECK constraint failed", `INSERT INTO categorias(nombre) VALUES ('Sugerencias')`)
}

func TestCorpusRejectsDuplicateCategoryName(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	expectErr(t, d, "UNIQUE constraint failed", `INSERT INTO categorias(nombre) VALUES ('Quejas')`)
}

func TestCorpusForeignKeyRestrict(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	mustExec(t, d, `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', 'hola', 1)`)

	// ON DELETE RESTRICT: a referenced category cannot be deleted.
	expectErr(t, d, "FOREIGN KEY constraint failed", `DELETE FROM categorias WHERE id = 1`)
	// Dangling reference rejected on insert.
	expectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', 'x', 99)`)
	// An unreferenced category can still be deleted.
	mustExec(t, d, `DELETE FROM categorias WHERE id = 5`)
}

func TestCorpusFechaGlob(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	mustExec(t, d, `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-02-29', 'ok', 1)`)
	for _, bad := range []string{"20240229", "2024/02/29", "2024-2-29", "24-02-29", "2024-02-29T10:00"} {
		expectErr(t, d, "CHECK constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, 'x', 1)`, bad)
	}
}

func TestCorpusNotNullColumns(t *testing.T) {
	d, _ := newCorpus(t)
	seedCategorias(t, d)
	expectErr(t, d, "NOT NULL constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', NULL, 1)`)
	expectErr(t, d, "NOT NULL constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (NULL, 'x', 1)`)
	expectErr(t, d, "NOT NULL constraint failed", `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES ('2024-01-01', 'x', NULL)`)
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

// ADR-007: both files use the rollback journal (journal_mode=DELETE), never WAL,
// so that VACUUMed images can later be embedded and deserialized.
func TestJournalModeDeleteOnBothFiles(t *testing.T) {
	for name, mk := range map[string]func(*testing.T) (*sql.DB, string){"corpus": newCorpus, "resultados": newResultados} {
		d, path := mk(t)
		var mode string
		if err := d.QueryRow(`PRAGMA journal_mode`).Scan(&mode); err != nil {
			t.Fatal(err)
		}
		if mode != "delete" {
			t.Errorf("%s: journal_mode = %q, want delete", name, mode)
		}
		if err := Finalize(d); err != nil {
			t.Fatalf("%s: Finalize: %v", name, err)
		}
		d.Close()

		raw, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		// SQLite header bytes 18/19 are the file-format write/read versions: 1 = legacy rollback journal, 2 = WAL.
		if raw[18] != 1 || raw[19] != 1 {
			t.Errorf("%s: header write/read version = %d/%d, want 1/1 (not WAL)", name, raw[18], raw[19])
		}
		for _, suffix := range []string{"-wal", "-shm", "-journal"} {
			if _, err := os.Stat(path + suffix); err == nil {
				t.Errorf("%s: unexpected sidecar file %s", name, path+suffix)
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
	mustExec(t, d, `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (10433, '{}', 2)`)
	mustExec(t, d, `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('factura', 'factur(a|as)', 1)`)
	mustExec(t, d, `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('cobro', 'cobro', 1)`)
}

// ADR-006: SQLite treats NULLs as distinct in a PRIMARY KEY, so the composite PK alone
// does not stop duplicate global rows; the partial unique index does.
func TestFrecuenciasGlobalRowIsUnique(t *testing.T) {
	d, _ := newResultados(t)
	seedResultados(t, d)
	mustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 10)`)
	expectErr(t, d, "UNIQUE constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 3)`)

	// Different lexeme or different corrida: another global row is fine.
	mustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (2, 1, NULL, 4)`)
	mustExec(t, d, `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (10433, '{}', 1)`)
	mustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 2, NULL, 7)`)
}

func TestFrecuenciasPerCategoryRowsKeepCompositePK(t *testing.T) {
	d, _ := newResultados(t)
	seedResultados(t, d)
	mustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 1, 6)`)
	mustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 2, 4)`)
	expectErr(t, d, "UNIQUE constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 1, 9)`)
	// A global row coexists with per-category rows of the same (lexema, corrida).
	mustExec(t, d, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, NULL, 10)`)
}

func TestResultadosCheckAndForeignKeyConstraints(t *testing.T) {
	d, _ := newResultados(t)
	expectErr(t, d, "CHECK constraint failed", `INSERT INTO corridas(semilla, parametros, total_lexemas) VALUES (10433, '{}', 0)`)
	seedResultados(t, d)
	expectErr(t, d, "CHECK constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 1, 0)`)
	expectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (99, 1, 1, 1)`)
	expectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('x', 'x', 42)`)
	expectErr(t, d, "UNIQUE constraint failed", `INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES ('factura', 'f', 1)`)

	mustExec(t, d, `INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (1, 'factura', 3)`)
	expectErr(t, d, "UNIQUE constraint failed", `INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (1, 'factura', 1)`)
	expectErr(t, d, "CHECK constraint failed", `INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (1, 'facturas', 0)`)
}

func TestOpenReadOnlyRejectsWrites(t *testing.T) {
	d, path := newCorpus(t)
	seedCategorias(t, d)
	if err := Finalize(d); err != nil {
		t.Fatal(err)
	}
	d.Close()

	ro, err := OpenReadOnly(path)
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
	if _, err := OpenReadOnly(filepath.Join(t.TempDir(), "nope.db")); err == nil {
		t.Fatal("OpenReadOnly on a missing file: expected error")
	}
}
