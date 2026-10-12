// Package casos is the single source of the test-case matrix of the Trabajo Escrito's Anexo II
// (test-case ID -> Go test -> measured result). The data is internal/assets/ui/casos_prueba.json;
// this package parses and validates it, renders docs/CASOS_PRUEBA.md and checks the two-way
// traceability against the *_test.go files (TE §6.3). The same JSON is embedded in the web page
// /casos-prueba, so the Markdown and the page cannot disagree.
package casos

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

// Estados valid for Caso.Estado.
const (
	EstadoPasa      = "pasa"
	EstadoPendiente = "pendiente"
)

// Matriz is the whole document.
type Matriz struct {
	Titulo     string     `json:"titulo"`
	Nota       string     `json:"nota"`
	Verificado Verificado `json:"verificado"`
	Niveles    []Nivel    `json:"niveles"`
	Casos      []Caso     `json:"casos"`
}

// Verificado records the last run of make verify that backs the "pasa" states. It holds no commit
// hash: hashes differ between the local repository and the GitHub export.
type Verificado struct {
	Fecha   string `json:"fecha"`
	Comando string `json:"comando"`
}

// Nivel is one group of the matrix (TE §6.3.2 to §6.3.8).
type Nivel struct {
	Nivel   int    `json:"nivel"`
	Nombre  string `json:"nombre"`
	Seccion string `json:"seccion"`
}

// Caso is one row. Extension marks a case of the project that is not in the Anexo II.
type Caso struct {
	ID         string   `json:"id"`
	IDAnexo    string   `json:"id_anexo"`
	Nivel      int      `json:"nivel"`
	Extension  bool     `json:"extension,omitempty"`
	Componente string   `json:"componente"`
	Tipo       string   `json:"tipo"`
	Entrada    string   `json:"entrada"`
	Criterio   string   `json:"criterio"`
	Pruebas    []Prueba `json:"pruebas"`
	Resultado  string   `json:"resultado"`
	Estado     string   `json:"estado"`
}

// Prueba is a Go test function and the file (relative to the module root) that defines it.
type Prueba struct {
	Funcion string `json:"funcion"`
	Archivo string `json:"archivo"`
}

// Parse reads and validates the JSON document: unknown fields, duplicate IDs, undefined levels,
// empty fields and states other than "pasa" or "pendiente" are errors.
func Parse(raw []byte) (*Matriz, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.DisallowUnknownFields()
	var m Matriz
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("casos: %w", err)
	}
	if err := m.validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Matriz) validate() error {
	if m.Titulo == "" || m.Verificado.Fecha == "" || m.Verificado.Comando == "" {
		return errors.New("casos: titulo and verificado (fecha, comando) are required")
	}
	levels := map[int]bool{}
	for _, n := range m.Niveles {
		if n.Nombre == "" || n.Seccion == "" || levels[n.Nivel] {
			return fmt.Errorf("casos: invalid or duplicate level %d", n.Nivel)
		}
		levels[n.Nivel] = true
	}
	seen := map[string]bool{}
	for _, c := range m.Casos {
		if c.ID == "" || seen[c.ID] {
			return fmt.Errorf("casos: empty or duplicate id %q", c.ID)
		}
		seen[c.ID] = true
		if !levels[c.Nivel] {
			return fmt.Errorf("casos: %s uses undefined level %d", c.ID, c.Nivel)
		}
		if c.Componente == "" || c.Tipo == "" || c.Entrada == "" || c.Criterio == "" || c.Resultado == "" {
			return fmt.Errorf("casos: %s has an empty field", c.ID)
		}
		if c.IDAnexo == "" && !c.Extension {
			return fmt.Errorf("casos: %s has no id_anexo and is not marked as extension", c.ID)
		}
		if c.Estado != EstadoPasa && c.Estado != EstadoPendiente {
			return fmt.Errorf("casos: %s has estado %q, want pasa or pendiente", c.ID, c.Estado)
		}
		if len(c.Pruebas) == 0 {
			return fmt.Errorf("casos: %s lists no Go tests", c.ID)
		}
		for _, p := range c.Pruebas {
			if p.Funcion == "" || p.Archivo == "" {
				return fmt.Errorf("casos: %s has a test without funcion or archivo", c.ID)
			}
		}
	}
	return nil
}

