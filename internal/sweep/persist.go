package sweep

import (
	"database/sql"
	"fmt"
	"os"

	"autoatend/internal/db"
	"autoatend/internal/db/build"
	"autoatend/internal/lexer"
	"autoatend/internal/regex"
)

// corridaID is the id of the only run a resultados.db holds (forma_canonica is
// UNIQUE and the I-2 query needs one global row per lexeme, so one file is one run).
const corridaID = 1

// persist builds a new resultados.db at path in one transaction: the run, then
// every lexeme in canonical order with its global frequency row, its non-zero
// per-category rows (ascending categoria_id) and its variants ordered by form.
// Then Finalize (VACUUM). On any error the partial file is removed; if the
// file already exists nothing is created and nothing is removed.
// It returns how many lexemes had to fall back to their canonical form
// because their generated ER failed the round-trip (see validatedER).
func persist(path string, entries []*lexer.Entry, seed int64, parametros []byte) (fallbacks int, err error) {
	d, err := build.CreateResultados(path)
	if err != nil {
		return 0, err
	}
	defer func() {
		d.Close()
		if err != nil {
			os.Remove(path)
		}
	}()
	if fallbacks, err = insertAll(d, entries, seed, parametros); err != nil {
		return 0, err
	}
	return fallbacks, build.Finalize(d)
}

// validatedER is the TE Etapa 6 check: the generated ER is compiled with the
// own engine and must accept every observed variant of the lexeme. If it does
// not, the lexeme falls back to the canonical form (escaped, as the generator
// would write it) and fallback is true.
func validatedER(canonical string, forms []string, er string) (string, bool) {
	if m, err := regex.Compile(er); err == nil {
		ok := true
		for _, f := range forms {
			if !m.Match(f) {
				ok = false
				break
			}
		}
		if ok {
			return er, false
		}
	}
	return lexer.GenerateER([]string{canonical}), true
}

func insertAll(d *sql.DB, entries []*lexer.Entry, seed int64, parametros []byte) (fallbacks int, err error) {
	tx, err := d.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	if _, err := tx.Exec(`INSERT INTO corridas(id, semilla, parametros, total_lexemas) VALUES (?, ?, ?, ?)`,
		corridaID, seed, string(parametros), len(entries)); err != nil {
		return 0, fmt.Errorf("sweep: insert corrida: %w", err)
	}
	insLex, err := tx.Prepare(`INSERT INTO lexemas(forma_canonica, expresion_regular, corrida_id) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer insLex.Close()
	insFrec, err := tx.Prepare(`INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer insFrec.Close()
	insVar, err := tx.Prepare(`INSERT INTO variantes(lexema_id, forma, frecuencia) VALUES (?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer insVar.Close()

	for _, e := range entries {
		variants := e.SortedVariants()
		forms := make([]string, len(variants))
		for i, v := range variants {
			forms[i] = v.Form
		}
		canon := e.Canonical()
		er, fellBack := validatedER(canon, forms, lexer.GenerateER(forms))
		if fellBack {
			fallbacks++
		}
		res, err := insLex.Exec(canon, er, corridaID)
		if err != nil {
			return 0, fmt.Errorf("sweep: insert lexema %q: %w", canon, err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return 0, err
		}
		if _, err := insFrec.Exec(id, corridaID, nil, e.Global); err != nil {
			return 0, fmt.Errorf("sweep: insert frecuencia global of %q: %w", canon, err)
		}
		for i, n := range e.PerCat {
			if n == 0 {
				continue
			}
			if _, err := insFrec.Exec(id, corridaID, db.Categories[i].ID, n); err != nil {
				return 0, fmt.Errorf("sweep: insert frecuencia of %q, categoria %d: %w", canon, db.Categories[i].ID, err)
			}
		}
		for _, v := range variants {
			if _, err := insVar.Exec(id, v.Form, v.N); err != nil {
				return 0, fmt.Errorf("sweep: insert variante %q: %w", v.Form, err)
			}
		}
	}
	return fallbacks, tx.Commit()
}
