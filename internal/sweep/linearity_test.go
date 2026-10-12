package sweep_test

import (
	"math"
	"strings"
	"testing"
	"time"

	"autoatend/internal/db"
	"autoatend/internal/sweep"
)

// FitLine is ordinary least squares; the expected numbers are computed by hand.
func TestCPP04_FitLineExact(t *testing.T) {
	// y = 2 + 3x on four points: perfect fit.
	a, b, r2 := sweep.FitLine([]float64{0, 1, 2, 3}, []float64{2, 5, 8, 11})
	if math.Abs(a-2) > 1e-12 || math.Abs(b-3) > 1e-12 || math.Abs(r2-1) > 1e-12 {
		t.Errorf("exact line: intercept %v slope %v R2 %v, want 2, 3, 1", a, b, r2)
	}
	// x = 1,2,3,4 and y = 1,3,2,4: slope 0.8, intercept 0.5, R2 = 0.64.
	a, b, r2 = sweep.FitLine([]float64{1, 2, 3, 4}, []float64{1, 3, 2, 4})
	if math.Abs(a-0.5) > 1e-12 || math.Abs(b-0.8) > 1e-12 || math.Abs(r2-0.64) > 1e-12 {
		t.Errorf("noisy line: intercept %v slope %v R2 %v, want 0.5, 0.8, 0.64", a, b, r2)
	}
	// A quadratic growth must fit worse than a line: this is what the invariant detects.
	_, _, r2 = sweep.FitLine([]float64{1, 10, 30, 50, 73}, []float64{1, 100, 900, 2500, 5329})
	if r2 >= 0.99 {
		t.Errorf("quadratic data has R2 %v, want < 0.99", r2)
	}
}

func TestCPP04_MeasureLinearityFixture(t *testing.T) {
	corpus := fixtureCorpus(t)
	l, err := sweep.MeasureLinearity(corpus, []int{4, 10, 20}, 3, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if len(l.Points) != 3 {
		t.Fatalf("%d points, want 3", len(l.Points))
	}
	for i, n := range []int{4, 10, 20} {
		p := l.Points[i]
		if p.N != n || len(p.Runs) != 3 || p.Median <= 0 {
			t.Errorf("point %d = %+v, want N=%d with 3 positive runs", i, p, n)
		}
		lo, hi := p.Runs[0], p.Runs[0]
		for _, r := range p.Runs {
			lo, hi = min(lo, r), max(hi, r)
		}
		if p.Median < lo || p.Median > hi {
			t.Errorf("median %v outside [%v, %v]", p.Median, lo, hi)
		}
	}
	if l.Predict(20) <= 0 && l.Slope > 0 {
		t.Errorf("Predict(20) = %v with a positive slope", l.Predict(20))
	}
	md := l.Markdown("equipo de prueba", "2026-10-12")
	for _, s := range []string{"equipo de prueba", "2026-10-12", "R²", "| 20 |"} {
		if !strings.Contains(md, s) {
			t.Errorf("Markdown lacks %q:\n%s", s, md)
		}
	}
}

// The subcorpus is the first N comments by id, copied to a temporary corpus.db: the committed
// data/corpus.db (whose SHA is registered) is never touched.
func TestCPP04_SubcorpusIsPrefixByID(t *testing.T) {
	corpus := fixtureCorpus(t)
	work := t.TempDir()
	path, err := sweep.WriteSubcorpus(corpus, 7, work)
	if err != nil {
		t.Fatal(err)
	}
	d, err := db.OpenReadOnly(path)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	var n int
	var last string
	if err := d.QueryRow(`SELECT COUNT(*), MAX(texto) FROM comentarios`).Scan(&n, &last); err != nil {
		t.Fatal(err)
	}
	if n != 7 {
		t.Errorf("subcorpus has %d rows, want 7", n)
	}
	var maxFecha string
	if err := d.QueryRow(`SELECT MAX(fecha) FROM comentarios`).Scan(&maxFecha); err != nil {
		t.Fatal(err)
	}
	if maxFecha != fixtureRows[6].fecha {
		t.Errorf("last fecha %s, want %s", maxFecha, fixtureRows[6].fecha)
	}
}

func TestCPP04_RejectsBadSizes(t *testing.T) {
	corpus := fixtureCorpus(t)
	if _, err := sweep.MeasureLinearity(corpus, []int{10, 999}, 1, t.TempDir()); err == nil {
		t.Error("a size larger than the corpus must be an error")
	}
	if _, err := sweep.MeasureLinearity(corpus, []int{10}, 1, t.TempDir()); err == nil {
		t.Error("fewer than 2 sizes cannot be fitted")
	}
	if _, err := sweep.MeasureLinearity(corpus, []int{5, 10}, 0, t.TempDir()); err == nil {
		t.Error("runs < 1 must be an error")
	}
}

// I-4 (TE §6.3.3 and §6.3.8) over the real corpus: T(N) for N = 1k, 10k, 30k, 50k and 73k (median
// of 3 runs), least-squares line with positive slope, R2 >= 0.99 and T(73k) < 60 s.
func TestCPP04_I4_Linealidad(t *testing.T) {
	if testing.Short() {
		t.Skip("I-4 sweeps five subcorpora three times each; skipped in -short mode")
	}
	l, err := sweep.MeasureLinearity(corpusPath, []int{1000, 10000, 30000, 50000, 73000}, 3, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range l.Points {
		t.Logf("N=%6d  T=%v  runs=%v", p.N, p.Median, p.Runs)
	}
	t.Logf("T(N) = %.4f + %.3e*N  seconds; R2 = %.5f; T(73k) measured %v", l.Intercept, l.Slope, l.R2, l.Points[len(l.Points)-1].Median)
	if l.Slope <= 0 {
		t.Errorf("slope %v, want > 0", l.Slope)
	}
	if l.R2 < 0.99 {
		t.Errorf("R2 = %.5f, want >= 0.99", l.R2)
	}
	if last := l.Points[len(l.Points)-1]; last.N != 73000 || last.Median >= 60*time.Second {
		t.Errorf("T(%d) = %v, want < 60 s", last.N, last.Median)
	}
}
