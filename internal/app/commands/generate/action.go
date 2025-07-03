package generate

import (
	"path"

	"github.com/Mad-Pixels/go-dyno/internal/app/flags"
	"github.com/Mad-Pixels/go-dyno/internal/generator"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
	"github.com/Mad-Pixels/go-dyno/internal/utils/writer"

	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) error {
	var (
		schemaPath = ctx.String(flags.LocalSchema.GetName())
		outputPath = ctx.String(flags.LocalOutputDir.GetName())
	)
	logger.Log.Debug().
		Str("schema", schemaPath).
		Str("output", outputPath).
		Msg("Starting code generation")

	g, err := generator.NewGenerator(schemaPath)
	if err != nil {
		return err
	}
	if err := g.Validate(); err != nil {
		return err
	}

	builder := g.NewRenderBuilder()
	if ctx.IsSet(flags.LocalPackageName.GetName()) {
		var (
			raw  = ctx.String(flags.LocalPackageName.GetName())
			safe = conv.ToLowerInlineCase(conv.ToSafeName(raw))
		)

		builder.WithPackageName(safe)
		logger.Log.Debug().
			Str("flag", flags.LocalPackageName.GetName()).
			Str("raw", raw).
			Str("safe", safe).
			Msg("Package name overridden via CLI flag")
	}
	if ctx.IsSet(flags.LocalFilename.GetName()) {
		var (
			raw  = ctx.String(flags.LocalFilename.GetName())
			safe = conv.ToLowerInlineCase(conv.ToSafeName(raw))
		)

		builder.WithFilename(safe)
		logger.Log.Debug().
			Str("flag", flags.LocalFilename.GetName()).
			Str("raw", raw).
			Str("safe", safe).
			Msg("Filename overridden via CLI flag")
	}

	var w writer.Writer
	switch outputPath {
	case "":
		w = writer.NewStdoutWriter()
		logger.Log.Debug().Msg("Using stdout writer")
	default:
		outputFilePath := path.Join(
			outputPath,
			builder.GetPackageName(),
			builder.GetFilename(),
		)
		w = writer.NewFileWriter(outputFilePath)
		logger.Log.Debug().
			Str("path", outputFilePath).
			Msg("Using file writer")
	}

	if err := w.Write([]byte(builder.Build())); err != nil {
		return logger.NewFailure("failed to write generated content", err).
			With("writer", w.Type()).
			With("schema", schemaPath)
	}
	logger.Log.Info().
		Str("schema", schemaPath).
		Str("table", g.TableName()).
		Str("package", builder.GetPackageName()).
		Str("filename", builder.GetFilename()).
		Str("writer", w.Type()).
		Msg("Code generated successfully")
	return nil
}
