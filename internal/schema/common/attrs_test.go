package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAttribute_GoType(t *testing.T) {
	tests := []struct {
		name     string
		attr     Attribute
		expected string
	}{
		// Default subtype behavior
		{"string_default", Attribute{Type: "S"}, "string"},
		{"number_default", Attribute{Type: "N"}, "float64"},
		{"boolean_default", Attribute{Type: "B"}, "bool"},
		{"unknown_default", Attribute{Type: "X"}, "any"},

		// Explicit subtype behavior
		{"string_explicit", Attribute{Type: "S", Subtype: SubtypeString}, "string"},
		{"int_explicit", Attribute{Type: "N", Subtype: SubtypeInt}, "int"},
		{"uint64_explicit", Attribute{Type: "N", Subtype: SubtypeUint64}, "uint64"},
		{"big_int_explicit", Attribute{Type: "N", Subtype: SubtypeBigInt}, "*big.Int"},
		{"bool_explicit", Attribute{Type: "B", Subtype: SubtypeBool}, "bool"},
		{"time_explicit", Attribute{Type: "S", Subtype: SubtypeTime}, "time.Time"},
		{"uuid_explicit", Attribute{Type: "S", Subtype: SubtypeUUID}, "uuid.UUID"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.attr.GoType())
		})
	}
}

func TestAttribute_ZeroValue(t *testing.T) {
	tests := []struct {
		name     string
		attr     Attribute
		expected string
	}{
		// Default subtype behavior
		{"string_default", Attribute{Type: "S"}, `""`},
		{"number_default", Attribute{Type: "N"}, "0.0"},
		{"boolean_default", Attribute{Type: "B"}, "false"},
		{"unknown_default", Attribute{Type: "X"}, "nil"},

		// Explicit subtype behavior
		{"int_explicit", Attribute{Type: "N", Subtype: SubtypeInt}, "0"},
		{"uint64_explicit", Attribute{Type: "N", Subtype: SubtypeUint64}, "0"},
		{"float32_explicit", Attribute{Type: "N", Subtype: SubtypeFloat32}, "0.0"},
		{"big_int_explicit", Attribute{Type: "N", Subtype: SubtypeBigInt}, "big.NewInt(0)"},
		{"decimal_explicit", Attribute{Type: "N", Subtype: SubtypeDecimal}, "decimal.Zero"},
		{"bool_explicit", Attribute{Type: "B", Subtype: SubtypeBool}, "false"},
		{"time_explicit", Attribute{Type: "S", Subtype: SubtypeTime}, "time.Time{}"},
		{"uuid_explicit", Attribute{Type: "S", Subtype: SubtypeUUID}, "uuid.UUID{}"},
		{"bytes_explicit", Attribute{Type: "BS", Subtype: SubtypeBytes}, "nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.attr.ZeroValue())
		})
	}
}

