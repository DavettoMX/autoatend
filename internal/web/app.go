// Package web is the AutoAtend web service (consigna e, TE §6.4): it extracts the
// embedded artifacts to a temporary directory, opens the databases read-only,
// serves the five endpoints on 127.0.0.1 and runs the --selftest battery.
//
// Console messages are plain ASCII on purpose: the cmd.exe code page of the exam
// machine is not UTF-8 and would garble accents. JSON and HTML are UTF-8.
package web

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sync/atomic"

	"autoatend/internal/assets"
	"autoatend/internal/db"
	"autoatend/internal/ml"
)

// Default network range (TE §6.4.1).
const (
	DefaultFirstPort = 8471
	DefaultLastPort  = 8480
)

// Options configure Boot and Run. The zero value is the production configuration.
type Options struct {
	// Modelo selects the model: "principal" (default, modelo.json) or
	// "alternativo" (modelo_alternativo.json, the TE §6.4.2 replica).
	Modelo string
	// Version is reported by /api/health and the banner.
	Version string
	// Assets holds the artifacts (data/<DataFiles...>, ui/index.html); default assets.FS.
	Assets fs.FS
	// NoBrowser disables opening the browser.
	NoBrowser bool
	// Selftest runs the headless checks instead of serving.
	Selftest bool
	// FirstPort and LastPort bound the port search; default 8471..8480.
	FirstPort, LastPort int
	// OpenBrowser opens a URL; default is the per-OS launcher.
	OpenBrowser func(url string) error
}

func (o Options) withDefaults() Options {
	if o.Modelo == "" {
		o.Modelo = "principal"
	}
	if o.Version == "" {
		o.Version = "dev"
	}
	if o.Assets == nil {
		o.Assets = assets.FS
	}
	if o.FirstPort == 0 {
		o.FirstPort = DefaultFirstPort
	}
	if o.LastPort == 0 {
		o.LastPort = DefaultLastPort
	}
	if o.OpenBrowser == nil {
		o.OpenBrowser = openBrowser
	}
	return o
}

func (o Options) modelFile() (string, error) {
	switch o.Modelo {
	case "principal":
		return "modelo.json", nil
	case "alternativo":
		return "modelo_alternativo.json", nil
	}
	return "", fmt.Errorf("valor %q no valido para -modelo (use principal o alternativo)", o.Modelo)
}

// App is a started service: the extracted artifacts, the read-only databases and the model.
type App struct {
	// Dir is the process-private temporary directory holding the extracted artifacts.
	Dir string
	// Model is the loaded classifier.
	Model *ml.Model

	opts       Options
	corpus     *sql.DB
	resultados *sql.DB
	erStmt     *sql.Stmt
	ui         []byte
	catalogo   []byte
	gold       []byte
	categorias []db.Category
	lexemas    int
	orphans    []int64
	port       atomic.Int32
	closed     atomic.Bool
}

