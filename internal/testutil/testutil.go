// Package testutil holds DB-agnostic helpers shared by test files that exercise
// SQLite. It is imported only from _test.go files.
package testutil

import (
	"database/sql"
	"strings"
	"testing"
)

// MustExec runs q with args on d and fails the test on error.
func MustExec(t *testing.T, d *sql.DB, q string, args ...any) {
	t.Helper()
	if _, err := d.Exec(q, args...); err != nil {
		t.Fatalf("exec %q: %v", q, err)
	}
}

// ExpectErr runs q with args on d and fails the test unless it errors with a
// message containing wantSubstr.
func ExpectErr(t *testing.T, d *sql.DB, wantSubstr, q string, args ...any) {
	t.Helper()
	_, err := d.Exec(q, args...)
	if err == nil {
		t.Fatalf("exec %q: expected error containing %q, got nil", q, wantSubstr)
	}
	if !strings.Contains(err.Error(), wantSubstr) {
		t.Fatalf("exec %q: error %q does not contain %q", q, err, wantSubstr)
	}
}

// Scalar runs q, which must return one integer, and fails the test unless it equals want.
func Scalar(t *testing.T, d *sql.DB, want int, q string) {
	t.Helper()
	var got int
	if err := d.QueryRow(q).Scan(&got); err != nil {
		t.Fatalf("%s: %v", q, err)
	}
	if got != want {
		t.Errorf("%s = %d, want %d", q, got, want)
	}
}
