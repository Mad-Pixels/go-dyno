package common

import "fmt"

// Attribute defines a basic DynamoDB attribute with a name and type.
// Supported types: "S" (string), "N" (number), "B" (boolean).
type Attribute struct {
	// Name is the logical name of the attribute as defined in the schema.
	Name string `json:"name"`

	// Type is the DynamoDB type of the attribute: "S", "N", or "B".
	Type string `json:"type"`

	// Subtype defines the specific Go type to generate (optional)
	Subtype AttributeSubtype `json:"subtype,omitempty"`
}

// GoType returns the Go type for this attribute.
func (a Attribute) GoType() string {
	if a.Subtype != SubtypeDefault {
		return a.Subtype.GoType()
	}

	switch a.Type {
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

// ZeroValue returns the zero value expression for this attribute.
// Uses subtype if specified, otherwise falls back to DynamoDB type defaults.
//
// Examples:
//
//	attr := Attribute{Type: "S"}                           → `""`
//	attr := Attribute{Type: "N", Subtype: SubtypeInt}      → "0"
//	attr := Attribute{Type: "N", Subtype: SubtypeBigInt}   → "big.NewInt(0)"
//	attr := Attribute{Type: "N", Subtype: SubtypeDefault}  → "0.0" (fallback)
func (a Attribute) ZeroValue() string {
	if a.Subtype != SubtypeDefault {
		return a.Subtype.ZeroValue()
	}

	switch a.Type {
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

// Validate checks if the attribute configuration is valid.
// Returns an error if the DynamoDB type and subtype are incompatible.
func (a Attribute) Validate() error {
	if a.Name == "" {
		return fmt.Errorf("attribute name cannot be empty")
	}

	validTypes := map[string]bool{
		"S":    true,
		"N":    true,
		"B":    true,
		"BOOL": true,
		"BS":   true,
		"SS":   true,
		"NS":   true,
		"L":    true,
		"M":    true,
		"NULL": true,
	}
	if !validTypes[a.Type] {
		return fmt.Errorf("attribute '%s': invalid DynamoDB type '%s'", a.Name, a.Type)
	}

	if err := a.Subtype.Validate(a.Type); err != nil {
		return fmt.Errorf("attribute '%s': %w", a.Name, err)
	}
	return nil
}

// CompositeKeyPart represents a part of a composite key.
// It can either be a constant string (e.g., "user") or an attribute (e.g., "user_id").
type CompositeKeyPart struct {
	// IsConstant indicates whether the part is a literal constant or a reference to an attribute.
	IsConstant bool

	// Value is either the constant string or the attribute name.
	Value string
}
