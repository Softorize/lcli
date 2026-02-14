package command

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Softorize/lcli/internal/config"
)

// runAuthLogout handles the auth logout subcommand.
func runAuthLogout(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("auth logout", flag.ContinueOnError)
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	tokenPath := filepath.Join(config.ConfigDir(), "tokens.json")

	if err := os.Remove(tokenPath); err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(deps.Stderr, "No stored credentials found.\n")
			return nil
		}
		return fmt.Errorf("auth logout: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Credentials removed.\n")
	return nil
}
