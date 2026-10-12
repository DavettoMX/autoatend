// Package sweep is the simulator pipeline (consigna b, TE §6.2 and §6.1.3): it
// sweeps corpus.db once, builds the symbol table with the lexer, persists the
// catalogue in resultados.db (corridas, lexemas, frecuencias, variantes),
// verifies the file by re-opening it, and exports the document frequency of
// the training split for the offline trainer. The run is deterministic: no
// timestamps, no map iteration while writing, no goroutines, so two runs over
// the same corpus produce byte-identical files (invariant I-3).
package sweep

import (
	"bufio"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"autoatend/internal/corpus"
	"autoatend/internal/db"
	"autoatend/internal/lexer"
)

// Options configures one sweep.
type Options struct {
	CorpusPath string // corpus.db to read
	OutPath    string // resultados.db to create (must not exist)
	DFPath     string // df_train.json to create (must not exist); empty = do not write
	Seed       int64  // recorded in corridas.semilla
}

// Result is everything a run measured.
type Result struct {
	Comentarios int           // comments swept
	TrainDocs   int           // comments of the training split
	Cutoff      string        // first fecha of the test split (db.SplitCutoff)
	Lexemas     int           // lexemes in the catalogue
	Tokens      int           // catalogued word occurrences (after the stoplist)
	ZeroDF      int           // lexemes that never occur in a training comment
	SHA256      string        // of the finished resultados.db
	Elapsed     time.Duration // wall time of the whole run
}

// Run is the whole pipeline: read -> count -> persist -> verify -> export DF.
// On any error nothing is left behind: a resultados.db or df_train.json that
// the run created is removed, and a file that already existed is never touched.
func Run(opts Options) (res Result, err error) {
	start := time.Now()
	if opts.DFPath != "" {
		if _, serr := os.Stat(opts.DFPath); serr == nil {
			return res, fmt.Errorf("sweep: refusing to overwrite existing file %s", opts.DFPath)
		} else if !errors.Is(serr, os.ErrNotExist) {
			return res, serr
		}
	}

	src, err := db.OpenReadOnly(opts.CorpusPath)
	if err != nil {
		return res, err
	}
	defer src.Close()

	sw, err := count(src)
	if err != nil {
		return res, err
	}
	res.Comentarios, res.TrainDocs, res.Cutoff, res.Tokens = sw.comentarios, sw.trainDocs, sw.cutoff, sw.tokens

	entries := sw.table.Entries()
	if len(entries) == 0 {
		return res, errors.New("sweep: the corpus yields no lexemes")
	}
	res.Lexemas = len(entries)

	params, err := marshalParametros(opts.Seed, sw.comentarios, sw.cutoff)
	if err != nil {
		return res, err
	}
	if err := persist(opts.OutPath, entries, opts.Seed, params); err != nil {
		return res, err
	}
	defer func() {
		if err != nil {
			os.Remove(opts.OutPath)
		}
	}()
	if err = verify(opts.OutPath, opts.CorpusPath, len(entries)); err != nil {
		return res, err
	}

	var dfs []dfEntry
	for _, e := range entries {
		if n := sw.df[e.Stem]; n > 0 {
			dfs = append(dfs, dfEntry{Lexema: e.Canonical(), DF: n})
		} else {
			res.ZeroDF++
		}
	}
	if opts.DFPath != "" {
		if err = writeDF(opts.DFPath, sw.cutoff, sw.trainDocs, dfs); err != nil {
			return res, err
		}
	}

	raw, err := os.ReadFile(opts.OutPath)
	if err != nil {
		return res, err
	}
	sum := sha256.Sum256(raw)
	res.SHA256 = hex.EncodeToString(sum[:])
	res.Elapsed = time.Since(start)
	return res, nil
}

// swept is the outcome of the single pass over corpus.db.
type swept struct {
	table       *lexer.Table
	df          map[string]int // document frequency of train, by stem
	comentarios int
	trainDocs   int
	cutoff      string
	tokens      int
}

