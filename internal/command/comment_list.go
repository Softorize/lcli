package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/toto/lcli/internal/output"
)

// runCommentList handles the comment list subcommand.
func runCommentList(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("comment list", flag.ContinueOnError)
	postURN := fs.String("post", "", "Post URN to list comments for (required)")
	count := fs.Int("count", 10, "Number of comments to retrieve")
	start := fs.Int("start", 0, "Pagination start index")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *postURN == "" {
		return fmt.Errorf("comment list: --post is required")
	}

	if err := requireAuth(deps.Comments); err != nil {
		return err
	}

	ctx := context.Background()
	list, err := deps.Comments.List(ctx, *postURN, *start, *count)
	if err != nil {
		return fmt.Errorf("comment list: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"ID", "Author", "Text", "Created"}
		rows := make([][]string, 0, len(list.Elements))
		for _, c := range list.Elements {
			rows = append(rows, []string{
				c.ID,
				c.Author,
				truncate(c.Text, 50),
				c.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(list)
}
