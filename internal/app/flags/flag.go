// Package flags provides centralized CLI flag definitions for GoDyno commands.
//
// Flag is a wrapper around cli.Flag that provides helper methods for consistent
// flag handling across the application.
//
// local.go contains local flags used internally by commands with automatic
// environment variable support (GODYNO_ prefix).
package flags

import "github.com/urfave/cli/v2"

// Flag wraps a cli.Flag with additional helper methods.
type Flag struct {
	// Object is the underlying urfave/cli/v2 flag implementation
	Object cli.Flag
}

// GetName returns the primary name of the flag (the first name in the Names() slice).
// This is used for consistent flag name resolution in command actions.
//
// Example:
//
//	flagName := flags.LocalSchema.GetName()  // Returns "schema"
//	value := ctx.String(flagName)
func (f Flag) GetName() string {
	return f.Object.Names()[0]
}
