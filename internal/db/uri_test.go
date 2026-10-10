package db

import (
	"net/url"
	"testing"
)

// A Windows absolute path (C:\...) must become file:///C:/..., never file://C:%5C...:
// SQLite reads everything between "file://" and the next "/" as the URI authority
// and rejects it ("invalid uri authority").
func TestFileURIWindowsDrivePath(t *testing.T) {
	got := fileURI(`C:/Users/exam/corpus.db`, []string{"foreign_keys(1)"}, url.Values{"mode": {"ro"}})
	want := "file:///C:/Users/exam/corpus.db?_pragma=foreign_keys%281%29&mode=ro"
	if got != want {
		t.Fatalf("fileURI = %q, want %q", got, want)
	}
}

func TestFileURIUnixPathUnchanged(t *testing.T) {
	if got, want := fileURI("/tmp/corpus.db", nil, nil), "file:///tmp/corpus.db"; got != want {
		t.Fatalf("fileURI = %q, want %q", got, want)
	}
}
