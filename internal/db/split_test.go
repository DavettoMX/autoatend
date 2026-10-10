package db_test

import (
	"database/sql"
	"fmt"
	"math/rand"
	"path/filepath"
	"sort"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
)

type fixtureRow struct {
	fecha string
	cat   int
}

// fixtureRows returns rows whose dates are NOT monotonic in insertion order and
// that contain many same-day ties, so that ordering by fecha differs from ordering by id.
func fixtureRows(seed int64) []fixtureRow {
	counts := []int{300, 200, 200, 200, 100}
	rng := rand.New(rand.NewSource(seed))
	var rows []fixtureRow
	for i, n := range counts {
		for j := 0; j < n; j++ {
			rows = append(rows, fixtureRow{
				fecha: fmt.Sprintf("2024-%02d-%02d", 1+rng.Intn(12), 1+rng.Intn(28)),
				cat:   i + 1,
			})
		}
	}
	rng.Shuffle(len(rows), func(i, j int) { rows[i], rows[j] = rows[j], rows[i] })
	return rows
}

// loadFixture builds a corpus.db with the five categories and rows as comentarios.
func loadFixture(t *testing.T, rows []fixtureRow) *sql.DB {
	t.Helper()
	d, err := build.CreateCorpus(filepath.Join(t.TempDir(), "corpus.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { d.Close() })
	if err := build.InsertCategorias(d); err != nil {
		t.Fatal(err)
	}
	tx, err := d.Begin()
	if err != nil {
		t.Fatal(err)
	}
	for i, r := range rows {
		if _, err := tx.Exec(`INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, ?, ?)`, r.fecha, fmt.Sprintf("texto %d", i), r.cat); err != nil {
			t.Fatal(err)
		}
	}
	if err := tx.Commit(); err != nil {
		t.Fatal(err)
	}
	return d
}

// oracleCutoff restates the rule independently: sort every fecha ascending; the
// cut-off is the fecha at 0-based position floor(80% of N).
func oracleCutoff(rows []fixtureRow) string {
	dates := make([]string, len(rows))
	for i, r := range rows {
		dates[i] = r.fecha
	}
	sort.Strings(dates)
	return dates[len(dates)*80/100]
}

func TestSplitCutoffFollowsGlobalTemporalRule(t *testing.T) {
	rows := fixtureRows(1)
	d := loadFixture(t, rows)
	got, err := db.SplitCutoff(d)
	if err != nil {
		t.Fatalf("SplitCutoff: %v", err)
	}
	if want := oracleCutoff(rows); got != want {
		t.Fatalf("SplitCutoff = %s, want %s", got, want)
	}
}

// TE §5.2.2: every training comment is earlier than every test comment, across
// all categories, and the training set never exceeds 80%.
func TestSplitTrainStrictlyPrecedesTest(t *testing.T) {
	for _, seed := range []int64{1, 5, 11} {
		d := loadFixture(t, fixtureRows(seed))
		cut, err := db.SplitCutoff(d)
		if err != nil {
			t.Fatal(err)
		}
		var maxTrain, minTest string
		var train, total int
		if err := d.QueryRow(`
			SELECT MAX(CASE WHEN fecha <  ? THEN fecha END),
			       MIN(CASE WHEN fecha >= ? THEN fecha END),
			       SUM(fecha < ?), COUNT(*)
			FROM comentarios`, cut, cut, cut).Scan(&maxTrain, &minTest, &train, &total); err != nil {
			t.Fatal(err)
		}
		if maxTrain >= minTest {
			t.Errorf("seed %d: max(train.fecha)=%s is not before min(test.fecha)=%s", seed, maxTrain, minTest)
		}
		if train > total*80/100 || train == 0 {
			t.Errorf("seed %d: train = %d of %d, want in (0, 80%%]", seed, train, total)
		}
	}
}

func TestSplitCutoffIsDeterministic(t *testing.T) {
	a, err := db.SplitCutoff(loadFixture(t, fixtureRows(7)))
	if err != nil {
		t.Fatal(err)
	}
	b, err := db.SplitCutoff(loadFixture(t, fixtureRows(7)))
	if err != nil {
		t.Fatal(err)
	}
	if a != b {
		t.Fatalf("cut-off differs between identical corpora: %s vs %s", a, b)
	}
}

func TestSplitCutoffRejectsEmptyCorpus(t *testing.T) {
	if _, err := db.SplitCutoff(loadFixture(t, nil)); err == nil {
		t.Fatal("SplitCutoff on an empty corpus: expected error")
	}
}
