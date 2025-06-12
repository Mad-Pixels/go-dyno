// Package conv provides string transformation and sanitization utilities
// used in code generation and template rendering.
//
// It includes helpers to:
//   - Convert arbitrary strings to safe Go identifiers (e.g. ToSafeName)
//   - Transform strings to CamelCase, snake_case variants
//   - Normalize and validate type names for Go and DynamoDB
//   - Handle reserved Go keywords and fix invalid characters
//   - Perform partial slicing and type detection for code templates
//
// This package is primarily used to ensure all identifiers generated from schema inputs
// are valid, readable, and collision-free in generated Go code.
package conv

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
