package corpus

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// ambiguousFile is data/corpus_ambiguos.json: the ids of the comments whose base
// sentence, as a bag of content tokens, also occurs in another category (K2). It lives
// outside corpus.db so the schema stays the literal Anexo I.
type ambiguousFile struct {
	Criterio string `json:"criterio"`
	Corpus   string `json:"corpus_sha256"`
	Total    int    `json:"total"`
	IDs      []int  `json:"ids"`
}

// writeAmbiguous writes the subset; ids are 1-based row positions, which are the
// comentarios ids because Persist inserts in generation order.
func writeAmbiguous(path string, cs []Comment, sha string) error {
	cats := map[string]uint8{}
	keys := make([]string, len(cs))
	for i, c := range cs {
		keys[i] = bagKey(c.Base)
		cats[keys[i]] |= 1 << (c.CategoriaID - 1)
	}
	f := ambiguousFile{
		Criterio: "K2: la oración base limpia, como bolsa de tokens sin la stoplist v1, aparece también en otra categoría (PLAN.md §3)",
		Corpus:   sha,
		IDs:      []int{},
	}
	for i := range cs {
		if m := cats[keys[i]]; m&(m-1) != 0 {
			f.IDs = append(f.IDs, i+1)
		}
	}
	f.Total = len(f.IDs)
	b, err := json.Marshal(f)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}
