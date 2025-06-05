package common

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttributeSubtype_String(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected string
	}{
		// String subtypes
		{"string", SubtypeString, "string"},

		// Signed integer subtypes
		{"int", SubtypeInt, "int"},
		{"int8", SubtypeInt8, "int8"},
		{"int16", SubtypeInt16, "int16"},
		{"int32", SubtypeInt32, "int32"},
		{"int64", SubtypeInt64, "int64"},

		// Unsigned integer subtypes
		{"uint", SubtypeUint, "uint"},
		{"uint8", SubtypeUint8, "uint8"},
		{"uint16", SubtypeUint16, "uint16"},
		{"uint32", SubtypeUint32, "uint32"},
		{"uint64", SubtypeUint64, "uint64"},

		// Floating point subtypes
		{"float32", SubtypeFloat32, "float32"},
		{"float64", SubtypeFloat64, "float64"},

		// Arbitrary precision subtypes
		{"big_int", SubtypeBigInt, "big_int"},
		{"decimal", SubtypeDecimal, "decimal"},

		// Boolean subtypes
		{"bool", SubtypeBool, "bool"},

		// Future extensibility
		{"bytes", SubtypeBytes, "bytes"},
		{"time", SubtypeTime, "time"},
		{"uuid", SubtypeUUID, "uuid"},

		// Default
		{"default", SubtypeDefault, "default"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subtype.String())
		})
	}
}

func TestAttributeSubtype_GoType(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected string
	}{
		// String subtypes
		{"string", SubtypeString, "string"},

		// Signed integer subtypes
		{"int", SubtypeInt, "int"},
		{"int8", SubtypeInt8, "int8"},
		{"int16", SubtypeInt16, "int16"},
		{"int32", SubtypeInt32, "int32"},
		{"int64", SubtypeInt64, "int64"},

		// Unsigned integer subtypes
		{"uint", SubtypeUint, "uint"},
		{"uint8", SubtypeUint8, "uint8"},
		{"uint16", SubtypeUint16, "uint16"},
		{"uint32", SubtypeUint32, "uint32"},
		{"uint64", SubtypeUint64, "uint64"},

		// Floating point subtypes
		{"float32", SubtypeFloat32, "float32"},
		{"float64", SubtypeFloat64, "float64"},

		// Arbitrary precision subtypes
		{"big_int", SubtypeBigInt, "*big.Int"},
		{"decimal", SubtypeDecimal, "*decimal.Decimal"},

		// Boolean subtypes
		{"bool", SubtypeBool, "bool"},

		// Future extensibility
		{"bytes", SubtypeBytes, "[]byte"},
		{"time", SubtypeTime, "time.Time"},
		{"uuid", SubtypeUUID, "uuid.UUID"},

		// Default
		{"default", SubtypeDefault, "any"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subtype.GoType())
		})
	}
}

func TestAttributeSubtype_ZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected string
	}{
		// String subtypes
		{"string", SubtypeString, `""`},

		// Integer subtypes (all return "0")
		{"int", SubtypeInt, "0"},
		{"int8", SubtypeInt8, "0"},
		{"uint64", SubtypeUint64, "0"},

		// Floating point subtypes
		{"float32", SubtypeFloat32, "0.0"},
		{"float64", SubtypeFloat64, "0.0"},

		// Boolean subtypes
		{"bool", SubtypeBool, "false"},

		// Arbitrary precision subtypes
		{"big_int", SubtypeBigInt, "big.NewInt(0)"},
		{"decimal", SubtypeDecimal, "decimal.Zero"},

		// Future extensibility
		{"bytes", SubtypeBytes, "nil"},
		{"time", SubtypeTime, "time.Time{}"},
		{"uuid", SubtypeUUID, "uuid.UUID{}"},

		// Default
		{"default", SubtypeDefault, "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subtype.ZeroValue())
		})
	}
}

