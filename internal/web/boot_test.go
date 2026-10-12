package web_test

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"autoatend/internal/web"
)

// Boot extracts the embedded databases into a temporary directory of its own and
// Close removes it (PLAN.md §2, Arranque).
func TestBootExtraeAlTemporalYCloseLoBorra(t *testing.T) {
	a, err := web.Boot(web.Options{Version: "test"})
	if err != nil {
		t.Fatal(err)
	}
	dir := a.Dir
	if !strings.HasPrefix(dir, os.TempDir()) {
		t.Errorf("working dir %q is not under the user temp dir %q", dir, os.TempDir())
	}
	for _, n := range []string{"corpus.db", "resultados.db", "modelo.json", "modelo_alternativo.json", "gold5.json"} {
		if st, err := os.Stat(filepath.Join(dir, n)); err != nil || st.Size() == 0 {
			t.Errorf("%s not extracted: %v", n, err)
		}
	}
	if a.Model == nil || a.Model.Familia != "rls" {
		t.Errorf("default model = %+v, want the principal (rls)", a.Model)
	}
	if err := a.Close(); err != nil {
		t.Errorf("Close: %v", err)
	}
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("temp dir %s still exists after Close (err=%v)", dir, err)
	}
	if err := a.Close(); err != nil {
		t.Errorf("second Close: %v", err)
	}
}

func TestBootModeloAlternativo(t *testing.T) {
	a := bootWith(t, web.Options{Modelo: "alternativo"})
	if a.Model.Familia != "nbm" {
		t.Errorf("alternative model familia = %q, want nbm", a.Model.Familia)
	}
	p := bootWith(t, web.Options{Modelo: "principal"})
	if p.Model.Familia != "rls" {
		t.Errorf("principal model familia = %q, want rls", p.Model.Familia)
	}
	if _, err := web.Boot(web.Options{Modelo: "otro"}); err == nil || !strings.Contains(err.Error(), "otro") {
		t.Errorf("unknown -modelo: err = %v, want an error naming the value", err)
	}
}

// Any startup problem is fatal with a diagnostic, and leaves no temp dir behind.
func TestBootErroresSonFatalesYLimpian(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	t.Setenv("TEMP", tmp)
	t.Setenv("TMP", tmp)

	missing := fstest.MapFS{}
	if _, err := web.Boot(web.Options{Assets: missing}); err == nil {
		t.Error("empty assets: expected an error")
	}

	badModel := mutatedAssets(t, "modelo.json", func(b []byte) []byte { return []byte(`{"version": 99}`) })
	_, err := web.Boot(web.Options{Assets: badModel})
	if err == nil || !strings.Contains(err.Error(), "modelo.json") {
		t.Errorf("corrupted model: err = %v, want a diagnostic naming modelo.json", err)
	}

	notSQLite := mutatedAssets(t, "resultados.db", func(b []byte) []byte { return []byte("this is not a database file at all") })
	if _, err := web.Boot(web.Options{Assets: notSQLite}); err == nil {
		t.Error("corrupted resultados.db: expected an error")
	}

	if left, _ := filepath.Glob(filepath.Join(tmp, "*")); len(left) != 0 {
		t.Errorf("failed boots left files behind: %v", left)
	}
}

// TE §6.1.3: a categoria_id of resultados.db missing from corpus.categorias is fatal at startup.
func TestBootIntegridadCruzadaFallaConHuerfanos(t *testing.T) {
	orphan := mutatedAssets(t, "resultados.db", func(b []byte) []byte {
		return editSQLite(t, b, `INSERT INTO frecuencias(lexema_id, corrida_id, categoria_id, valor) VALUES (1, 1, 9, 1)`)
	})
	_, err := web.Boot(web.Options{Assets: orphan})
	if err == nil || !strings.Contains(err.Error(), "categoria_id") || !strings.Contains(err.Error(), "9") {
		t.Fatalf("err = %v, want a diagnostic naming the dangling categoria_id 9", err)
	}
}

