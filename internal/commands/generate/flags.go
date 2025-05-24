package generate

import (
	"fmt"
	"strings"

	godyno "github.com/Mad-Pixels/go-dyno"

	"github.com/urfave/cli/v2"
)

func getFlagCfgValue(ctx *cli.Context) string {
	return ctx.String(flagCfg)
}

func getFlagDestValue(ctx *cli.Context) string {
	return ctx.String(flagDest)
}

func flags() []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:  flagCfg,
			Usage: "Set 'JSON' config for generate goLang objects.",
			Aliases: []string{
				"c",
			},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper(flagCfg)),
			},
			Required: true,
		},
		&cli.StringFlag{
			Name:  flagDest,
			Usage: "Set destination filepath.",
			Aliases: []string{
				"d",
			},
			EnvVars: []string{
				fmt.Sprintf("%s_%s", godyno.EnvPrefix, strings.ToUpper(flagDest)),
			},
			Required: true,
		},
	}
}
