// Package main define an entrypoint for cli.
package main

import (
	"os"

	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/Mad-Pixels/go-dyno/internal/app/commands/generate"
	"github.com/Mad-Pixels/go-dyno/internal/app/commands/validate"
	"github.com/Mad-Pixels/go-dyno/internal/logger"

	"github.com/rs/zerolog"
	cli "github.com/urfave/cli/v2"
)

func init() {
	logger.Init()
}

func main() {
	app := &cli.App{
		Name:    godyno.Name,
		Usage:   godyno.Usage,
		Version: godyno.Version,

		Commands: []*cli.Command{
			generate.Command(),
			validate.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		if failure, ok := err.(*logger.Failure); ok {
			failure.Log(zerolog.ErrorLevel)
		} else {
			logger.Log.Error().Msg(err.Error())
		}
	}
}
