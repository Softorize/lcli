package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/Softorize/lcli/internal/output"
)

// runPostList handles the post list subcommand.
func runPostList(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("post list", flag.ContinueOnError)
	count := fs.Int("count", 10, "Number of posts to retrieve")
	start := fs.Int("start", 0, "Pagination start index")
	author := fs.String("author", "me", "Author URN (defaults to 'me')")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := requireAuth(deps.Posts); err != nil {
		return err
	}

	ctx := context.Background()
	list, err := deps.Posts.ListByAuthor(ctx, *author, *start, *count)
	if err != nil {
		return fmt.Errorf("post list: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"ID", "Text", "Visibility", "Created"}
		rows := make([][]string, 0, len(list.Elements))
		for _, p := range list.Elements {
			rows = append(rows, []string{
				p.ID,
				truncate(p.Text, 50),
				p.Visibility,
				p.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(list)
}
