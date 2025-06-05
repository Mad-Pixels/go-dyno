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

// TestSchemaStructsTemplate validates the generation of DynamoDB schema structures and item types.
// This template generates Go structures that represent DynamoDB table items and schema metadata
// with full support for AttributeSubtype system.
func TestSchemaStructsTemplate(t *testing.T) {
	// Test with simple table - basic structure without subtypes
	t.Run("simple_table_default_types", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "usertable",
			TableName:   "UserTable",
			HashKey:     "user_id",
			RangeKey:    "created_at",
			Attributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "is_active", Type: "B"},
			},
			CommonAttributes: []common.Attribute{
				{Name: "updated_at", Type: "N"},
			},
			AllAttributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "is_active", Type: "B"},
				{Name: "updated_at", Type: "N"},
			},
			SecondaryIndexes: []common.SecondaryIndex{},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
		testBasicStructure(t, rendered, templateMap)
		testDefaultGoTypes(t, rendered, templateMap)
	})

	// Test with complex table using explicit subtypes
	t.Run("complex_table_with_subtypes", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "ordertable",
			TableName:   "OrderTable",
			HashKey:     "order_id",
			RangeKey:    "created_at",
			Attributes: []common.Attribute{
				{Name: "order_id", Type: "S", Subtype: common.SubtypeString},
				{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
				{Name: "user_id", Type: "N", Subtype: common.SubtypeUint64},
				{Name: "status_code", Type: "N", Subtype: common.SubtypeUint8},
				{Name: "is_paid", Type: "B", Subtype: common.SubtypeBool},
			},
			CommonAttributes: []common.Attribute{
				{Name: "price_cents", Type: "N", Subtype: common.SubtypeBigInt},
				{Name: "discount_rate", Type: "N", Subtype: common.SubtypeFloat32},
				{Name: "request_id", Type: "S", Subtype: common.SubtypeUUID},
				{Name: "created_time", Type: "S", Subtype: common.SubtypeTime},
			},
			AllAttributes: []common.Attribute{
				{Name: "order_id", Type: "S", Subtype: common.SubtypeString},
				{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
				{Name: "user_id", Type: "N", Subtype: common.SubtypeUint64},
				{Name: "status_code", Type: "N", Subtype: common.SubtypeUint8},
				{Name: "is_paid", Type: "B", Subtype: common.SubtypeBool},
				{Name: "price_cents", Type: "N", Subtype: common.SubtypeBigInt},
				{Name: "discount_rate", Type: "N", Subtype: common.SubtypeFloat32},
				{Name: "request_id", Type: "S", Subtype: common.SubtypeUUID},
				{Name: "created_time", Type: "S", Subtype: common.SubtypeTime},
			},
			SecondaryIndexes: []common.SecondaryIndex{
				{
					Name:           "UserOrdersIndex",
					HashKey:        "user_id",
					RangeKey:       "created_at",
					ProjectionType: "ALL",
				},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
		testBasicStructure(t, rendered, templateMap)
		testSubtypeGoTypes(t, rendered, templateMap)
		testSecondaryIndexes(t, rendered, templateMap)
	})

	// Test with mixed subtypes and defaults
	t.Run("mixed_subtypes_and_defaults", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "mixedtable",
			TableName:   "MixedTable",
			HashKey:     "id",
			RangeKey:    "",
			Attributes: []common.Attribute{
				{Name: "id", Type: "S"},
				{Name: "count", Type: "N", Subtype: common.SubtypeInt},
				{Name: "score", Type: "N"},
				{Name: "enabled", Type: "B", Subtype: common.SubtypeBool},
			},
			CommonAttributes: []common.Attribute{
				{Name: "timestamp", Type: "N", Subtype: common.SubtypeUint64},
				{Name: "metadata", Type: "S"},
			},
			AllAttributes: []common.Attribute{
				{Name: "id", Type: "S"},
				{Name: "count", Type: "N", Subtype: common.SubtypeInt},
				{Name: "score", Type: "N"},
				{Name: "enabled", Type: "B", Subtype: common.SubtypeBool},
				{Name: "timestamp", Type: "N", Subtype: common.SubtypeUint64},
				{Name: "metadata", Type: "S"},
			},
			SecondaryIndexes: []common.SecondaryIndex{},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
		testBasicStructure(t, rendered, templateMap)
		testMixedTypes(t, rendered)
	})

	// Test edge cases with uncommon subtypes
	t.Run("edge_case_subtypes", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "edgetable",
			TableName:   "EdgeTable",
			HashKey:     "id",
			RangeKey:    "",
			Attributes: []common.Attribute{
				{Name: "id", Type: "S", Subtype: common.SubtypeUUID},
				{Name: "data", Type: "BS", Subtype: common.SubtypeBytes},
				{Name: "created_at", Type: "S", Subtype: common.SubtypeTime},
			},
			CommonAttributes: []common.Attribute{
				{Name: "balance", Type: "N", Subtype: common.SubtypeDecimal},
				{Name: "tiny_num", Type: "N", Subtype: common.SubtypeInt8},
				{Name: "huge_num", Type: "N", Subtype: common.SubtypeBigInt},
			},
			AllAttributes: []common.Attribute{
				{Name: "id", Type: "S", Subtype: common.SubtypeUUID},
				{Name: "data", Type: "BS", Subtype: common.SubtypeBytes},
				{Name: "created_at", Type: "S", Subtype: common.SubtypeTime},
				{Name: "balance", Type: "N", Subtype: common.SubtypeDecimal},
				{Name: "tiny_num", Type: "N", Subtype: common.SubtypeInt8},
				{Name: "huge_num", Type: "N", Subtype: common.SubtypeBigInt},
			},
			SecondaryIndexes: []common.SecondaryIndex{},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
		testBasicStructure(t, rendered, templateMap)
		testEdgeCaseTypes(t, rendered, templateMap)
	})
}

