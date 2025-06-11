package templatefunc

import (
	"strings"
	"unicode"
)

var (
	reservedWords = map[string]bool{
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
)

func toCamelCase(s string) string {
	var (
		res     strings.Builder
		capNext = true
	)

	for _, r := range s {
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
