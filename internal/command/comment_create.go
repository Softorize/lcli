package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/Softorize/lcli/internal/model"
)

// runCommentCreate handles the comment create subcommand.
func runCommentCreate(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("comment create", flag.ContinueOnError)
	postURN := fs.String("post", "", "Post URN to comment on (required)")
	text := fs.String("text", "", "Comment text content (required)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *postURN == "" {
		return fmt.Errorf("comment create: --post is required")
	}
	if *text == "" {
		return fmt.Errorf("comment create: --text is required")
	}

	if err := requireAuth(deps.Comments); err != nil {
		return err
	}

	ctx := context.Background()
	req := &model.CreateCommentRequest{
		PostURN: *postURN,
		Text:    *text,
	}

	comment, err := deps.Comments.Create(ctx, req)
	if err != nil {
		return fmt.Errorf("comment create: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Comment created: %s\n", comment.ID)
	return nil
}