func TestListenFallback(t *testing.T) {
	busy, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer busy.Close()
	first := busy.Addr().(*net.TCPAddr).Port

	ln, err := web.Listen(first, first+9)
	if err != nil {
		t.Fatalf("Listen: %v", err)
	}
	defer ln.Close()
	got := ln.Addr().(*net.TCPAddr)
	if got.Port <= first || got.Port > first+9 {
		t.Errorf("bound port %d, want in (%d, %d]", got.Port, first, first+9)
	}
	if !got.IP.IsLoopback() || got.IP.String() != "127.0.0.1" {
		t.Errorf("bound %v, want 127.0.0.1 only", got.IP)
	}
}

func TestListenSinPuertoLibre(t *testing.T) {
	var held []net.Listener
	defer func() {
		for _, l := range held {
			l.Close()
		}
	}()
	// Reserve a small consecutive range, then ask for exactly that range.
	var first int
	for attempt := 0; attempt < 50 && len(held) == 0; attempt++ {
		l, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		first = l.Addr().(*net.TCPAddr).Port
		held = append(held, l)
		for p := first + 1; p <= first+2; p++ {
			l2, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", itoa(p)))
			if err != nil {
				for _, h := range held {
					h.Close()
				}
				held = nil
				break
			}
			held = append(held, l2)
		}
	}
	if len(held) != 3 {
		t.Skip("could not reserve three consecutive ports")
	}
	if ln, err := web.Listen(first, first+2); err == nil {
		ln.Close()
		t.Fatal("expected an error when every port of the range is busy")
	} else if !strings.Contains(err.Error(), itoa(first)) {
		t.Errorf("err = %v, want the tried range in the message", err)
	}
}

func itoa(n int) string { return strconv.Itoa(n) }

func TestBanner(t *testing.T) {
	b := web.Banner("http://127.0.0.1:8472", "1.0.0")
	if !strings.Contains(b, "http://127.0.0.1:8472") || !strings.Contains(b, "AutoAtend") || !strings.Contains(b, "1.0.0") {
		t.Errorf("banner = %q", b)
	}
}

