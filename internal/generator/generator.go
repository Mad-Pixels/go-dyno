// Package generator provides ...
package generator

import (
	"path/filepath"

	"github.com/Mad-Pixels/go-dyno/internal/generator/schema"
	"github.com/Mad-Pixels/go-dyno/internal/templ"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
)

func Load(path string) (*schema.Schema, error) {
	spec, err := schema.NewSchema(path)
	if err != nil {
		return nil, err
	}
	if err := spec.Validate(); err != nil {
		return nil, err
	}
	return spec, nil
}

func Generate(config *Config) error {
	spec, err := Load(config.schemaPath)
	if err != nil {
		return err
	}

	g := filepath.Join(config.outputDir, spec.PackageName(), spec.Filename())
	if err := utils.IsFileOrCreate(g); err != nil {
		return err
	}

	schemaMap := v2.TemplateMap{
		PackageName:      spec.PackageName(),
		TableName:        spec.TableName(),
		HashKey:          spec.HashKey(),
		RangeKey:         spec.RangeKey(),
		Attributes:       spec.Attributes(),
		CommonAttributes: spec.CommonAttributes(),
		AllAttributes:    spec.AllAttributes(),
		SecondaryIndexes: spec.SecondaryIndexes(),
	}

	res := templ.MustParseTemplateFormattedToString(v2.CodeTemplate, schemaMap)
	if err := utils.WriteToFile(g, []byte(res)); err != nil {
		return err
	}

	return nil

}
