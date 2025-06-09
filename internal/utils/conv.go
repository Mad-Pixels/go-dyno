package utils

import (
	"fmt"
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

// ToGolangBaseType maps a DynamoDB attribute to the corresponding Go base type.
// Uses the attribute's subtype if specified, otherwise uses default mapping.
//
// Examples:
//
//	attr := Attribute{Type: "S"}
//	ToGolangBaseType(attr) → "string"
//
//	attr := Attribute{Type: "N", Subtype: SubtypeInt}
//	ToGolangBaseType(attr) → "int"
//
//	attr := Attribute{Type: "SS"}
//	ToGolangBaseType(attr) → "[]string"
func ToGolangBaseType(attr common.Attribute) string {
	// Handle Set types first
	switch attr.Type {
	case "SS":
		return "[]string"
	case "NS":
		return "[]int"
	case "BS":
		return "[][]byte"
	default:
		return attr.GoType()
	}
}

// ToGolangZeroType returns the zero value as a string literal for a DynamoDB attribute.
// Handles Set types properly.
//
// Examples:
//
//	attr := Attribute{Type: "S"}                           → `""`
//	attr := Attribute{Type: "N", Subtype: SubtypeInt}      → "0"
//	attr := Attribute{Type: "SS"}                          → "nil"
func ToGolangZeroType(attr common.Attribute) string {
	// Handle Set types first
	switch attr.Type {
	case "SS", "NS", "BS":
		return "nil"
	default:
		return attr.ZeroValue()
	}
}

// IsNumericAttr returns true if the attribute represents a numeric type
func IsNumericAttr(attr common.Attribute) bool {
	if attr.Subtype != common.SubtypeDefault {
		return attr.Subtype.IsNumeric()
	}
	// Fallback for default types
	return attr.Type == "N"
}

// IsIntegerAttr returns true if the attribute represents an integer type
func IsIntegerAttr(attr common.Attribute) bool {
	if attr.Subtype != common.SubtypeDefault {
		return attr.Subtype.IsInteger()
	}
	// Fallback for default types - default N is now int
	return attr.Type == "N"
}

// ToGolangAttrType looks up a specific attribute in the provided list and
// returns its mapped Go base type.
//
// Example:
//
//	attrs := []common.Attribute{
//	  {Name: "id", Type: "S"},
//	  {Name: "count", Type: "N"},
//	  {Name: "tags", Type: "SS"},
//	}
//	ToGolangAttrType("count", attrs)   → "int"
//	ToGolangAttrType("tags", attrs)    → "[]string"
//	ToGolangAttrType("missing", attrs) → "any"
func ToGolangAttrType(attrName string, attributes []common.Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return ToGolangBaseType(attr)
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

// ToDynamoDBStructTag returns the appropriate dynamodbav struct tag for the attribute.
// For NS/SS/BS types it adds the required set tags to ensure proper marshaling.
func ToDynamoDBStructTag(attr common.Attribute) string {
	switch attr.Type {
	case "NS":
		return fmt.Sprintf(`dynamodbav:"%s,numberset"`, attr.Name)
	case "SS":
		return fmt.Sprintf(`dynamodbav:"%s,stringset"`, attr.Name)
	case "BS":
		return fmt.Sprintf(`dynamodbav:"%s,binaryset"`, attr.Name)
	default:
		return fmt.Sprintf(`dynamodbav:"%s"`, attr.Name)
	}
}
