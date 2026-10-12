package web_test

import (
	"context"
	"encoding/json"
	"net"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"autoatend/internal/web"
)

func runSelftest(t *testing.T, opts web.Options) (code int, stdout, stderr string, elapsed time.Duration) {
	t.Helper()
	opts.Selftest = true
	var out, errb syncBuffer
	start := time.Now()
	code = web.Run(context.Background(), opts, &out, &errb)
	return code, out.String(), errb.String(), time.Since(start)
}

// TE §6.3.9 / §6.4.3: good artifacts -> exit 0 in < 3 s, one diagnostic line per check.
func TestSelftestExit0(t *testing.T) {
	code, out, errs, elapsed := runSelftest(t, web.Options{Version: "test"})
	t.Logf("selftest: exit %d in %v\n%s", code, elapsed, out)
	if code != 0 {
		t.Fatalf("exit %d\nstdout:\n%s\nstderr:\n%s", code, out, errs)
	}
	if errs != "" {
		t.Errorf("stderr is not empty on success: %q", errs)
	}
	for _, want := range []string{"corpus.db", "integridad", "catalogo", "gold5", "puerto"} {
		if !strings.Contains(out, want) {
			t.Errorf("selftest output has no line about %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "FALLO") {
		t.Errorf("unexpected failure line:\n%s", out)
	}
	if !testing.Short() && elapsed > 3*time.Second {
		t.Errorf("selftest took %v, want < 3 s", elapsed)
	}
}

func TestSelftestModeloAlternativo(t *testing.T) {
	code, out, errs, _ := runSelftest(t, web.Options{Modelo: "alternativo"})
	if code != 0 {
		t.Fatalf("exit %d\n%s\n%s", code, out, errs)
	}
}

// A gold comment with the wrong expected label fails the check and the exit code is 1.
func TestSelftestGold5Corrupto(t *testing.T) {
	bad := mutatedAssets(t, "gold5.json", func(b []byte) []byte {
		var g []map[string]any
		if err := json.Unmarshal(b, &g); err != nil {
			t.Fatal(err)
		}
		g[0]["categoria_id"] = float64(int(g[0]["categoria_id"].(float64))%5 + 1) // a different label
		out, _ := json.Marshal(g)
		return out
	})
	code, out, errs, _ := runSelftest(t, web.Options{Assets: bad})
	if code != 1 {
		t.Fatalf("exit %d, want 1\nstdout:\n%s\nstderr:\n%s", code, out, errs)
	}
	if !strings.Contains(errs, "gold5") {
		t.Errorf("stderr does not name the failed check:\n%s", errs)
	}
	// The other checks still ran and passed.
	if !strings.Contains(out, "corpus.db") || !strings.Contains(out, "catalogo") {
		t.Errorf("remaining checks did not report:\n%s", out)
	}
}

func TestSelftestGold5ConLongitudIncorrecta(t *testing.T) {
	bad := mutatedAssets(t, "gold5.json", func(b []byte) []byte { return []byte(`[]`) })
	if code, _, errs, _ := runSelftest(t, web.Options{Assets: bad}); code != 1 || !strings.Contains(errs, "gold5") {
		t.Errorf("exit %d, stderr %q; want 1 naming gold5", code, errs)
	}
}

func TestSelftestCorpusConConteoIncorrecto(t *testing.T) {
	bad := mutatedAssets(t, "corpus.db", func(b []byte) []byte {
		return editSQLite(t, b, `DELETE FROM comentarios WHERE id = (SELECT MIN(id) FROM comentarios)`)
	})
	code, _, errs, _ := runSelftest(t, web.Options{Assets: bad})
	if code != 1 || !strings.Contains(errs, "72999") || !strings.Contains(errs, "73000") {
		t.Errorf("exit %d, stderr %q; want 1 with found 72999 and expected 73000", code, errs)
	}
}

func TestSelftestERInvalidaDelCatalogo(t *testing.T) {
	bad := mutatedAssets(t, "resultados.db", func(b []byte) []byte {
		return editSQLite(t, b, `UPDATE lexemas SET expresion_regular = 'ab(c' WHERE id = 1`)
	})
	code, _, errs, _ := runSelftest(t, web.Options{Assets: bad})
	if code != 1 || !strings.Contains(errs, "catalogo") {
		t.Errorf("exit %d, stderr %q; want 1 naming the catalogue check", code, errs)
	}
}

// The catalogue ER must accept the variants it summarizes (I-1).
func TestSelftestERQueNoAceptaSusVariantes(t *testing.T) {
	bad := mutatedAssets(t, "resultados.db", func(b []byte) []byte {
		return editSQLite(t, b, `UPDATE lexemas SET expresion_regular = 'zzz' WHERE id = (SELECT lexema_id FROM variantes LIMIT 1)`)
	})
	code, _, errs, _ := runSelftest(t, web.Options{Assets: bad})
	if code != 1 || !strings.Contains(errs, "catalogo") {
		t.Errorf("exit %d, stderr %q; want 1 naming the catalogue check", code, errs)
	}
}

func TestSelftestArranqueFallido(t *testing.T) {
	code, _, errs, _ := runSelftest(t, web.Options{Assets: fstest.MapFS{}})
	if code != 1 || errs == "" {
		t.Errorf("exit %d, stderr %q; want 1 with a diagnostic", code, errs)
	}
}

// With the whole port range busy the port check fails.
func TestSelftestPuertoOcupado(t *testing.T) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer l.Close()
	p := l.Addr().(*net.TCPAddr).Port
	code, _, errs, _ := runSelftest(t, web.Options{FirstPort: p, LastPort: p})
	if code != 1 || !strings.Contains(errs, "puerto") {
		t.Errorf("exit %d, stderr %q; want 1 naming the port check", code, errs)
	}
	// With a free port next to the busy one the fallback makes the check pass.
	code, out, errs, _ := runSelftest(t, web.Options{FirstPort: p, LastPort: p + 9})
	if code != 0 {
		t.Errorf("fallback: exit %d\n%s\n%s", code, out, errs)
	}
}
