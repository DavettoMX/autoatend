package sweep

import (
	"fmt"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
)

// i2Query is the frequency-conservation query of TE §6.3.3, verbatim: it
// returns every lexeme whose global frequency differs from the sum over
// categories (and so also detects a duplicated global row).
const i2Query = `
	SELECT l.forma_canonica,
	       SUM(CASE WHEN f.categoria_id IS NULL THEN f.valor ELSE 0 END) AS global,
	       SUM(CASE WHEN f.categoria_id IS NOT NULL THEN f.valor ELSE 0 END) AS suma_categorias
	FROM lexemas l
	JOIN frecuencias f ON f.lexema_id = l.id
	GROUP BY l.id
	HAVING global != suma_categorias`

// verify re-opens the finished resultados.db read-only and checks it
// independently of the in-memory counters: the I-2 query returns no rows,
// there is exactly one run, the lexeme count matches, no lexeme lacks variants
// or frequencies, and no categoria_id dangles against corpus.db.
func verify(path, corpusPath string, wantLexemas int) error {
	d, err := db.OpenReadOnly(path)
	if err != nil {
		return err
	}
	defer d.Close()

	checks := []struct {
		name  string
		query string
		want  int
	}{
		{"I-2 conservation (TE §6.3.3)", `SELECT COUNT(*) FROM (` + i2Query + `)`, 0},
		{"corridas rows", `SELECT COUNT(*) FROM corridas`, 1},
		{"lexemas rows", `SELECT COUNT(*) FROM lexemas`, wantLexemas},
		{"lexemas without variants", `SELECT COUNT(*) FROM lexemas WHERE id NOT IN (SELECT lexema_id FROM variantes)`, 0},
		{"lexemas without frequencies", `SELECT COUNT(*) FROM lexemas WHERE id NOT IN (SELECT lexema_id FROM frecuencias)`, 0},
	}
	for _, c := range checks {
		var got int
		if err := d.QueryRow(c.query).Scan(&got); err != nil {
			return fmt.Errorf("sweep: verify %s: %w", c.name, err)
		}
		if got != c.want {
			return fmt.Errorf("sweep: verify %s: got %d, want %d", c.name, got, c.want)
		}
	}
	orphans, err := build.CheckCategoriaIntegrity(d, corpusPath)
	if err != nil {
		return err
	}
	if len(orphans) != 0 {
		return fmt.Errorf("sweep: verify categoria_id: dangling ids %v", orphans)
	}
	return nil
}
