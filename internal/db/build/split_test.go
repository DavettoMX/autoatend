package build

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"testing"

	"autoatend/internal/db"
	"autoatend/internal/testutil"
)

type fixtureRow struct {
	fecha string
	cat   int
}

// fixtureRows returns rows whose dates are NOT monotonic in insertion order and
// that contain many same-day ties, so that "ORDER BY fecha, id" differs from "ORDER BY id".
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

func loadFixture(t *testing.T, rows []fixtureRow) *sql.DB {
	t.Helper()
	d, _ := newCorpus(t)
	seedCategorias(t, d)
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

func assignments(t *testing.T, d *sql.DB) map[int64]string {
	t.Helper()
	rows, err := d.Query(`SELECT comentario_id, conjunto FROM particiones ORDER BY comentario_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	out := map[int64]string{}
	for rows.Next() {
		var id int64
		var c string
		if err := rows.Scan(&id, &c); err != nil {
			t.Fatal(err)
		}
		out[id] = c
	}
	return out
}

// oracleSplit is an independent restatement of the rule: per category, order by
// (fecha, id); the oldest 80% (integer floor) are train, the rest test.
func oracleSplit(t *testing.T, d *sql.DB) map[int64]string {
	t.Helper()
	type rec struct {
		id    int64
		fecha string
		cat   int
	}
	rows, err := d.Query(`SELECT id, fecha, categoria_id FROM comentarios`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	byCat := map[int][]rec{}
	for rows.Next() {
		var r rec
		if err := rows.Scan(&r.id, &r.fecha, &r.cat); err != nil {
			t.Fatal(err)
		}
		byCat[r.cat] = append(byCat[r.cat], r)
	}
	out := map[int64]string{}
	for _, rs := range byCat {
		sort.Slice(rs, func(i, j int) bool {
			if rs[i].fecha != rs[j].fecha {
				return rs[i].fecha < rs[j].fecha
			}
			return rs[i].id < rs[j].id
		})
		nTrain := len(rs) * 80 / 100
		for i, r := range rs {
			if i < nTrain {
				out[r.id] = db.SplitTrain
			} else {
				out[r.id] = db.SplitTest
			}
		}
	}
	return out
}

func TestSplitFollowsTemporalStratifiedRule(t *testing.T) {
	d := loadFixture(t, fixtureRows(1))
	if err := ApplySplit(d); err != nil {
		t.Fatalf("ApplySplit: %v", err)
	}
	got, want := assignments(t, d), oracleSplit(t, d)
	if len(got) != 1000 {
		t.Fatalf("assigned %d comments, want 1000 (every comment exactly once)", len(got))
	}
	for id, w := range want {
		if got[id] != w {
			t.Fatalf("comment %d: got %q, want %q", id, got[id], w)
		}
	}
}

func TestSplitIsDeterministic(t *testing.T) {
	a := loadFixture(t, fixtureRows(7))
	b := loadFixture(t, fixtureRows(7))
	if err := ApplySplit(a); err != nil {
		t.Fatal(err)
	}
	if err := ApplySplit(b); err != nil {
		t.Fatal(err)
	}
	ga, gb := assignments(t, a), assignments(t, b)
	if len(ga) == 0 || len(ga) != len(gb) {
		t.Fatalf("assignment sizes differ: %d vs %d", len(ga), len(gb))
	}
	for id, v := range ga {
		if gb[id] != v {
			t.Fatalf("comment %d differs between identical builds: %q vs %q", id, v, gb[id])
		}
	}
}

func TestSplitIsFixedOnce(t *testing.T) {
	d := loadFixture(t, fixtureRows(3))
	if err := ApplySplit(d); err != nil {
		t.Fatal(err)
	}
	before := assignments(t, d)
	if err := ApplySplit(d); err == nil {
		t.Fatal("second ApplySplit: expected error (the split is fixed once, never re-split)")
	}
	after := assignments(t, d)
	if len(before) != len(after) {
		t.Fatal("assignment changed after refused re-split")
	}
}

func TestSplitStratificationWithinOnePercent(t *testing.T) {
	d := loadFixture(t, fixtureRows(5))
	if err := ApplySplit(d); err != nil {
		t.Fatal(err)
	}
	rows, err := d.Query(`
		SELECT c.categoria_id,
		       SUM(p.conjunto = 'train'), COUNT(*)
		FROM comentarios c JOIN particiones p ON p.comentario_id = c.id
		GROUP BY c.categoria_id ORDER BY c.categoria_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	var totalTrain, total int
	type cat struct{ train, n int }
	var cats []cat
	for rows.Next() {
		var id, train, n int
		if err := rows.Scan(&id, &train, &n); err != nil {
			t.Fatal(err)
		}
		if share := float64(train) / float64(n); math.Abs(share-0.80) > 0.01 {
			t.Errorf("category %d: train share %.4f not within 80%% ± 1%%", id, share)
		}
		cats = append(cats, cat{train, n})
		totalTrain += train
		total += n
	}
	// The class mix of the train partition must match the class mix of the corpus.
	for i, c := range cats {
		mixTrain := float64(c.train) / float64(totalTrain)
		mixAll := float64(c.n) / float64(total)
		if math.Abs(mixTrain-mixAll) > 0.01 {
			t.Errorf("category %d: train class share %.4f vs corpus %.4f differ by > 1%%", i+1, mixTrain, mixAll)
		}
	}
	if len(cats) != 5 {
		t.Fatalf("categories seen = %d, want 5", len(cats))
	}
}

func TestSplitHasNoTemporalLeakage(t *testing.T) {
	d := loadFixture(t, fixtureRows(11))
	if err := ApplySplit(d); err != nil {
		t.Fatal(err)
	}
	rows, err := d.Query(`
		SELECT c.categoria_id,
		       MAX(CASE WHEN p.conjunto = 'train' THEN c.fecha END),
		       MIN(CASE WHEN p.conjunto = 'test'  THEN c.fecha END)
		FROM comentarios c JOIN particiones p ON p.comentario_id = c.id
		GROUP BY c.categoria_id`)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()
	seen := 0
	for rows.Next() {
		var cat int
		var maxTrain, minTest string
		if err := rows.Scan(&cat, &maxTrain, &minTest); err != nil {
			t.Fatal(err)
		}
		if maxTrain > minTest {
			t.Errorf("category %d: max(train.fecha)=%s > min(test.fecha)=%s", cat, maxTrain, minTest)
		}
		seen++
	}
	if seen != 5 {
		t.Fatalf("categories checked = %d, want 5", seen)
	}
}

func TestParticionesConstraints(t *testing.T) {
	d := loadFixture(t, fixtureRows(1)[:10])
	testutil.ExpectErr(t, d, "CHECK constraint failed", `INSERT INTO particiones(comentario_id, conjunto) VALUES (1, 'val')`)
	testutil.ExpectErr(t, d, "FOREIGN KEY constraint failed", `INSERT INTO particiones(comentario_id, conjunto) VALUES (999, 'train')`)
	testutil.MustExec(t, d, `INSERT INTO particiones(comentario_id, conjunto) VALUES (1, 'train')`)
	testutil.ExpectErr(t, d, "UNIQUE constraint failed", `INSERT INTO particiones(comentario_id, conjunto) VALUES (1, 'test')`)
}
