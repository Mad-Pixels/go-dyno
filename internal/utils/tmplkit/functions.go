package tmplkit

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/Mad-Pixels/go-dyno/internal/generator/attribute"
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

// TrimLeftN returns a substring of the input string `s` with the first `start` characters removed.
//
// If `start` is greater than or equal to the length of the string, it returns an empty string.
// This function is useful in templates and code generation where controlled slicing is needed.
//
// Examples:
//
//	TrimLeftN("[]int", 2)       → "int"
//	TrimLeftN("##Hello", 2)     → "Hello"
//	TrimLeftN("GoLang", 0)      → "GoLang"
//	TrimLeftN("short", 10)      → ""
func TrimLeftN(s string, start int) string {
	if start >= len(s) {
		return ""
	}
	return s[start:]
}

// IsFloatType checks whether the provided Go type name is a floating-point type.
//
// This is commonly used in code generation scenarios to determine how numeric types
// should be handled (e.g., when converting to DynamoDB number sets).
//
// Examples:
//
//	IsFloatType("float32") → true
//	IsFloatType("float64") → true
//	IsFloatType("int64")   → false
//	IsFloatType("double")  → false
func IsFloatType(goType string) bool {
	return goType == "float32" || goType == "float64"
}

// GetUsedNumericSetTypes returns a list of distinct Go slice types used for DynamoDB "NS" (Number Set) attributes.
//
// It inspects the provided list of attributes and collects all unique slice types based on
// their numeric subtype (e.g., "int", "float64") as `[]<type>` strings. If no subtype is specified,
// it defaults to `[]int`.
//
// This function is useful for code generation, particularly when producing type-switch logic
// for marshaling numeric sets.
//
// Examples:
//
//	Input: []Attribute{
//		{Type: "NS", Subtype: Subtype{Kind: "float64"}},
//		{Type: "NS", Subtype: Subtype{Kind: "int"}},
//	}
//	Output: []string{"[]float64", "[]int"}
func GetUsedNumericSetTypes(attributes []attribute.Attribute) []string {
	typesSet := make(map[string]bool)

	for _, attr := range attributes {
		if attr.Type == "NS" && attr.Subtype != attribute.SubtypeDefault {
			typesSet["[]"+attr.Subtype.GoType()] = true
		} else if attr.Type == "NS" {
			typesSet["[]int"] = true
		}
	}

	var types []string
	for t := range typesSet {
		types = append(types, t)
	}
	return types
}

// ToDynamoDBStructTag generates the appropriate `dynamodbav` struct tag for a given attribute.
//
// For set types (NS, SS, BS), it appends the corresponding set directive (e.g., `numberset`)
// to ensure proper encoding when using the AWS SDK for Go v2. For all other types,
// it generates a standard tag with just the attribute name.
//
// Examples:
//
//	ToDynamoDBStructTag(Attribute{Name: "ids", Type: "NS"}) → `dynamodbav:"ids,numberset"`
//	ToDynamoDBStructTag(Attribute{Name: "tags", Type: "SS"}) → `dynamodbav:"tags,stringset"`
//	ToDynamoDBStructTag(Attribute{Name: "data", Type: "S"})  → `dynamodbav:"data"`
func ToDynamoDBStructTag(attr attribute.Attribute) string {
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

// IsIntegerAttr returns true if the given attribute is considered an integer type.
//
// It checks the attribute's subtype (if defined) using its IsInteger method.
// If no subtype is specified, it falls back to treating DynamoDB type "N" (Number)
// as an integer by default.
//
// This is useful when generating code that requires knowledge of whether
// a numeric attribute is an integer (e.g., for marshaling, validation, or formatting).
//
// Examples:
//
//	IsIntegerAttr(Attribute{Type: "N", Subtype: "int"})     → true
//	IsIntegerAttr(Attribute{Type: "N", Subtype: "float64"}) → false
//	IsIntegerAttr(Attribute{Type: "N"})                     → true
//	IsIntegerAttr(Attribute{Type: "S"})                     → false
func IsIntegerAttr(attr attribute.Attribute) bool {
	if attr.Subtype != attribute.SubtypeDefault {
		return attr.Subtype.IsInteger()
	}
	// Fallback for default types - default N is now int
	return attr.Type == "N"
}

// IsNumericAttr returns true if the given attribute is considered a numeric type.
//
// It checks the attribute's subtype (if defined) using its IsNumeric method.
// If no subtype is specified, it defaults to treating DynamoDB type "N" (Number)
// as numeric.
//
// This function is useful for code generation tasks that require distinguishing
// numeric attributes (integers or floats) from other types.
//
// Examples:
//
//	IsNumericAttr(Attribute{Type: "N", Subtype: "int"})     → true
//	IsNumericAttr(Attribute{Type: "N", Subtype: "float64"}) → true
//	IsNumericAttr(Attribute{Type: "N"})                     → true
//	IsNumericAttr(Attribute{Type: "S"})                     → false
func IsNumericAttr(attr attribute.Attribute) bool {
	if attr.Subtype != attribute.SubtypeDefault {
		return attr.Subtype.IsNumeric()
	}
	// Fallback for default types
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
func ToGolangAttrType(attrName string, attributes []attribute.Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return ToGolangBaseType(attr)
		}
	}
	return "any"
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
func ToGolangBaseType(attr attribute.Attribute) string {
	// Handle Set types first
	switch attr.Type {
	case "SS":
		return "[]string"
	case "NS":
		if attr.Subtype != attribute.SubtypeDefault {
			return "[]" + attr.Subtype.GoType()
		}
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
func ToGolangZeroType(attr attribute.Attribute) string {
	// Handle Set types first
	switch attr.Type {
	case "SS", "NS", "BS":
		return "nil"
	default:
		return attr.ZeroValue()
	}
}
