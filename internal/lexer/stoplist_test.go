package lexer

import (
	"encoding/json"
	"os"
	"sort"
	"testing"
)

// policy is the ADR-003 stoplist policy (PLAN.md §4) written out independently
// of stoplist/v1.json, class by class.
var policy = map[string][]string{
	"articulo":     {"el", "la", "los", "las", "un", "una", "unos", "unas", "lo", "al", "del"},
	"preposicion":  {"de", "en", "con", "por", "a", "ante", "tras", "durante", "mediante", "según", "sin", "hasta", "desde", "contra", "hacia", "segun", "x"},
	"conjuncion":   {"y", "e", "o", "u", "ni", "que", "pero", "sino", "aunque", "porque", "pues", "q"},
	"interjeccion": {"ah", "eh", "ay", "uf", "olé", "ole"},
}

func TestCPU11_StoplistFiltersDe(t *testing.T) {
	class, ok := Lookup("de")
	if !ok || class != "preposicion" {
		t.Fatalf(`Lookup("de") = %q, %v; want preposicion, true`, class, ok)
	}
}

func TestCPU12_StoplistKeepsBajo(t *testing.T) {
	if class, ok := Lookup("bajo"); ok {
		t.Fatalf(`Lookup("bajo") = %q, true; the homograph must be retained`, class)
	}
}

func TestCPU13_StoplistKeepsServicio(t *testing.T) {
	if class, ok := Lookup("servicio"); ok {
		t.Fatalf(`Lookup("servicio") = %q, true; a content word must be retained`, class)
	}
}

func TestCPU14_StoplistPolicy(t *testing.T) {
	// Every form of the policy is filtered, with its class.
	for class, forms := range policy {
		for _, f := range forms {
			got, ok := Lookup(f)
			if !ok || got != class {
				t.Errorf("Lookup(%q) = %q, %v; want %q, true", f, got, ok, class)
			}
		}
	}
	// Surface variants of filtered forms stay in the class of the original.
	for variant, original := range map[string]string{"segun": "según", "ole": "olé"} {
		cv, _ := Lookup(variant)
		co, _ := Lookup(original)
		if cv == "" || cv != co {
			t.Errorf("%q is in class %q but %q is in %q", variant, cv, original, co)
		}
	}
	if c, _ := Lookup("q"); c != "conjuncion" {
		t.Errorf(`"q" class = %q, want conjuncion`, c)
	}
	if c, _ := Lookup("x"); c != "preposicion" {
		t.Errorf(`"x" class = %q, want preposicion`, c)
	}
	if c, _ := Lookup("hacia"); c != "preposicion" {
		t.Errorf(`"hacia" class = %q, want preposicion`, c)
	}
	// Homographs of the TE and words outside the four classes are retained.
	for _, f := range []string{"bajo", "para", "sobre", "entre", "como", "más", "solo", "aun",
		"él", "qué", "sí", "té", "se", "no", "este", "eso"} {
		if c, ok := Lookup(f); ok {
			t.Errorf("Lookup(%q) = %q, true; must be retained", f, c)
		}
	}
}

// The embedded file has exactly the policy: nothing more, nothing less.
func TestStoplistFileMatchesPolicy(t *testing.T) {
	raw, err := os.ReadFile("stoplist/v1.json")
	if err != nil {
		t.Fatal(err)
	}
	var f struct {
		Version string            `json:"version"`
		Entries map[string]string `json:"entries"`
	}
	if err := json.Unmarshal(raw, &f); err != nil {
		t.Fatal(err)
	}
	if f.Version != "v1" || Version() != "v1" {
		t.Errorf("version = %q / %q, want v1", f.Version, Version())
	}
	want := map[string]string{}
	for class, forms := range policy {
		for _, form := range forms {
			want[form] = class
		}
	}
	var extra, missing, wrong []string
	for form, class := range f.Entries {
		wc, ok := want[form]
		switch {
		case !ok:
			extra = append(extra, form)
		case wc != class:
			wrong = append(wrong, form)
		}
	}
	for form := range want {
		if _, ok := f.Entries[form]; !ok {
			missing = append(missing, form)
		}
	}
	sort.Strings(extra)
	sort.Strings(missing)
	sort.Strings(wrong)
	if len(extra)+len(missing)+len(wrong) > 0 {
		t.Errorf("v1.json differs from the policy: extra %v, missing %v, wrong class %v", extra, missing, wrong)
	}
	if len(f.Entries) != 46 {
		t.Errorf("entries = %d, want 46", len(f.Entries))
	}
	// Every entry is already in normalized form, or Analyze could never match it.
	for form := range f.Entries {
		if Normalize(form) != form {
			t.Errorf("entry %q is not in normalized form (%q)", form, Normalize(form))
		}
	}
}
