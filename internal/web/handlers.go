package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

const (
	// maxBody is the largest accepted request body (TE §6.3.7: a 1 MB text is rejected with 413).
	maxBody = 1 << 20
	// defaultPageSize and maxPageSize bound /api/lexemas pages.
	defaultPageSize = 100
	maxPageSize     = 500
	// maxPrefix bounds the q filter, in bytes.
	maxPrefix = 100
	// maxRune is the largest code point; prefix + maxRune is above every string with that prefix.
	maxRune = "\U0010FFFF"
)

// Handler returns the HTTP handler of the five endpoints of TE §6.4.1. Every
// database access is a parameterized query; the text of a comment never reaches SQL.
func (a *App) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/{$}", only(http.MethodGet, a.serveUI))
	mux.HandleFunc("/api/clasificar", only(http.MethodPost, a.clasificar))
	mux.HandleFunc("/api/lexemas", only(http.MethodGet, a.lexemasPage))
	mux.HandleFunc("/api/categorias", only(http.MethodGet, a.categoriasList))
	mux.HandleFunc("/api/health", only(http.MethodGet, a.health))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		writeError(w, http.StatusNotFound, "Ruta no encontrada: "+r.URL.Path)
	})
	return recoverer(mux)
}

// only restricts a handler to one method, answering 405 with a JSON error otherwise.
func only(method string, h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.Header().Set("Allow", method)
			writeError(w, http.StatusMethodNotAllowed, fmt.Sprintf("Metodo %s no permitido; use %s.", r.Method, method))
			return
		}
		h(w, r)
	}
}

// recoverer turns a panic into a 500 JSON error instead of dropping the connection.
func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				writeError(w, http.StatusInternalServerError, "Error interno del servidor.")
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		status = http.StatusInternalServerError
		buf.Reset()
		buf.WriteString(`{"error":"Error interno del servidor."}` + "\n")
	}
	h := w.Header()
	h.Set("Content-Type", "application/json; charset=utf-8")
	h.Set("Cache-Control", "no-store")
	h.Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// serveUI answers GET / with the embedded single-page UI. The CSP forbids every
// network access except the same-origin API, so nothing is ever fetched from outside.
func (a *App) serveUI(w http.ResponseWriter, r *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Cache-Control", "no-store")
	h.Set("X-Content-Type-Options", "nosniff")
	h.Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; script-src 'unsafe-inline'; connect-src 'self'; base-uri 'none'; form-action 'none'")
	w.Write(a.ui)
}

type clasificarIn struct {
	Texto *string `json:"texto"`
}

type lexemaClave struct {
	Lexema           string  `json:"lexema"`
	ExpresionRegular string  `json:"expresion_regular"`
	Contribucion     float64 `json:"contribucion"`
}

type clasificarOut struct {
	Categoria    string             `json:"categoria"`
	Confianzas   map[string]float64 `json:"confianzas"`
	LexemasClave []lexemaClave      `json:"lexemas_clave"`
	Vacio        bool               `json:"vacio"`
}

// clasificar implements POST /api/clasificar (TE §6.4.1, §6.3.7).
func (a *App) clasificar(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBody)
	var in clasificarIn
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		var tooBig *http.MaxBytesError
		switch {
		case errors.As(err, &tooBig):
			writeError(w, http.StatusRequestEntityTooLarge, fmt.Sprintf("El cuerpo de la solicitud supera el maximo de %d bytes (1 MB).", maxBody))
		case errors.Is(err, io.EOF):
			writeError(w, http.StatusBadRequest, `El cuerpo esta vacio; envie un JSON con el campo "texto".`)
		default:
			writeError(w, http.StatusBadRequest, `JSON no valido; se espera {"texto": "..."} con el texto como cadena.`)
		}
		return
	}
	if in.Texto == nil || strings.TrimSpace(*in.Texto) == "" {
		writeError(w, http.StatusBadRequest, `El campo "texto" es obligatorio y no puede estar vacio.`)
		return
	}

	p := a.Model.Predict(*in.Texto)
	out := clasificarOut{
		Categoria:    p.Categoria,
		Confianzas:   make(map[string]float64, len(p.Confianzas)),
		LexemasClave: make([]lexemaClave, 0, len(p.LexemasClave)),
		Vacio:        p.Vacio,
	}
	for k, c := range a.Model.Categorias {
		out.Confianzas[c.Nombre] = p.Confianzas[k]
	}
	for _, lc := range p.LexemasClave {
		var er string
		// The regular expression comes from the catalogue of resultados.db (consigna b).
		if err := a.erStmt.QueryRow(lc.Lexema).Scan(&er); err != nil {
			writeError(w, http.StatusInternalServerError, "No se pudo leer el catalogo de expresiones regulares.")
			return
		}
		out.LexemasClave = append(out.LexemasClave, lexemaClave{Lexema: lc.Lexema, ExpresionRegular: er, Contribucion: lc.Contribucion})
	}
	writeJSON(w, http.StatusOK, out)
}

type lexemaOut struct {
	FormaCanonica    string         `json:"forma_canonica"`
	ExpresionRegular string         `json:"expresion_regular"`
	FrecuenciaGlobal int            `json:"frecuencia_global"`
	Frecuencias      map[string]int `json:"frecuencias"`
	Variantes        int            `json:"variantes"`
}

type lexemasOut struct {
	Total   int         `json:"total"`
	Pagina  int         `json:"pagina"`
	Tamanio int         `json:"tamanio"`
	Lexemas []lexemaOut `json:"lexemas"`
}

