// Package generator provides the core logic for converting DynamoDB schema definitions
// into Go source code using structured templates.
//
// It orchestrates the code generation process by:
//   - Loading and validating a schema definition from JSON (via the schema package)
//   - Parsing attributes, indexes, and configuration options
//   - Rendering code using predefined templates
//   - Writing generated code to disk in a safe and configurable way
//
// Subpackages include:
//   - schema: schema parsing, validation, and introspection
//   - attribute: attribute typing, subtype handling, and Go type mapping
//   - index: index type validation and composite key resolution
//
// The generator supports customization via the ConfigBuilder API, allowing for dry-run,
// verbose logging, output file overrides, and scoping.
//
// This package is intended to be used internally by code generation tools.
package generator

import (
	"github.com/Mad-Pixels/go-dyno/internal/generator/schema"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmpl"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
)

type Generator struct {
	schemaPath string
	schema     *schema.Schema
}

// NewGenerator object.
func NewGenerator(schemaPath string) (*Generator, error) {
	s, err := schema.NewSchema(schemaPath)
	if err != nil {
		return nil, err
	}
	return &Generator{
		schemaPath: schemaPath,
		schema:     s,
	}, nil
}

// Validate schema.
func (g *Generator) Validate() error {
	if err := g.schema.Validate(); err != nil {
		return err
	}
	logger.Log.Info().Str("schema", g.schemaPath).Msg("Schema is valid")
	return g.schema.Validate()
}

// RenderContent ...
func (g *Generator) RenderContent(packageFl string) string {
	pkg := g.schema.PackageName()
	if packageFl != "" {
		pkg = conv.ToLowerInlineCase(conv.ToSafeName(packageFl))
	}

	tmplMap := v2.TemplateMap{
		PackageName:      pkg,
		TableName:        g.schema.TableName(),
		HashKey:          g.schema.HashKey(),
		RangeKey:         g.schema.RangeKey(),
		Attributes:       g.schema.Attributes(),
		CommonAttributes: g.schema.CommonAttributes(),
		AllAttributes:    g.schema.AllAttributes(),
		SecondaryIndexes: g.schema.SecondaryIndexes(),
	}
	logger.Log.Debug().Any("data", tmplMap).Msg("Template map prepared")
	return tmpl.MustParseTemplateFormattedToString(v2.CodeTemplate, tmplMap)
}

func (g *Generator) FileName() string {
	if g.schema != nil {
		return g.schema.Filename()
	}
	return ""
}

func (g *Generator) PackageName() string {
	if g.schema != nil {
		return g.schema.PackageName()
	}
	return ""
}
