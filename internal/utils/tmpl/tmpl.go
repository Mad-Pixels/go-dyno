// Package tmpl provides utility functions for rendering Go text templates with
// built-in helpers for code generation.
//
// It supports:
//   - Safe rendering of templates with panic-on-error semantics
//   - Automatic formatting using gofumpt and goimports for production-ready Go code
//   - Common template functions such as Join, CamelCase conversion, DynamoDB struct tags,
//     Go type resolution for attributes, and more
//
// Typical use cases include:
//   - Generating Go source code from schema definitions
//   - Producing templates with dynamic attribute, type, and tag handling
package tmpl

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/Mad-Pixels/go-dyno/internal/generator/attribute"
	"github.com/Mad-Pixels/go-dyno/internal/generator/mode"
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils/conv"
	"github.com/rs/zerolog"
	"golang.org/x/tools/imports"
	"mvdan.cc/gofumpt/format"
)

// MustParseTemplate renders the given Go text template `tmpl` into buffer `b`
// using the provided `vars`. If parsing or execution fails, it logs and exits.
//
// This function provides built-in helper functions for templates:
// - Join
// - ToUpperCamelCase
// - ToLowerCamelCase
// - ToGolangBaseType
// - ToGolangZeroType
// - ToGolangAttrType
// - ToSafeName
//
// This is intended for internal code generation or templating purposes.
//
// Example:
//
//	var b bytes.Buffer
//	tmpl := "Hello, {{ .Name }}! Your key: {{ ToSafeName .Key }}"
//	vars := map[string]string{"Name": "GoDyno", "Key": "1value"}
//	utils.MustParseTemplate(&b, tmpl, vars)
//	fmt.Println(b.String())
//
// Output:
//
//	Hello, GoDyno! Your key: x1value
func MustParseTemplate(b *bytes.Buffer, tmpl string, vars any) {
	renderTemplate(b, tmpl, vars, false)
}

// MustParseTemplateFormatted renders the given Go text template `tmpl` into buffer `b`
// using the provided `vars` and automatically formats the result using gofumpt.
// If parsing, execution, or formatting fails, it logs and exits.
//
// This function provides the same built-in helper functions as MustParseTemplate
// and additionally ensures the generated Go code is properly formatted with:
// - Correct indentation and spacing
// - Aligned struct fields and tags
// - Proper import grouping and sorting
// - Stricter formatting rules than go fmt
//
// This is intended for Go code generation that must be production-ready.
//
// Example:
//
//	var b bytes.Buffer
//	tmpl := "type User struct {\n{{ .Field }} {{ .Type }}\n}"
//	vars := map[string]string{"Field": "Name", "Type": "string"}
//	utils.MustParseTemplateFormatted(&b, tmpl, vars)
//	fmt.Println(b.String())
//
// Output:
//
//	type User struct {
//		Name string
//	}
func MustParseTemplateFormatted(b *bytes.Buffer, tmpl string, vars any) {
	renderTemplate(b, tmpl, vars, true)
}

// MustParseTemplateToString renders the provided template with variables and
// returns the result as a string. Panics fatally on any error.
//
// It's a convenience wrapper around MustParseTemplate.
//
// Example:
//
//	tmpl := "Field: {{ .Field }}, Type: {{ ToGolangBaseType . }}"
//	output := utils.MustParseTemplateToString(tmpl, map[string]string{
//		"Field": "age",
//		"Type":  "N",
//	})
//	fmt.Println(output)
//
// Output:
//
//	Field: age, Type: int
func MustParseTemplateToString(tmpl string, vars any) string {
	var b bytes.Buffer
	MustParseTemplate(&b, tmpl, vars)
	return b.String()
}

// MustParseTemplateFormattedToString renders the provided template with variables,
// formats the result using gofumpt, and returns it as a string. Panics fatally on any error.
//
// It's a convenience wrapper around MustParseTemplateFormatted.
//
// Example:
//
//	tmpl := "type {{ .Name }} struct {\n{{ .Field }} {{ .Type }}\n}"
//	output := utils.MustParseTemplateFormattedToString(tmpl, map[string]interface{}{
//		"Name":  "User",
//		"Field": "Email",
//		"Type":  "string",
//	})
//	fmt.Println(output)
//
// Output:
//
//	type User struct {
//		Email string
//	}
func MustParseTemplateFormattedToString(tmpl string, vars any) string {
	var b bytes.Buffer
	MustParseTemplateFormatted(&b, tmpl, vars)
	return b.String()
}

// renderTemplate is the internal implementation for template rendering with optional formatting
func renderTemplate(b *bytes.Buffer, tmpl string, vars any, shouldFormat bool) {
	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"Join":                   strings.Join,
			"ToUpper":                strings.ToUpper,
			"ToUpperCamelCase":       conv.ToUpperCamelCase,
			"ToLowerCamelCase":       conv.ToLowerCamelCase,
			"ToGolangBaseType":       attribute.ToGolangBaseType,
			"ToGolangZeroType":       attribute.ToGolangZeroType,
			"ToGolangAttrType":       attribute.ToGolangAttrType,
			"ToSafeName":             conv.ToSafeName,
			"IsNumericAttr":          attribute.IsNumericAttr,
			"IsIntegerAttr":          attribute.IsIntegerAttr,
			"ToDynamoDBStructTag":    attribute.ToDynamoDBStructTag,
			"GetUsedNumericSetTypes": attribute.GetUsedNumericSetTypes,
			"IsFloatType":            conv.IsFloatType,
			"Slice":                  conv.TrimLeftN,
			"IsALL":                  mode.IsALL,
			"IsMIN":                  mode.IsMIN,
			"IsMode":                 mode.IsMode,
		},
	).
		Parse(tmpl)
	if err != nil {
		logger.NewFailure("internal: failed to create template", err).
			Log(zerolog.FatalLevel)
		os.Exit(1)
	}

	if err = t.Execute(b, vars); err != nil {
		logger.NewFailure("internal: failed to write template data", err).
			Log(zerolog.FatalLevel)
		os.Exit(1)
	}

	// Apply formatting if requested
	if shouldFormat {
		formatted, err := format.Source(b.Bytes(), format.Options{})
		if err != nil {
			logger.NewFailure("internal: failed to format generated code with gofumpt", err).
				Log(zerolog.FatalLevel)
			os.Exit(1)
		}
		imported, err := imports.Process("", formatted, &imports.Options{
			Comments:  true,
			TabWidth:  8,
			TabIndent: true,
		})
		if err != nil {
			logger.NewFailure("internal: failed to process imports with goimports", err).
				Log(zerolog.FatalLevel)
			os.Exit(1)
		}

		b.Reset()
		b.Write(imported)
	}
}
