package command

import (
	"errors"
	"os/exec"
	"runtime"

	"github.com/Softorize/lcli/internal/output"
)

// errNotAuthenticated is returned when a command requires authentication
// but no valid token is available.
var errNotAuthenticated = errors.New("not authenticated — run 'lcli auth login' first")

// newPrinter creates an output.Printer from a format string flag value.
// It writes to deps.Stdout with the parsed format.
func newPrinter(deps *Deps, fmtStr string) (*output.Printer, error) {
	f, err := output.ParseFormat(fmtStr)
	if err != nil {
		return nil, err
	}
	return output.NewPrinter(deps.Stdout, f), nil
}

// openBrowser attempts to open the given URL in the default browser.
// Errors are silently ignored since this is a best-effort convenience.
func openBrowser(url string) {
	switch runtime.GOOS {
	case "darwin":
		_ = exec.Command("open", url).Start()
	case "linux":
		_ = exec.Command("xdg-open", url).Start()
	case "windows":
		_ = exec.Command("cmd", "/c", "start", url).Start()
	}
}

// truncate shortens s to maxLen characters, appending "..." if truncated.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
