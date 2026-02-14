// lcli is a command-line tool for interacting with the LinkedIn API.
package main

import (
	"fmt"
	"os"

	"github.com/toto/lcli/internal/client"
	"github.com/toto/lcli/internal/command"
	"github.com/toto/lcli/internal/config"
	"github.com/toto/lcli/internal/linkedin"
	"github.com/toto/lcli/internal/output"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	deps := &command.Deps{
		Cfg:    cfg,
		Output: output.NewPrinter(os.Stdout, output.FormatTable),
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	if err := initServices(cfg, deps); err != nil {
		// Non-fatal: services will be nil and commands that need
		// auth will return an appropriate error.
		fmt.Fprintf(os.Stderr, "warning: %v\n", err)
	}

	return command.Run(os.Args[1:], deps)
}

func initServices(cfg *config.Config, deps *command.Deps) error {
	token, err := config.LoadToken()
	if err != nil {
		return fmt.Errorf("load token: %w", err)
	}

	if token == nil || !token.Valid() {
		return nil
	}

	cli := client.New(token.AccessToken, cfg.APIVersion)

	deps.Profile = linkedin.NewProfileService(cli)
	deps.Posts = linkedin.NewPostService(cli)
	deps.Comments = linkedin.NewCommentService(cli)
	deps.Reactions = linkedin.NewReactionService(cli)
	deps.Media = linkedin.NewMediaService(cli)
	deps.Orgs = linkedin.NewOrgService(cli)
	deps.Analytics = linkedin.NewAnalyticsService(cli)

	return nil
}
