package flags

import "github.com/urfave/cli/v2"

type Flag struct {
	Object cli.Flag
}

// GetName of the current flag.
func (f Flag) GetName() string {
	return f.Object.Names()[0]
}
