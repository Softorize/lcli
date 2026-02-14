package command

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/Softorize/lcli/internal/auth"
	"github.com/Softorize/lcli/internal/config"
)

// runAuthLogin handles the auth login subcommand.
func runAuthLogin(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("auth login", flag.ContinueOnError)
	port := fs.Int("port", 8484, "Local port for OAuth callback server")
	timeout := fs.Int("timeout", 120, "Timeout in seconds for login flow")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("auth login: %w", err)
	}

	if cfg.ClientID == "" || cfg.ClientSecret == "" {
		return fmt.Errorf("auth login: credentials not configured — run 'lcli config setup' first")
	}

	authenticator := auth.NewAuthenticator(cfg)

	state, err := auth.RandomState()
	if err != nil {
		return fmt.Errorf("auth login: %w", err)
	}

	srv := auth.NewCallbackServer(*port)
	url := authenticator.AuthorizationURL(state)

	fmt.Fprintf(deps.Stderr, "Open this URL in your browser:\n\n  %s\n\nWaiting for callback...\n", url)
	openBrowser(url)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout)*time.Second)
	defer cancel()

	code, returnedState, err := srv.Start(ctx)
	if err != nil {
		return fmt.Errorf("auth login: %w", err)
	}

	if returnedState != state {
		return fmt.Errorf("auth login: state mismatch — possible CSRF attack")
	}

	token, err := authenticator.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("auth login: %w", err)
	}

	if err := config.SaveToken(token); err != nil {
		return fmt.Errorf("auth login: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Authentication successful!\n")
	return nil
}
