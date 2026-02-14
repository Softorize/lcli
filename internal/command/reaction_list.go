package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/Softorize/lcli/internal/output"
)

// runReactionList handles the reaction list subcommand.
func runReactionList(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("reaction list", flag.ContinueOnError)
	count := fs.Int("count", 10, "Number of reactions to retrieve")
	start := fs.Int("start", 0, "Pagination start index")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("reaction list: post URN argument is required")
	}

	if err := requireAuth(deps.Reactions); err != nil {
		return err
	}

	urn := fs.Arg(0)
	ctx := context.Background()

	list, err := deps.Reactions.List(ctx, urn, *start, *count)
	if err != nil {
		return fmt.Errorf("reaction list: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Actor", "Type", "Created"}
		rows := make([][]string, 0, len(list.Elements))
		for _, r := range list.Elements {
			rows = append(rows, []string{
				r.Actor,
				string(r.Type),
				r.CreatedAt.Format("2006-01-02 15:04"),
			})
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(list)
}
