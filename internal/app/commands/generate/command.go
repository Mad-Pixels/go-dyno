// Package generate provides a CLI command for generate GoLang code from JSON schema.
package generate

import (
	godyno "github.com/Mad-Pixels/go-dyno"
	"github.com/Mad-Pixels/go-dyno/internal/app/flags"
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmpl"

	cli "github.com/urfave/cli/v2"
)

var (
	name  = "generate"
	usage = "generate static golang code from JSON schema"
)

type tmplUsage struct {
	Command   string
	EnvPrefix string

	FlagSchemaPath string
	FlagMode       string
}

// Command entrypoint.
func Command() *cli.Command {
	usageText := tmpl.MustParseTemplateToString(
		usageTemplate,
		tmplUsage{
			Command:   name,
			EnvPrefix: godyno.EnvPrefix,

			FlagSchemaPath: flags.LocalSchema.GetName(),
			FlagMode:       flags.LocalGenerateMode.GetName(),
		},
	)

	return &cli.Command{
		Name:      name,
		Usage:     usage,
		UsageText: usageText,
		Action:    action,

		Flags: []cli.Flag{
			flags.LocalSchema.Object,
			flags.LocalOutputDir.Object,
			flags.LocalFilename.Object,
			flags.LocalPackageName.Object,
			flags.LocalGenerateMode.Object,
		},
	}
}
