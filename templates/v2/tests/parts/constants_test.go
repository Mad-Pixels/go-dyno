package parts

import (
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestConstantsTemplate validates the generation of DynamoDB table and attribute constants.
// This template generates compile-time safe constants for table names, indexes, and columns.
func TestConstantsTemplate(t *testing.T) {
	// Test with simple table - basic constants generation
	// Example: user table with minimal indexes and attributes
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
			SecondaryIndexes: []common.SecondaryIndex{
				{
					Name:           "EmailIndex",
					HashKey:        "email",
					ProjectionType: "KEYS_ONLY",
				},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.ConstantsTemplate, templateMap)
		testConstantsContent(t, rendered, templateMap)
	})

	// Test with complex table - multiple indexes with different projection types
	// Example: order table with composite keys and various projection strategies
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
				{Name: "user_id", Type: "S"},
			},
			CommonAttributes: []common.Attribute{
				{Name: "updated_at", Type: "N"},
				{Name: "total_amount", Type: "N"},
			},
			AllAttributes: []common.Attribute{
				{Name: "order_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "status", Type: "S"},
				{Name: "user_id", Type: "S"},
				{Name: "updated_at", Type: "N"},
				{Name: "total_amount", Type: "N"},
			},
			SecondaryIndexes: []common.SecondaryIndex{
				{
					Name:           "StatusIndex",
					HashKey:        "status",
					RangeKey:       "created_at",
					ProjectionType: "ALL",
				},
				{
					Name:           "UserOrdersIndex",
					HashKey:        "user_id",
					RangeKey:       "created_at",
					ProjectionType: "KEYS_ONLY",
				},
				{
					Name:             "StatusSummaryIndex",
					HashKey:          "status",
					ProjectionType:   "INCLUDE",
					NonKeyAttributes: []string{"total_amount", "user_id"},
				},
			},
		}
		rendered := utils.MustParseTemplateToString(v2.ConstantsTemplate, templateMap)
		testConstantsContent(t, rendered, templateMap)
	})
}

