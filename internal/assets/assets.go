// Package assets embeds the artifacts the single binary ships with (TE §7): the two
// SQLite databases, the two models, the five gold comments and the web UI.
//
// go:embed cannot reach outside the package directory, so the artifacts of data/
// are copied into internal/assets/data/ by "make assets" (never edited by hand and
// not versioned: they are byproducts of data/). The UI source is ui/index.html and
// is versioned here. A build that has not run "make assets" fails with
// "pattern data/...: no matching files found"; every make target that compiles Go
// depends on assets.
package assets

import "embed"

// FS holds data/<DataFiles...>, ui/index.html, ui/catalogo.html and ui/api.html.
//
//go:embed data/corpus.db data/resultados.db data/modelo.json data/modelo_alternativo.json data/gold5.json ui/index.html ui/catalogo.html ui/api.html
var FS embed.FS

// DataFiles are the names of the embedded artifacts under data/, in the order the
// startup writes them out.
var DataFiles = []string{"corpus.db", "resultados.db", "modelo.json", "modelo_alternativo.json", "gold5.json"}

// UIPath is the embedded single-page UI.
const UIPath = "ui/index.html"

// CatalogoPath is the embedded catalogue page (a readable view of /api/lexemas).
const CatalogoPath = "ui/catalogo.html"

// ApiDocsPath is the embedded API documentation page served at GET /api.
const ApiDocsPath = "ui/api.html"
