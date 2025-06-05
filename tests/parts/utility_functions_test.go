package parts

import (
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUtilityFunctionsTemplate validates the UtilityFunctionsTemplate rendering.
// This template generates a set of helper functions for DynamoDB operations:
// - BatchPutItems, PutItem for item marshaling
// - BoolToInt, IntToBool for boolean conversions
// - ExtractFromDynamoDBStreamEvent with comprehensive type handling
// - IsFieldModified, GetBoolFieldChanged for stream processing
// - CreateKey, CreateKeyFromItem for key generation
// - CreateTriggerHandler for Lambda event handling
// - ConvertMapToAttributeValues for dynamic conversions
func TestUtilityFunctionsTemplate(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "usertable",
		TableName:   "UserTable",
		HashKey:     "user_id",
		RangeKey:    "created_at",
		AllAttributes: []common.Attribute{
			{Name: "user_id", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "status", Type: "S"},
			{Name: "updated_at", Type: "N"},
			{Name: "is_active", Type: "B"},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

	// Test that the rendered code is valid Go syntax with all required imports
	// Example: parsing with complete import set should succeed without errors
	t.Run("go_syntax_valid", func(t *testing.T) {
		src := "package test\n\nimport (\n" +
			"\t\"context\"\n" +
			"\t\"encoding/json\"\n" +
			"\t\"fmt\"\n" +
			"\t\"strconv\"\n" +
			"\t\"time\"\n" +
			"\t\"encoding/base64\"\n" +
			"\t\"math/big\"\n" +
			"\t\"github.com/google/uuid\"\n" +
			"\t\"github.com/shopspring/decimal\"\n" +
			"\t\"github.com/aws/aws-lambda-go/events\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
			")\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "utils.go", src, parser.ParseComments)
		require.NoError(t, err, "Rendered UtilityFunctionsTemplate should be valid Go syntax")
	})

	// Test that all core utility functions are present in generated code
	// Example: BatchPutItems, PutItem, BoolToInt, etc.
	t.Run("core_functions_present", func(t *testing.T) {
		coreFunctions := []string{
			"func BatchPutItems(",
			"func PutItem(",
			"func BoolToInt(",
			"func IntToBool(",
			"func ExtractFromDynamoDBStreamEvent(",
			"func IsFieldModified(",
			"func GetBoolFieldChanged(",
			"func CreateKey(",
			"func CreateKeyFromItem(",
			"func CreateTriggerHandler(",
			"func ConvertMapToAttributeValues(",
		}
		for _, fn := range coreFunctions {
			assert.Contains(t, rendered, fn, "Should contain core function: %s", fn)
		}
	})

	// Test that ExtractFromDynamoDBStreamEvent handles default DynamoDB types
	// Example: string attributes use val.String(), number attributes use strconv.ParseFloat
	t.Run("extract_default_type_handling", func(t *testing.T) {
		// Default string handling
		assert.Contains(t, rendered, "val.String()", "Should handle string attributes with val.String()")

		// Default number handling (should use float64 as default)
		assert.Contains(t, rendered, "strconv.ParseFloat(val.Number(), 64)", "Should handle number attributes with ParseFloat")

		// Default boolean handling
		assert.Contains(t, rendered, "val.Boolean()", "Should handle boolean attributes with val.Boolean()")
	})

	// Test that ConvertMapToAttributeValues covers all basic type conversions
	// Example: string, float64, bool, nil, map, slice, and default cases
	t.Run("convert_map_type_coverage", func(t *testing.T) {
		typeConversions := []string{
			"case string:",
			"case float64:",
			"case bool:",
			"case nil:",
			"case map[string]interface{}:",
			"case []interface{}:",
			"default:",
		}
		for _, typeCase := range typeConversions {
			assert.Contains(t, rendered, typeCase, "ConvertMapToAttributeValues should handle: %s", typeCase)
		}
	})

	// Test that boolean conversion utilities work correctly
	// Example: BoolToInt converts true→1, false→0; IntToBool converts 0→false, non-zero→true
	t.Run("boolean_conversion_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "if b {\n        return 1\n    }\n    return 0",
			"BoolToInt should convert bool to int correctly")
		assert.Contains(t, rendered, "return i != 0",
			"IntToBool should convert int to bool correctly")
	})
}

