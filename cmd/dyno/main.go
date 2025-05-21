package main

import (
	"os"

	"github.com/Mad-Pixels/go-dyno/internal/commands/generate"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/rs/zerolog"

	godyno "github.com/Mad-Pixels/go-dyno"
	cli "github.com/urfave/cli/v2"
)

func main() {
	logger.Init()
	app := &cli.App{
		Name:    godyno.Name,
		Usage:   godyno.Usage,
		Version: godyno.Version,

		Commands: []*cli.Command{
			generate.Command(),
		},
	}

	logger.Log.Info().Msg("ASDAD")

	if err := app.Run(os.Args); err != nil {
		if failure, ok := err.(*logger.Failure); ok {
			failure.Log(zerolog.ErrorLevel)
		} else {
			logger.Log.Error().Msg(err.Error())
		}
	}
}
