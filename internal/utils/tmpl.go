package utils

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/Mad-Pixels/go-dyno/internal/logger"

	"github.com/rs/zerolog"
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
	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"Join":             strings.Join,
			"ToUpper":          strings.ToUpper,
			"ToUpperCamelCase": ToUpperCamelCase,
			"ToLowerCamelCase": ToLowerCamelCase,
			"ToGolangBaseType": ToGolangBaseType,
			"ToGolangZeroType": ToGolangZeroType,
			"ToGolangAttrType": ToGolangAttrType,
			"ToSafeName":       ToSafeName,
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
}

// MustParseTemplateToString renders the provided template with variables and
// returns the result as a string. Panics fatally on any error.
//
// It's a convenience wrapper around MustParseTemplate.
//
// Example:
//
//	tmpl := "Field: {{ .Field }}, Type: {{ ToGolangBaseType .Type }}"
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