// Boot extracts the embedded artifacts to a fresh directory under the user's
// temporary directory (no administrator rights needed), opens corpus.db and
// resultados.db with mode=ro&immutable=1, checks the cross-file categoria_id
// integrity with ATTACH (TE §6.1.3) and loads the model. Any problem is an error
// that leaves nothing behind; Close removes the directory.
func Boot(opts Options) (app *App, err error) {
	opts = opts.withDefaults()
	modelName, err := opts.modelFile()
	if err != nil {
		return nil, err
	}
	dir, err := os.MkdirTemp("", "autoatend-")
	if err != nil {
		return nil, fmt.Errorf("no se pudo crear el directorio temporal: %w", err)
	}
	a := &App{Dir: dir, opts: opts}
	defer func() {
		if err != nil {
			a.Close()
			app = nil
		}
	}()

	for _, name := range assets.DataFiles {
		b, rerr := fs.ReadFile(opts.Assets, "data/"+name)
		if rerr != nil {
			return nil, fmt.Errorf("artefacto embebido %s: %w", name, rerr)
		}
		if name == "gold5.json" {
			a.gold = b
		}
		if werr := os.WriteFile(filepath.Join(dir, name), b, 0o600); werr != nil {
			return nil, fmt.Errorf("no se pudo escribir %s en %s: %w", name, dir, werr)
		}
	}
	if a.ui, err = fs.ReadFile(opts.Assets, assets.UIPath); err != nil {
		return nil, fmt.Errorf("artefacto embebido %s: %w", assets.UIPath, err)
	}
	if a.catalogo, err = fs.ReadFile(opts.Assets, assets.CatalogoPath); err != nil {
		return nil, fmt.Errorf("artefacto embebido %s: %w", assets.CatalogoPath, err)
	}

	if a.corpus, err = db.OpenReadOnly(filepath.Join(dir, "corpus.db")); err != nil {
		return nil, fmt.Errorf("corpus.db: %w", err)
	}
	if a.resultados, err = db.OpenReadOnly(filepath.Join(dir, "resultados.db")); err != nil {
		return nil, fmt.Errorf("resultados.db: %w", err)
	}
	if err = a.loadCategorias(); err != nil {
		return nil, err
	}
	if a.orphans, err = db.CheckCategoriaIntegrity(a.resultados, filepath.Join(dir, "corpus.db")); err != nil {
		return nil, fmt.Errorf("integridad cruzada de categoria_id (ATTACH): %w", err)
	}
	if len(a.orphans) > 0 {
		return nil, fmt.Errorf("integridad cruzada: resultados.db tiene categoria_id sin fila en corpus.db (categorias): %v", a.orphans)
	}
	if a.Model, err = ml.Load(filepath.Join(dir, modelName)); err != nil {
		return nil, err
	}
	if err = a.resultados.QueryRow(`SELECT COUNT(*) FROM lexemas`).Scan(&a.lexemas); err != nil {
		return nil, fmt.Errorf("resultados.db: contando lexemas: %w", err)
	}
	if a.lexemas == 0 {
		return nil, fmt.Errorf("resultados.db: el catalogo de lexemas esta vacio")
	}
	if a.erStmt, err = a.resultados.Prepare(`SELECT expresion_regular FROM lexemas WHERE forma_canonica = ?`); err != nil {
		return nil, fmt.Errorf("resultados.db: %w", err)
	}
	return a, nil
}

// loadCategorias reads corpus.categorias and requires exactly the five labels of the case.
func (a *App) loadCategorias() error {
	rows, err := a.corpus.Query(`SELECT id, nombre FROM categorias ORDER BY id`)
	if err != nil {
		return fmt.Errorf("corpus.db: leyendo categorias: %w", err)
	}
	defer rows.Close()
	for rows.Next() {
		var c db.Category
		if err := rows.Scan(&c.ID, &c.Nombre); err != nil {
			return err
		}
		a.categorias = append(a.categorias, c)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if len(a.categorias) != len(db.Categories) {
		return fmt.Errorf("corpus.db: %d categorias, el caso define %d", len(a.categorias), len(db.Categories))
	}
	for i, c := range db.Categories {
		if a.categorias[i] != c {
			return fmt.Errorf("corpus.db: categoria %d es %+v, el caso define %+v", i, a.categorias[i], c)
		}
	}
	return nil
}

// Close closes the databases and removes the temporary directory. It is safe to call twice.
func (a *App) Close() error {
	if !a.closed.CompareAndSwap(false, true) {
		return nil
	}
	if a.erStmt != nil {
		a.erStmt.Close()
	}
	if a.corpus != nil {
		a.corpus.Close()
	}
	if a.resultados != nil {
		a.resultados.Close()
	}
	return os.RemoveAll(a.Dir)
}

// SetPort records the port the service listens on, reported by /api/health.
func (a *App) SetPort(p int) { a.port.Store(int32(p)) }

// Port is the port set by SetPort (0 before listening).
func (a *App) Port() int { return int(a.port.Load()) }
