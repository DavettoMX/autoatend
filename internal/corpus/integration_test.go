package corpus

import (
	"math"
	"path/filepath"
	"testing"
	"time"

	"autoatend/internal/db"
	"autoatend/internal/testutil"
)

// Level 4a integration: gencorpus -> corpus.db at full case size (seed 10433, 73,000 rows).
// The assertions query the database directly, independently of the generator's own statistics.
func TestLevel4aGencorpusToCorpusDB(t *testing.T) {
	pendingTemplates(t)
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

	// Global temporal split (TE §5.2.2). Independent oracle: 100 comments per day,
	// so the first 80% are exactly the first Total*80/100/PerDay days.
	cfg := DefaultConfig()
	start, err := time.Parse("2006-01-02", cfg.Start)
	if err != nil {
		t.Fatal(err)
	}
	cut := start.AddDate(0, 0, cfg.Total*80/100/cfg.PerDay).Format("2006-01-02")
	if got, err := db.SplitCutoff(d); err != nil || got != cut {
		t.Fatalf("SplitCutoff = %q (%v), want %s", got, err, cut)
	}
	var maxTrain, minTest string
	if err := d.QueryRow(`SELECT MAX(fecha) FROM comentarios WHERE fecha < ?`, cut).Scan(&maxTrain); err != nil {
		t.Fatal(err)
	}
	if err := d.QueryRow(`SELECT MIN(fecha) FROM comentarios WHERE fecha >= ?`, cut).Scan(&minTest); err != nil {
		t.Fatal(err)
	}
	if maxTrain >= minTest {
		t.Errorf("temporal leakage: max(train.fecha)=%s >= min(test.fecha)=%s", maxTrain, minTest)
	}
	testutil.Scalar(t, d, 14600, `SELECT COUNT(*) FROM comentarios WHERE fecha >= '`+cut+`'`)

	// Stratification: per category ~80% train, and the class mix of both sets
	// matches the corpus, within ±1%.
	srows, err := d.Query(`
		SELECT categoria_id, SUM(fecha < ?), SUM(fecha >= ?), COUNT(*)
		FROM comentarios GROUP BY categoria_id`, cut, cut)
	if err != nil {
		t.Fatal(err)
	}
	defer srows.Close()
	for srows.Next() {
		var cat, train, test, n int
		if err := srows.Scan(&cat, &train, &test, &n); err != nil {
			t.Fatal(err)
		}
		if share := float64(train) / float64(n); math.Abs(share-0.80) > 0.01 {
			t.Errorf("category %d train share %.4f, want 0.80 ± 0.01", cat, share)
		}
		all := float64(n) / 73000
		if mix := float64(train) / 58400; math.Abs(mix-all) > 0.01 {
			t.Errorf("category %d train mix %.4f vs corpus %.4f", cat, mix, all)
		}
		if mix := float64(test) / 14600; math.Abs(mix-all) > 0.01 {
			t.Errorf("category %d test mix %.4f vs corpus %.4f", cat, mix, all)
		}
	}
}