func TestAttribute_Validate(t *testing.T) {
	tests := []struct {
		name        string
		attr        Attribute
		expectError bool
		errorMsg    string
	}{
		// Valid attributes
		{
			name:        "valid_string_default",
			attr:        Attribute{Name: "username", Type: "S"},
			expectError: false,
		},
		{
			name:        "valid_number_with_int_subtype",
			attr:        Attribute{Name: "count", Type: "N", Subtype: SubtypeInt},
			expectError: false,
		},
		{
			name:        "valid_number_with_uint64_subtype",
			attr:        Attribute{Name: "user_id", Type: "N", Subtype: SubtypeUint64},
			expectError: false,
		},
		{
			name:        "valid_number_with_big_int_subtype",
			attr:        Attribute{Name: "price", Type: "N", Subtype: SubtypeBigInt},
			expectError: false,
		},
		{
			name:        "valid_boolean_default",
			attr:        Attribute{Name: "is_active", Type: "B"},
			expectError: false,
		},
		{
			name:        "valid_boolean_with_bool_subtype",
			attr:        Attribute{Name: "is_premium", Type: "B", Subtype: SubtypeBool},
			expectError: false,
		},
		{
			name:        "valid_boolean_with_int_subtype",
			attr:        Attribute{Name: "status", Type: "B", Subtype: SubtypeInt},
			expectError: false,
		},
		{
			name:        "valid_string_with_time_subtype",
			attr:        Attribute{Name: "created_at", Type: "S", Subtype: SubtypeTime},
			expectError: false,
		},
		{
			name:        "valid_string_with_uuid_subtype",
			attr:        Attribute{Name: "id", Type: "S", Subtype: SubtypeUUID},
			expectError: false,
		},

		// Invalid attributes - empty name
		{
			name:        "empty_name",
			attr:        Attribute{Name: "", Type: "S"},
			expectError: true,
			errorMsg:    "attribute name cannot be empty",
		},

		// Invalid attributes - invalid DynamoDB type
		{
			name:        "invalid_dynamodb_type",
			attr:        Attribute{Name: "test", Type: "INVALID"},
			expectError: true,
			errorMsg:    "invalid DynamoDB type 'INVALID'",
		},

		// Invalid attributes - incompatible subtype with DynamoDB type
		{
			name:        "string_subtype_with_number_type",
			attr:        Attribute{Name: "invalid", Type: "N", Subtype: SubtypeString},
			expectError: true,
			errorMsg:    "not compatible with DynamoDB type 'N'",
		},
		{
			name:        "numeric_subtype_with_string_type",
			attr:        Attribute{Name: "invalid", Type: "S", Subtype: SubtypeInt},
			expectError: true,
			errorMsg:    "not compatible with DynamoDB type 'S'",
		},
		{
			name:        "float_subtype_with_boolean_type",
			attr:        Attribute{Name: "invalid", Type: "B", Subtype: SubtypeFloat64},
			expectError: true,
			errorMsg:    "not compatible with DynamoDB type 'B'",
		},
		{
			name:        "bytes_subtype_with_string_type",
			attr:        Attribute{Name: "invalid", Type: "S", Subtype: SubtypeBytes},
			expectError: true,
			errorMsg:    "not compatible with DynamoDB type 'S'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.attr.Validate()

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

func TestAttribute_ValidateAllDynamoDBTypes(t *testing.T) {
	validTypes := []string{"S", "N", "B", "BOOL", "BS", "SS", "NS", "L", "M", "NULL"}

	for _, dynamoType := range validTypes {
		t.Run("type_"+dynamoType, func(t *testing.T) {
			attr := Attribute{
				Name: "test_attr",
				Type: dynamoType,
			}

			err := attr.Validate()
			assert.NoError(t, err, "DynamoDB type %s should be valid", dynamoType)
		})
	}
}

func TestAttribute_SubtypeCompatibilityMatrix(t *testing.T) {
	// Test compatibility matrix between DynamoDB types and subtypes
	compatibilityTests := []struct {
		dynamoType string
		subtype    AttributeSubtype
		compatible bool
	}{
		// String type compatibility
		{"S", SubtypeString, true},
		{"S", SubtypeTime, true},
		{"S", SubtypeUUID, true},
		{"S", SubtypeInt, false},
		{"S", SubtypeBool, false},
		{"S", SubtypeBytes, false},

		// Number type compatibility
		{"N", SubtypeInt, true},
		{"N", SubtypeInt8, true},
		{"N", SubtypeInt16, true},
		{"N", SubtypeInt32, true},
		{"N", SubtypeInt64, true},
		{"N", SubtypeUint, true},
		{"N", SubtypeUint8, true},
		{"N", SubtypeUint16, true},
		{"N", SubtypeUint32, true},
		{"N", SubtypeUint64, true},
		{"N", SubtypeFloat32, true},
		{"N", SubtypeFloat64, true},
		{"N", SubtypeBigInt, true},
		{"N", SubtypeDecimal, true},
		{"N", SubtypeString, false},
		{"N", SubtypeBool, false},
		{"N", SubtypeBytes, false},

		// Boolean type compatibility (stored as Number in DynamoDB)
		{"B", SubtypeBool, true},
		{"B", SubtypeInt, true},
		{"B", SubtypeUint32, true},
		{"B", SubtypeFloat64, false},
		{"B", SubtypeString, false},
		{"B", SubtypeBytes, false},
	}

	for _, tt := range compatibilityTests {
		testName := tt.dynamoType + "_with_" + tt.subtype.String()
		t.Run(testName, func(t *testing.T) {
			attr := Attribute{
				Name:    "test",
				Type:    tt.dynamoType,
				Subtype: tt.subtype,
			}

			err := attr.Validate()
			if tt.compatible {
				assert.NoError(t, err, "Should be compatible: %s with %s", tt.dynamoType, tt.subtype.String())
			} else {
				assert.Error(t, err, "Should be incompatible: %s with %s", tt.dynamoType, tt.subtype.String())
			}
		})
	}
}

func TestAttribute_EdgeCases(t *testing.T) {
	t.Run("whitespace_name", func(t *testing.T) {
		attr := Attribute{Name: "   ", Type: "S"}
		err := attr.Validate()
		assert.NoError(t, err)
	})

	t.Run("unicode_name", func(t *testing.T) {
		attr := Attribute{Name: "用户名", Type: "S"}
		err := attr.Validate()
		assert.NoError(t, err)
	})

	t.Run("very_long_name", func(t *testing.T) {
		longName := string(make([]byte, 1000))
		for i := range longName {
			longName = longName[:i] + "a" + longName[i+1:]
		}
		attr := Attribute{Name: longName, Type: "S"}
		err := attr.Validate()
		assert.NoError(t, err)
	})

	t.Run("case_sensitive_type", func(t *testing.T) {
		attr := Attribute{Name: "test", Type: "s"} // lowercase
		err := attr.Validate()
		assert.Error(t, err) // Should fail because DynamoDB types are case-sensitive
	})
}

// Example usage test
func TestAttribute_RealWorldExamples(t *testing.T) {
	examples := []struct {
		name string
		attr Attribute
		desc string
	}{
		{
			name: "user_id_as_uint64",
			attr: Attribute{Name: "user_id", Type: "N", Subtype: SubtypeUint64},
			desc: "User ID stored as unsigned 64-bit integer",
		},
		{
			name: "price_as_big_int",
			attr: Attribute{Name: "price_cents", Type: "N", Subtype: SubtypeBigInt},
			desc: "Price in cents using arbitrary precision",
		},
		{
			name: "timestamp_as_time",
			attr: Attribute{Name: "created_at", Type: "S", Subtype: SubtypeTime},
			desc: "ISO timestamp stored as string, parsed as time.Time",
		},
		{
			name: "uuid_identifier",
			attr: Attribute{Name: "request_id", Type: "S", Subtype: SubtypeUUID},
			desc: "UUID stored as string",
		},
		{
			name: "feature_flag",
			attr: Attribute{Name: "is_premium", Type: "B", Subtype: SubtypeBool},
			desc: "Boolean feature flag",
		},
		{
			name: "counter",
			attr: Attribute{Name: "view_count", Type: "N", Subtype: SubtypeUint32},
			desc: "View counter as unsigned 32-bit integer",
		},
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			err := example.attr.Validate()
			require.NoError(t, err, "Real-world example should be valid: %s", example.desc)

			goType := example.attr.GoType()
			zeroValue := example.attr.ZeroValue()

			assert.NotEmpty(t, goType, "Should have a Go type")
			assert.NotEmpty(t, zeroValue, "Should have a zero value")

			t.Logf("Example: %s", example.desc)
			t.Logf("  Attribute: %+v", example.attr)
			t.Logf("  Go Type: %s", goType)
			t.Logf("  Zero Value: %s", zeroValue)
		})
	}
}
