package sweep

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
)

// This file implements invariant I-4 (TE §6.3.3 and §6.3.8, CPP04): the sweep time must grow
// linearly with the number of comments. The subcorpora are the first N comments of the real
// corpus by id, copied to a temporary corpus.db, so the committed data/corpus.db (whose SHA is
// registered) is never regenerated or touched.

// LinearityPoint is the timing of the sweep over the first N comments.
type LinearityPoint struct {
	N      int
	Runs   []time.Duration // wall time of each run, in execution order
	Median time.Duration
}

// Linearity is the outcome of MeasureLinearity: the timings and the least-squares line
// T(N) = Intercept + Slope*N, in seconds, fitted on the medians.
type Linearity struct {
	Points    []LinearityPoint
	Intercept float64 // seconds
	Slope     float64 // seconds per comment
	R2        float64
}

// Predict is the fitted time for n comments.
func (l Linearity) Predict(n int) time.Duration {
	return time.Duration((l.Intercept + l.Slope*float64(n)) * float64(time.Second))
}

// FitLine is the ordinary least-squares line y = intercept + slope*x and its coefficient of
// determination R2 = 1 - SSres/SStot (1 when y is constant and the line is exact).
func FitLine(xs, ys []float64) (intercept, slope, r2 float64) {
	n := float64(len(xs))
	var sx, sy float64
	for i := range xs {
		sx += xs[i]
		sy += ys[i]
	}
	mx, my := sx/n, sy/n
	var sxx, sxy, syy float64
	for i := range xs {
		dx, dy := xs[i]-mx, ys[i]-my
		sxx += dx * dx
		sxy += dx * dy
		syy += dy * dy
	}
	slope = sxy / sxx
	intercept = my - slope*mx
	var ssRes float64
	for i := range xs {
		e := ys[i] - (intercept + slope*xs[i])
		ssRes += e * e
	}
	if syy == 0 {
		return intercept, slope, 1
	}
	return intercept, slope, 1 - ssRes/syy
}

// WriteSubcorpus copies the first n comments of corpusPath (by id) into a new corpus.db inside
// dir and returns its path. Ids restart at 1 and keep the original order, so the subcorpus is a
// chronological prefix and db.SplitCutoff applies to it as to any corpus.
func WriteSubcorpus(corpusPath string, n int, dir string) (path string, err error) {
	src, err := db.OpenReadOnly(corpusPath)
	if err != nil {
		return "", err
	}
	defer src.Close()

	path = filepath.Join(dir, fmt.Sprintf("corpus_%d.db", n))
	d, err := build.CreateCorpus(path)
	if err != nil {
		return "", err
	}
	defer func() {
		d.Close()
		if err != nil {
			os.Remove(path)
		}
	}()
	if err = build.InsertCategorias(d); err != nil {
		return "", err
	}
	if err = copyPrefix(src, d, n); err != nil {
		return "", err
	}
	if err = build.Finalize(d); err != nil {
		return "", err
	}
	return path, nil
}

