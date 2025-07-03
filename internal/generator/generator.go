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
	"path/filepath"

	"github.com/Mad-Pixels/go-dyno/internal/generator/schema"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/fs"
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmpl"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
)

// Load reads, parses, and validates a DynamoDB schema definition from the given JSON file path.
//
// It performs the following steps:
//  1. Parses the JSON schema into a structured Schema object.
//  2. Validates all attributes, keys, and index definitions for correctness.
//  3. Returns the fully validated *schema.Schema or an error if invalid.
//
// This function is typically used as the first step in the code generation pipeline.
//
// Example:
//
//	schema, err := generator.Load("schema/user-posts.json")
//	if err != nil {
//	    log.Fatal(err)
//	}
func Load(path string) (*schema.Schema, error) {
	spec, err := schema.NewSchema(path)
	if err != nil {
		return nil, err
	}
	if err := spec.Validate(); err != nil {
		return nil, err
	}

	logger.Log.Info().Str("path", path).Msg("Schema was loaded and validated")
	return spec, nil
}

// Render generates Go source code from the given validated DynamoDB schema.
//
// It uses a predefined Go template and a structured template map containing schema metadata,
// including table name, keys, attributes, and secondary indexes.
//
// Returns the formatted Go source code as a string, or an error if rendering fails.
//
// This function assumes the schema has already passed validation.
//
// Example:
//
//	s, _ := generator.Load("schema.json")
//	code, err := generator.Render(s)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(code)
func Render(spec *schema.Schema) (string, error) {
	if spec == nil {
		return "", logger.NewFailure("schema is nil", nil)
	}

	tmplMap := v2.TemplateMap{
		PackageName:      spec.PackageName(),
		TableName:        spec.TableName(),
		HashKey:          spec.HashKey(),
		RangeKey:         spec.RangeKey(),
		Attributes:       spec.Attributes(),
		CommonAttributes: spec.CommonAttributes(),
		AllAttributes:    spec.AllAttributes(),
		SecondaryIndexes: spec.SecondaryIndexes(),
	}

	logger.Log.Debug().Any("data", tmplMap).Msg("Template map prepared")
	return tmpl.MustParseTemplateFormattedToString(v2.CodeTemplate, tmplMap), nil
}

// Generate performs the full code generation pipeline based on the given configuration.
//
// It performs the following steps:
//  1. Loads and validates the schema from the specified path.
//  2. Renders the Go source code using templates.
//  3. Writes the generated code to the target output file, creating the file and directories if needed.
//
// This function is intended to be called by CLI tools or automation processes.
//
// Example:
//
//	builder := generator.NewConfigBuilder("schema.json", "./out").WithVerbose(true)
//	config, _ := builder.Build()
//	err := generator.Generate(config)
//	if err != nil {
//	    log.Fatal(err)
//	}
func Generate(config *Config) error {
	s, err := Load(config.schemaPath)
	if err != nil {
		return err
	}

	c, err := Render(s)
	if err != nil {
		return err
	}

	f := filepath.Join(config.outputDir, s.PackageName(), s.Filename())
	if err := fs.IsFileOrCreate(f); err != nil {
		return err
	}
	if err := fs.WriteToFile(f, []byte(c)); err != nil {
		return err
	}

	logger.Log.Info().Str("path", f).Msg("File created")
	return nil
}
