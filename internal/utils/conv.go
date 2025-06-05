package utils

import (
	"strings"
	"unicode"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
)

// ToUpperCamelCase converts a string into UpperCamelCase format,
// ensuring the result is safe for use in Go identifiers.
//
// Examples:
//
//	ToUpperCamelCase("user_id")         → "UserId"
//	ToUpperCamelCase("user-name")       → "UserName"
//	ToUpperCamelCase("1type")           → "X1type"
//	ToUpperCamelCase("full#access")     → "FullAccess"
//	ToUpperCamelCase("!@#special-case") → "XxxSpecialCase"
func ToUpperCamelCase(s string) string {
	res := ToSafeName(toCamelCase(s))
	return strings.ToUpper(res[:1]) + res[1:]
}

// ToLowerCamelCase converts a string into lowerCamelCase format,
// ensuring the result is safe for use in Go identifiers.
//
// Examples:
//
//	ToLowerCamelCase("user_id")  → "userId"
//	ToLowerCamelCase("Type")     → "type"
//	ToLowerCamelCase("1invalid") → "x1invalid"
func ToLowerCamelCase(s string) string {
	res := ToSafeName(toCamelCase(s))
	return strings.ToLower(res[:1]) + res[1:]
}

// ToLowerInlineCase converts a string to lowercase without underscores.
//
// Examples:
//
//	ToLowerInlineCase("user_id")          → "userid"
//	ToLowerInlineCase("snake_case_value") → "snakecasevalue"
func ToLowerInlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToLower(res)
}

// ToUpperinlineCase converts a string to uppercase without underscores.
//
// Examples:
//
//	ToUpperinlineCase("user_id")   → "USERID"
//	ToUpperinlineCase("Api_Token") → "APITOKEN"
func ToUpperinlineCase(s string) string {
	res := strings.ReplaceAll(ToSafeName(s), "_", "")
	return strings.ToUpper(res)
}

// ToSafeName sanitizes any string into a Go-safe identifier:
// - replaces invalid characters with underscores
// - prefixes with 'x' if it starts with a number or is a reserved keyword
//
// Examples:
//
//	ToSafeName("1test")       → "x1test"
//	ToSafeName("type")        → "xtype"
//	ToSafeName("hello_world") → "hello_world"
//	ToSafeName("$$$abc")      → "abc"
func ToSafeName(s string) string {
	s = strings.TrimFunc(s, func(r rune) bool {
		return (r < 'A' || r > 'Z') &&
			(r < 'a' || r > 'z') &&
			(r < '0' || r > '9')
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
		}
	}
	switch {
	case b.Len() == 0:
		return "xxx"
	case unicode.IsDigit(rune(b.String()[0])) || reservedWords[strings.ToLower(b.String())]:
		return "x" + b.String()
	default:
		return b.String()
	}
}

// ToGolangBaseType maps a DynamoDB type to the corresponding Go base type.
//
// Examples:
//
//	ToGolangBaseType("S")       → "string"
//	ToGolangBaseType("N")       → "int"
//	ToGolangBaseType("BOOL")    → "bool"
//	ToGolangBaseType("SS")      → "[]string"
//	ToGolangBaseType("NS")      → "[]int"
//	ToGolangBaseType("UNKNOWN") → "any"
func ToGolangBaseType(dynamoType string) string {
	switch dynamoType {
	case "S":
		return "string"
	case "N":
		return "int"
	case "BOOL":
		return "bool"
	case "SS":
		return "[]string"
	case "NS":
		return "[]int"
	default:
		return "any"
	}
}

// ToGolangZeroType returns the zero value as a string literal for a DynamoDB type.
//
// Examples:
//
//	ToGolangZeroType("S")    → `""`
//	ToGolangZeroType("N")    → "0"
//	ToGolangZeroType("BOOL") → "false"
//	ToGolangZeroType("SS")   → "nil"
//	ToGolangZeroType("NS")   → "nil"
//	ToGolangZeroType("X")    → "nil"
func ToGolangZeroType(dynamoType string) string {
	switch dynamoType {
	case "S":
		return `""`
	case "N":
		return "0"
	case "BOOL":
		return "false"
	case "SS", "NS":
		return "nil"
	default:
		return "nil"
	}
}

// ToGolangAttrType looks up a specific attribute in the provided list and
// returns its mapped Go base type.
//
// Example:
//
//	attrs := []common.Attribute{
//	  {Name: "id", Type: "S"},
//	  {Name: "count", Type: "N"},
//	  {Name: "is_active", Type: "BOOL"},
//	  {Name: "tags", Type: "SS"},
//	  {Name: "scores", Type: "NS"},
//	}
//	ToGolangAttrType("count", attrs)     → "int"
//	ToGolangAttrType("is_active", attrs) → "bool"
//	ToGolangAttrType("tags", attrs)      → "[]string"
//	ToGolangAttrType("scores", attrs)    → "[]int"
//	ToGolangAttrType("missing", attrs)   → "any"
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
