package command

import "fmt"

// runAuth dispatches to auth subcommands: login, logout, status.
func runAuth(args []string, deps *Deps) error {
	if len(args) == 0 {
		printAuthUsage(deps)
		return nil
	}

	switch args[0] {
	case "login":
		return runAuthLogin(args[1:], deps)
	case "logout":
		return runAuthLogout(args[1:], deps)
	case "status":
		return runAuthStatus(args[1:], deps)
	case "-help", "--help", "-h":
		printAuthUsage(deps)
		return nil
	default:
		return fmt.Errorf("auth: unknown subcommand %q", args[0])
	}
}

// printAuthUsage writes auth command help text.
func printAuthUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli auth <subcommand> [flags]

Subcommands:
  login     Authenticate with LinkedIn via OAuth
  logout    Remove stored credentials
  status    Show current authentication status

Use "lcli auth <subcommand> -help" for more information.
`)
}
