package command

import (
	"flag"
	"fmt"

	"github.com/toto/lcli/internal/config"
)

// runConfig dispatches to config subcommands: setup.
func runConfig(args []string, deps *Deps) error {
	if len(args) == 0 {
		printConfigUsage(deps)
		return nil
	}

	switch args[0] {
	case "setup":
		return runConfigSetup(args[1:], deps)
	case "-help", "--help", "-h":
		printConfigUsage(deps)
		return nil
	default:
		return fmt.Errorf("config: unknown subcommand %q", args[0])
	}
}

// printConfigUsage writes config command help text.
func printConfigUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli config <subcommand> [flags]

Subcommands:
  setup     Configure LinkedIn app credentials

Use "lcli config <subcommand> -help" for more information.
`)
}

// runConfigSetup handles the config setup subcommand.
func runConfigSetup(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("config setup", flag.ContinueOnError)
	clientID := fs.String("client-id", "", "LinkedIn app client ID (required)")
	clientSecret := fs.String("client-secret", "", "LinkedIn app client secret (required)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *clientID == "" {
		return fmt.Errorf("config setup: --client-id is required")
	}
	if *clientSecret == "" {
		return fmt.Errorf("config setup: --client-secret is required")
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("config setup: %w", err)
	}

	cfg.ClientID = *clientID
	cfg.ClientSecret = *clientSecret

	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("config setup: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Configuration saved to %s\n", config.ConfigDir())
	return nil
}