// testBasicStructure validates the basic structure of generated schema structs
func testBasicStructure(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	t.Run("go_syntax_valid", func(t *testing.T) {
		testCode := "package test\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Generated schema structs should be valid Go syntax")
	})

	t.Run("required_structs_present", func(t *testing.T) {
		requiredStructs := []string{
			"type DynamoSchema struct",
			"type Attribute struct",
			"type CompositeKeyPart struct",
			"type SecondaryIndex struct",
			"type SchemaItem struct",
		}

		for _, structDef := range requiredStructs {
			assert.Contains(t, rendered, structDef,
				"Should contain required struct definition: %s", structDef)
		}
	})

	t.Run("schema_item_has_all_attributes", func(t *testing.T) {
		for _, attr := range templateMap.AllAttributes {
			expectedField := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			expectedTag := "`dynamodbav:\"" + attr.Name + "\"`"

			assert.Contains(t, rendered, expectedField,
				"SchemaItem should contain field for attribute: %s", attr.Name)
			assert.Contains(t, rendered, expectedTag,
				"SchemaItem should contain dynamodb tag for attribute: %s", attr.Name)
		}
	})

	t.Run("table_schema_properly_initialized", func(t *testing.T) {
		assert.Contains(t, rendered, "var TableSchema = DynamoSchema{",
			"Should declare TableSchema variable")

		assert.Contains(t, rendered, "TableName: \""+templateMap.TableName+"\"",
			"TableSchema should have correct table name")

		assert.Contains(t, rendered, "HashKey:   \""+templateMap.HashKey+"\"",
			"TableSchema should have correct hash key")

		assert.Contains(t, rendered, "RangeKey:  \""+templateMap.RangeKey+"\"",
			"TableSchema should have correct range key")
	})
}

// testDefaultGoTypes validates default Go type mappings (without explicit subtypes)
func testDefaultGoTypes(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	t.Run("default_go_types_for_dynamodb_types", func(t *testing.T) {
		defaultTypeMapping := map[string]string{
			"S": "string",
			"N": "float64",
			"B": "bool",
		}

		for _, attr := range templateMap.AllAttributes {
			// Only test attributes without explicit subtypes
			if attr.Subtype == common.SubtypeDefault {
				expectedType, exists := defaultTypeMapping[attr.Type]
				if exists {
					fieldName := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
					expectedField := fieldName + " " + expectedType

					assert.Contains(t, rendered, expectedField,
						"Field %s should have default Go type %s for DynamoDB type %s",
						fieldName, expectedType, attr.Type)
				}
			}
		}
	})
}

