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
	logger.Log.Debug().
		Str("schema", schemaPath).
		Msg("Starting schema validation")

	g, err := generator.NewGenerator(schemaPath)
	if err != nil {
		return err
	}
	if err := g.Validate(); err != nil {
		return err
	}

	logger.Log.Info().
		Str("schema", schemaPath).
		Str("table", g.TableName()).
		Str("package", g.PackageName()).
		Msg("Schema validation completed successfully")
	return nil
}
