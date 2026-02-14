package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/Softorize/lcli/internal/output"
)

// runPostGet handles the post get subcommand.
func runPostGet(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("post get", flag.ContinueOnError)
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("post get: post URN argument is required")
	}

	if err := requireAuth(deps.Posts); err != nil {
		return err
	}

	urn := fs.Arg(0)
	ctx := context.Background()

	post, err := deps.Posts.Get(ctx, urn)
	if err != nil {
		return fmt.Errorf("post get: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Field", "Value"}
		rows := [][]string{
			{"ID", post.ID},
			{"Author", post.Author},
			{"Text", truncate(post.Text, 80)},
			{"Visibility", post.Visibility},
			{"State", post.LifecycleState},
			{"Created", post.CreatedAt.Format("2006-01-02 15:04")},
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(post)
}
