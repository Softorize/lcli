package command

import (
	"flag"
	"fmt"
	"strings"
	"time"

	"github.com/Softorize/lcli/internal/config"
)

// runAuthStatus handles the auth status subcommand.
func runAuthStatus(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("auth status", flag.ContinueOnError)
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	token, err := config.LoadToken()
	if err != nil {
		return fmt.Errorf("auth status: %w", err)
	}

	if token == nil {
		fmt.Fprintf(deps.Stdout, "Status: not authenticated\n")
		fmt.Fprintf(deps.Stdout, "Run 'lcli auth login' to authenticate.\n")
		return nil
	}

	status := "valid"
	if !token.Valid() {
		status = "expired"
	}

	fmt.Fprintf(deps.Stdout, "Status:  %s\n", status)
	fmt.Fprintf(deps.Stdout, "Expires: %s\n", token.ExpiresAt.Format(time.RFC3339))
	if len(token.Scopes) > 0 {
		fmt.Fprintf(deps.Stdout, "Scopes:  %s\n", strings.Join(token.Scopes, ", "))
	}

	return nil
}
