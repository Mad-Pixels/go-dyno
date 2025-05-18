package main

import (
	"os"

	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

func init() {
	godyno.Logger.SetOutput(os.Stdout)
	godyno.Logger.SetLevel(logrus.InfoLevel)
	godyno.Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: false,
		ForceColors:   true,
	})
}

func main() {
	app := &cli.App{
		Name:     godyno.Name,
		Usage:    godyno.Usage,
		Commands: commands.Commands(),
	}
	if err := app.Run(os.Args); err != nil {
		godyno.Logger.Fatal(err)
	}
}
