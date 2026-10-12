package web_test

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"autoatend/internal/assets"
	"autoatend/internal/casos"
	"autoatend/internal/web"
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
	for _, path := range []string{"/", "/catalogo", "/api", "/casos-prueba"} {
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
		`<a href="/casos-prueba">Casos de prueba</a>`,
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

// dataBlock extracts the JSON that the server injects into /casos-prueba.
func dataBlock(t *testing.T, body string) string {
	t.Helper()
	const open = `<script type="application/json" id="casos-datos">`
	i := strings.Index(body, open)
	if i < 0 {
		t.Fatal("casos-prueba page lacks the injected JSON block")
	}
	rest := body[i+len(open):]
	j := strings.Index(rest, "</script>")
	if j < 0 {
		t.Fatal("the injected JSON block is not closed")
	}
	return rest[:j]
}

// GET /casos-prueba: the Anexo II matrix, served like the other pages (exact path, CSP, no external
// resources, textContent only) with the JSON of the matrix injected by the server. There is no
// extra /api route: the five endpoints of TE §6.4.1 are unchanged.
func TestUI_CasosPruebaPagina(t *testing.T) {
	h := bootApp(t).Handler()
	rec := do(t, h, http.MethodGet, "/casos-prueba", nil)
	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Errorf("Content-Type = %q", ct)
	}
	csp := rec.Header().Get("Content-Security-Policy")
	for _, d := range []string{"default-src 'none'", "connect-src 'self'"} {
		if !strings.Contains(csp, d) {
			t.Errorf("Content-Security-Policy %q lacks %q", csp, d)
		}
	}
	// The injected data may quote URLs as text (CPE03, CPA05); only the page itself is checked.
	body := strings.Replace(rec.Body.String(), dataBlock(t, rec.Body.String()), "", 1)
	for _, bad := range []string{"http://", "https://", "//cdn", "<link", "@import", "@font-face", "<script src", "url(", "innerHTML", `target="_blank"`, "{{CASOS_JSON}}", "consigna"} {
		if strings.Contains(body, bad) {
			t.Errorf("casos.html contains %q", bad)
		}
	}
	for _, s := range []string{
		`<html lang="es"`, `<meta name="viewport"`, `<main`, `href="/"`, "Volver al clasificador",
		"<h1>Casos de prueba</h1>", `<label for="nivel"`, `<select id="nivel"`, "<table", "prefers-color-scheme",
		"aria-expanded", "textContent", "Anexo II",
	} {
		if !strings.Contains(body, s) {
			t.Errorf("casos.html lacks %q", s)
		}
	}
	if strings.Count(body, `href="/`) > 2 {
		t.Error("casos.html links to more than the way back (and the API docs)")
	}
}

// The injected JSON is exactly the embedded matrix, safe inside a script element, and carries every
// ID of the Anexo II with the tests that cover it.
func TestUI_CasosPruebaDatosInyectados(t *testing.T) {
	h := bootApp(t).Handler()
	body := do(t, h, http.MethodGet, "/casos-prueba", nil).Body.String()
	data := dataBlock(t, body)
	if strings.ContainsAny(data, "<>&") {
		t.Error("the injected JSON contains an unescaped <, > or & (it could close the script element)")
	}
	m, err := casos.Parse([]byte(data))
	if err != nil {
		t.Fatalf("the injected JSON is not a valid matrix: %v", err)
	}
	raw, _ := assets.FS.ReadFile(assets.CasosDataPath)
	want, err := casos.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	if len(m.Casos) != len(want.Casos) || len(m.Casos) < 30 {
		t.Fatalf("injected %d cases, embedded file has %d (want at least 30)", len(m.Casos), len(want.Casos))
	}
	var compact, original any
	if err := json.Unmarshal([]byte(data), &compact); err != nil {
		t.Fatal(err)
	}
	if err := json.Unmarshal(raw, &original); err != nil {
		t.Fatal(err)
	}
	a, _ := json.Marshal(compact)
	b, _ := json.Marshal(original)
	if string(a) != string(b) {
		t.Error("the injected JSON differs from ui/casos_prueba.json")
	}
	for _, c := range want.Casos {
		if !strings.Contains(data, `"`+c.ID+`"`) {
			t.Errorf("the page lacks %s", c.ID)
		}
	}
}

// A hostile string in the matrix cannot break out of the script element.
func TestUI_CasosPruebaEscapaElJSON(t *testing.T) {
	hostile := "</script><img src=x onerror=alert(1)> & fin"
	fsys := mutatedFile(t, assets.CasosDataPath, func(b []byte) []byte {
		return []byte(strings.Replace(string(b), "«servicio» no está en la stoplist y se conserva.", hostile, 1))
	})
	a := bootWith(t, web.Options{Version: "test", Assets: fsys})
	body := do(t, a.Handler(), http.MethodGet, "/casos-prueba", nil).Body.String()
	if strings.Contains(body, "<img src=x") || strings.Contains(body, "</script><img") {
		t.Error("the hostile string reached the page unescaped")
	}
	if !strings.Contains(dataBlock(t, body), `\u003c/script\u003e`) {
		t.Error("the hostile string is not escaped as \\u003c inside the data block")
	}
}

// The matrix page is an exact path and read-only, and it adds no API route.
func TestUI_CasosPruebaRutaExactaYSinNuevaRutaAPI(t *testing.T) {
	h := bootApp(t).Handler()
	if rec := do(t, h, http.MethodPost, "/casos-prueba", nil); rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST /casos-prueba status %d, want 405", rec.Code)
	}
	for _, p := range []string{"/casos-prueba/", "/casos-prueba/x", "/api/casos", "/api/casos-prueba"} {
		if rec := do(t, h, http.MethodGet, p, nil); rec.Code != http.StatusNotFound {
			t.Errorf("GET %s status %d, want 404", p, rec.Code)
		}
	}
}

// A missing or corrupt matrix stops the start-up instead of serving an empty page.
func TestBoot_CasosPruebaAusenteOCorrupto(t *testing.T) {
	for name, edit := range map[string]func([]byte) []byte{
		"not JSON": func([]byte) []byte { return []byte("no es json") },
		"empty":    func([]byte) []byte { return nil },
	} {
		if a, err := web.Boot(web.Options{Version: "test", Assets: mutatedFile(t, assets.CasosDataPath, edit)}); err == nil {
			a.Close()
			t.Errorf("%s: Boot accepted it", name)
		}
	}
	if a, err := web.Boot(web.Options{Version: "test", Assets: mutatedFile(t, assets.CasosPath, func([]byte) []byte { return nil })}); err == nil {
		a.Close()
		t.Error("an empty casos.html was accepted")
	}
}

// The summary of the page is short: the technical note stays in docs/CASOS_PRUEBA.md, and the
// verification line shows the Spanish date without any commit hash.
func TestUI_CasosPruebaResumenSinNotaNiCommit(t *testing.T) {
	body := do(t, bootApp(t).Handler(), http.MethodGet, "/casos-prueba", nil).Body.String()
	page := strings.Replace(body, dataBlock(t, body), "", 1)
	for _, bad := range []string{`id="nota"`, "matriz.nota", "commit", ".commit"} {
		if strings.Contains(page, bad) {
			t.Errorf("casos.html contains %q", bad)
		}
	}
	for _, want := range []string{`"Última verificación: "`, "octubre", "diciembre", " de "} {
		if !strings.Contains(page, want) {
			t.Errorf("casos.html lacks %q (Spanish date in the verification line)", want)
		}
	}
	if strings.Contains(dataBlock(t, body), "commit") {
		t.Error("the injected matrix carries a commit")
	}
}