// TestUtilityFunctionsTemplate_WithSubtypes validates subtype handling in utility functions.
// This test ensures that ExtractFromDynamoDBStreamEvent properly processes all subtype variations
// including signed/unsigned integers, floating point, arbitrary precision, and special types.
func TestUtilityFunctionsTemplate_WithSubtypes(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "subtypetable",
		TableName:   "SubtypeTable",
		HashKey:     "user_id",
		RangeKey:    "created_at",
		AllAttributes: []common.Attribute{
			{Name: "user_id", Type: "N", Subtype: common.SubtypeUint64},
			{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
			{Name: "tiny_count", Type: "N", Subtype: common.SubtypeInt8},
			{Name: "big_count", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
			{Name: "balance", Type: "N", Subtype: common.SubtypeDecimal},
			{Name: "score", Type: "N", Subtype: common.SubtypeFloat32},
			{Name: "rating", Type: "N", Subtype: common.SubtypeFloat64},
			{Name: "timestamp", Type: "S", Subtype: common.SubtypeTime},
			{Name: "request_id", Type: "S", Subtype: common.SubtypeUUID},
			{Name: "data", Type: "BS", Subtype: common.SubtypeBytes},
			{Name: "is_active", Type: "B", Subtype: common.SubtypeBool},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

	// Test that all integer subtype extractions are properly handled
	// Example: int8 uses ParseInt with bit size 8, uint64 uses ParseUint with bit size 64
	t.Run("integer_subtype_extraction", func(t *testing.T) {
		integerExtractions := []struct {
			subtype  string
			expected string
		}{
			{"int8", "strconv.ParseInt(val.Number(), 10, 8)"},
			{"int64", "strconv.ParseInt(val.Number(), 10, 64)"},
			{"uint32", "strconv.ParseUint(val.Number(), 10, 32)"},
			{"uint64", "strconv.ParseUint(val.Number(), 10, 64)"},
		}

		for _, test := range integerExtractions {
			assert.Contains(t, rendered, test.expected,
				"Should contain %s extraction logic", test.subtype)
		}
	})

	// Test that floating point subtype extractions are properly handled
	// Example: float32 uses ParseFloat with bit size 32, float64 uses bit size 64
	t.Run("float_subtype_extraction", func(t *testing.T) {
		floatExtractions := []struct {
			subtype  string
			expected string
		}{
			{"float32", "strconv.ParseFloat(val.Number(), 32)"},
			{"float64", "strconv.ParseFloat(val.Number(), 64)"},
		}

		for _, test := range floatExtractions {
			assert.Contains(t, rendered, test.expected,
				"Should contain %s extraction logic", test.subtype)
		}
	})

	// Test that arbitrary precision and special types are properly handled
	// Example: *big.Int uses SetString, *decimal.Decimal uses NewFromString
	t.Run("special_subtype_extraction", func(t *testing.T) {
		specialExtractions := []struct {
			subtype  string
			expected string
		}{
			{"*big.Int", "new(big.Int).SetString(val.Number(), 10)"},
			{"*decimal.Decimal", "decimal.NewFromString(val.Number())"},
			{"time.Time", "time.Parse(time.RFC3339, val.String())"},
			{"uuid.UUID", "uuid.Parse(val.String())"},
			{"[]byte", "base64.StdEncoding.DecodeString(val.String())"},
		}

		for _, test := range specialExtractions {
			assert.Contains(t, rendered, test.expected,
				"Should contain %s extraction logic", test.subtype)
		}
	})

	// Test that type casting is properly applied for smaller integer types
	// Example: int8 result should be cast to int8, uint32 result should be cast to uint32
	t.Run("integer_type_casting", func(t *testing.T) {
		typeCasts := []struct {
			subtype  string
			expected string
		}{
			{"int8", "int8(n)"},
			{"uint32", "uint32(n)"},
			{"float32", "float32(f)"},
		}

		for _, test := range typeCasts {
			assert.Contains(t, rendered, test.expected,
				"Should contain type casting for %s", test.subtype)
		}
	})
}

// TestUtilityFunctionsTemplate_DefaultTypeFallback validates fallback logic for default types.
// This test ensures that appropriate fallback branches are generated when using default types
// without explicit subtypes, and that complex DynamoDB types are handled gracefully.
func TestUtilityFunctionsTemplate_DefaultTypeFallback(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "fallbacktable",
		TableName:   "FallbackTable",
		HashKey:     "id",
		RangeKey:    "",
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},            // Default string
			{Name: "count", Type: "N"},         // Default number - should use fallback
			{Name: "active", Type: "B"},        // Default boolean
			{Name: "tags", Type: "SS"},         // String Set - complex type
			{Name: "metadata", Type: "M"},      // Map - complex type
			{Name: "items", Type: "L"},         // List - complex type
			{Name: "null_field", Type: "NULL"}, // NULL type
		},
	}

	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

	// Test that fallback logic is present for unknown numeric types
	// Example: when no specific subtype is defined, should fall back to float64
	t.Run("numeric_fallback_logic", func(t *testing.T) {
		// For default N type without subtype, should use ParseFloat as fallback
		assert.Contains(t, rendered, "strconv.ParseFloat(val.Number(), 64)",
			"Should use ParseFloat as numeric fallback for default N type")
	})

	// Test that complex DynamoDB types are handled with JSON fallback
	// Example: SS, NS, L, M types should be marshaled to JSON string
	t.Run("complex_type_handling", func(t *testing.T) {
		complexTypeHandling := []string{
			"Complex types (String Set, Number Set, List, Map)",
			"json.Marshal(val)",
			"string(jsonData)",
		}

		for _, expected := range complexTypeHandling {
			assert.Contains(t, rendered, expected,
				"Should handle complex types with JSON: %s", expected)
		}
	})

	// Test that NULL type is properly handled
	// Example: NULL type should leave field as zero value with appropriate comment
	t.Run("null_type_handling", func(t *testing.T) {
		assert.Contains(t, rendered, "NULL type - leave field as zero value",
			"Should handle NULL type appropriately")
	})

	// Test that completely unknown types have final fallback
	// Example: any unrecognized type should be stored as string representation
	t.Run("unknown_type_fallback", func(t *testing.T) {
		// Since we don't have truly unknown types in this test, we check that the fallback logic exists
		// The template should have fallback for unexpected types
		assert.Contains(t, rendered, "val.String()",
			"Should use val.String() for fallback cases")
	})
}

