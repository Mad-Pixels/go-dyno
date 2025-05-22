package generate

import (
	"path/filepath"

	"github.com/Mad-Pixels/go-dyno/internal/schema"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"

	"github.com/urfave/cli/v2"
)

func action(ctx *cli.Context) (err error) {
	var (
		cfgFl  = getFlagCfgValue(ctx)
		destFl = getFlagDestValue(ctx)
	)

	if err = utils.IsFileOrError(cfgFl); err != nil {
		return err
	}
	dynamoSchema, err := schema.LoadSchema(cfgFl)
	if err != nil {
		return err
	}
	genFilepath := filepath.Join(destFl, dynamoSchema.Directory(), dynamoSchema.Filename())
	if err = utils.IsFileOrCreate(
		genFilepath,
	); err != nil {
		return err
	}

	if err = process(dynamoSchema, genFilepath); err != nil {
		utils.RemovePath(destFl)
	}
	return err
}

func process(dynamoSchema *schema.DynamoSchema, p string) error {
	schemaMap := v2.TemplateMapV2{
		PackageName:      dynamoSchema.PackageName(),
		TableName:        dynamoSchema.TableName(),
		HashKey:          dynamoSchema.HashKey(),
		RangeKey:         dynamoSchema.RangeKey(),
		Attributes:       dynamoSchema.Attributes(),
		CommonAttributes: dynamoSchema.CommonAttributes(),
		AllAttributes:    dynamoSchema.AllAttributes(),
		SecondaryIndexes: dynamoSchema.SecondaryIndexes(),
	}

	res := utils.MustParseTemplateToString(v2.CodeTemplate, schemaMap)
	if err := utils.WriteToFile(p, []byte(res)); err != nil {
		return err
	}

	return nil
}
