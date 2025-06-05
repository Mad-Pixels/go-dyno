package utils

import (
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/stretchr/testify/assert"
)

func TestToSafeName(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "   -Привет-123_foo-",
			expected: "x123_foo",
		},
		{
			input:    "$$$",
			expected: "xxx",
		},
		{
			input:    "9lives",
			expected: "x9lives",
		},
		{
			input:    "Hello-World",
			expected: "Hello_World",
		},
		{
			input:    "_foo_bar_",
			expected: "foo_bar",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ToSafeName(tt.input), "input: %q", tt.input)
	}
}

func TestToUpperCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "   -Привет-123_foo-",
			expected: "X123Foo",
		},
		{
			input:    "$$$",
			expected: "Xxx",
		},
		{
			input:    "9lives",
			expected: "X9lives",
		},
		{
			input:    "Hello-World",
			expected: "HelloWorld",
		},
		{
			input:    "_foo_bar_",
			expected: "FooBar",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ToUpperCamelCase(tt.input), "input: %q", tt.input)
	}
}

func TestToLowerCase(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    "   -Привет-123_foo-",
			expected: "x123Foo",
		},
		{
			input:    "$$$",
			expected: "xxx",
		},
		{
			input:    "9lives",
			expected: "x9lives",
		},
		{
			input:    "Hello-World",
			expected: "helloWorld",
		},
		{
			input:    "_foo_bar_",
			expected: "fooBar",
		},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.expected, ToLowerCamelCase(tt.input), "input: %q", tt.input)
	}
}

func TestToGolangZeroType_WithSubtypes(t *testing.T) {
	tests := []struct {
		name     string
		attr     common.Attribute
		expected string
	}{
		// Default behavior
		{"string_default", common.Attribute{Type: "S"}, `""`},
		{"number_default", common.Attribute{Type: "N"}, "0.0"},
		{"boolean_default", common.Attribute{Type: "B"}, "false"},

		// Subtype behavior
		{"int_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeInt}, "0"},
		{"uint64_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeUint64}, "0"},
		{"float32_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeFloat32}, "0.0"},
		{"big_int_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeBigInt}, "big.NewInt(0)"},
		{"decimal_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeDecimal}, "decimal.Zero"},
		{"bool_subtype", common.Attribute{Type: "B", Subtype: common.SubtypeBool}, "false"},
		{"time_subtype", common.Attribute{Type: "S", Subtype: common.SubtypeTime}, "time.Time{}"},
		{"uuid_subtype", common.Attribute{Type: "S", Subtype: common.SubtypeUUID}, "uuid.UUID{}"},
		{"bytes_subtype", common.Attribute{Type: "BS", Subtype: common.SubtypeBytes}, "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGolangZeroType(tt.attr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestToGolangBaseType_WithSubtypes(t *testing.T) {
	tests := []struct {
		name     string
		attr     common.Attribute
		expected string
	}{
		// Default subtype behavior (existing behavior should remain)
		{"string_default", common.Attribute{Type: "S"}, "string"},
		{"number_default", common.Attribute{Type: "N"}, "float64"},
		{"boolean_default", common.Attribute{Type: "B"}, "bool"},
		{"unknown_default", common.Attribute{Type: "UNKNOWN"}, "any"},

		// Signed integer subtypes
		{"int_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeInt}, "int"},
		{"int8_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeInt8}, "int8"},
		{"int16_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeInt16}, "int16"},
		{"int32_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeInt32}, "int32"},
		{"int64_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeInt64}, "int64"},

		// Unsigned integer subtypes
		{"uint_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeUint}, "uint"},
		{"uint8_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeUint8}, "uint8"},
		{"uint16_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeUint16}, "uint16"},
		{"uint32_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeUint32}, "uint32"},
		{"uint64_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeUint64}, "uint64"},

		// Floating point subtypes
		{"float32_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeFloat32}, "float32"},
		{"float64_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeFloat64}, "float64"},

		// Arbitrary precision subtypes
		{"big_int_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeBigInt}, "*big.Int"},
		{"decimal_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeDecimal}, "*decimal.Decimal"},

		// String subtypes
		{"string_explicit", common.Attribute{Type: "S", Subtype: common.SubtypeString}, "string"},
		{"time_subtype", common.Attribute{Type: "S", Subtype: common.SubtypeTime}, "time.Time"},
		{"uuid_subtype", common.Attribute{Type: "S", Subtype: common.SubtypeUUID}, "uuid.UUID"},

		// Boolean subtypes
		{"bool_explicit", common.Attribute{Type: "B", Subtype: common.SubtypeBool}, "bool"},

		// Binary subtypes
		{"bytes_subtype", common.Attribute{Type: "BS", Subtype: common.SubtypeBytes}, "[]byte"},

		// Mixed scenarios - DynamoDB type vs subtype
		{"number_type_with_string_subtype", common.Attribute{Type: "N", Subtype: common.SubtypeString}, "string"},
		{"string_type_with_time_subtype", common.Attribute{Type: "S", Subtype: common.SubtypeTime}, "time.Time"},
		{"boolean_type_with_int_subtype", common.Attribute{Type: "B", Subtype: common.SubtypeInt}, "int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGolangBaseType(tt.attr)
			assert.Equal(t, tt.expected, result,
				"Attribute{Type: %q, Subtype: %s} should return Go type %q",
				tt.attr.Type, tt.attr.Subtype.String(), tt.expected)
		})
	}
}

func TestToGolangBaseType_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		attr     common.Attribute
		expected string
		note     string
	}{
		{
			name:     "empty_type",
			attr:     common.Attribute{Type: ""},
			expected: "any",
			note:     "Empty type should default to 'any'",
		},
		{
			name:     "subtype_default_explicit",
			attr:     common.Attribute{Type: "N", Subtype: common.SubtypeDefault},
			expected: "float64",
			note:     "Explicit SubtypeDefault should use DynamoDB type mapping",
		},
		{
			name:     "subtype_zero_value",
			attr:     common.Attribute{Type: "S", Subtype: 0},
			expected: "string",
			note:     "Zero value subtype should use DynamoDB type mapping",
		},
		{
			name:     "case_sensitive_type",
			attr:     common.Attribute{Type: "s"},
			expected: "any",
			note:     "DynamoDB types are case-sensitive",
		},
		{
			name:     "numeric_type_with_explicit_string_subtype",
			attr:     common.Attribute{Type: "N", Subtype: common.SubtypeString},
			expected: "string",
			note:     "Subtype overrides DynamoDB type (even if incompatible)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToGolangBaseType(tt.attr)
			assert.Equal(t, tt.expected, result, tt.note)
		})
	}
}
