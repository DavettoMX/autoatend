package corpus

import (
	"math"
	"path/filepath"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/testutil"
)

// Level 4a integration: gencorpus -> corpus.db at full case size (seed 10433, 73,000 rows).
// The assertions query the database directly, independently of the generator's own statistics.
func TestLevel4aGencorpusToCorpusDB(t *testing.T) {
	if testing.Short() {
		t.Skip("full-size generation skipped in -short mode")
	}
	dir := t.TempDir()
	opts := func(name string) Options {
		return Options{Config: DefaultConfig(), OutPath: filepath.Join(dir, name+".db")}
	}
	first, err := Run(opts("a"))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	second, err := Run(opts("b"))
	if err != nil {
		t.Fatalf("second run: %v", err)
	}
	// Groundwork for I-3: two builds with the same seed are byte-identical.
	if first.Persisted.SHA256 != second.Persisted.SHA256 {
		t.Errorf("corpus.db differs between identical builds: %s vs %s", first.Persisted.SHA256, second.Persisted.SHA256)
	}

	d, err := db.OpenReadOnly(filepath.Join(dir, "a.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()

	// Exact target counts.
	wantCounts := map[int]int{1: 21900, 2: 14600, 3: 14600, 4: 14600, 5: 7300}
	rows, err := d.Query(`SELECT categoria_id, COUNT(*) FROM comentarios GROUP BY categoria_id ORDER BY categoria_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	total := 0
	got := map[int]int{}
	for rows.Next() {
		var cat, n int
		if err := rows.Scan(&cat, &n); err != nil {
			t.Fatal(err)
		}
		got[cat] = n
		total += n
	}
	if total != 73000 {
		t.Errorf("total rows = %d, want 73000", total)
	}
	for cat, want := range wantCounts {
		if got[cat] != want {
			t.Errorf("category %d: %d rows, want %d", cat, got[cat], want)
		}
	}

	// Distribution 30/20/20/20/10 within ±1%.
	wantShare := map[int]float64{1: 0.30, 2: 0.20, 3: 0.20, 4: 0.20, 5: 0.10}
	for cat, want := range wantShare {
		if share := float64(got[cat]) / float64(total); math.Abs(share-want) > 0.01 {
			t.Errorf("category %d share %.4f, want %.2f ± 0.01", cat, share, want)
		}
	}

	// Zero empty texts.
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM comentarios WHERE TRIM(texto) = ''`)
	// 100 comments per day over 730 days.
	testutil.Scalar(t, d, 730, `SELECT COUNT(DISTINCT fecha) FROM comentarios`)
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM (SELECT fecha FROM comentarios GROUP BY fecha HAVING COUNT(*) != 100)`)

	// Structural integrity.
	var ic string
	if err := d.QueryRow(`PRAGMA integrity_check`).Scan(&ic); err != nil || ic != "ok" {
		t.Errorf("integrity_check = %q (%v)", ic, err)
	}
	testutil.Scalar(t, d, 0, `SELECT COUNT(*) FROM pragma_foreign_key_check`)

	// The persisted split covers every row and is stratified and leak-free at full size.
	testutil.Scalar(t, d, 73000, `SELECT COUNT(*) FROM particiones`)
	srows, err := d.Query(`
		SELECT c.categoria_id, SUM(p.conjunto = 'train'), COUNT(*),
		       MAX(CASE WHEN p.conjunto = 'train' THEN c.fecha END),
		       MIN(CASE WHEN p.conjunto = 'test'  THEN c.fecha END)
		FROM comentarios c JOIN particiones p ON p.comentario_id = c.id
		GROUP BY c.categoria_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer srows.Close()
	for srows.Next() {
		var cat, train, n int
		var maxTrain, minTest string
		if err := srows.Scan(&cat, &train, &n, &maxTrain, &minTest); err != nil {
			t.Fatal(err)
		}
		if share := float64(train) / float64(n); math.Abs(share-0.80) > 0.01 {
			t.Errorf("category %d train share %.4f", cat, share)
		}
		if maxTrain > minTest {
			t.Errorf("category %d: temporal leakage %s > %s", cat, maxTrain, minTest)
		}
	}
}
