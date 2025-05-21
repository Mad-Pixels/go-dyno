package utils

import (
	"regexp"
	"strings"
	"unicode"
)

func ToUpperCamelCase(s string) string {
	return toCamelCase(s)
}

func ToLowerCamelCase(s string) string {
	res := toCamelCase(s)
	return strings.ToLower(res[:1]) + res[1:]
}

func ToSafeName(s string) string {
	s = unsupportedSymbols.ReplaceAllString(s, "_")

	switch {
	case s == "":
		return "x"
	case unicode.IsDigit(rune(s[0])):
		s = "x" + s
	}

	if reservedWords[strings.ToLower(s)] {
		s = s + "x"
	}
	return s
}

func ToGolangBaseType(dynamoType string) string {
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

func ToGolangZeroType(dynamoType string) string {
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

func ToGolangAttrType(attrName string, attributes []Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return ToGolangBaseType(attr.Type)
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
