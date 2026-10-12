package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"autoatend/internal/regex"
)

// ExpectedComentarios is the size of corpus.db fixed by the case (Anexo I, 73,000 comments).
const ExpectedComentarios = 73000

// maxReported bounds the detail lines of a failing check.
const maxReported = 5

// check is one selftest verification: it returns a one-line summary on success
// or an error whose message is the diagnostic.
type check struct {
	name string
	run  func() (string, error)
}

// Selftest runs the headless verification battery of TE §6.3.9 / §6.4.3 on the
// started app: counts and integrity of corpus.db, validity of every regular
// expression of the catalogue, the five gold comments and the port. Every check
// prints one line: "ok" lines to stdout and "FALLO" lines to stderr. It reports
// whether all passed; all checks run even after a failure.
func (a *App) Selftest(stdout, stderr io.Writer) bool {
	checks := []check{
		{"corpus.db", a.checkCorpus},
		{"integridad", a.checkIntegridad},
		{"catalogo", a.checkCatalogo},
		{"gold5", a.checkGold5},
		{"puerto", a.checkPuerto},
	}
	ok := true
	for _, c := range checks {
		summary, err := c.run()
		if err != nil {
			ok = false
			fmt.Fprintf(stderr, "FALLO  %s: %v\n", c.name, err)
			continue
		}
		fmt.Fprintf(stdout, "ok     %s: %s\n", c.name, summary)
	}
	if ok {
		fmt.Fprintf(stdout, "selftest: todas las verificaciones pasaron (version %s)\n", a.opts.Version)
	} else {
		fmt.Fprintln(stderr, "selftest: FALLO")
	}
	return ok
}

// checkCorpus: record count and categories of corpus.db, and no comment with a dangling categoria_id.
func (a *App) checkCorpus() (string, error) {
	var n, cats, dangling int
	if err := a.corpus.QueryRow(`SELECT COUNT(*) FROM comentarios`).Scan(&n); err != nil {
		return "", fmt.Errorf("contando comentarios: %w", err)
	}
	if err := a.corpus.QueryRow(`SELECT COUNT(*) FROM categorias`).Scan(&cats); err != nil {
		return "", fmt.Errorf("contando categorias: %w", err)
	}
	if n != ExpectedComentarios {
		return "", fmt.Errorf("%d comentarios, se esperaban %d", n, ExpectedComentarios)
	}
	if cats != len(a.categorias) {
		return "", fmt.Errorf("%d categorias, se esperaban %d", cats, len(a.categorias))
	}
	if err := a.corpus.QueryRow(`SELECT COUNT(*) FROM comentarios c LEFT JOIN categorias k ON k.id = c.categoria_id WHERE k.id IS NULL`).Scan(&dangling); err != nil {
		return "", err
	}
	if dangling != 0 {
		return "", fmt.Errorf("%d comentarios con categoria_id inexistente", dangling)
	}
	return fmt.Sprintf("%d comentarios en %d categorias", n, cats), nil
}

// checkIntegridad: SQLite structural check of both files plus the cross-file categoria_id check done at startup.
func (a *App) checkIntegridad() (string, error) {
	for _, f := range []struct {
		name string
		d    *sql.DB
	}{{"corpus.db", a.corpus}, {"resultados.db", a.resultados}} {
		var res string
		if err := f.d.QueryRow(`PRAGMA quick_check`).Scan(&res); err != nil {
			return "", fmt.Errorf("%s: %w", f.name, err)
		}
		if res != "ok" {
			return "", fmt.Errorf("%s: PRAGMA quick_check = %q", f.name, res)
		}
	}
	if len(a.orphans) != 0 {
		return "", fmt.Errorf("categoria_id huerfanos: %v", a.orphans)
	}
	return "quick_check ok en corpus.db y resultados.db; categoria_id de frecuencias existe en categorias (ATTACH, 0 huerfanos)", nil
}

