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
		`/api/lexemas?orden=variantes`,
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

// GET /catalogo: a readable table of the catalogue (consigna b) built from /api/lexemas,
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
