package validate

import (
	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/Mad-Pixels/go-dyno/internal/app/flags"
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmpl"

	cli "github.com/urfave/cli/v2"
)

var (
	name  = "validate"
	usage = "validate JSON schema"
)

type tmplUsage struct {
	Command   string
	EnvPrefix string

	FlagSchemaPath string
}

// Command entrypoint.
func Command() *cli.Command {
	usageText := tmpl.MustParseTemplateToString(
		usageTemplate,
		tmplUsage{
			Command:   name,
			EnvPrefix: godyno.EnvPrefix,

			FlagSchemaPath: flags.LocalSchema.GetName(),
		},
	)

	return &cli.Command{
		Name:      name,
		Usage:     usage,
		UsageText: usageText,
		Action:    action,

		Flags: []cli.Flag{
			flags.LocalSchema.Object,
		},
	}
}
