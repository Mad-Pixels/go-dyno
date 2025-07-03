package validate

import (
	"github.com/Mad-Pixels/go-dyno/internal/app/flags"
	"github.com/Mad-Pixels/go-dyno/internal/generator"
	"github.com/Mad-Pixels/go-dyno/internal/logger"

	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) (err error) {
	var (
		schemaPath = ctx.String(flags.LocalSchema.GetName())
	)

	g, err := generator.NewGenerator(schemaPath)
	if err != nil {
		return err
	}
	if err := g.Validate(); err != nil {
		return err
	}

	logger.Log.Info().Str("path", schemaPath).Msg("schema is valid")
	return nil
}
