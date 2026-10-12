package web

import (
	"fmt"
	"net"
	"os/exec"
	"runtime"
	"strconv"
)

// Listen binds the first free port of first..last on the loopback interface only
// (TE §6.4.1): the service is never reachable from the network, so no firewall
// rule or administrator permission is involved.
func Listen(first, last int) (net.Listener, error) {
	var lastErr error
	for p := first; p <= last; p++ {
		ln, err := net.Listen("tcp", net.JoinHostPort("127.0.0.1", strconv.Itoa(p)))
		if err == nil {
			return ln, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("ningun puerto libre en 127.0.0.1:%d-%d: %w", first, last, lastErr)
}

// Banner is the console text printed once the service is listening.
func Banner(url, version string) string {
	return fmt.Sprintf("AutoAtend %s - clasificador de comentarios de TelArg\n"+
		"Servicio listo en %s\n"+
		"Abra esa direccion en el navegador. Para terminar presione Ctrl+C o cierre esta ventana.\n", version, url)
}

// BrowserCommand is the command that opens url in the default browser on goos
// (windows, darwin, linux); ok is false for any other system.
func BrowserCommand(goos, url string) (argv []string, ok bool) {
	switch goos {
	case "windows":
		return []string{"rundll32", "url.dll,FileProtocolHandler", url}, true
	case "darwin":
		return []string{"open", url}, true
	case "linux":
		return []string{"xdg-open", url}, true
	}
	return nil, false
}

// openBrowser launches the default browser. It is the only place that reads
// runtime.GOOS (CLAUDE.md, reglas duras). It does not wait for the browser.
func openBrowser(url string) error {
	argv, ok := BrowserCommand(runtime.GOOS, url)
	if !ok {
		return fmt.Errorf("sistema %s sin lanzador de navegador", runtime.GOOS)
	}
	cmd := exec.Command(argv[0], argv[1:]...)
	if err := cmd.Start(); err != nil {
		return err
	}
	go cmd.Wait()
	return nil
}
