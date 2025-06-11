package attribute

import (
	"github.com/Mad-Pixels/go-dyno/internal/logger"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
)

var (
	// validTypes lists all supported DynamoDB types
	validTypes = map[string]bool{
		// Scalar types
		"S":    true,
		"N":    true,
		"B":    true,
		"BOOL": true,

		// Set types
		"SS": true,
		"NS": true,
		"BS": true,

		// Document and special types
		"L":    true,
		"M":    true,
		"NULL": true,
	}
)

// Attribute defines a DynamoDB attribute with a name, DynamoDB type, and optional Go subtype.
type Attribute struct {
	// Name is the logical name of the attribute as defined in the schema.
	Name string `json:"name"`

	// Type is the DynamoDB type of the attribute: "S", "N", etc...
	Type string `json:"type"`

	// Subtype defines the specific Go type to generate. Optional.
	Subtype attributeSubtype `json:"subtype,omitempty"`
}

// GoType return the Go type for this attribute.
func (a Attribute) GoType() string {
	if !a.Subtype.IsDefault() {
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
func (a Attribute) ZeroValue() string {
	if !a.Subtype.IsDefault() {
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
func (a Attribute) Validate() error {
	if a.Name == "" {
		return logger.NewFailure("attribute name cannot be empty", nil)
	}
	if !validTypes[a.Type] {
		return logger.NewFailure("invalid attribute type", nil).
			With("name", a.Name).
			With("type", a.Type).
			With("available", utils.AvailableKeys(validTypes))
	}
	return a.Subtype.Validate(a.Type)
}