// count reads every comment once, in id order: Analyze, Insert into the
// symbol table, and, for training comments (fecha < cut-off), add one to the
// document frequency of each distinct lexeme of the comment.
func count(src *sql.DB) (*swept, error) {
	sw := &swept{table: lexer.NewTable(), df: make(map[string]int)}
	var err error
	if sw.cutoff, err = db.SplitCutoff(src); err != nil {
		return nil, err
	}
	if err = src.QueryRow(`SELECT COUNT(*) FROM comentarios`).Scan(&sw.comentarios); err != nil {
		return nil, err
	}
	rows, err := src.Query(`SELECT texto, categoria_id, fecha < ? FROM comentarios ORDER BY id`, sw.cutoff)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	seen := make(map[string]struct{})
	n := 0
	for rows.Next() {
		var texto string
		var cat int
		var train bool
		if err := rows.Scan(&texto, &cat, &train); err != nil {
			return nil, err
		}
		n++
		clear(seen)
		for _, form := range lexer.Analyze(texto) {
			e, err := sw.table.Insert(form, cat)
			if err != nil {
				return nil, fmt.Errorf("sweep: comment %d: %w", n, err)
			}
			sw.tokens++
			if train {
				if _, dup := seen[e.Stem]; !dup {
					seen[e.Stem] = struct{}{}
					sw.df[e.Stem]++
				}
			}
		}
		if train {
			sw.trainDocs++
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if n != sw.comentarios {
		return nil, fmt.Errorf("sweep: read %d comments, COUNT(*) says %d", n, sw.comentarios)
	}
	return sw, nil
}

type dfEntry struct {
	Lexema string
	DF     int
}

func writeDF(path, cutoff string, trainDocs int, dfs []dfEntry) (err error) {
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); err == nil {
			err = cerr
		}
		if err != nil {
			os.Remove(path)
		}
	}()
	w := bufio.NewWriter(f)
	corte, _ := json.Marshal(cutoff)
	fmt.Fprintf(w, "{\"corte\":%s,\"documentos_train\":%d,\"df\":[", corte, trainDocs)
	for i, e := range dfs {
		sep := "\n"
		if i > 0 {
			sep = ",\n"
		}
		lex, _ := json.Marshal(e.Lexema)
		fmt.Fprintf(w, "%s{\"lexema\":%s,\"df\":%d}", sep, lex, e.DF)
	}
	if _, err = io.WriteString(w, "\n]}\n"); err != nil {
		return err
	}
	return w.Flush()
}

func marshalParametros(seed int64, comentarios int, cutoff string) ([]byte, error) {
	n := corpus.DefaultNoise
	return json.Marshal(parametros{
		Semilla:     seed,
		Comentarios: comentarios,
		Ruido: ruido{
			Tasa: n.Rate, Tilde: n.TildeDrop, Abreviaturas: n.Abbrev,
			Erratas: n.Typo, Mayusculas: n.CaseJitter, Puntuacion: n.Punct,
		},
		Stoplist: lexer.Version(),
		Tokenizador: tokenizador{
			Alfabeto: lexer.Alphabet,
			Digitos:  "0-9",
			Numeros:  "reconocidos por el AFD, no catalogados",
		},
		Stemmer:    lexer.StemVersion,
		CorteTrain: cutoff,
	})
}

// parametros is the content of corridas.parametros (TE §6.1.4). The field
// order is fixed, so the serialized JSON is reproducible.
type parametros struct {
	Semilla     int64       `json:"semilla"`
	Comentarios int         `json:"comentarios"`
	Ruido       ruido       `json:"ruido"`
	Stoplist    string      `json:"stoplist"`
	Tokenizador tokenizador `json:"tokenizador"`
	Stemmer     string      `json:"stemmer"`
	CorteTrain  string      `json:"corte_train"`
}

// ruido mirrors corpus.NoiseConfig, the calibrated configuration the corpus was generated with.
type ruido struct {
	Tasa         float64 `json:"tasa"`
	Tilde        float64 `json:"tilde"`
	Abreviaturas float64 `json:"abreviaturas"`
	Erratas      float64 `json:"erratas"`
	Mayusculas   float64 `json:"mayusculas"`
	Puntuacion   float64 `json:"puntuacion"`
}

type tokenizador struct {
	Alfabeto string `json:"alfabeto"`
	Digitos  string `json:"digitos"`
	Numeros  string `json:"numeros"`
}
