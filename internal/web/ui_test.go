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
