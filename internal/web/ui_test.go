package web_test

import (
	"net/http"
	"strings"
	"testing"

	"autoatend/internal/assets"
)

// GET /: the single embedded HTML page, system fonts, no CDN (TE §6.4.1).
func TestUI_PaginaUnicaSinRecursosExternos(t *testing.T) {
	h := bootApp(t).Handler()
	rec := do(t, h, http.MethodGet, "/", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
	want, _ := assets.FS.ReadFile(assets.UIPath)
	if rec.Body.String() != string(want) {
		t.Error("GET / does not serve the embedded ui/index.html")
	}
	csp := rec.Header().Get("Content-Security-Policy")
	for _, d := range []string{"default-src 'none'", "connect-src 'self'"} {
		if !strings.Contains(csp, d) {
			t.Errorf("Content-Security-Policy %q lacks %q", csp, d)
		}
	}

	html := rec.Body.String()
	// Nothing is fetched from outside the page: no absolute URLs, no external scripts, styles, fonts or imports.
	for _, bad := range []string{"http://", "https://", "//cdn", "<link", "@import", "@font-face", "<script src", "<img src=\"http", "url("} {
		if strings.Contains(html, bad) {
			t.Errorf("index.html contains %q (external or forbidden resource)", bad)
		}
	}
}

func TestUI_ContenidoYAccesibilidad(t *testing.T) {
	html := bootApp(t).Handler()
	rec := do(t, html, http.MethodGet, "/", nil)
	body := rec.Body.String()
	for _, want := range []string{
		`<html lang="es"`,
		`<meta name="viewport"`,
		`<main`,
		`<label for="texto"`,
		`<textarea id="texto"`,
		`<button`,
		`Clasificar`,
		`aria-live="polite"`,
		`/api/clasificar`,
		`lexemas_clave`,
		`expresion_regular`,
		`confianzas`,
		`vacio`,
		`prefers-color-scheme`,
		`system-ui`,
	} {
		if !strings.Contains(body, want) {
			t.Errorf("index.html lacks %q", want)
		}
	}
	// Results are rendered with textContent: the page never builds HTML from API data.
	if strings.Contains(body, "innerHTML") {
		t.Error("index.html uses innerHTML")
	}
}

// GET /catalogo: a readable table of the catalogue built from /api/lexemas,
// with the same no-external-resources and textContent rules as the main page.
func TestUI_CatalogoPagina(t *testing.T) {
	h := bootApp(t).Handler()
	rec := do(t, h, http.MethodGet, "/catalogo", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
	want, _ := assets.FS.ReadFile(assets.CatalogoPath)
	if len(want) == 0 || rec.Body.String() != string(want) {
		t.Fatal("GET /catalogo does not serve the embedded ui/catalogo.html")
	}
	if csp := rec.Header().Get("Content-Security-Policy"); !strings.Contains(csp, "default-src 'none'") {
		t.Errorf("Content-Security-Policy %q", csp)
	}
	body := rec.Body.String()
	for _, bad := range []string{"http://", "https://", "//cdn", "<link", "@import", "@font-face", "<script src", "url(", "innerHTML"} {
		if strings.Contains(body, bad) {
			t.Errorf("catalogo.html contains %q", bad)
		}
	}
	for _, s := range []string{`<html lang="es"`, `<main`, `/api/lexemas`, `min_frecuencia`, `orden`, `<table`, `<label`, `href="/"`} {
		if !strings.Contains(body, s) {
			t.Errorf("catalogo.html lacks %q", s)
		}
	}
}

// The main page links to the readable catalogue and documents both keyboard shortcuts.
func TestUI_EnlaceAlCatalogoYAtajos(t *testing.T) {
	body := do(t, bootApp(t).Handler(), http.MethodGet, "/", nil).Body.String()
	for _, s := range []string{`href="/catalogo`, `⌘`} {
		if !strings.Contains(body, s) {
			t.Errorf("index.html lacks %q", s)
		}
	}
}

// Every page stays in the same tab: no target="_blank" anywhere in the UI, and every
// page reachable from the UI offers a visible way back to the classifier.
func TestUI_SinPestanasNuevas(t *testing.T) {
	h := bootApp(t).Handler()
	for _, path := range []string{"/", "/catalogo", "/api"} {
		body := do(t, h, http.MethodGet, path, nil).Body.String()
		if strings.Contains(body, `target="_blank"`) {
			t.Errorf("%s contains target=\"_blank\"", path)
		}
		if path != "/" && !strings.Contains(body, "Volver al clasificador") {
			t.Errorf("%s lacks the link \"Volver al clasificador\"", path)
		}
	}
}

// The footer of the main page names the two secondary pages, never the loopback address,
// and never links to raw JSON (a browser tab on raw JSON has no way back).
func TestUI_PieDeLaPaginaPrincipal(t *testing.T) {
	body := do(t, bootApp(t).Handler(), http.MethodGet, "/", nil).Body.String()
	for _, s := range []string{
		`<a href="/catalogo">Catálogo de expresiones regulares</a>`,
		`<a href="/api">Documentación de la API</a>`,
	} {
		if !strings.Contains(body, s) {
			t.Errorf("index.html lacks %q", s)
		}
	}
	for _, bad := range []string{"127.0.0.1", "consigna", `href="/api/`} {
		if strings.Contains(body, bad) {
			t.Errorf("index.html contains %q", bad)
		}
	}
}

// The catalogue shows its JSON in an in-page collapsible panel instead of a link.
func TestUI_CatalogoPanelJSONSinConsigna(t *testing.T) {
	body := do(t, bootApp(t).Handler(), http.MethodGet, "/catalogo", nil).Body.String()
	for _, s := range []string{"<details", "<summary>Ver datos en JSON</summary>", "JSON.stringify(", "<pre"} {
		if !strings.Contains(body, s) {
			t.Errorf("catalogo.html lacks %q", s)
		}
	}
	for _, bad := range []string{"consigna", `id="json" href`, `href="/api/`} {
		if strings.Contains(body, bad) {
			t.Errorf("catalogo.html contains %q", bad)
		}
	}
}

// GET /api: the documentation page of the five endpoints of TE §6.4.1, served like the other pages.
func TestUI_DocumentacionAPI(t *testing.T) {
	h := bootApp(t).Handler()
	rec := do(t, h, http.MethodGet, "/api", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
	want, _ := assets.FS.ReadFile(assets.ApiDocsPath)
	if len(want) == 0 || rec.Body.String() != string(want) {
		t.Fatal("GET /api does not serve the embedded ui/api.html")
	}
	csp := rec.Header().Get("Content-Security-Policy")
	for _, d := range []string{"default-src 'none'", "connect-src 'self'"} {
		if !strings.Contains(csp, d) {
			t.Errorf("Content-Security-Policy %q lacks %q", csp, d)
		}
	}
	body := rec.Body.String()
	for _, bad := range []string{"http://", "https://", "//cdn", "<link", "@import", "@font-face", "<script src", "url(", "innerHTML", `target="_blank"`, "consigna", "Estas cinco rutas", "Trabajo Escrito", "solo lectura", "no requiere internet"} {
		if strings.Contains(body, bad) {
			t.Errorf("api.html contains %q", bad)
		}
	}
	for _, s := range []string{
		`<html lang="es"`, `<main`, `href="/"`, "Volver al clasificador",
		"POST /api/clasificar", "GET /api/lexemas", "GET /api/categorias", "GET /api/health", "GET /",
		"Probar",
		`<p class="lead">Servicio web HTTP REST local de AutoAtend. Las respuestas son JSON en UTF-8.</p>`,
	} {
		if !strings.Contains(body, s) {
			t.Errorf("api.html lacks %q", s)
		}
	}
}

// Every endpoint documents its error responses as one block per status code, shaped like the
// 200 block (badge heading, one explanatory line, its own example body) and carrying the real
// messages of the handlers; there is no table of codes.
func TestUI_APIRespuestasDeError(t *testing.T) {
	body := do(t, bootApp(t).Handler(), http.MethodGet, "/api", nil).Body.String()
	if n := strings.Count(body, "<h3>Respuestas de error</h3>"); n < 4 {
		t.Errorf("api.html has %d \"Respuestas de error\" subheadings, want one per endpoint (>= 4)", n)
	}
	if n := strings.Count(body, `<span class="cod ok">200</span>`); n < 4 {
		t.Errorf("api.html has %d 200 badges, want one per endpoint (>= 4)", n)
	}
	// 3 (clasificar) + 2 (lexemas) + 1 (categorias) + 1 (health) + 404 and 500 in the common section.
	if n := strings.Count(body, `<h3><span class="cod err">`); n < 9 {
		t.Errorf("api.html has %d error blocks with a badge heading, want >= 9", n)
	}
	for _, code := range []string{"400", "405", "413", "404", "500"} {
		if !strings.Contains(body, `<h3><span class="cod err">`+code+`</span>`) {
			t.Errorf("api.html lacks an error block headed by the %s badge", code)
		}
	}
	for _, s := range []string{
		`{"error": "El campo \"texto\" es obligatorio y no puede estar vacio."}`,
		`{"error": "El cuerpo esta vacio; envie un JSON con el campo \"texto\"."}`,
		`{"error": "JSON no valido; se espera {\"texto\": \"...\"} con el texto como cadena."}`,
		`{"error": "El cuerpo de la solicitud supera el maximo de 1048576 bytes (1 MB)."}`,
		`{"error": "Metodo GET no permitido; use POST."}`,
		`{"error": "Metodo POST no permitido; use GET."}`,
		`{"error": "El parametro \"pagina\" debe ser un entero &gt;= 1."}`,
		`{"error": "Ruta no encontrada: /api/otra"}`,
		`{"error": "No se pudo leer el catalogo de lexemas."}`,
		"Falla interna al leer el catálogo. No ocurre en operación normal: las bases de datos se extraen y verifican al arrancar (y con --selftest).",
	} {
		if !strings.Contains(body, s) {
			t.Errorf("api.html lacks %q", s)
		}
	}
	for _, bad := range []string{"Errores: 400", `<th scope="col">Código</th>`, `<th scope="col">Cuándo</th>`} {
		if strings.Contains(body, bad) {
			t.Errorf("api.html still contains %q (errors must be blocks, not a table)", bad)
		}
	}
}

// /api is exact: the data routes below it keep working and unknown subpaths stay 404.
func TestUI_RutaAPIExacta(t *testing.T) {
	h := bootApp(t).Handler()
	if rec := do(t, h, http.MethodGet, "/api/health", nil); !strings.HasPrefix(rec.Header().Get("Content-Type"), "application/json") {
		t.Errorf("/api/health Content-Type = %q", rec.Header().Get("Content-Type"))
	}
	if rec := do(t, h, http.MethodGet, "/api/otra", nil); rec.Code != http.StatusNotFound {
		t.Errorf("/api/otra status %d, want 404", rec.Code)
	}
	if rec := do(t, h, http.MethodPost, "/api", nil); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /api status %d, want 405", rec.Code)
	}
}
