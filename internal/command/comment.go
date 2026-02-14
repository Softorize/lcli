package command

import "fmt"

// runComment dispatches to comment subcommands: create, list, delete.
func runComment(args []string, deps *Deps) error {
	if len(args) == 0 {
		printCommentUsage(deps)
		return nil
	}

	switch args[0] {
	case "create":
		return runCommentCreate(args[1:], deps)
	case "list":
		return runCommentList(args[1:], deps)
	case "delete":
		return runCommentDelete(args[1:], deps)
	case "-help", "--help", "-h":
		printCommentUsage(deps)
		return nil
	default:
		return fmt.Errorf("comment: unknown subcommand %q", args[0])
	}
}

// printCommentUsage writes comment command help text.
func printCommentUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli comment <subcommand> [flags]

Subcommands:
  create    Add a comment to a post
  list      List comments on a post
  delete    Delete a comment

Use "lcli comment <subcommand> -help" for more information.
`)
}
