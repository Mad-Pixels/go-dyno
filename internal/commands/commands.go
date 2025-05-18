package commands

import (
	"github.com/Mad-Pixels/go-dyno/internal/commands/generate"

	"github.com/urfave/cli/v2"
)

func Commands() []*cli.Command {
	return []*cli.Command{
		generate.Command(),
	}
}
