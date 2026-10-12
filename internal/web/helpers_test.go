package web_test

import (
	"bytes"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"autoatend/internal/assets"
	"autoatend/internal/db"
	"autoatend/internal/web"
)

// bootApp starts the application on the embedded artifacts and removes its
// temporary directory when the test ends.
func bootApp(t *testing.T) *web.App {
	t.Helper()
	return bootWith(t, web.Options{Version: "test"})
}

func bootWith(t *testing.T, opts web.Options) *web.App {
	t.Helper()
	a, err := web.Boot(opts)
	if err != nil {
		t.Fatalf("Boot: %v", err)
	}
	t.Cleanup(func() { a.Close() })
	return a
}

// do sends a request to the handler and returns the recorded response.
func do(t *testing.T, h http.Handler, method, target string, body io.Reader) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, target, body)
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec
}

func postClasificar(t *testing.T, h http.Handler, texto string) *httptest.ResponseRecorder {
	t.Helper()
	b, err := json.Marshal(map[string]string{"texto": texto})
	if err != nil {
		t.Fatal(err)
	}
	return do(t, h, http.MethodPost, "/api/clasificar", bytes.NewReader(b))
}

// decode parses the recorded JSON body into v.
func decode(t *testing.T, rec *httptest.ResponseRecorder, v any) {
	t.Helper()
	if err := json.Unmarshal(rec.Body.Bytes(), v); err != nil {
		t.Fatalf("body is not the expected JSON (%v): %.300s", err, rec.Body.String())
	}
}

func sha(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	h := sha256.Sum256(b)
	return hex.EncodeToString(h[:])
}

// goldEntry is one record of gold5.json.
type goldEntry struct {
	ID          int    `json:"id"`
	Texto       string `json:"texto"`
	CategoriaID int    `json:"categoria_id"`
	Categoria   string `json:"categoria"`
}

func loadGold(t *testing.T) []goldEntry {
	t.Helper()
	raw, err := assets.FS.ReadFile("data/gold5.json")
	if err != nil {
		t.Fatal(err)
	}
	var g []goldEntry
	if err := json.Unmarshal(raw, &g); err != nil {
		t.Fatal(err)
	}
	if len(g) != 5 {
		t.Fatalf("gold5.json has %d entries, want 5", len(g))
	}
	return g
}

// mutatedAssets returns a copy of the embedded artifacts in which the data file
// name is replaced by edit(contents).
func mutatedAssets(t *testing.T, name string, edit func([]byte) []byte) fs.FS {
	t.Helper()
	m := fstest.MapFS{}
	for _, n := range append([]string{assets.UIPath}, prefixed("data/", assets.DataFiles)...) {
		b, err := assets.FS.ReadFile(n)
		if err != nil {
			t.Fatal(err)
		}
		if n == "data/"+name {
			b = edit(b)
		}
		m[n] = &fstest.MapFile{Data: b}
	}
	return m
}

func prefixed(p string, names []string) []string {
	out := make([]string, len(names))
	for i, n := range names {
		out[i] = p + n
	}
	return out
}

// editSQLite applies stmts (no arguments) to a copy of the SQLite image b and
// returns the modified image.
func editSQLite(t *testing.T, b []byte, stmts ...string) []byte {
	t.Helper()
	path := filepath.Join(t.TempDir(), "edit.db")
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
	uri, err := db.FileURI(path, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	d, err := sql.Open("sqlite", uri)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range stmts {
		if _, err := d.Exec(s); err != nil {
			t.Fatalf("%s: %v", s, err)
		}
	}
	d.Close()
	out, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return out
}

// syncBuffer is a bytes.Buffer safe for a writer goroutine and a polling reader.
type syncBuffer struct {
	mu sync.Mutex
	b  strings.Builder
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}
