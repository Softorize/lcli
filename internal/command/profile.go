package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/toto/lcli/internal/model"
	"github.com/toto/lcli/internal/output"
)

// runProfile dispatches to profile subcommands: me, view.
func runProfile(args []string, deps *Deps) error {
	if len(args) == 0 {
		printProfileUsage(deps)
		return nil
	}

	switch args[0] {
	case "me":
		return runProfileMe(args[1:], deps)
	case "view":
		return runProfileView(args[1:], deps)
	case "-help", "--help", "-h":
		printProfileUsage(deps)
		return nil
	default:
		return fmt.Errorf("profile: unknown subcommand %q", args[0])
	}
}

// printProfileUsage writes profile command help text.
func printProfileUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli profile <subcommand> [flags]

Subcommands:
  me      Show your own LinkedIn profile
  view    View another user's profile

Use "lcli profile <subcommand> -help" for more information.
`)
}

// runProfileMe handles the profile me subcommand.
func runProfileMe(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("profile me", flag.ContinueOnError)
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if err := requireAuth(deps.Profile); err != nil {
		return err
	}

	ctx := context.Background()
	profile, err := deps.Profile.Me(ctx)
	if err != nil {
		return fmt.Errorf("profile me: %w", err)
	}

	return printProfile(deps, *outputFmt, profile)
}

// runProfileView handles the profile view subcommand.
func runProfileView(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("profile view", flag.ContinueOnError)
	id := fs.String("id", "", "LinkedIn profile ID to view")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *id == "" {
		return fmt.Errorf("profile view: --id is required")
	}

	if err := requireAuth(deps.Profile); err != nil {
		return err
	}

	ctx := context.Background()
	profile, err := deps.Profile.GetByID(ctx, *id)
	if err != nil {
		return fmt.Errorf("profile view: %w", err)
	}

	return printProfile(deps, *outputFmt, profile)
}

// printProfile renders a profile in the requested format.
func printProfile(deps *Deps, fmtStr string, p *model.Profile) error {
	printer, err := newPrinter(deps, fmtStr)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Field", "Value"}
		rows := [][]string{
			{"ID", p.ID},
			{"Name", p.FirstName + " " + p.LastName},
			{"Headline", p.Headline},
			{"Vanity", p.Vanity},
			{"Email", p.Email},
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(p)
}
