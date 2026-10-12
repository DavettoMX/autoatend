package assets_test

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"autoatend/internal/assets"
)

// The embedded files are byte-for-byte copies of data/ (make assets): they are never
// edited by hand, so they cannot drift from the artifacts that the rest of the
// repository versions and tests.
func TestEmbeddedDataEqualsDataDir(t *testing.T) {
	for _, name := range assets.DataFiles {
		got, err := assets.FS.ReadFile("data/" + name)
		if err != nil {
			t.Fatalf("embedded data/%s: %v", name, err)
		}
		want, err := os.ReadFile(filepath.Join("..", "..", "data", name))
		if err != nil {
			t.Fatal(err)
		}
		if len(got) == 0 || !bytes.Equal(got, want) {
			t.Errorf("data/%s: embedded copy (%d bytes) differs from ../../data/%s (%d bytes); run make assets", name, len(got), name, len(want))
		}
	}
}

func TestEmbeddedDataFilesList(t *testing.T) {
	want := []string{"corpus.db", "resultados.db", "modelo.json", "modelo_alternativo.json", "gold5.json"}
	if len(assets.DataFiles) != len(want) {
		t.Fatalf("DataFiles = %v, want %v", assets.DataFiles, want)
	}
	for i, n := range want {
		if assets.DataFiles[i] != n {
			t.Fatalf("DataFiles = %v, want %v", assets.DataFiles, want)
		}
	}
}

func TestEmbeddedUIAndJSON(t *testing.T) {
	ui, err := assets.FS.ReadFile(assets.UIPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(ui, []byte("<!DOCTYPE html>")) {
		t.Error("ui/index.html is not an HTML document")
	}
	for _, n := range []string{"modelo.json", "modelo_alternativo.json", "gold5.json"} {
		raw, _ := assets.FS.ReadFile("data/" + n)
		if !json.Valid(raw) {
			t.Errorf("data/%s is not valid JSON", n)
		}
	}
}

// The test-case matrix page and its JSON source are embedded; the page carries the placeholder that
// the server replaces with the JSON, and the JSON is a valid matrix.
func TestEmbeddedCasosDePrueba(t *testing.T) {
	page, err := assets.FS.ReadFile(assets.CasosPath)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(page, []byte("<!DOCTYPE html>")) || !bytes.Contains(page, []byte(assets.CasosPlaceholder)) {
		t.Errorf("ui/casos.html is not an HTML document with the %s placeholder", assets.CasosPlaceholder)
	}
	raw, err := assets.FS.ReadFile(assets.CasosDataPath)
	if err != nil {
		t.Fatal(err)
	}
	if !json.Valid(raw) {
		t.Error("ui/casos_prueba.json is not valid JSON")
	}
}
