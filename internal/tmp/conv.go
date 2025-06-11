package tmp

import (
	"fmt"

	"github.com/Mad-Pixels/go-dyno/internal/generator/attribute"
)

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

// IsNumericAttr returns true if the attribute represents a numeric type
func IsNumericAttr(attr attribute.Attribute) bool {
	if attr.Subtype != attribute.SubtypeDefault {
		return attr.Subtype.IsNumeric()
	}
	// Fallback for default types
	return attr.Type == "N"
}

// IsIntegerAttr returns true if the attribute represents an integer type
func IsIntegerAttr(attr attribute.Attribute) bool {
	if attr.Subtype != attribute.SubtypeDefault {
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
func ToGolangAttrType(attrName string, attributes []attribute.Attribute) string {
	for _, attr := range attributes {
		if attr.Name == attrName {
			return ToGolangBaseType(attr)
		}
	}
	return "any"
}

// ToDynamoDBStructTag returns the appropriate dynamodbav struct tag for the attribute.
// For NS/SS/BS types it adds the required set tags to ensure proper marshaling.
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

// GetUsedNumericTypes возвращает все используемые числовые типы в схеме
func GetUsedNumericTypes(attributes []attribute.Attribute) []string {
	typesSet := make(map[string]bool)

	for _, attr := range attributes {
		if attr.Type == "N" && attr.Subtype != attribute.SubtypeDefault {
			typesSet[attr.Subtype.GoType()] = true
		} else if attr.Type == "N" {
			typesSet["int"] = true // default
		}
	}

	var types []string
	for t := range typesSet {
		types = append(types, t)
	}
	return types
}

// GetUsedNumericSetTypes возвращает все используемые числовые set типы в схеме
func GetUsedNumericSetTypes(attributes []attribute.Attribute) []string {
	typesSet := make(map[string]bool)

	for _, attr := range attributes {
		if attr.Type == "NS" && attr.Subtype != attribute.SubtypeDefault {
			typesSet["[]"+attr.Subtype.GoType()] = true
		} else if attr.Type == "NS" {
			typesSet["[]int"] = true // default
		}
	}

	var types []string
	for t := range typesSet {
		types = append(types, t)
	}
	return types
}
