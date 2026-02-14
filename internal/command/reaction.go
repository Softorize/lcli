package command

import "fmt"

// runReaction dispatches to reaction subcommands: like, unlike, list.
func runReaction(args []string, deps *Deps) error {
	if len(args) == 0 {
		printReactionUsage(deps)
		return nil
	}

	switch args[0] {
	case "like":
		return runReactionLike(args[1:], deps)
	case "unlike":
		return runReactionUnlike(args[1:], deps)
	case "list":
		return runReactionList(args[1:], deps)
	case "-help", "--help", "-h":
		printReactionUsage(deps)
		return nil
	default:
		return fmt.Errorf("reaction: unknown subcommand %q", args[0])
	}
}

// printReactionUsage writes reaction command help text.
func printReactionUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli reaction <subcommand> [flags]

Subcommands:
  like      React to a post
  unlike    Remove a reaction from a post
  list      List reactions on a post

Use "lcli reaction <subcommand> -help" for more information.
`)
}