// TestUtilityFunctionsTemplate_StreamEventHandling validates DynamoDB Stream event processing.
// This test ensures that stream-related utilities handle various event types correctly
// and provide proper trigger handler functionality.
func TestUtilityFunctionsTemplate_StreamEventHandling(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "streamtable",
		TableName:   "StreamTable",
		HashKey:     "id",
		RangeKey:    "timestamp",
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "timestamp", Type: "N"},
			{Name: "status", Type: "S"},
			{Name: "is_processed", Type: "B"},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

	// Test that CreateTriggerHandler supports all DynamoDB stream event types
	// Example: INSERT, MODIFY, REMOVE events should be handled with appropriate callbacks
	t.Run("trigger_handler_event_types", func(t *testing.T) {
		eventTypes := []string{
			"case \"INSERT\":",
			"case \"MODIFY\":",
			"case \"REMOVE\":",
		}

		for _, eventType := range eventTypes {
			assert.Contains(t, rendered, eventType,
				"CreateTriggerHandler should handle event type: %s", eventType)
		}
	})

	// Test that field modification detection works correctly
	// Example: IsFieldModified should check MODIFY events and compare old/new images
	t.Run("field_modification_detection", func(t *testing.T) {
		modificationChecks := []string{
			"dbEvent.EventName != \"MODIFY\"",
			"dbEvent.Change.OldImage == nil || dbEvent.Change.NewImage == nil",
			"oldVal.String() != newVal.String()",
		}

		for _, check := range modificationChecks {
			assert.Contains(t, rendered, check,
				"IsFieldModified should include check: %s", check)
		}
	})

	// Test that boolean field change detection is properly implemented
	// Example: GetBoolFieldChanged should detect false→true transitions
	t.Run("boolean_field_change_detection", func(t *testing.T) {
		booleanChangeLogic := []string{
			"oldValue := false",
			"newValue := false",
			"oldVal.Boolean()",
			"newVal.Boolean()",
			"!oldValue && newValue",
		}

		for _, logic := range booleanChangeLogic {
			assert.Contains(t, rendered, logic,
				"GetBoolFieldChanged should include logic: %s", logic)
		}
	})
}

// TestUtilityFunctionsTemplateFormatting validates that the rendered template is gofmt-compliant.
// This ensures that generated code follows Go formatting standards and can be used directly
// in production without additional formatting steps.
func TestUtilityFunctionsTemplateFormatting(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "formattable",
		TableName:   "FormatTable",
		HashKey:     "user_id",
		RangeKey:    "created_at",
		AllAttributes: []common.Attribute{
			{Name: "user_id", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "status", Type: "S"},
			{Name: "updated_at", Type: "N"},
			{Name: "is_active", Type: "B"},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

	fullCode := "package test\n\nimport (\n" +
		"\t\"context\"\n" +
		"\t\"encoding/json\"\n" +
		"\t\"fmt\"\n" +
		"\t\"strconv\"\n" +
		"\t\"time\"\n" +
		"\t\"encoding/base64\"\n" +
		"\t\"math/big\"\n" +
		"\t\"github.com/google/uuid\"\n" +
		"\t\"github.com/shopspring/decimal\"\n" +
		"\t\"github.com/aws/aws-lambda-go/events\"\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue\"\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
		")\n\n" +
		rendered + "\n"

	_, err := format.Source([]byte(fullCode))
	assert.NoError(t, err, "UtilityFunctionsTemplate should be gofmt-compliant")
}
