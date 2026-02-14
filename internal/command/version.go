package command

import "fmt"

// runVersion prints the build version, commit, and date.
func runVersion(deps *Deps) error {
	fmt.Fprintf(deps.Stdout, "lcli %s\n", BuildVersion)
	fmt.Fprintf(deps.Stdout, "  commit: %s\n", BuildCommit)
	fmt.Fprintf(deps.Stdout, "  built:  %s\n", BuildDate)
	return nil
}