// lexemasOrder maps the orden parameter to a fixed ORDER BY clause (never built from input).
var lexemasOrder = map[string]string{
	"":           "l.forma_canonica",
	"frecuencia": "glob DESC, l.forma_canonica",
	"variantes":  "nvar DESC, l.forma_canonica",
}

// lexemasPage implements GET /api/lexemas?pagina=N&tamanio=M[&q=prefijo][&orden=frecuencia|variantes].
func (a *App) lexemasPage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	pagina, err := intParam(q.Get("pagina"), 1)
	if err != nil || pagina < 1 {
		writeError(w, http.StatusBadRequest, `El parametro "pagina" debe ser un entero >= 1.`)
		return
	}
	tamanio, err := intParam(q.Get("tamanio"), defaultPageSize)
	if err != nil || tamanio < 1 {
		writeError(w, http.StatusBadRequest, `El parametro "tamanio" debe ser un entero >= 1.`)
		return
	}
	if tamanio > maxPageSize {
		tamanio = maxPageSize
	}
	order, ok := lexemasOrder[q.Get("orden")]
	if !ok {
		writeError(w, http.StatusBadRequest, `El parametro "orden" debe ser "frecuencia" o "variantes".`)
		return
	}
	prefix := strings.ToLower(q.Get("q"))
	if len(prefix) > maxPrefix {
		writeError(w, http.StatusBadRequest, `El parametro "q" es demasiado largo.`)
		return
	}
	// Prefix filter as a range on the unique index: prefix <= forma < prefix + max code point.
	lo, hi := prefix, prefix+maxRune

	out := lexemasOut{Pagina: pagina, Tamanio: tamanio, Lexemas: []lexemaOut{}}
	if err := a.resultados.QueryRow(`SELECT COUNT(*) FROM lexemas WHERE forma_canonica >= ? AND forma_canonica < ?`, lo, hi).Scan(&out.Total); err != nil {
		writeError(w, http.StatusInternalServerError, "No se pudo leer el catalogo de lexemas.")
		return
	}
	rows, err := a.resultados.Query(`
		SELECT l.id, l.forma_canonica, l.expresion_regular,
		       COALESCE((SELECT f.valor FROM frecuencias f WHERE f.lexema_id = l.id AND f.categoria_id IS NULL), 0) AS glob,
		       (SELECT COUNT(*) FROM variantes v WHERE v.lexema_id = l.id) AS nvar
		FROM lexemas l
		WHERE l.forma_canonica >= ? AND l.forma_canonica < ?
		ORDER BY `+order+`
		LIMIT ? OFFSET ?`, lo, hi, tamanio, (pagina-1)*tamanio)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "No se pudo leer el catalogo de lexemas.")
		return
	}
	defer rows.Close()
	index := map[int64]int{}
	var ids []any
	for rows.Next() {
		var id int64
		var l lexemaOut
		if err := rows.Scan(&id, &l.FormaCanonica, &l.ExpresionRegular, &l.FrecuenciaGlobal, &l.Variantes); err != nil {
			writeError(w, http.StatusInternalServerError, "No se pudo leer el catalogo de lexemas.")
			return
		}
		l.Frecuencias = make(map[string]int, len(a.categorias))
		for _, c := range a.categorias {
			l.Frecuencias[c.Nombre] = 0
		}
		index[id] = len(out.Lexemas)
		ids = append(ids, id)
		out.Lexemas = append(out.Lexemas, l)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "No se pudo leer el catalogo de lexemas.")
		return
	}
	rows.Close()

	if len(ids) > 0 {
		names := make(map[int]string, len(a.categorias))
		for _, c := range a.categorias {
			names[c.ID] = c.Nombre
		}
		fr, err := a.resultados.Query(`SELECT lexema_id, categoria_id, valor FROM frecuencias
			WHERE categoria_id IS NOT NULL AND lexema_id IN (?`+strings.Repeat(",?", len(ids)-1)+`)`, ids...)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "No se pudo leer las frecuencias.")
			return
		}
		defer fr.Close()
		for fr.Next() {
			var lid int64
			var cid, v int
			if err := fr.Scan(&lid, &cid, &v); err != nil {
				writeError(w, http.StatusInternalServerError, "No se pudo leer las frecuencias.")
				return
			}
			if name, ok := names[cid]; ok {
				out.Lexemas[index[lid]].Frecuencias[name] = v
			}
		}
		if err := fr.Err(); err != nil {
			writeError(w, http.StatusInternalServerError, "No se pudo leer las frecuencias.")
			return
		}
	}
	writeJSON(w, http.StatusOK, out)
}

func intParam(s string, def int) (int, error) {
	if s == "" {
		return def, nil
	}
	return strconv.Atoi(s)
}

// categoriasList implements GET /api/categorias: the five labels of the case with their ids.
func (a *App) categoriasList(w http.ResponseWriter, r *http.Request) {
	type cat struct {
		ID     int    `json:"id"`
		Nombre string `json:"nombre"`
	}
	out := make([]cat, len(a.categorias))
	for i, c := range a.categorias {
		out[i] = cat{c.ID, c.Nombre}
	}
	writeJSON(w, http.StatusOK, out)
}

// health implements GET /api/health: exactly {"ok","version","lexemas","puerto"}.
func (a *App) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, struct {
		OK      bool   `json:"ok"`
		Version string `json:"version"`
		Lexemas int    `json:"lexemas"`
		Puerto  int    `json:"puerto"`
	}{true, a.opts.Version, a.lexemas, a.Port()})
}
