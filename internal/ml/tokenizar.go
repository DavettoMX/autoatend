package ml

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"autoatend/internal/db"
)

// LoadCatalogue builds the Lexicon of every canonical form of the lexemas table
// of a resultados.db (opened read-only).
func LoadCatalogue(path string) (*Lexicon, error) {
	d, err := db.OpenReadOnly(path)
	if err != nil {
		return nil, err
	}
	defer d.Close()
	rows, err := d.Query(`SELECT forma_canonica FROM lexemas ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("ml: reading catalogue %s: %w", path, err)
	}
	defer rows.Close()
	var canon []string
	for rows.Next() {
		var c string
		if err := rows.Scan(&c); err != nil {
			return nil, err
		}
		canon = append(canon, c)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return NewLexicon(canon)
}

type tokenizarIn struct {
	ID    json.RawMessage `json:"id"`
	Texto string          `json:"texto"`
}

type tokenizarOut struct {
	ID      json.RawMessage `json:"id"`
	Lexemas []string        `json:"lexemas"`
}

// RunTokenizar is the logic of cmd/tokenizar (CPA06, TE §6.3.6): it reads JSON
// lines {"id":..,"texto":..} from r and writes, for each one, the JSON line
// {"id":..,"lexemas":[..]} with the catalogue lexemes of the text. The id is
// echoed verbatim. Blank lines are skipped; a malformed line is an error that
// names its line number.
func RunTokenizar(r io.Reader, w io.Writer, lex *Lexicon) error {
	br := bufio.NewReaderSize(r, 1<<20)
	bw := bufio.NewWriter(w)
	for n := 1; ; n++ {
		line, rerr := br.ReadBytes('\n')
		if len(bytes.TrimSpace(line)) > 0 {
			if err := tokenizarLine(line, n, bw, lex); err != nil {
				return err
			}
		}
		if rerr == io.EOF {
			break
		}
		if rerr != nil {
			return rerr
		}
	}
	return bw.Flush()
}

func tokenizarLine(line []byte, n int, w io.Writer, lex *Lexicon) error {
	var in tokenizarIn
	if err := json.Unmarshal(line, &in); err != nil {
		return fmt.Errorf("tokenizar: line %d: %w", n, err)
	}
	if in.ID == nil {
		return fmt.Errorf("tokenizar: line %d: missing id", n)
	}
	out, err := json.Marshal(tokenizarOut{ID: in.ID, Lexemas: lex.Lexemes(in.Texto)})
	if err != nil {
		return err
	}
	out = append(out, '\n')
	_, err = w.Write(out)
	return err
}
