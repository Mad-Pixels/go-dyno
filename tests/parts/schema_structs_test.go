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
// This template generates Go structs that represent DynamoDB table items and schema metadata.
func TestSchemaStructsTemplate(t *testing.T) {
	// Test with simple table - only hash key and basic attributes
	// Example: user table with id (hash key) and email
	t.Run("simple_table", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "usertable",
			TableName:   "UserTable",
			HashKey:     "user_id",
			RangeKey:    "",
			Attributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "email", Type: "S"},
			},
			CommonAttributes: []common.Attribute{
				{Name: "created_at", Type: "N"},
			},
			AllAttributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "email", Type: "S"},
				{Name: "created_at", Type: "N"},
			},
			SecondaryIndexes: []common.SecondaryIndex{},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
		testSchemaStructsContent(t, rendered, templateMap)
	})

	// Test with complex table - hash + range keys, GSI with composite keys
	// Example: order table with composite keys and secondary indexes
	t.Run("complex_table", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "ordertable",
			TableName:   "OrderTable",
			HashKey:     "order_id",
			RangeKey:    "created_at",
			Attributes: []common.Attribute{
				{Name: "order_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "status", Type: "S"},
			},
			CommonAttributes: []common.Attribute{
				{Name: "updated_at", Type: "N"},
			},
			AllAttributes: []common.Attribute{
				{Name: "order_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "status", Type: "S"},
				{Name: "updated_at", Type: "N"},
			},
			SecondaryIndexes: []common.SecondaryIndex{
				{
					Name:           "StatusIndex",
					HashKey:        "status",
					RangeKey:       "created_at",
					ProjectionType: "ALL",
				},
				{
					Name:             "UserStatusIndex",
					HashKey:          "user#status",
					RangeKey:         "created_at",
					ProjectionType:   "INCLUDE",
					NonKeyAttributes: []string{"order_total", "shipping_address"},
					HashKeyParts: []common.CompositeKeyPart{
						{IsConstant: false, Value: "user_id"},
						{IsConstant: false, Value: "status"},
					},
					RangeKeyParts: []common.CompositeKeyPart{
						{IsConstant: false, Value: "created_at"},
					},
				},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.SchemaStructsTemplate, templateMap)
		testSchemaStructsContent(t, rendered, templateMap)
	})
}

// testSchemaStructsContent validates the content of rendered schema structs template
func testSchemaStructsContent(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	// Test that generated code has valid Go syntax
	// Example: should parse without errors when combined with package declaration
	t.Run("go_syntax_valid", func(t *testing.T) {
		testCode := "package test\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Generated schema structs should be valid Go syntax")
	})

	// Test that all required struct types are generated
	// Example: DynamoSchema, Attribute, SecondaryIndex, SchemaItem, CompositeKeyPart
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

	// Test that SchemaItem struct contains all table attributes as fields
	// Example: UserId string `dynamodbav:"user_id"`, Email string `dynamodbav:"email"`
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

	// Test that TableSchema variable is properly initialized with template values
	// Example: TableName: "UserTable", HashKey: "user_id", RangeKey: "created_at"
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

	// Test that all table attributes are included in TableSchema.Attributes array
	// Example: {Name: "user_id", Type: "S"}, {Name: "email", Type: "S"}
	t.Run("table_schema_attributes_complete", func(t *testing.T) {
		for _, attr := range templateMap.Attributes {
			expectedAttr := "{Name: \"" + attr.Name + "\", Type: \"" + attr.Type + "\"}"
			assert.Contains(t, rendered, expectedAttr,
				"TableSchema.Attributes should contain: %s", expectedAttr)
		}
	})

	// Test that all common attributes are included in TableSchema.CommonAttributes array
	// Example: {Name: "created_at", Type: "N"}, {Name: "updated_at", Type: "N"}
	t.Run("table_schema_common_attributes_complete", func(t *testing.T) {
		for _, attr := range templateMap.CommonAttributes {
			expectedAttr := "{Name: \"" + attr.Name + "\", Type: \"" + attr.Type + "\"}"
			assert.Contains(t, rendered, expectedAttr,
				"TableSchema.CommonAttributes should contain: %s", expectedAttr)
		}
	})

	// Test that all secondary indexes are properly defined in TableSchema.SecondaryIndexes
	// Example: Name: "StatusIndex", HashKey: "status", RangeKey: "created_at"
	t.Run("table_schema_secondary_indexes_complete", func(t *testing.T) {
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

	// Test that Go field types match DynamoDB attribute types correctly
	// Example: "S" -> string, "N" -> int, "B" -> bool
	t.Run("correct_go_types_for_dynamodb_types", func(t *testing.T) {
		typeMapping := map[string]string{
			"S": "string",
			"N": "int",
			"B": "bool",
		}

		for _, attr := range templateMap.AllAttributes {
			expectedType, exists := typeMapping[attr.Type]
			if exists {
				fieldName := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
				expectedField := fieldName + " " + expectedType
				assert.Contains(t, rendered, expectedField,
					"Field %s should have Go type %s for DynamoDB type %s",
					fieldName, expectedType, attr.Type)
			}
		}
	})
}

// TestSchemaStructsTemplateFormatting validates that the schema structs template follows Go formatting standards.
// This ensures the generated code is gofmt-compliant and doesn't need additional formatting.
func TestSchemaStructsTemplateFormatting(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "testtable",
		TableName:   "TestTable",
		HashKey:     "id",
		RangeKey:    "sort_key",
		Attributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
		},
		CommonAttributes: []common.Attribute{
			{Name: "created_at", Type: "N"},
		},
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "created_at", Type: "N"},
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

	if _, err := format.Source([]byte(fullCode)); err != nil {
		t.Fatalf("SchemaStructsTemplate is not gofmt-compliant: %v", err)
	}
}