// checkCatalogo: every regular expression of the catalogue compiles with the project's engine and
// accepts every variant it summarizes (I-1); every vocabulary lexeme of the model is in the catalogue.
func (a *App) checkCatalogo() (string, error) {
	rows, err := a.resultados.Query(`
		SELECT l.id, l.forma_canonica, l.expresion_regular, v.forma
		FROM lexemas l LEFT JOIN variantes v ON v.lexema_id = l.id
		ORDER BY l.id, v.id`)
	if err != nil {
		return "", fmt.Errorf("leyendo el catalogo: %w", err)
	}
	defer rows.Close()

	var (
		curID                      int64 = -1
		nfa                        *regex.NFA
		canonical                  = map[string]bool{}
		lexemes, variants, problem int
		details                    []string
	)
	report := func(format string, args ...any) {
		problem++
		if len(details) < maxReported {
			details = append(details, fmt.Sprintf(format, args...))
		}
	}
	for rows.Next() {
		var id int64
		var forma, er string
		var variante *string
		if err := rows.Scan(&id, &forma, &er, &variante); err != nil {
			return "", err
		}
		if id != curID {
			curID = id
			lexemes++
			canonical[forma] = true
			var cerr error
			if nfa, cerr = regex.Compile(er); cerr != nil {
				nfa = nil
				report("la ER de %q no compila: %v", forma, cerr)
			}
		}
		if variante != nil && nfa != nil {
			variants++
			if !nfa.Match(*variante) {
				report("la ER %q de %q no acepta su variante %q", er, forma, *variante)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return "", err
	}
	missing := 0
	for _, w := range a.Model.Vocabulario {
		if !canonical[w] {
			missing++
			if missing <= 3 {
				report("el lexema %q del vocabulario del modelo no esta en el catalogo", w)
			}
		}
	}
	if problem > 0 {
		return "", fmt.Errorf("%d problemas: %s", problem, strings.Join(details, "; "))
	}
	return fmt.Sprintf("%d lexemas, %d expresiones regulares compiladas y %d variantes aceptadas (I-1); vocabulario del modelo (%d) contenido en el catalogo",
		lexemes, lexemes, variants, len(a.Model.Vocabulario)), nil
}

// checkGold5: the five embedded gold comments are classified with their expected category.
func (a *App) checkGold5() (string, error) {
	var gold []struct {
		ID          int    `json:"id"`
		Texto       string `json:"texto"`
		CategoriaID int    `json:"categoria_id"`
	}
	if err := json.Unmarshal(a.gold, &gold); err != nil {
		return "", fmt.Errorf("gold5.json no es valido: %w", err)
	}
	if len(gold) != 5 {
		return "", fmt.Errorf("gold5.json tiene %d comentarios, se esperaban 5", len(gold))
	}
	var bad []string
	for _, g := range gold {
		if p := a.Model.Predict(g.Texto); p.CategoriaID != g.CategoriaID {
			bad = append(bad, fmt.Sprintf("comentario %d: se predijo la categoria %d (%s), se esperaba %d", g.ID, p.CategoriaID, p.Categoria, g.CategoriaID))
		}
	}
	if len(bad) > 0 {
		return "", fmt.Errorf("%d de 5 mal clasificados: %s", len(bad), strings.Join(bad, "; "))
	}
	return fmt.Sprintf("5 de 5 comentarios de oro clasificados con su categoria esperada (modelo %s/%s)", a.Model.Familia, a.Model.Estrategia), nil
}

// checkPuerto: the primary port is free or the fallback range yields a free port.
func (a *App) checkPuerto() (string, error) {
	ln, err := Listen(a.opts.FirstPort, a.opts.LastPort)
	if err != nil {
		return "", err
	}
	defer ln.Close()
	p := ln.Addr().(*net.TCPAddr).Port
	if p == a.opts.FirstPort {
		return fmt.Sprintf("127.0.0.1:%d libre", p), nil
	}
	return fmt.Sprintf("127.0.0.1:%d ocupado; el fallback asigna 127.0.0.1:%d", a.opts.FirstPort, p), nil
}