// testConstantsContent validates the content of rendered constants template
func testConstantsContent(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	// Test that generated code has valid Go syntax
	// Example: should parse without errors as valid Go constants and variables
	t.Run("go_syntax_valid", func(t *testing.T) {
		testCode := "package test\n\n" + rendered + "\n"
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Generated constants should be valid Go syntax")
	})

	// Test that TableName constant is properly defined
	// Example: const TableName = "UserTable"
	t.Run("table_name_constant_present", func(t *testing.T) {
		expectedTableName := "TableName = \"" + templateMap.TableName + "\""
		assert.Contains(t, rendered, expectedTableName,
			"Should contain TableName constant with correct value")

		assert.Contains(t, rendered, "const (",
			"TableName should be declared in a const block")
	})

	// Test that all secondary index constants are generated
	// Example: const IndexEmailIndex = "EmailIndex", const IndexStatusIndex = "StatusIndex"
	t.Run("index_constants_present", func(t *testing.T) {
		for _, idx := range templateMap.SecondaryIndexes {
			expectedIndexConstant := "Index" + idx.Name + " = \"" + idx.Name + "\""
			assert.Contains(t, rendered, expectedIndexConstant,
				"Should contain index constant for: %s", idx.Name)
		}
	})

	// Test that all attribute column constants are generated with safe names
	// Example: const ColumnUserId = "user_id", const ColumnCreatedAt = "created_at"
	t.Run("column_constants_present", func(t *testing.T) {
		for _, attr := range templateMap.AllAttributes {
			safeName := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			expectedColumnConstant := "Column" + safeName + " = \"" + attr.Name + "\""
			assert.Contains(t, rendered, expectedColumnConstant,
				"Should contain column constant for attribute: %s", attr.Name)
		}
	})

	// Test that AttributeNames slice contains all table attributes
	// Example: AttributeNames = []string{"user_id", "email", "created_at"}
	t.Run("attribute_names_slice_complete", func(t *testing.T) {
		assert.Contains(t, rendered, "AttributeNames = []string{",
			"Should declare AttributeNames slice")

		for _, attr := range templateMap.AllAttributes {
			expectedEntry := "\"" + attr.Name + "\""
			assert.Contains(t, rendered, expectedEntry,
				"AttributeNames should contain: %s", attr.Name)
		}
	})

	// Test that IndexProjections map is properly structured
	// Example: IndexProjections = map[string][]string{"EmailIndex": {...}}
	t.Run("index_projections_map_structure", func(t *testing.T) {
		assert.Contains(t, rendered, "IndexProjections = map[string][]string{",
			"Should declare IndexProjections map")

		for _, idx := range templateMap.SecondaryIndexes {
			expectedMapKey := "\"" + idx.Name + "\": {"
			assert.Contains(t, rendered, expectedMapKey,
				"IndexProjections should contain entry for index: %s", idx.Name)
		}
	})

	// Test that index projections contain correct attributes based on projection type
	// Example: ALL projection includes all attributes, KEYS_ONLY includes only keys
	t.Run("index_projections_content_correct", func(t *testing.T) {
		for _, idx := range templateMap.SecondaryIndexes {
			switch idx.ProjectionType {
			case "ALL":
				for _, attr := range templateMap.AllAttributes {
					expectedAttr := "\"" + attr.Name + "\""
					assert.Contains(t, rendered, expectedAttr,
						"Index %s with ALL projection should include attribute: %s",
						idx.Name, attr.Name)
				}

			case "KEYS_ONLY":
				expectedHashKey := "\"" + idx.HashKey + "\""
				assert.Contains(t, rendered, expectedHashKey,
					"Index %s with KEYS_ONLY should include hash key: %s",
					idx.Name, idx.HashKey)

				if idx.RangeKey != "" {
					expectedRangeKey := "\"" + idx.RangeKey + "\""
					assert.Contains(t, rendered, expectedRangeKey,
						"Index %s with KEYS_ONLY should include range key: %s",
						idx.Name, idx.RangeKey)
				}

			case "INCLUDE":
				expectedHashKey := "\"" + idx.HashKey + "\""
				assert.Contains(t, rendered, expectedHashKey,
					"Index %s with INCLUDE should include hash key: %s",
					idx.Name, idx.HashKey)

				if idx.RangeKey != "" {
					expectedRangeKey := "\"" + idx.RangeKey + "\""
					assert.Contains(t, rendered, expectedRangeKey,
						"Index %s with INCLUDE should include range key: %s",
						idx.Name, idx.RangeKey)
				}

				for _, nonKeyAttr := range idx.NonKeyAttributes {
					expectedAttr := "\"" + nonKeyAttr + "\""
					assert.Contains(t, rendered, expectedAttr,
						"Index %s with INCLUDE should include non-key attribute: %s",
						idx.Name, nonKeyAttr)
				}
			}
		}
	})

	// Test that constant names follow Go naming conventions
	// Example: ColumnUserId not column_user_id, IndexEmailIndex not indexEmailIndex
	t.Run("constant_naming_conventions", func(t *testing.T) {
		for _, idx := range templateMap.SecondaryIndexes {
			expectedPattern := "Index" + idx.Name + " ="
			assert.Contains(t, rendered, expectedPattern,
				"Index constant should follow IndexXxx naming: %s", idx.Name)
		}
		for _, attr := range templateMap.AllAttributes {
			safeName := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			expectedPattern := "Column" + safeName + " ="
			assert.Contains(t, rendered, expectedPattern,
				"Column constant should follow ColumnXxx naming for: %s", attr.Name)
		}
	})

	// Test that no duplicate constants are generated
	// Example: ensure ColumnUserId appears only once, not multiple times
	t.Run("no_duplicate_constants", func(t *testing.T) {
		lines := strings.Split(rendered, "\n")
		constDeclarations := make(map[string]int)

		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(trimmed, " = ") && !strings.HasPrefix(trimmed, "//") {
				parts := strings.Split(trimmed, " = ")
				if len(parts) >= 2 {
					constName := strings.TrimSpace(parts[0])
					constDeclarations[constName]++
				}
			}
		}
		for constName, count := range constDeclarations {
			assert.Equal(t, 1, count,
				"Constant %s should be declared exactly once, found %d times",
				constName, count)
		}
	})
}

// TestConstantsTemplateFormatting validates that the constants template follows Go formatting standards.
// This ensures the generated code is gofmt‐compliant and doesn't need additional formatting.
func TestConstantsTemplateFormatting(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "testtable",
		TableName:   "TestTable",
		HashKey:     "id",
		RangeKey:    "sort_key",
		Attributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "status", Type: "S"},
		},
		CommonAttributes: []common.Attribute{
			{Name: "created_at", Type: "N"},
			{Name: "updated_at", Type: "N"},
		},
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "status", Type: "S"},
			{Name: "created_at", Type: "N"},
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
				Name:             "StatusSummaryIndex",
				HashKey:          "status",
				ProjectionType:   "INCLUDE",
				NonKeyAttributes: []string{"updated_at"},
			},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.ConstantsTemplate, templateMap)
	fullCode := "package test\n\n" + rendered + "\n"

	if _, err := format.Source([]byte(fullCode)); err != nil {
		t.Fatalf("Generated code is not gofmt-compliant: %v", err)
	}
}
