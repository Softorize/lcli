package command

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/toto/lcli/internal/output"
)

// runAnalytics dispatches to analytics subcommands: post, views.
func runAnalytics(args []string, deps *Deps) error {
	if len(args) == 0 {
		printAnalyticsUsage(deps)
		return nil
	}

	switch args[0] {
	case "post":
		return runAnalyticsPost(args[1:], deps)
	case "views":
		return runAnalyticsViews(args[1:], deps)
	case "-help", "--help", "-h":
		printAnalyticsUsage(deps)
		return nil
	default:
		return fmt.Errorf("analytics: unknown subcommand %q", args[0])
	}
}

// printAnalyticsUsage writes analytics command help text.
func printAnalyticsUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli analytics <subcommand> [flags]

Subcommands:
  post      View analytics for a specific post
  views     View your profile view count

Use "lcli analytics <subcommand> -help" for more information.
`)
}

// runAnalyticsPost handles the analytics post subcommand.
func runAnalyticsPost(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("analytics post", flag.ContinueOnError)
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("analytics post: post URN argument is required")
	}

	if err := requireAuth(deps.Analytics); err != nil {
		return err
	}

	urn := fs.Arg(0)
	ctx := context.Background()

	stats, err := deps.Analytics.PostAnalytics(ctx, urn)
	if err != nil {
		return fmt.Errorf("analytics post: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Metric", "Value"}
		rows := make([][]string, 0, len(stats))
		for k, v := range stats {
			rows = append(rows, []string{k, fmt.Sprintf("%v", v)})
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(stats)
}

// runAnalyticsViews handles the analytics views subcommand.
func runAnalyticsViews(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("analytics views", flag.ContinueOnError)
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := requireAuth(deps.Analytics); err != nil {
		return err
	}

	ctx := context.Background()
	count, err := deps.Analytics.ProfileViews(ctx)
	if err != nil {
		return fmt.Errorf("analytics views: %w", err)
	}

	printer, err := newPrinter(deps, *outputFmt)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Metric", "Value"}
		rows := [][]string{
			{"Network Size", strconv.Itoa(count)},
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(map[string]int{"networkSize": count})
}
