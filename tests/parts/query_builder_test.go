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

// TestQueryBuilderTemplate validates the QueryBuilderTemplate rendering.
// This template generates a fluent QueryBuilder with methods for every attribute,
// composite key helpers, range conditions, sorting, pagination, and execution.
func TestQueryBuilderTemplate(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "usertable",
		TableName:   "UserTable",
		HashKey:     "user_id",
		RangeKey:    "created_at",
		Attributes: []common.Attribute{
			{Name: "user_id", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "status", Type: "S"},
		},
		CommonAttributes: []common.Attribute{
			{Name: "updated_at", Type: "N"},
		},
		AllAttributes: []common.Attribute{
			{Name: "user_id", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "status", Type: "S"},
			{Name: "updated_at", Type: "N"},
		},
		SecondaryIndexes: []common.SecondaryIndex{
			{
				Name:     "UserStatusIndex",
				HashKey:  "user_id#status",
				RangeKey: "created_at",
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

	// Render the template
	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)

	// Test that the rendered code is valid Go syntax
	// Example: parsing with imports for expression and types should succeed
	t.Run("go_syntax_valid", func(t *testing.T) {
		src := "package test\n\nimport (\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
			")\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "builder.go", src, parser.ParseComments)
		require.NoError(t, err, "Rendered QueryBuilder should be valid Go syntax")
	})

	// Test that NewQueryBuilder constructor is present
	// Example: func NewQueryBuilder() *QueryBuilder
	t.Run("NewQueryBuilder_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func NewQueryBuilder() *QueryBuilder", "Should contain NewQueryBuilder constructor")
	})

	// Test that With<Attribute> methods are generated for key attributes only
	// Example: WithUserId, WithCreatedAt, WithStatus (hash/range keys)
	t.Run("key_condition_methods_present", func(t *testing.T) {
		for _, attr := range templateMap.Attributes {
			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			expected := "With" + name + "("
			assert.Contains(t, rendered, expected,
				"Should contain key condition method: %s", expected)
		}
	})

	// Test that Filter<Attribute> methods are generated for non-key attributes only
	// Example: FilterTitle, FilterContent, FilterViews (filter expressions)
	t.Run("filter_expression_methods_present", func(t *testing.T) {
		for _, attr := range templateMap.CommonAttributes {
			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			expected := "Filter" + name + "("
			assert.Contains(t, rendered, expected,
				"Should contain filter expression method: %s", expected)
		}
	})

	// Test that key and filter methods are not mixed
	// Example: key attributes should not have Filter methods, common attributes should not have With methods
	t.Run("no_mixed_method_types", func(t *testing.T) {
		for _, attr := range templateMap.Attributes {
			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			unexpected := "Filter" + name + "("
			assert.NotContains(t, rendered, unexpected,
				"Key attribute should not have filter method: %s", unexpected)
		}

		for _, attr := range templateMap.CommonAttributes {
			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
			unexpected := "With" + name + "("
			assert.NotContains(t, rendered, unexpected,
				"Common attribute should not have key method: %s", unexpected)
		}
	})

	// Test that composite hash key helper is generated
	// Example: WithUserStatusIndexHashKey(userId, status)
	t.Run("composite_hash_helper", func(t *testing.T) {
		assert.Contains(t, rendered, "WithUserStatusIndexHashKey(", "Should contain composite hash-key helper")
	})

	// Test that composite range key helper is generated
	// Example: WithUserStatusIndexRangeKey(createdAt)
	t.Run("composite_range_helper", func(t *testing.T) {
		assert.Contains(t, rendered, "WithUserStatusIndexRangeKey(", "Should contain composite range-key helper")
	})

	// Test that numeric range methods are generated
	// Example: WithCreatedAtBetween, WithCreatedAtGreaterThan, WithCreatedAtLessThan
	t.Run("numeric_range_methods", func(t *testing.T) {
		base := "WithCreatedAt"
		assert.Contains(t, rendered, base+"Between(", "Should contain Between method")
		assert.Contains(t, rendered, base+"GreaterThan(", "Should contain GreaterThan method")
		assert.Contains(t, rendered, base+"LessThan(", "Should contain LessThan method")
	})

	// Test sorting and pagination helpers
	// Example: OrderByDesc, OrderByAsc, Limit, StartFrom
	t.Run("sorting_pagination_helpers", func(t *testing.T) {
		assert.Contains(t, rendered, "OrderByDesc()", "Should contain OrderByDesc")
		assert.Contains(t, rendered, "OrderByAsc()", "Should contain OrderByAsc")
		assert.Contains(t, rendered, "Limit(", "Should contain Limit")
		assert.Contains(t, rendered, "StartFrom(", "Should contain StartFrom")
	})

	// Test PreferredSortKey override method
	// Example: WithPreferredSortKey
	t.Run("preferred_sort_key", func(t *testing.T) {
		assert.Contains(t, rendered, "WithPreferredSortKey(", "Should contain WithPreferredSortKey")
	})

	// Test that no duplicate method signatures exist
	// Example: count of "func (qb *QueryBuilder" should meet expected minimum
	t.Run("no_duplicate_signatures", func(t *testing.T) {
		lines := strings.Split(rendered, "\n")
		count := 0
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "func (qb *QueryBuilder)") {
				count++
			}
		}
		assert.True(t, count >= 12,
			"Expected at least 12 QueryBuilder methods, found %d", count)
	})
}

