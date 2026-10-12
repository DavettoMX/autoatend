package main

import (
	"bytes"
	"context"
	"net"
	"os/exec"
	"strings"
	"sync"
	"testing"
	"time"
)

type syncBuf struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

func TestSelftestFlagExit0(t *testing.T) {
	for _, flagForm := range []string{"--selftest", "-selftest"} {
		var out, errb bytes.Buffer
		if code := run(context.Background(), []string{flagForm}, &out, &errb); code != 0 {
			t.Fatalf("%s: exit %d\n%s\n%s", flagForm, code, out.String(), errb.String())
		}
	}
	var out, errb bytes.Buffer
	if code := run(context.Background(), []string{"--selftest", "-modelo", "alternativo"}, &out, &errb); code != 0 {
		t.Fatalf("alternativo: exit %d\n%s\n%s", code, out.String(), errb.String())
	}
}

func TestFlagsInvalidos(t *testing.T) {
	var out, errb bytes.Buffer
	if code := run(context.Background(), []string{"-modelo", "otro"}, &out, &errb); code != 2 || !strings.Contains(errb.String(), "otro") {
		t.Errorf("-modelo otro: exit %d, stderr %q; want 2 naming the value", code, errb.String())
	}
	errb.Reset()
	if code := run(context.Background(), []string{"-no-existe"}, &out, &errb); code != 2 {
		t.Errorf("unknown flag: exit %d, want 2", code)
	}
}

// CPA05 on the real ports: with 8471 busy the service starts on the next one and the banner shows it.
func TestCPA05_Puerto8471Ocupado(t *testing.T) {
	busy, err := net.Listen("tcp", "127.0.0.1:8471")
	if err != nil {
		t.Skipf("8471 is not bindable here (already in use): %v", err)
	}
	defer busy.Close()

	ctx, cancel := context.WithCancel(context.Background())
	var out, errb syncBuf
	done := make(chan int, 1)
	go func() { done <- run(ctx, []string{"-no-browser"}, &out, &errb) }()

	deadline := time.Now().Add(10 * time.Second)
	for !strings.Contains(out.String(), "Servicio listo en") && time.Now().Before(deadline) {
		time.Sleep(10 * time.Millisecond)
	}
	banner := out.String()
	if strings.Contains(banner, "127.0.0.1:8471") || !strings.Contains(banner, "http://127.0.0.1:8472") {
		t.Errorf("banner = %q, want the URL on 8472", banner)
	}
	cancel()
	if code := <-done; code != 0 {
		t.Errorf("exit %d: %s", code, errb.String())
	}
}

// The runtime must not link the generation-only package (CLAUDE.md, reglas duras).
func TestNoDependeDeDBBuild(t *testing.T) {
	out, err := exec.Command("go", "list", "-deps", ".").CombinedOutput()
	if err != nil {
		t.Fatalf("go list -deps: %v\n%s", err, out)
	}
	deps := string(out)
	if !strings.Contains(deps, "autoatend/internal/web") {
		t.Fatalf("go list -deps does not list the web package:\n%s", deps)
	}
	for _, forbidden := range []string{"autoatend/internal/db/build", "autoatend/internal/sweep", "autoatend/internal/corpus"} {
		if strings.Contains(deps, forbidden+"\n") {
			t.Errorf("cmd/autoatend depends on %s", forbidden)
		}
	}
}
