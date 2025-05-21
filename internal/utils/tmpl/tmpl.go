package tmpl

import (
	"bytes"
	"os"
	"regexp"
	"strings"
	"text/template"
	"unicode"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/schema"

	"github.com/rs/zerolog"
)

func MustParseTemplate(b *bytes.Buffer, tmpl string, vars any) {
	t, err := template.New("tmpl").Funcs(
		template.FuncMap{
			"Join":             strings.Join,
			"ToUpperCamelCase": fmToUpperCamelCase,
			"ToLowerCamelCase": fmToLowerCamelCase,
			"ToSafeName":       fmToSafeName,
			"ToGolangBaseType": fmToGolangBaseType,
			"ToGolangZeroType": fmToGolangZeroType,
			"ToGolangAttrType": fmToGolangAttrType,
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

func fmToUpperCamelCase(s string) string {
	return toCamelCase(s)
}

func fmToLowerCamelCase(s string) string {
	res := toCamelCase(s)
	return strings.ToLower(res[:1]) + res[1:]
}

func fmToSafeName(s string) string {
	s = unsupportedSymbols.ReplaceAllString(s, "_")

	switch {
	case s == "":
		return "_empty"
	case unicode.IsDigit(rune(s[0])):
		s = "_" + s
	}

	if reservedWords[strings.ToLower(s)] {
		s = s + "_"
	}
	return s

}

func fmToGolangBaseType(dynamoType string) string {
	switch dynamoType {
	case "S":
		return "string"
	case "N":
		return "int"
	case "B":
		return "bool"
	default:
		return "any"
	}
}

func fmToGolangZeroType(dynamoType string) string {
	switch dynamoType {
	case "S":
		return `""`
	case "N":
		return "0"
	case "B":
		return "false"
	default:
		return "nil"
	}
}

func fmToGolangAttrType(attrName string, attributes []schema.Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return fmToGolangBaseType(attr.Type)
		}
	}
	return "any"
}

var (
	unsupportedSymbols = regexp.MustCompile(`[^a-zA-Z0-9_]`)
)

var reservedWords = map[string]bool{
	"break":       true,
	"continue":    true,
	"return":      true,
	"fallthrough": true,
	"goto":        true,

	"if":      true,
	"else":    true,
	"for":     true,
	"range":   true,
	"switch":  true,
	"case":    true,
	"default": true,
	"select":  true,

	"var":       true,
	"const":     true,
	"type":      true,
	"struct":    true,
	"interface": true,
	"map":       true,
	"chan":      true,
	"func":      true,
	"package":   true,
	"import":    true,
	"defer":     true,

	"any": true,

	"go": true,
}

func toCamelCase(s string) string {
	var (
		res     strings.Builder
		capNext = true
	)

	for _, r := range s {
		switch {
		case r == '_' || r == '-':
			capNext = true
		case capNext:
			res.WriteRune(unicode.ToUpper(r))
			capNext = false
		default:
			res.WriteRune(r)
		}
	}
	return res.String()
}
