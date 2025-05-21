package generate

import (
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmpl"

	cli "github.com/urfave/cli/v2"
)

var (
	name  = "gen"
	usage = "generate static golang code from config"

	flagCfg  = "cfg"
	flagDest = "dest"
)

type tmplUsage struct {
}

func Command() *cli.Command {
	usageText := tmpl.MustParseTemplateToString(
		usageTemplate,
		tmplUsage{},
	)

	return &cli.Command{
		Name:      name,
		Usage:     usage,
		UsageText: usageText,
		Action:    action,
		Flags:     flags(),
	}
}
