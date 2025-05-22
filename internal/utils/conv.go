package utils

import (
	"strings"
	"unicode"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
)

func ToUpperCamelCase(s string) string {
	res := ToSafeName(toCamelCase(s))
	return strings.ToUpper(res[:1]) + res[1:]
}

func ToLowerCamelCase(s string) string {
	res := ToSafeName(toCamelCase(s))
	return strings.ToLower(res[:1]) + res[1:]
}

func ToLowerInlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToLower(res)
}

func ToUpperinlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToUpper(res)
}

func ToSafeName(s string) string {
	s = strings.TrimFunc(s, func(r rune) bool {
		return !((r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9'))
	})

	var b strings.Builder
	for _, r := range s {
		switch {
		case (r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(r >= '0' && r <= '9'):
			b.WriteRune(r)
		default:
			b.WriteRune('_')

			//case r == '-' || r == '_':
			//		b.WriteRune('_')
		}
	}

	result := b.String()
	if result == "" {
		result = "xxx"
	}
	if unicode.IsDigit(rune(result[0])) {
		result = "x" + result
	}
	return result
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

func ToGolangAttrType(attrName string, attributes []common.Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return ToGolangBaseType(attr.Type)
		}
	}
	return "any"
}

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

	for _, r := range strings.ToLower(s) {
		switch {
		case r == '_' || r == '-' || r == '#':
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
