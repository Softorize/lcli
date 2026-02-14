package command

import (
	"context"
	"flag"
	"fmt"
)

// runPostDelete handles the post delete subcommand.
func runPostDelete(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("post delete", flag.ContinueOnError)
	confirm := fs.Bool("confirm", false, "Skip confirmation prompt")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("post delete: post URN argument is required")
	}

	urn := fs.Arg(0)
	if !*confirm {
		fmt.Fprintf(deps.Stderr, "Delete post %s? Use --confirm to skip this prompt.\n", urn)
		return nil
	}

	if err := requireAuth(deps.Posts); err != nil {
		return err
	}

	ctx := context.Background()
	if err := deps.Posts.Delete(ctx, urn); err != nil {
		return fmt.Errorf("post delete: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Post %s deleted.\n", urn)
	return nil
}
