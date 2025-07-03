package generate

import (
	"path"

	"github.com/Mad-Pixels/go-dyno/internal/app/flags"
	"github.com/Mad-Pixels/go-dyno/internal/generator"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/writer"

	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) (err error) {
	var (
		w writer.Writer

		genPackageName = ""
		genFileName    = ""

		schemaPath = ctx.String(flags.LocalSchema.GetName())
		outputPath = ctx.String(flags.LocalOutputDir.GetName())
	)

	g, err := generator.NewGenerator(schemaPath)
	if err != nil {
		return err
	}
	if err := g.Validate(); err != nil {
		return err
	}

	genFileName = g.FileName()
	if flags.LocalFilename.Object.IsSet() {
		genFileName = ctx.String(flags.LocalFilename.GetName())
	}

	genPackageName = g.PackageName()
	if flags.LocalPackageName.Object.IsSet() {
		genPackageName = ctx.String(flags.LocalPackageName.GetName())
	}

	switch outputPath {
	case "":
		w = writer.NewStdoutWriter()
	default:
		w = writer.NewFileWriter(path.Join(outputPath, genPackageName, genFileName))
	}
	w.Write([]byte(g.RenderContent(genPackageName)))

	logger.Log.Info().Str("path", schemaPath).Msg("schema is valid")
	return nil
}