// testSubtypeGoTypes validates Go types with explicit subtypes
func testSubtypeGoTypes(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	t.Run("explicit_subtype_go_types", func(t *testing.T) {
		for _, attr := range templateMap.AllAttributes {
			if attr.Subtype != common.SubtypeDefault {
				fieldName := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
				expectedType := attr.GoType()
				expectedField := fieldName + " " + expectedType

				assert.Contains(t, rendered, expectedField,
					"Field %s should have subtype Go type %s (subtype: %s, DynamoDB type: %s)",
					fieldName, expectedType, attr.Subtype.String(), attr.Type)
			}
		}
	})

	t.Run("specific_subtype_validations", func(t *testing.T) {
		subtypeTests := map[string]string{
			"OrderId":      "string",    // SubtypeString
			"CreatedAt":    "int64",     // SubtypeInt64
			"UserId":       "uint64",    // SubtypeUint64
			"StatusCode":   "uint8",     // SubtypeUint8
			"IsPaid":       "bool",      // SubtypeBool
			"PriceCents":   "*big.Int",  // SubtypeBigInt
			"DiscountRate": "float32",   // SubtypeFloat32
			"RequestId":    "uuid.UUID", // SubtypeUUID
			"CreatedTime":  "time.Time", // SubtypeTime
		}

		for fieldName, expectedType := range subtypeTests {
			expectedField := fieldName + " " + expectedType
			assert.Contains(t, rendered, expectedField,
				"Should contain field declaration: %s", expectedField)
		}
	})
}

// testMixedTypes validates mixing of default and explicit subtypes
func testMixedTypes(t *testing.T, rendered string) {
	t.Helper()

	t.Run("mixed_default_and_explicit_types", func(t *testing.T) {
		expectedFields := map[string]string{
			"Id":        "string",  // default S -> string
			"Count":     "int",     // explicit SubtypeInt
			"Score":     "float64", // default N -> float64
			"Enabled":   "bool",    // explicit SubtypeBool
			"Timestamp": "uint64",  // explicit SubtypeUint64
			"Metadata":  "string",  // default S -> string
		}

		for fieldName, expectedType := range expectedFields {
			expectedField := fieldName + " " + expectedType
			assert.Contains(t, rendered, expectedField,
				"Mixed types: should contain field declaration: %s", expectedField)
		}
	})
}

// testEdgeCaseTypes validates uncommon but valid subtypes
func testEdgeCaseTypes(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	t.Run("edge_case_subtype_mappings", func(t *testing.T) {
		edgeCaseFields := map[string]string{
			"Id":        "uuid.UUID",        // SubtypeUUID
			"Data":      "[]byte",           // SubtypeBytes
			"CreatedAt": "time.Time",        // SubtypeTime
			"Balance":   "*decimal.Decimal", // SubtypeDecimal
			"TinyNum":   "int8",             // SubtypeInt8
			"HugeNum":   "*big.Int",         // SubtypeBigInt
		}

		for fieldName, expectedType := range edgeCaseFields {
			expectedField := fieldName + " " + expectedType
			assert.Contains(t, rendered, expectedField,
				"Edge case: should contain field declaration: %s", expectedField)
		}
	})

	t.Run("special_type_comments", func(t *testing.T) {
		// Validate that comments reflect the correct types
		commentTests := []struct {
			field   string
			comment string
		}{
			{"Id", "uuid.UUID"},
			{"Data", "[]byte"},
			{"CreatedAt", "time.Time"},
			{"Balance", "*decimal.Decimal"},
			{"HugeNum", "*big.Int"},
		}

		for _, test := range commentTests {
			expectedComment := "// DynamoDB type: " + findAttributeType(templateMap.AllAttributes, test.field) + " -> Go type: " + test.comment
			assert.Contains(t, rendered, expectedComment,
				"Should contain correct type comment for %s", test.field)
		}
	})
}

// testSecondaryIndexes validates secondary index generation
func testSecondaryIndexes(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	t.Run("secondary_indexes_structure", func(t *testing.T) {
		for _, idx := range templateMap.SecondaryIndexes {
			assert.Contains(t, rendered, "Name:           \""+idx.Name+"\"",
				"TableSchema.SecondaryIndexes should contain index: %s", idx.Name)

			assert.Contains(t, rendered, "HashKey:        \""+idx.HashKey+"\"",
				"Index %s should have correct hash key", idx.Name)

			if idx.RangeKey != "" {
				assert.Contains(t, rendered, "RangeKey:       \""+idx.RangeKey+"\"",
					"Index %s should have correct range key", idx.Name)
			}

			assert.Contains(t, rendered, "ProjectionType: \""+idx.ProjectionType+"\"",
				"Index %s should have correct projection type", idx.Name)
		}
	})
}

