package command

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/toto/lcli/internal/output"
)

// runOrgFollowers handles the org followers subcommand.
func runOrgFollowers(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("org followers", flag.ContinueOnError)
	orgURN := fs.String("org", "", "Organization URN (required)")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *orgURN == "" {
		return fmt.Errorf("org followers: --org is required")
	}

	if err := requireAuth(deps.Orgs); err != nil {
		return err
	}

	ctx := context.Background()
	stats, err := deps.Orgs.FollowerStats(ctx, *orgURN)
	if err != nil {
		return fmt.Errorf("org followers: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Metric", "Count"}
		rows := [][]string{
			{"Total", strconv.Itoa(stats.TotalCount)},
			{"Organic", strconv.Itoa(stats.OrganicCount)},
			{"Paid", strconv.Itoa(stats.PaidCount)},
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(stats)
}
