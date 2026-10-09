package build

import (
	"database/sql"
	"errors"

	"autoatend/internal/db"
)

// ApplySplit materializes the deterministic 80/20 temporal-stratified split in
// the particiones table: per category, comentarios are ordered by (fecha, id)
// and the oldest floor(80%) go to train, the newest remainder to test.
// The split is fixed once: it errors if an assignment already exists.
func ApplySplit(d *sql.DB) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var existing int
	if err := tx.QueryRow(`SELECT COUNT(*) FROM particiones`).Scan(&existing); err != nil {
		return err
	}
	if existing > 0 {
		return errors.New("db: split already applied; it is fixed once and never re-split")
	}
	if _, err := tx.Exec(`
		INSERT INTO particiones(comentario_id, conjunto)
		SELECT id,
		       CASE WHEN rn <= cnt * ? / 100 THEN ? ELSE ? END
		FROM (
		    SELECT id,
		           ROW_NUMBER() OVER (PARTITION BY categoria_id ORDER BY fecha, id) AS rn,
		           COUNT(*)     OVER (PARTITION BY categoria_id)                    AS cnt
		    FROM comentarios
		)
		ORDER BY id`, db.TrainPercent, db.SplitTrain, db.SplitTest); err != nil {
		return err
	}
	return tx.Commit()
}
