package generate

import (
	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/Mad-Pixels/go-dyno/internal/utils"

	cli "github.com/urfave/cli/v2"
)

var (
	name  = "gen"
	usage = "generate static golang code from config"

	flagCfg  = "cfg"
	flagDest = "dest"
)

type tmplUsage struct {
	Command     string
	FlagCfg     string
	FlagDst     string
	EnvPrefix   string
	ExampleJSON string
}

// Command ...
func Command() *cli.Command {
	usageText := utils.MustParseTemplateToString(
		usageTemplate,
		tmplUsage{
			Command:     name,
			FlagCfg:     flagCfg,
			FlagDst:     flagCfg,
			EnvPrefix:   godyno.EnvPrefix,
			ExampleJSON: "dynamo_db_description.json",
		},
	)

	return &cli.Command{
		Name:      name,
		Usage:     usage,
		UsageText: usageText,
		Action:    action,
		Flags:     flags(),
	}
}