// TestSchemaStructsTemplateFormatting validates that the schema structs template follows Go formatting standards
func TestSchemaStructsTemplateFormatting(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "testtable",
		TableName:   "TestTable",
		HashKey:     "id",
		RangeKey:    "sort_key",
		Attributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "count", Type: "N", Subtype: common.SubtypeUint32},
		},
		CommonAttributes: []common.Attribute{
			{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
		},
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "count", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
		},
		SecondaryIndexes: []common.SecondaryIndex{
			{
				Name:           "CreatedAtIndex",
				HashKey:        "created_at",
				ProjectionType: "KEYS_ONLY",
			},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
	fullCode := "package test\n\n" + rendered + "\n"

	_, err := format.Source([]byte(fullCode))
	assert.NoError(t, err, "SchemaStructsTemplate should be gofmt-compliant")
}

// TestAttributeSubtypeIntegration validates that AttributeSubtype system works end-to-end
func TestAttributeSubtypeIntegration(t *testing.T) {
	t.Run("comprehensive_subtype_integration", func(t *testing.T) {
		// Test all supported subtypes in one template
		allSubtypes := []common.Attribute{
			// String subtypes
			{Name: "name", Type: "S", Subtype: common.SubtypeString},
			{Name: "uuid_field", Type: "S", Subtype: common.SubtypeUUID},
			{Name: "time_field", Type: "S", Subtype: common.SubtypeTime},

			// Signed integer subtypes
			{Name: "tiny_int", Type: "N", Subtype: common.SubtypeInt8},
			{Name: "small_int", Type: "N", Subtype: common.SubtypeInt16},
			{Name: "medium_int", Type: "N", Subtype: common.SubtypeInt32},
			{Name: "big_int_std", Type: "N", Subtype: common.SubtypeInt64},
			{Name: "default_int", Type: "N", Subtype: common.SubtypeInt},

			// Unsigned integer subtypes
			{Name: "tiny_uint", Type: "N", Subtype: common.SubtypeUint8},
			{Name: "small_uint", Type: "N", Subtype: common.SubtypeUint16},
			{Name: "medium_uint", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "big_uint", Type: "N", Subtype: common.SubtypeUint64},
			{Name: "default_uint", Type: "N", Subtype: common.SubtypeUint},

			// Floating point subtypes
			{Name: "float_32", Type: "N", Subtype: common.SubtypeFloat32},
			{Name: "float_64", Type: "N", Subtype: common.SubtypeFloat64},

			// Arbitrary precision subtypes
			{Name: "big_integer", Type: "N", Subtype: common.SubtypeBigInt},
			{Name: "decimal_num", Type: "N", Subtype: common.SubtypeDecimal},

			// Boolean subtypes
			{Name: "flag", Type: "B", Subtype: common.SubtypeBool},

			// Binary subtypes
			{Name: "binary_data", Type: "BS", Subtype: common.SubtypeBytes},
		}

		templateMap := v2.TemplateMap{
			PackageName:      "comprehensive",
			TableName:        "ComprehensiveTable",
			HashKey:          "name",
			RangeKey:         "",
			Attributes:       allSubtypes[:10],
			CommonAttributes: allSubtypes[10:],
			AllAttributes:    allSubtypes,
			SecondaryIndexes: []common.SecondaryIndex{},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)

		// Validate that all subtypes are correctly mapped
		expectedMappings := map[string]string{
			"Name":        "string",
			"UuidField":   "uuid.UUID",
			"TimeField":   "time.Time",
			"TinyInt":     "int8",
			"SmallInt":    "int16",
			"MediumInt":   "int32",
			"BigIntStd":   "int64",
			"DefaultInt":  "int",
			"TinyUint":    "uint8",
			"SmallUint":   "uint16",
			"MediumUint":  "uint32",
			"BigUint":     "uint64",
			"DefaultUint": "uint",
			"Float32":     "float32",
			"Float64":     "float64",
			"BigInteger":  "*big.Int",
			"DecimalNum":  "*decimal.Decimal",
			"Flag":        "bool",
			"BinaryData":  "[]byte",
		}

		for fieldName, expectedType := range expectedMappings {
			expectedField := fieldName + " " + expectedType
			assert.Contains(t, rendered, expectedField,
				"Comprehensive test: should contain field %s with type %s", fieldName, expectedType)
		}

		// Validate syntax is still correct with all types
		testCode := "package test\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Comprehensive subtype test should generate valid Go syntax")
	})
}

// Helper function to find attribute type by field name
func findAttributeType(attrs []common.Attribute, fieldName string) string {
	targetName := utils.ToLowerCamelCase(fieldName)
	for _, attr := range attrs {
		if utils.ToLowerCamelCase(utils.ToSafeName(attr.Name)) == targetName {
			return attr.Type
		}
	}
	return "UNKNOWN"
}
