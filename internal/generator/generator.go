// Package generator converts DynamoDB schema definitions into Go source code.
//
// Core workflow:
//   - Load and validate JSON schema
//   - Parse attributes, indexes, and configuration
//   - Render Go code using templates
//   - Support for customization via Builder pattern
//
// Subpackages:
//   - schema: JSON parsing and validation
//   - attribute: DynamoDB type mapping to Go types
//   - index: secondary index handling
package generator

import (
	"github.com/Mad-Pixels/go-dyno/internal/generator/schema"
)

// Generator orchestrates the code generation process from DynamoDB schema to Go code.
type Generator struct {
	schemaPath string
	schema     *schema.Schema
}

// NewGenerator creates a new generator instance from a schema file path.
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

// FileName returns the default output filename based on schema.
func (g *Generator) FileName() string {
	if g.schema != nil {
		return g.schema.Filename()
	}
	return ""
}

// PackageName returns the Go package name derived from schema.
func (g *Generator) PackageName() string {
	if g.schema != nil {
		return g.schema.PackageName()
	}
	return ""
}

// TableName returns the DynamoDB table name from schema.
func (g *Generator) TableName() string {
	if g.schema != nil {
		return g.schema.TableName()
	}
	return ""
}

// NewRenderBuilder creates a new builder instance.
func (g *Generator) NewRenderBuilder() *RenderBuilder {
	return &RenderBuilder{
		generator: g,
	}
}

// Validate performs comprehensive schema validation.
func (g *Generator) Validate() error {
	if err := g.schema.Validate(); err != nil {
		return err
	}
	return g.schema.Validate()
}
