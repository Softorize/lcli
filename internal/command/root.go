package command

import "fmt"

// Version information set via ldflags at build time.
var (
	// BuildVersion is the semantic version of the build.
	BuildVersion = "dev"
	// BuildCommit is the git commit hash of the build.
	BuildCommit = "unknown"
	// BuildDate is the date the binary was built.
	BuildDate = "unknown"
)

// Run is the main dispatch function that routes to the appropriate subcommand.
func Run(args []string, deps *Deps) error {
	if len(args) == 0 {
		printUsage(deps)
		return nil
	}

	cmd := args[0]
	sub := args[1:]

	switch cmd {
	case "help", "-help", "--help", "-h":
		printUsage(deps)
		return nil
	case "auth":
		return runAuth(sub, deps)
	case "config":
		return runConfig(sub, deps)
	case "profile":
		return runProfile(sub, deps)
	case "post":
		return runPost(sub, deps)
	case "comment":
		return runComment(sub, deps)
	case "reaction":
		return runReaction(sub, deps)
	case "media":
		return runMedia(sub, deps)
	case "org":
		return runOrg(sub, deps)
	case "analytics":
		return runAnalytics(sub, deps)
	case "version":
		return runVersion(deps)
	default:
		return fmt.Errorf("unknown command: %s\nRun \"lcli help\" for usage", cmd)
	}
}

// printUsage writes the top-level help text to stdout.
func printUsage(deps *Deps) {
	fmt.Fprint(deps.Stdout, `lcli - LinkedIn CLI tool

Usage:
  lcli <command> [flags]

Commands:
  auth        Authenticate with LinkedIn (login, logout, status)
  config      Configure client credentials (setup)
  profile     View LinkedIn profiles
  post        Create, list, and manage posts
  comment     Manage comments on posts
  reaction    Like and react to posts
  media       Upload images and videos
  org         Manage organization pages
  analytics   View post and profile analytics
  version     Print version information

Use "lcli <command> -help" for more information about a command.
`)
}
