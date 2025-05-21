package utils

import (
	"bytes"
	"os"
	"strings"
	"text/template"

	"github.com/Mad-Pixels/go-dyno/internal/logger"

	"github.com/rs/zerolog"
)

func MustParseTemplate(b *bytes.Buffer, tmpl string, vars any) {
	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"Join":             strings.Join,
			"ToUpperCamelCase": ToUpperCamelCase,
			"ToLowerCamelCase": ToLowerCamelCase,
			"ToSafeName":       ToSafeName,
			"ToGolangBaseType": ToGolangBaseType,
			"ToGolangZeroType": ToGolangZeroType,
			"ToGolangAttrType": ToGolangAttrType,
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

func MustParseTemplateToString(tmpl string, vars any) string {
	var b bytes.Buffer
	MustParseTemplate(&b, tmpl, vars)
	return b.String()
}
