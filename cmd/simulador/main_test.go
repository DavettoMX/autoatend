package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"autoatend/internal/db/build"
	"autoatend/internal/testutil"
)

func tinyCorpus(t *testing.T, dir string) string {
	t.Helper()
	p := filepath.Join(dir, "corpus.db")
	d, err := build.CreateCorpus(p)
	if err != nil {
		t.Fatal(err)
	}
	defer d.Close()
	if err := build.InsertCategorias(d); err != nil {
		t.Fatal(err)
	}
	for i, txt := range []string{"La factura llegó mal", "Cobraron la factura", "Sin señal", "Gracias por todo", "Excelente servicio"} {
		testutil.MustExec(t, d, `INSERT INTO comentarios(fecha, texto, categoria_id) VALUES (?, ?, ?)`,
			"2025-02-0"+string(rune('1'+i)), txt, i%5+1)
	}
	if err := build.Finalize(d); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestRun_WritesTheCatalogueAndReportsOneLine(t *testing.T) {
	dir := t.TempDir()
	corpus := tinyCorpus(t, dir)
	out, df := filepath.Join(dir, "r.db"), filepath.Join(dir, "df.json")
	var stdout, stderr bytes.Buffer
	code := run([]string{"-corpus", corpus, "-out", out, "-df", df, "-seed", "10433"}, &stdout, &stderr)
	if code != 0 {
		t.Fatalf("exit %d, stderr %q", code, stderr.String())
	}
	line := strings.TrimSpace(stdout.String())
	if strings.Count(line, "\n") != 0 || !strings.HasPrefix(line, "simulador: wrote "+out+" (") ||
		!strings.Contains(line, " lexemas, sha256 ") {
		t.Errorf("unexpected output line %q", line)
	}
	for _, p := range []string{out, df} {
		if _, err := os.Stat(p); err != nil {
			t.Errorf("%s not written: %v", p, err)
		}
	}
}

func TestRun_FailureHasPrefixAndExitCode1(t *testing.T) {
	dir := t.TempDir()
	corpus := tinyCorpus(t, dir)
	out, df := filepath.Join(dir, "r.db"), filepath.Join(dir, "df.json")
	args := []string{"-corpus", corpus, "-out", out, "-df", df}
	var stdout, stderr bytes.Buffer
	if code := run(args, &stdout, &stderr); code != 0 {
		t.Fatalf("first run exit %d: %s", code, stderr.String())
	}
	stdout.Reset()
	stderr.Reset()
	if code := run(args, &stdout, &stderr); code != 1 {
		t.Errorf("second run (files exist) exit %d, want 1", code)
	}
	if !strings.HasPrefix(stderr.String(), "simulador:") {
		t.Errorf("stderr %q lacks the simulador: prefix", stderr.String())
	}
	if stdout.Len() != 0 {
		t.Errorf("stdout %q on failure", stdout.String())
	}
}

func TestRun_BadFlagExitsNonZero(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"-nope"}, &stdout, &stderr); code == 0 {
		t.Error("an unknown flag must not exit 0")
	}
}
