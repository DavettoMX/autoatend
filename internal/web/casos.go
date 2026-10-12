package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"

	"autoatend/internal/assets"
)

// loadCasosPage reads ui/casos.html and ui/casos_prueba.json and returns the page with the matrix
// injected in place of assets.CasosPlaceholder. The page therefore needs no API route (the five
// endpoints of TE §6.4.1 stay as they are). The JSON is compacted by json.Marshal, which applies
// HTML escaping (less-than, greater-than and ampersand become \\u escapes), so no string of the
// matrix can close the script element it lives in. Anything missing or invalid is a start-up
// error rather than an empty page.
func loadCasosPage(fsys fs.FS) ([]byte, error) {
	page, err := fs.ReadFile(fsys, assets.CasosPath)
	if err != nil {
		return nil, fmt.Errorf("artefacto embebido %s: %w", assets.CasosPath, err)
	}
	raw, err := fs.ReadFile(fsys, assets.CasosDataPath)
	if err != nil {
		return nil, fmt.Errorf("artefacto embebido %s: %w", assets.CasosDataPath, err)
	}
	safe, err := json.Marshal(json.RawMessage(raw))
	if err != nil {
		return nil, fmt.Errorf("artefacto embebido %s: JSON no valido: %w", assets.CasosDataPath, err)
	}
	if !bytes.Contains(page, []byte(assets.CasosPlaceholder)) {
		return nil, fmt.Errorf("artefacto embebido %s: falta la marca %s", assets.CasosPath, assets.CasosPlaceholder)
	}
	return bytes.Replace(page, []byte(assets.CasosPlaceholder), safe, 1), nil
}
