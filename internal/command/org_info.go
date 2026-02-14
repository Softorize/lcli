package command

import (
	"context"
	"flag"
	"fmt"
	"strconv"

	"github.com/Softorize/lcli/internal/model"
	"github.com/Softorize/lcli/internal/output"
)

// runOrgInfo handles the org info subcommand.
func runOrgInfo(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("org info", flag.ContinueOnError)
	id := fs.String("id", "", "Organization ID")
	vanity := fs.String("vanity", "", "Organization vanity name")
	outputFmt := fs.String("output", "table", "Output format (json/table/yaml)")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if *id == "" && *vanity == "" {
		return fmt.Errorf("org info: --id or --vanity is required")
	}

	if err := requireAuth(deps.Orgs); err != nil {
		return err
	}

	ctx := context.Background()
	var org *model.Organization
	var err error

	if *id != "" {
		numID, parseErr := strconv.ParseInt(*id, 10, 64)
		if parseErr != nil {
			return fmt.Errorf("org info: invalid id %q: %w", *id, parseErr)
		}
		org, err = deps.Orgs.Get(ctx, numID)
	} else {
		org, err = deps.Orgs.GetByVanity(ctx, *vanity)
	}

	if err != nil {
		return fmt.Errorf("org info: %w", err)
	}

	return printOrg(deps, *outputFmt, org)
}

// printOrg renders an organization in the requested format.
func printOrg(deps *Deps, fmtStr string, org *model.Organization) error {
	printer, err := newPrinter(deps, fmtStr)
	if err != nil {
		return err
	}

	if printer.Format() == output.FormatTable {
		headers := []string{"Field", "Value"}
		rows := [][]string{
			{"ID", strconv.FormatInt(org.ID, 10)},
			{"Name", org.Name},
			{"Vanity", org.VanityName},
			{"Description", truncate(org.Description, 60)},
			{"Website", org.Website},
			{"Followers", strconv.Itoa(org.FollowerCount)},
		}
		return printer.PrintTable(headers, rows)
	}

	return printer.Print(org)
}
