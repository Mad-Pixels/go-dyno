package generate

import (
	"github.com/Mad-Pixels/go-dyno/internal/generator"

	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) (err error) {
	var (
		cfgFl = getFlagCfgValue(ctx)
		dstFl = getFlagDestValue(ctx)
	)

	cfg, err := generator.NewConfigBuilder(cfgFl, dstFl).Build()
	if err != nil {
		return err
	}
	return generator.Generate(cfg)
}
