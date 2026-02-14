package command

import (
	"context"
	"flag"
	"fmt"
)

// runCommentDelete handles the comment delete subcommand.
func runCommentDelete(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("comment delete", flag.ContinueOnError)
	confirm := fs.Bool("confirm", false, "Skip confirmation prompt")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("comment delete: comment URN argument is required")
	}

	urn := fs.Arg(0)
	if !*confirm {
		fmt.Fprintf(deps.Stderr, "Delete comment %s? Use --confirm to skip this prompt.\n", urn)
		return nil
	}

	if err := requireAuth(deps.Comments); err != nil {
		return err
	}

	ctx := context.Background()
	if err := deps.Comments.Delete(ctx, urn); err != nil {
		return fmt.Errorf("comment delete: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Comment %s deleted.\n", urn)
	return nil
}
