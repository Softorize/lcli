package command

import "fmt"

// runPost dispatches to post subcommands: create, list, get, delete.
func runPost(args []string, deps *Deps) error {
	if len(args) == 0 {
		printPostUsage(deps)
		return nil
	}

	switch args[0] {
	case "create":
		return runPostCreate(args[1:], deps)
	case "list":
		return runPostList(args[1:], deps)
	case "get":
		return runPostGet(args[1:], deps)
	case "delete":
		return runPostDelete(args[1:], deps)
	case "-help", "--help", "-h":
		printPostUsage(deps)
		return nil
	default:
		return fmt.Errorf("post: unknown subcommand %q", args[0])
	}
}

// printPostUsage writes post command help text.
func printPostUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `Usage: lcli post <subcommand> [flags]

Subcommands:
  create    Create a new LinkedIn post
  list      List your recent posts
  get       Get a single post by URN
  delete    Delete a post by URN

Use "lcli post <subcommand> -help" for more information.
`)
}
