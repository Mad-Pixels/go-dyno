package utils

import (
	"bytes"
	"strings"
	"text/template"

	"github.com/Mad-Pixels/go-dyno/internal/logger"

	"github.com/rs/zerolog"
)

func MustParseTemplate(b *bytes.Buffer, tmpl string, vars any) {
	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"Join":  strings.Join,
			"Upper": strings.ToUpper,
			"Lower": strings.ToLower,
		},
	).
		Parse(tmpl)
	if err != nil {
		logger.NewFailure("internal: failed to create template", err).
			Log(zerolog.FatalLevel)
	}

	if err = t.Execute(b, vars); err != nil {
		logger.NewFailure("internal: failed to write template data", err).
			Log(zerolog.FatalLevel)
	}
}

func MustParseTemplateToString(tmpl string, vars any) string {
	var b bytes.Buffer
	MustParseTemplate(&b, tmpl, vars)
	return b.String()
}