// Resumen counts the cases of the Anexo II, the project extensions and the cases that pass.
func (m *Matriz) Resumen() (anexo, extensiones, pasan int) {
	for _, c := range m.Casos {
		if c.Extension {
			extensiones++
		} else {
			anexo++
		}
		if c.Estado == EstadoPasa {
			pasan++
		}
	}
	return anexo, extensiones, pasan
}

// Markdown renders docs/CASOS_PRUEBA.md. The output is a pure function of the matrix.
func (m *Matriz) Markdown() string {
	var b strings.Builder
	anexo, ext, pasan := m.Resumen()
	fmt.Fprintf(&b, "# %s\n\n", m.Titulo)
	b.WriteString("> Generado por `make casos` desde `internal/assets/ui/casos_prueba.json`; no se edita a mano. La misma fuente alimenta la página `/casos-prueba` de la aplicación.\n\n")
	fmt.Fprintf(&b, "- Verificación: `%s`, %s.\n", m.Verificado.Comando, FechaES(m.Verificado.Fecha))
	fmt.Fprintf(&b, "- Resumen: %d casos del Anexo II y %d extensión del proyecto; %d de %d pasan.\n\n", anexo, ext, pasan, anexo+ext)
	fmt.Fprintf(&b, "%s\n", m.Nota)
	for _, n := range m.Niveles {
		var rows []Caso
		for _, c := range m.Casos {
			if c.Nivel == n.Nivel {
				rows = append(rows, c)
			}
		}
		if len(rows) == 0 {
			continue
		}
		fmt.Fprintf(&b, "\n## Nivel %d: %s (TE %s)\n\n", n.Nivel, n.Nombre, n.Seccion)
		b.WriteString("| ID | Componente | Tipo | Entrada | Criterio de aceptación | Pruebas Go | Resultado medido | Estado |\n")
		b.WriteString("|---|---|---|---|---|---|---|---|\n")
		for _, c := range rows {
			id := "**" + c.ID + "**"
			if c.IDAnexo != "" {
				id += " (" + c.IDAnexo + ")"
			} else if c.Extension {
				id += " (extensión, no está en el Anexo II)"
			}
			var tests []string
			for _, p := range c.Pruebas {
				tests = append(tests, "`"+p.Funcion+"` en `"+p.Archivo+"`")
			}
			fmt.Fprintf(&b, "| %s | %s | %s | %s | %s | %s | %s | %s |\n",
				id, cell(c.Componente), cell(c.Tipo), cell(c.Entrada), cell(c.Criterio),
				strings.Join(tests, "<br>"), cell(c.Resultado), cell(c.Estado))
		}
	}
	return b.String()
}

// cell makes s safe inside a Markdown table cell.
func cell(s string) string {
	s = strings.ReplaceAll(s, "|", `\|`)
	return strings.ReplaceAll(s, "\n", " ")
}

var meses = [...]string{"enero", "febrero", "marzo", "abril", "mayo", "junio", "julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

// FechaES renders an ISO date (2026-10-12) in Spanish ("12 de octubre de 2026"). Anything that is
// not a YYYY-MM-DD date is returned unchanged.
func FechaES(iso string) string {
	var y, m, d int
	if n, err := fmt.Sscanf(iso, "%4d-%2d-%2d", &y, &m, &d); err != nil || n != 3 || len(iso) != 10 || m < 1 || m > 12 || d < 1 || d > 31 {
		return iso
	}
	return fmt.Sprintf("%d de %s de %d", d, meses[m-1], y)
}
