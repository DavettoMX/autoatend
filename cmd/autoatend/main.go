// Command autoatend is the AutoAtend web service (consigna e, TE §6.4): a single
// static binary with the databases and the model embedded. It only parses flags
// and signals; the logic lives in internal/web.
//
// Build with: go build -ldflags "-X main.version=1.0.0" ./cmd/autoatend (make assets first).
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"autoatend/internal/web"
)

// version is set at link time with -ldflags "-X main.version=...".
var version = "dev"

func main() {
	// SIGINT/SIGTERM cancel the context: the server shuts down and the temp dir is removed.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	code := run(ctx, os.Args[1:], os.Stdout, os.Stderr)
	stop()
	os.Exit(code)
}

func run(ctx context.Context, args []string, stdout, stderr io.Writer) int {
	fs := flag.NewFlagSet("autoatend", flag.ContinueOnError)
	fs.SetOutput(stderr)
	modelo := fs.String("modelo", "principal", "modelo a servir: principal (regresion logistica) o alternativo (Naive Bayes)")
	noBrowser := fs.Bool("no-browser", false, "no abrir el navegador al iniciar")
	selftest := fs.Bool("selftest", false, "ejecutar las verificaciones headless y salir con 0 (exito) o 1 (fallo)")
	if err := fs.Parse(args); err != nil {
		return 2
	}
	if *modelo != "principal" && *modelo != "alternativo" {
		fmt.Fprintf(stderr, "autoatend: valor %q no valido para -modelo (use principal o alternativo)\n", *modelo)
		return 2
	}
	return web.Run(ctx, web.Options{
		Modelo:    *modelo,
		Version:   version,
		NoBrowser: *noBrowser,
		Selftest:  *selftest,
	}, stdout, stderr)
}