func TestAttributeSubtype_IsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected bool
	}{
		// Numeric types
		{"int", SubtypeInt, true},
		{"uint64", SubtypeUint64, true},
		{"float32", SubtypeFloat32, true},
		{"big_int", SubtypeBigInt, true},
		{"decimal", SubtypeDecimal, true},

		// Non-numeric types
		{"string", SubtypeString, false},
		{"bool", SubtypeBool, false},
		{"bytes", SubtypeBytes, false},
		{"default", SubtypeDefault, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subtype.IsNumeric())
		})
	}
}

func TestAttributeSubtype_IsUnsigned(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected bool
	}{
		// Unsigned types
		{"uint", SubtypeUint, true},
		{"uint64", SubtypeUint64, true},

		// Non-unsigned types
		{"int", SubtypeInt, false},
		{"float32", SubtypeFloat32, false},
		{"string", SubtypeString, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subtype.IsUnsigned())
		})
	}
}

func TestAttributeSubtype_IsInteger(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected bool
	}{
		// Integer types
		{"int", SubtypeInt, true},
		{"uint64", SubtypeUint64, true},

		// Non-integer types
		{"float32", SubtypeFloat32, false},
		{"string", SubtypeString, false},
		{"big_int", SubtypeBigInt, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.subtype.IsInteger())
		})
	}
}

func TestAttributeSubtype_Validate(t *testing.T) {
	tests := []struct {
		name        string
		subtype     AttributeSubtype
		dynamoType  string
		expectError bool
		errorMsg    string
	}{
		// Valid combinations for "S" (String)
		{"string_with_S", SubtypeString, "S", false, ""},
		{"time_with_S", SubtypeTime, "S", false, ""},
		{"default_with_S", SubtypeDefault, "S", false, ""},

		// Invalid combinations for "S"
		{"int_with_S", SubtypeInt, "S", true, "not compatible"},

		// Valid combinations for "N" (Number)
		{"int_with_N", SubtypeInt, "N", false, ""},
		{"uint64_with_N", SubtypeUint64, "N", false, ""},
		{"float64_with_N", SubtypeFloat64, "N", false, ""},

		// Invalid combinations for "N"
		{"string_with_N", SubtypeString, "N", true, "not compatible"},

		// Valid combinations for "B" (Boolean)
		{"bool_with_B", SubtypeBool, "B", false, ""},
		{"int_with_B", SubtypeInt, "B", false, ""},

		// Unknown DynamoDB type
		{"unknown_type", SubtypeString, "X", true, "unknown DynamoDB type"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.subtype.Validate(tt.dynamoType)

			if tt.expectError {
				assert.Error(t, err)
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestAttributeSubtype_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		subtype  AttributeSubtype
		expected string
	}{
		{"string", SubtypeString, `"string"`},
		{"uint64", SubtypeUint64, `"uint64"`},
		{"big_int", SubtypeBigInt, `"big_int"`},
		{"default", SubtypeDefault, `"default"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.subtype)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, string(data))
		})
	}
}

func TestAttributeSubtype_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected AttributeSubtype
	}{
		{"string", `"string"`, SubtypeString},
		{"uint64", `"uint64"`, SubtypeUint64},
		{"big_int", `"big_int"`, SubtypeBigInt},
		{"unknown", `"unknown"`, SubtypeDefault},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var subtype AttributeSubtype
			err := json.Unmarshal([]byte(tt.json), &subtype)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, subtype)
		})
	}
}

func TestAttributeSubtype_JSONRoundTrip(t *testing.T) {
	subtypes := []AttributeSubtype{
		SubtypeDefault, SubtypeString, SubtypeInt, SubtypeUint64,
		SubtypeFloat32, SubtypeBigInt, SubtypeBool,
	}

	for _, original := range subtypes {
		t.Run(original.String(), func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(original)
			require.NoError(t, err)

			// Unmarshal from JSON
			var unmarshaled AttributeSubtype
			err = json.Unmarshal(data, &unmarshaled)
			require.NoError(t, err)

			// Should be equal
			assert.Equal(t, original, unmarshaled)
		})
	}
}
