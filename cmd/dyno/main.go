package main

import (
	"os"

	"github.com/Mad-Pixels/go-dyno/internal/commands/generate"
	"github.com/Mad-Pixels/go-dyno/internal/logger"

	godyno "github.com/Mad-Pixels/go-dyno"
	cli "github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    godyno.Name,
		Usage:   godyno.Usage,
		Version: godyno.Version,

		Commands: []*cli.Command{
			generate.Command(),
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Log.Fatal(err)
	}
}