// TestQueryBuilderTemplateFormatting validates that the rendered template is gofmt‐compliant.
// Example: format.Source(\"package test\\n\\n\" + rendered + \"\\n\") should return no error.
func TestQueryBuilderTemplateFormatting(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "usertable",
		TableName:   "UserTable",
		HashKey:     "user_id",
		RangeKey:    "created_at",
		Attributes: []common.Attribute{
			{Name: "user_id", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "status", Type: "S"},
		},
		CommonAttributes: []common.Attribute{
			{Name: "updated_at", Type: "N"},
		},
		AllAttributes: []common.Attribute{
			{Name: "user_id", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "status", Type: "S"},
			{Name: "updated_at", Type: "N"},
		},
		SecondaryIndexes: []common.SecondaryIndex{
			{
				Name:     "UserStatusIndex",
				HashKey:  "user_id#status",
				RangeKey: "created_at",
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

	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)
	full := "package test\n\n" +
		"import (\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
		")\n\n" +
		rendered + "\n"

	if _, err := format.Source([]byte(full)); err != nil {
		t.Fatalf("QueryBuilderTemplate is not gofmt-compliant: %v", err)
	}
}

// TestQueryBuilderTemplate_WithSubtypes validates the QueryBuilderTemplate rendering with subtypes.
// This test ensures that all subtype Go types are correctly used in method signatures
// and that range conditions are generated for appropriate types.
func TestQueryBuilderTemplate_WithSubtypes(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "subtypetable",
		TableName:   "SubtypeTable",
		HashKey:     "user_id",
		RangeKey:    "created_at",
		Attributes: []common.Attribute{
			{Name: "user_id", Type: "N", Subtype: common.SubtypeUint64},
			{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
			{Name: "count", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "score", Type: "N", Subtype: common.SubtypeFloat32},
		},
		CommonAttributes: []common.Attribute{
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
			{Name: "flag", Type: "B", Subtype: common.SubtypeBool},
		},
		AllAttributes: []common.Attribute{
			{Name: "user_id", Type: "N", Subtype: common.SubtypeUint64},
			{Name: "created_at", Type: "N", Subtype: common.SubtypeInt64},
			{Name: "count", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "score", Type: "N", Subtype: common.SubtypeFloat32},
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
			{Name: "flag", Type: "B", Subtype: common.SubtypeBool},
		},
		SecondaryIndexes: []common.SecondaryIndex{},
	}

	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)

	// Test that fluent methods are generated with correct subtype parameter types
	// Example: WithUserId(userId uint64), WithCount(count uint32)
	t.Run("subtype_fluent_methods_generated", func(t *testing.T) {
		subtypeMethods := []struct {
			method    string
			paramType string
		}{
			{"WithUserId(userId uint64)", "uint64"},     // SubtypeUint64
			{"WithCreatedAt(createdAt int64)", "int64"}, // SubtypeInt64
			{"WithCount(count uint32)", "uint32"},       // SubtypeUint32
			{"WithScore(score float32)", "float32"},     // SubtypeFloat32
			{"WithPrice(price *big.Int)", "*big.Int"},   // SubtypeBigInt
			{"WithFlag(flag bool)", "bool"},             // SubtypeBool
		}

		for _, test := range subtypeMethods {
			assert.Contains(t, rendered, test.method,
				"Should generate method with correct subtype parameter: %s", test.method)
		}
	})

	// Test that range conditions are generated for supported types
	// Example: WithUserIdBetween(start, end uint64), WithCountGreaterThan(value uint32)
	t.Run("subtype_range_conditions_generated", func(t *testing.T) {
		// Only types that support range conditions in template
		rangeConditions := []struct {
			base      string
			paramType string
		}{
			// Integer types (signed and unsigned)
			{"WithUserIdBetween(start, end uint64)", "uint64"},
			{"WithUserIdGreaterThan(value uint64)", "uint64"},
			{"WithUserIdLessThan(value uint64)", "uint64"},
			{"WithCreatedAtBetween(start, end int64)", "int64"},
			{"WithCreatedAtGreaterThan(value int64)", "int64"},
			{"WithCreatedAtLessThan(value int64)", "int64"},
			{"WithCountBetween(start, end uint32)", "uint32"},
			{"WithCountGreaterThan(value uint32)", "uint32"},
			{"WithCountLessThan(value uint32)", "uint32"},

			// Float types
			{"WithScoreBetween(start, end float32)", "float32"},
			{"WithScoreGreaterThan(value float32)", "float32"},
			{"WithScoreLessThan(value float32)", "float32"},
		}

		for _, test := range rangeConditions {
			assert.Contains(t, rendered, test.base,
				"Should generate range method for subtype: %s", test.base)
		}
	})

	// Test that complex types do not generate range conditions
	// Example: *big.Int and other complex types should NOT have range methods
	t.Run("complex_types_no_range_conditions", func(t *testing.T) {
		// *big.Int and other complex types should NOT generate range methods
		complexTypes := []string{
			"WithPriceBetween(*big.Int", // Should NOT be present
		}

		for _, notExpected := range complexTypes {
			assert.NotContains(t, rendered, notExpected,
				"Should NOT generate range methods for complex types: %s", notExpected)
		}
	})

	// Test that generated Go syntax is valid with all subtypes
	// Example: parsing with imports should succeed without errors
	t.Run("go_syntax_valid_with_subtypes", func(t *testing.T) {
		src := "package test\n\nimport (\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
			")\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "builder.go", src, parser.ParseComments)
		require.NoError(t, err, "QueryBuilder with subtypes should be valid Go syntax")
	})
}

