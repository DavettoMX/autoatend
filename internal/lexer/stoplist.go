package lexer

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

// The stoplist is a versioned, closed table form -> class (ADR-003). Only
// functional words of the four closed classes (articles with their contractions,
// prepositions, conjunctions, interjections) are filtered. The homographs of
// the TE (bajo, para, sobre, entre, como, más, solo, aun) are retained, and so
// are se, no, este and eso. The lookup is by exact normalized surface form, with
// no accent folding: "segun", "ole", "q" and "x" are listed on their own.
//
//go:embed stoplist/v1.json
var stoplistJSON []byte

type stoplistFile struct {
	Version string            `json:"version"`
	Entries map[string]string `json:"entries"`
}

var stoplist = mustLoadStoplist(stoplistJSON)

func mustLoadStoplist(raw []byte) stoplistFile {
	var f stoplistFile
	if err := json.Unmarshal(raw, &f); err != nil {
		panic(fmt.Sprintf("lexer: stoplist/v1.json: %v", err))
	}
	return f
}

// Lookup reports the closed class of a normalized surface form, and whether
// the form is in the stoplist (that is, must be filtered).
func Lookup(form string) (class string, ok bool) {
	class, ok = stoplist.Entries[form]
	return class, ok
}

// Version is the version of the embedded stoplist, recorded in corridas.parametros.
func Version() string { return stoplist.Version }
