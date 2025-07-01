package attribute

import (
	"fmt"
)

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
func GetUsedNumericSetTypes(attributes []Attribute) []string {
	typesSet := make(map[string]bool)

	for _, attr := range attributes {
		if attr.Type == "NS" && attr.Subtype != SubtypeDefault {
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
func ToDynamoDBStructTag(attr Attribute) string {
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
func IsIntegerAttr(attr Attribute) bool {
	if attr.Subtype != SubtypeDefault {
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
func IsNumericAttr(attr Attribute) bool {
	if attr.Subtype != SubtypeDefault {
		return attr.Subtype.IsNumeric()
	}
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
func ToGolangAttrType(attrName string, attributes []Attribute) string {
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
func ToGolangBaseType(attr Attribute) string {
	switch attr.Type {
	case "SS":
		return "[]string"
	case "NS":
		if attr.Subtype != SubtypeDefault {
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
func ToGolangZeroType(attr Attribute) string {
	switch attr.Type {
	case "SS", "NS", "BS":
		return "nil"
	default:
		return attr.ZeroValue()
	}
}
