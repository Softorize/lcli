package command

import "fmt"

// runOrg dispatches to org subcommands: info, followers, stats.
func runOrg(args []string, deps *Deps) error {
	if len(args) == 0 {
		printOrgUsage(deps)
		return nil
	}

	switch args[0] {
	case "info":
		return runOrgInfo(args[1:], deps)
	case "followers":
		return runOrgFollowers(args[1:], deps)
	case "stats":
		return runOrgStats(args[1:], deps)
	case "-help", "--help", "-h":
		printOrgUsage(deps)
		return nil
	default:
		return fmt.Errorf("org: unknown subcommand %q", args[0])
	}
}

// printOrgUsage writes org command help text.
func printOrgUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli org <subcommand> [flags]

Subcommands:
  info        Get organization information
  followers   Get follower statistics
  stats       Get page statistics

Use "lcli org <subcommand> -help" for more information.
`)
}