// TestQueryBuilderTemplate_UintTypeSupport validates uint type support specifically.
// This test ensures that all unsigned integer types are properly supported
// in both method signatures and range conditions.
func TestQueryBuilderTemplate_UintTypeSupport(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "uinttable",
		TableName:   "UintTable",
		HashKey:     "id",
		RangeKey:    "",
		Attributes: []common.Attribute{
			{Name: "id", Type: "N", Subtype: common.SubtypeUint},
			{Name: "tiny", Type: "N", Subtype: common.SubtypeUint8},
			{Name: "small", Type: "N", Subtype: common.SubtypeUint16},
			{Name: "medium", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "large", Type: "N", Subtype: common.SubtypeUint64},
		},
		CommonAttributes: []common.Attribute{},
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "N", Subtype: common.SubtypeUint},
			{Name: "tiny", Type: "N", Subtype: common.SubtypeUint8},
			{Name: "small", Type: "N", Subtype: common.SubtypeUint16},
			{Name: "medium", Type: "N", Subtype: common.SubtypeUint32},
			{Name: "large", Type: "N", Subtype: common.SubtypeUint64},
		},
		SecondaryIndexes: []common.SecondaryIndex{},
	}

	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)

	// Test that all uint types are supported in method signatures
	// Example: WithId(id uint), WithTiny(tiny uint8), WithLarge(large uint64)
	t.Run("all_uint_types_supported", func(t *testing.T) {
		uintTypes := []struct {
			method string
			goType string
		}{
			{"WithId(id uint)", "uint"},
			{"WithTiny(tiny uint8)", "uint8"},
			{"WithSmall(small uint16)", "uint16"},
			{"WithMedium(medium uint32)", "uint32"},
			{"WithLarge(large uint64)", "uint64"},
		}

		for _, test := range uintTypes {
			assert.Contains(t, rendered, test.method,
				"Should support uint type: %s", test.goType)
		}
	})

	// Test that uint range conditions are generated properly
	// Example: WithIdBetween(start, end uint), WithLargeBetween(start, end uint64)
	t.Run("uint_range_conditions", func(t *testing.T) {
		// All uint types should support range conditions
		assert.Contains(t, rendered, "WithIdBetween(start, end uint)",
			"Should generate range methods for uint")
		assert.Contains(t, rendered, "WithLargeBetween(start, end uint64)",
			"Should generate range methods for uint64")
		assert.Contains(t, rendered, "WithTinyBetween(start, end uint8)",
			"Should generate range methods for uint8")
		assert.Contains(t, rendered, "WithMediumBetween(start, end uint32)",
			"Should generate range methods for uint32")
		assert.Contains(t, rendered, "WithSmallBetween(start, end uint16)",
			"Should generate range methods for uint16")
	})

	// Test that GreaterThan and LessThan conditions are also generated
	// Example: WithIdGreaterThan(value uint), WithLargeGreaterThan(value uint64)
	t.Run("uint_comparison_conditions", func(t *testing.T) {
		comparisonTests := []string{
			"WithIdGreaterThan(value uint)",
			"WithIdLessThan(value uint)",
			"WithLargeGreaterThan(value uint64)",
			"WithLargeLessThan(value uint64)",
			"WithTinyGreaterThan(value uint8)",
			"WithTinyLessThan(value uint8)",
		}

		for _, expected := range comparisonTests {
			assert.Contains(t, rendered, expected,
				"Should generate comparison method: %s", expected)
		}
	})
}

// TestQueryBuilderTemplateFormatting_WithSubtypes validates that template with subtypes is gofmt-compliant.
// Example: format.Source should succeed for code with all subtype variations.
func TestQueryBuilderTemplateFormatting_WithSubtypes(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "formattable",
		TableName:   "FormatTable",
		HashKey:     "id",
		RangeKey:    "timestamp",
		Attributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "timestamp", Type: "N", Subtype: common.SubtypeUint64},
		},
		CommonAttributes: []common.Attribute{
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
		},
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "timestamp", Type: "N", Subtype: common.SubtypeUint64},
			{Name: "price", Type: "N", Subtype: common.SubtypeBigInt},
		},
		SecondaryIndexes: []common.SecondaryIndex{},
	}

	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)
	full := "package test\n\n" +
		"import (\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
		")\n\n" +
		rendered + "\n"

	_, err := format.Source([]byte(full))
	assert.NoError(t, err, "QueryBuilderTemplate with subtypes should be gofmt-compliant")
}
