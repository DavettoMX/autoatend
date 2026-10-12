package web

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"
)

// Run starts the service (or the --selftest battery) and returns the process exit
// code: 0 on success and on a clean shutdown (ctx cancelled by SIGINT/SIGTERM),
// 1 on any startup error or failed check. The temporary directory is removed on
// every path.
func Run(ctx context.Context, opts Options, stdout, stderr io.Writer) int {
	opts = opts.withDefaults()
	app, err := Boot(opts)
	if err != nil {
		fmt.Fprintf(stderr, "autoatend: no se pudo iniciar: %v\n", err)
		return 1
	}
	defer app.Close()

	if opts.Selftest {
		if app.Selftest(stdout, stderr) {
			return 0
		}
		return 1
	}

	ln, err := Listen(opts.FirstPort, opts.LastPort)
	if err != nil {
		fmt.Fprintf(stderr, "autoatend: no se pudo iniciar: %v\n", err)
		return 1
	}
	tcp := ln.Addr().(*net.TCPAddr).Port
	app.SetPort(tcp)
	url := fmt.Sprintf("http://127.0.0.1:%d", tcp)

	srv := &http.Server{
		Handler:           app.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	serveErr := make(chan error, 1)
	go func() { serveErr <- srv.Serve(ln) }()

	fmt.Fprint(stdout, Banner(url, opts.Version))
	if !opts.NoBrowser {
		if err := opts.OpenBrowser(url); err != nil {
			fmt.Fprintf(stdout, "No se pudo abrir el navegador (%v); abra %s manualmente.\n", err, url)
		}
	}

	select {
	case err := <-serveErr:
		fmt.Fprintf(stderr, "autoatend: el servidor se detuvo: %v\n", err)
		return 1
	case <-ctx.Done():
	}
	sctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	if err := srv.Shutdown(sctx); err != nil && !errors.Is(err, context.DeadlineExceeded) {
		fmt.Fprintf(stderr, "autoatend: error al detener el servidor: %v\n", err)
	}
	srv.Close()
	fmt.Fprintln(stdout, "Servicio detenido.")
	return 0
}
