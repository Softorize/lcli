package command

import (
	"context"
	"flag"
	"fmt"

	"github.com/Softorize/lcli/internal/model"
)

// validReactionTypes lists the accepted reaction type values.
var validReactionTypes = map[string]bool{
	"LIKE":       true,
	"CELEBRATE":  true,
	"SUPPORT":    true,
	"LOVE":       true,
	"INSIGHTFUL": true,
	"FUNNY":      true,
}

// runReactionLike handles the reaction like subcommand.
func runReactionLike(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("reaction like", flag.ContinueOnError)
	reactionType := fs.String("type", "LIKE", "Reaction type: LIKE, CELEBRATE, SUPPORT, LOVE, INSIGHTFUL, FUNNY")
	actor := fs.String("actor", "me", "Actor URN (defaults to 'me')")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("reaction like: post URN argument is required")
	}

	if !validReactionTypes[*reactionType] {
		return fmt.Errorf("reaction like: invalid type %q", *reactionType)
	}

	if err := requireAuth(deps.Reactions); err != nil {
		return err
	}

	urn := fs.Arg(0)
	ctx := context.Background()

	err := deps.Reactions.React(ctx, *actor, urn, model.ReactionType(*reactionType))
	if err != nil {
		return fmt.Errorf("reaction like: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Reacted to %s with %s.\n", urn, *reactionType)
	return nil
}

// runReactionUnlike handles the reaction unlike subcommand.
func runReactionUnlike(args []string, deps *Deps) error {
	fs := flag.NewFlagSet("reaction unlike", flag.ContinueOnError)
	actor := fs.String("actor", "me", "Actor URN (defaults to 'me')")
	fs.SetOutput(deps.Stderr)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() < 1 {
		return fmt.Errorf("reaction unlike: post URN argument is required")
	}

	if err := requireAuth(deps.Reactions); err != nil {
		return err
	}

	urn := fs.Arg(0)
	ctx := context.Background()

	if err := deps.Reactions.Unreact(ctx, *actor, urn); err != nil {
		return fmt.Errorf("reaction unlike: %w", err)
	}

	fmt.Fprintf(deps.Stderr, "Reaction removed from %s.\n", urn)
	return nil
}
