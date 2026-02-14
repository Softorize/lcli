package command

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/toto/lcli/internal/output"
)

// runOrgStats handles the org stats subcommand.
func runOrgStats(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("org stats", flag.ContinueOnError)
	orgURN := fs.String("org", "", "Organization URN (required)")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *orgURN == "" {
		return fmt.Errorf("org stats: --org is required")
	}

	if err := requireAuth(deps.Orgs); err != nil {
		return err
	}

	ctx := context.Background()
	stats, err := deps.Orgs.PageStats(ctx, *orgURN)
	if err != nil {
		return fmt.Errorf("org stats: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Metric", "Value"}
		rows := [][]string{
			{"Views", strconv.Itoa(stats.Views)},
			{"Unique Visitors", strconv.Itoa(stats.UniqueVisitors)},
			{"Clicks", strconv.Itoa(stats.Clicks)},
			{"Period", stats.Period},
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(stats)
}
