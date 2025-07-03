package generator

import (
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmpl"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
)

// RenderBuilder provides a customizing code generation.
// Allows overriding schema defaults (package name, filename) via CLI flags.
type RenderBuilder struct {
	generator   *Generator
	packageName *string
	filename    *string
}

// WithPackageName overrides the package name with safe conversion.
func (rb *RenderBuilder) WithPackageName(name string) *RenderBuilder {
	if name != "" {
		cleanName := conv.ToLowerInlineCase(conv.ToSafeName(name))
		rb.packageName = &cleanName
	}
	return rb
}

// WithFilename overrides the filename with safe conversion.
func (rb *RenderBuilder) WithFilename(name string) *RenderBuilder {
	if name != "" {
		cleanName := conv.ToSafeName(name)
		rb.filename = &cleanName
	}
	return rb
}

// Build renders the final Go code using configured overrides.
func (rb *RenderBuilder) Build() string {
	var (
		tmplMap = rb.buildTemplateMap()
	)
	logger.Log.Debug().
		Any("data", tmplMap).
		Msg("Template map prepared")
	return tmpl.MustParseTemplateFormattedToString(v2.CodeTemplate, tmplMap)
}

// GetPackageName returns the final package name (override or schema default).
func (rb *RenderBuilder) GetPackageName() string {
	if rb.packageName != nil {
		return *rb.packageName
	}
	return rb.generator.schema.PackageName()
}

// GetFilename returns the final filename (override or schema default).
func (rb *RenderBuilder) GetFilename() string {
	if rb.filename != nil {
		return *rb.filename
	}
	return rb.generator.schema.Filename()
}

// buildTemplateMap creates template data with schema and overrides.
func (rb *RenderBuilder) buildTemplateMap() v2.TemplateMap {
	schema := rb.generator.schema

	return v2.TemplateMap{
		PackageName:      rb.getPackageName(),
		TableName:        schema.TableName(),
		HashKey:          schema.HashKey(),
		RangeKey:         schema.RangeKey(),
		Attributes:       schema.Attributes(),
		CommonAttributes: schema.CommonAttributes(),
		AllAttributes:    schema.AllAttributes(),
		SecondaryIndexes: schema.SecondaryIndexes(),
	}
}

// getPackageName internal helper for consistent package name resolution.
func (rb *RenderBuilder) getPackageName() string {
	if rb.packageName != nil {
		return *rb.packageName
	}
	return rb.generator.schema.PackageName()
}