func copyPrefix(src, dst *sql.DB, n int) error {
	rows, err := src.Query(`SELECT fecha, texto, categoria_id FROM comentarios ORDER BY id LIMIT ?`, n)
	if err != nil {
		return err
	}
	defer rows.Close()
	tx, err := dst.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	copied := 0
	for rows.Next() {
		var fecha, texto string
		var cat int
		if err := rows.Scan(&fecha, &texto, &cat); err != nil {
			return err
		}
		if _, err := stmt.Exec(fecha, texto, cat); err != nil {
			return err
		}
		copied++
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if copied != n {
		return fmt.Errorf("sweep: the corpus has only %d comments, want a subcorpus of %d", copied, n)
	}
	return tx.Commit()
}

// MeasureLinearity times Run over the first N comments of corpusPath for every N in sizes
// (runs times each, median kept) and fits a line through the medians. workDir receives the
// temporary subcorpora and sweep outputs, which are removed as they are no longer needed.
func MeasureLinearity(corpusPath string, sizes []int, runs int, workDir string) (Linearity, error) {
	var l Linearity
	if len(sizes) < 2 {
		return l, errors.New("sweep: a line needs at least 2 sizes")
	}
	if runs < 1 {
		return l, errors.New("sweep: runs must be >= 1")
	}
	for _, n := range sizes {
		if n < 1 {
			return l, fmt.Errorf("sweep: invalid subcorpus size %d", n)
		}
		p, err := measureSize(corpusPath, n, runs, workDir)
		if err != nil {
			return l, err
		}
		l.Points = append(l.Points, p)
	}
	xs := make([]float64, len(l.Points))
	ys := make([]float64, len(l.Points))
	for i, p := range l.Points {
		xs[i], ys[i] = float64(p.N), p.Median.Seconds()
	}
	l.Intercept, l.Slope, l.R2 = FitLine(xs, ys)
	return l, nil
}

func measureSize(corpusPath string, n, runs int, workDir string) (LinearityPoint, error) {
	p := LinearityPoint{N: n}
	sub, err := WriteSubcorpus(corpusPath, n, workDir)
	if err != nil {
		return p, err
	}
	defer os.Remove(sub)
	for r := 0; r < runs; r++ {
		out := filepath.Join(workDir, fmt.Sprintf("resultados_%d_%d.db", n, r))
		res, err := Run(Options{CorpusPath: sub, OutPath: out, Seed: 10433})
		os.Remove(out)
		if err != nil {
			return p, err
		}
		p.Runs = append(p.Runs, res.Elapsed)
	}
	p.Median = median(p.Runs)
	return p, nil
}

func median(ds []time.Duration) time.Duration {
	s := slices.Clone(ds)
	slices.Sort(s)
	if len(s)%2 == 1 {
		return s[len(s)/2]
	}
	return (s[len(s)/2-1] + s[len(s)/2]) / 2
}

// Markdown renders the report docs/reports/i4_linealidad.md. equipo and fecha describe the
// machine and the day of the measurement; the caller supplies them so that the measurement
// itself stays free of environment lookups.
func (l Linearity) Markdown(equipo, fecha string) string {
	var b strings.Builder
	last := l.Points[len(l.Points)-1]
	fmt.Fprintf(&b, "# I-4: linealidad del barrido (CPP04)\n\n")
	fmt.Fprintf(&b, "Invariante I-4 del TE §6.3.3 y §6.3.8: el tiempo de barrido crece linealmente con el número de comentarios, T(k·N) ≈ k·T(N). Umbrales: pendiente positiva, R² ≥ 0.99 y T(73k) < 60 s.\n\n")
	fmt.Fprintf(&b, "- Equipo: %s\n- Fecha: %s\n", equipo, fecha)
	fmt.Fprintf(&b, "- Método: subcorpus de los primeros N comentarios de `data/corpus.db` (por `id`), copiados a una base temporal; se mide `sweep.Run` completo (lectura, tabla de símbolos, persistencia y verificación de `resultados.db`); mediana de %d corridas por tamaño.\n\n", len(l.Points[0].Runs))
	fmt.Fprintf(&b, "| N | Mediana (s) | Mínimo (s) | Máximo (s) |\n|---:|---:|---:|---:|\n")
	for _, p := range l.Points {
		lo, hi := slices.Min(p.Runs), slices.Max(p.Runs)
		fmt.Fprintf(&b, "| %d | %.3f | %.3f | %.3f |\n", p.N, p.Median.Seconds(), lo.Seconds(), hi.Seconds())
	}
	fmt.Fprintf(&b, "\nRecta de mínimos cuadrados sobre las medianas: T(N) = %.4f + %.3e·N segundos.\n\n", l.Intercept, l.Slope)
	verdict := func(ok bool) string {
		if ok {
			return "cumple"
		}
		return "NO cumple"
	}
	fmt.Fprintf(&b, "| Criterio | Umbral | Medido | Resultado |\n|---|---|---|---|\n")
	fmt.Fprintf(&b, "| Pendiente | > 0 | %.3e s por comentario | %s |\n", l.Slope, verdict(l.Slope > 0))
	fmt.Fprintf(&b, "| R² | ≥ 0.99 | %.5f | %s |\n", l.R2, verdict(l.R2 >= 0.99))
	fmt.Fprintf(&b, "| T(%d) medido | < 60 s | %.3f s (ajustado por la recta: %.3f s) | %s |\n", last.N, last.Median.Seconds(), l.Predict(last.N).Seconds(), verdict(last.Median < 60*time.Second))
	return b.String()
}
