package common

import (
	"encoding/json"
	"fmt"
)

// AttributeSubtype defines the specific Go type for DynamoDB attributes.
type AttributeSubtype int

//revive:disable:exported
const (
	// Default (zero value) - use automatic mapping.
	SubtypeDefault AttributeSubtype = iota

	// String subtypes.
	SubtypeString

	// Numeric subtypes.
	SubtypeInt
	SubtypeInt8
	SubtypeInt16
	SubtypeInt32
	SubtypeInt64
	SubtypeFloat32
	SubtypeFloat64
	SubtypeUint
	SubtypeUint8
	SubtypeUint16
	SubtypeUint32
	SubtypeUint64

	// Boolean subtypes.
	SubtypeBool
)

// String returns the string representation of AttributeSubtype.
// JSON marshal/unmarshal (equal with JSON schema).
func (s AttributeSubtype) String() string {
	switch s {
	// String subtypes
	case SubtypeString:
		return "string"

	// Numeric subtypes
	case SubtypeInt:
		return "int"
	case SubtypeInt8:
		return "int8"
	case SubtypeInt16:
		return "int16"
	case SubtypeInt32:
		return "int32"
	case SubtypeInt64:
		return "int64"
	case SubtypeFloat32:
		return "float32"
	case SubtypeFloat64:
		return "float64"
	case SubtypeUint:
		return "uint"
	case SubtypeUint8:
		return "uint8"
	case SubtypeUint16:
		return "uint16"
	case SubtypeUint32:
		return "uint32"
	case SubtypeUint64:
		return "uint64"

	// Boolean subtypes
	case SubtypeBool:
		return "bool"

	// default
	default:
		return "default"
	}
}

// GoType returns the Go type string for code generation.
// Represent in template generation.
// type SchemaItem struct { Price *big.Int }
func (s AttributeSubtype) GoType() string {
	switch s {
	// String subtypes
	case SubtypeString:
		return "string"

	// Numeric subtypes
	case SubtypeInt:
		return "int"
	case SubtypeInt8:
		return "int8"
	case SubtypeInt16:
		return "int16"
	case SubtypeInt32:
		return "int32"
	case SubtypeInt64:
		return "int64"
	case SubtypeFloat32:
		return "float32"
	case SubtypeFloat64:
		return "float64"
	case SubtypeUint:
		return "uint"
	case SubtypeUint8:
		return "uint8"
	case SubtypeUint16:
		return "uint16"
	case SubtypeUint32:
		return "uint32"
	case SubtypeUint64:
		return "uint64"

	// Boolean subtypes
	case SubtypeBool:
		return "bool"

	// default
	default:
		return "any"
	}
}

// MarshalJSON converts AttributeSubtype to JSON string
func (s AttributeSubtype) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON converts JSON string to AttributeSubtype
func (s *AttributeSubtype) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}

	switch str {
	// String subtypes
	case "string":
		*s = SubtypeString

	// Numeric subtypes
	case "int":
		*s = SubtypeInt
	case "int8":
		*s = SubtypeInt8
	case "int16":
		*s = SubtypeInt16
	case "int32":
		*s = SubtypeInt32
	case "int64":
		*s = SubtypeInt64
	case "float32":
		*s = SubtypeFloat32
	case "float64":
		*s = SubtypeFloat64
	case "uint":
		*s = SubtypeUint
	case "uint8":
		*s = SubtypeUint8
	case "uint16":
		*s = SubtypeUint16
	case "uint32":
		*s = SubtypeUint32
	case "uint64":
		*s = SubtypeUint64

	// Boolean subtypes
	case "bool":
		*s = SubtypeBool

	default:
		*s = SubtypeDefault
	}
	return nil
}

// IsNumeric returns true if the subtype represents a numeric Go type.
func (s AttributeSubtype) IsNumeric() bool {
	switch s {
	case SubtypeInt, SubtypeInt8, SubtypeInt16, SubtypeInt32, SubtypeInt64,
		SubtypeUint, SubtypeUint8, SubtypeUint16, SubtypeUint32, SubtypeUint64,
		SubtypeFloat32, SubtypeFloat64:
		return true
	default:
		return false
	}
}

// IsUnsigned returns true if the subtype represents an unsigned integer type.
func (s AttributeSubtype) IsUnsigned() bool {
	switch s {
	case SubtypeUint, SubtypeUint8, SubtypeUint16, SubtypeUint32, SubtypeUint64:
		return true
	default:
		return false
	}
}

// IsInteger returns true if the subtype represents an integer type.
func (s AttributeSubtype) IsInteger() bool {
	switch s {
	case SubtypeInt, SubtypeInt8, SubtypeInt16, SubtypeInt32, SubtypeInt64,
		SubtypeUint, SubtypeUint8, SubtypeUint16, SubtypeUint32, SubtypeUint64:
		return true
	default:
		return false
	}
}

// Validate checks if the subtype is compatible with the given DynamoDB type.
func (s AttributeSubtype) Validate(dynamoType string) error {
	if s == SubtypeDefault {
		return nil
	}

	switch dynamoType {
	case "S":
		if s != SubtypeString {
			return fmt.Errorf("subtype %s is not compatible with DynamoDB type 'S'", s.String())
		}
	case "N":
		if !s.IsNumeric() {
			return fmt.Errorf("subtype %s is not compatible with DynamoDB type 'N'", s.String())
		}
	case "NS":
		if !s.IsNumeric() {
			return fmt.Errorf("subtype %s is not compatible with DynamoDB type 'NS'", s.String())
		}
	default:
		return fmt.Errorf("unknown DynamoDB type: %s", dynamoType)
	}
	return nil
}

// ZeroValue return GoLang zero value which equal current type.
func (s AttributeSubtype) ZeroValue() string {
	switch s {
	case SubtypeString:
		return `""`
	case SubtypeInt, SubtypeInt8, SubtypeInt16, SubtypeInt32, SubtypeInt64,
		SubtypeUint, SubtypeUint8, SubtypeUint16, SubtypeUint32, SubtypeUint64:
		return "0"
	case SubtypeFloat32, SubtypeFloat64:
		return "0.0"
	case SubtypeBool:
		return "false"
	default:
		return "nil"
	}
}