// CPA05: with the primary port busy the server starts on the next one and the banner shows its URL.
func TestCPA05_FallbackDePuerto(t *testing.T) {
	busy, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer busy.Close()
	first := busy.Addr().(*net.TCPAddr).Port

	ctx, cancel := context.WithCancel(context.Background())
	var out, errb syncBuffer
	done := make(chan int, 1)
	go func() {
		done <- web.Run(ctx, web.Options{Version: "test", NoBrowser: true, FirstPort: first, LastPort: first + 9}, &out, &errb)
	}()

	url := waitForURL(t, &out)
	if !strings.HasSuffix(url, ":"+itoa(first+1)) && !strings.HasSuffix(url, ":"+itoa(first+2)) {
		t.Errorf("banner URL %q: want the next free port after %d", url, first)
	}
	if strings.HasSuffix(url, ":"+itoa(first)) {
		t.Errorf("banner announces the busy port: %s", url)
	}
	resp, err := http.Get(url + "/api/health")
	if err != nil {
		t.Fatal(err)
	}
	var h map[string]any
	json.NewDecoder(resp.Body).Decode(&h)
	resp.Body.Close()
	if !strings.HasSuffix(url, ":"+itoa(int(h["puerto"].(float64)))) {
		t.Errorf("/api/health puerto = %v, banner URL = %s", h["puerto"], url)
	}

	cancel()
	select {
	case code := <-done:
		if code != 0 {
			t.Errorf("exit code %d after shutdown, stderr: %s", code, errb.String())
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run did not return after the context was cancelled")
	}
	if _, err := http.Get(url + "/api/health"); err == nil {
		t.Error("server still answers after shutdown")
	}
}

// Run removes its temp dir on exit (SIGINT/SIGTERM cancel the context in cmd/autoatend).
func TestRunLimpiaElTemporalAlSalir(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("TMPDIR", tmp)
	t.Setenv("TEMP", tmp)
	t.Setenv("TMP", tmp)

	ctx, cancel := context.WithCancel(context.Background())
	var out, errb syncBuffer
	done := make(chan int, 1)
	go func() {
		done <- web.Run(ctx, web.Options{NoBrowser: true, FirstPort: 18471, LastPort: 18480}, &out, &errb)
	}()
	waitForURL(t, &out)
	if entries, _ := os.ReadDir(tmp); len(entries) == 0 {
		t.Error("no working directory while running")
	}
	cancel()
	if code := <-done; code != 0 {
		t.Fatalf("exit %d: %s", code, errb.String())
	}
	if entries, _ := os.ReadDir(tmp); len(entries) != 0 {
		t.Errorf("temp dir not cleaned: %v", entries)
	}
}

func TestRunErrorDeArranqueEsFatal(t *testing.T) {
	var out, errb syncBuffer
	code := web.Run(context.Background(), web.Options{Assets: fstest.MapFS{}, NoBrowser: true}, &out, &errb)
	if code != 1 {
		t.Errorf("exit %d, want 1", code)
	}
	if !strings.Contains(errb.String(), "autoatend") || len(errb.String()) < 20 {
		t.Errorf("stderr = %q, want a startup diagnostic", errb.String())
	}
}

func TestRunAbreElNavegadorConLaURLEfectiva(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var out, errb syncBuffer
	opened := make(chan string, 1)
	done := make(chan int, 1)
	go func() {
		done <- web.Run(ctx, web.Options{FirstPort: 18481, LastPort: 18490,
			OpenBrowser: func(u string) error { opened <- u; return nil }}, &out, &errb)
	}()
	url := waitForURL(t, &out)
	select {
	case u := <-opened:
		if u != url {
			t.Errorf("browser opened %q, banner says %q", u, url)
		}
	case <-time.After(3 * time.Second):
		t.Error("browser was not opened")
	}
	cancel()
	<-done

	// -no-browser (NoBrowser) never calls the opener.
	ctx2, cancel2 := context.WithCancel(context.Background())
	var out2, err2 syncBuffer
	done2 := make(chan int, 1)
	go func() {
		done2 <- web.Run(ctx2, web.Options{NoBrowser: true, FirstPort: 18481, LastPort: 18490,
			OpenBrowser: func(u string) error { t.Error("opener called despite NoBrowser"); return errors.New("x") }}, &out2, &err2)
	}()
	waitForURL(t, &out2)
	cancel2()
	<-done2
}

func TestBrowserCommand(t *testing.T) {
	for _, c := range []struct {
		goos string
		want []string
	}{
		{"windows", []string{"rundll32", "url.dll,FileProtocolHandler", "http://127.0.0.1:8471"}},
		{"darwin", []string{"open", "http://127.0.0.1:8471"}},
		{"linux", []string{"xdg-open", "http://127.0.0.1:8471"}},
	} {
		got, ok := web.BrowserCommand(c.goos, "http://127.0.0.1:8471")
		if !ok || strings.Join(got, " ") != strings.Join(c.want, " ") {
			t.Errorf("%s: %v (%v), want %v", c.goos, got, ok, c.want)
		}
	}
	if _, ok := web.BrowserCommand("plan9", "http://x"); ok {
		t.Error("unknown OS must not claim a command")
	}
}

// waitForURL polls the banner until it prints the URL and returns it.
func waitForURL(t *testing.T, out *syncBuffer) string {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		for _, f := range strings.Fields(out.String()) {
			if strings.HasPrefix(f, "http://127.0.0.1:") {
				return strings.TrimRight(f, ".,)")
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("no banner URL after 10 s; output: %q", out.String())
	return ""
}
