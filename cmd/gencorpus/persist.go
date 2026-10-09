package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"

	"autoatend/internal/db"
)

// Persist builds a new corpus.db at path: the five categories, then every
// comment in chronological order (so ids follow time) in one bulk-insert
// transaction, then the 80/20 temporal-stratified split, then Finalize
// (journal_mode=DELETE + VACUUM). On any error the partial file is removed.
func Persist(path string, cs []Comment) (err error) {
	d, err := db.CreateCorpus(path)
	if err != nil {
		return err
	}
	defer func() {
		d.Close()
		if err != nil {
			os.Remove(path)
		}
	}()

	if err = db.InsertCategorias(d); err != nil {
		return err
	}
	if err = bulkInsert(d, cs); err != nil {
		return err
	}
	if err = db.ApplySplit(d); err != nil {
		return err
	}
	return db.Finalize(d)
}

func bulkInsert(d *sql.DB, cs []Comment) error {
	tx, err := d.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare(`INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()
	for i, c := range cs {
		if _, err := stmt.Exec(c.Fecha, c.Texto, c.CategoriaID); err != nil {
			return fmt.Errorf("gencorpus: insert row %d: %w", i, err)
		}
	}
	return tx.Commit()
}

// PersistedStats are measured by re-opening the finished file read-only,
// independently of the in-memory generation statistics.
type PersistedStats struct {
	Rows         int
	PerCategory  [5]int
	Train, Test  [5]int
	MaxTrainDate [5]string
	MinTestDate  [5]string
	Empty        int
	FirstDate    string
	LastDate     string
	SHA256       string
}

func verifyCorpus(path string) (PersistedStats, error) {
	var p PersistedStats
	d, err := db.OpenReadOnly(path)
	if err != nil {
		return p, err
	}
	defer d.Close()

	rows, err := d.Query(`
		SELECT c.categoria_id, COUNT(*),
		       SUM(p.conjunto = 'train'), SUM(p.conjunto = 'test'),
		       COALESCE(MAX(CASE WHEN p.conjunto = 'train' THEN c.fecha END), ''),
		       COALESCE(MIN(CASE WHEN p.conjunto = 'test'  THEN c.fecha END), '')
		FROM comentarios c JOIN particiones p ON p.comentario_id = c.id
		GROUP BY c.categoria_id ORDER BY c.categoria_id`)
	if err != nil {
		return p, err
	}
	defer rows.Close()
	for rows.Next() {
		var cat, n, train, test int
		var maxTrain, minTest string
		if err := rows.Scan(&cat, &n, &train, &test, &maxTrain, &minTest); err != nil {
			return p, err
		}
		if cat < 1 || cat > 5 {
			return p, fmt.Errorf("gencorpus: unexpected categoria_id %d", cat)
		}
		i := cat - 1
		p.PerCategory[i], p.Train[i], p.Test[i] = n, train, test
		p.MaxTrainDate[i], p.MinTestDate[i] = maxTrain, minTest
		p.Rows += n
	}
	if err := rows.Err(); err != nil {
		return p, err
	}
	if err := d.QueryRow(`SELECT COUNT(*) FROM comentarios WHERE TRIM(texto) = ''`).Scan(&p.Empty); err != nil {
		return p, err
	}
	if err := d.QueryRow(`SELECT MIN(fecha), MAX(fecha) FROM comentarios`).Scan(&p.FirstDate, &p.LastDate); err != nil {
		return p, err
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return p, err
	}
	sum := sha256.Sum256(raw)
	p.SHA256 = hex.EncodeToString(sum[:])
	return p, nil
}
