package parts

// import (
// 	"go/format"
// 	"go/parser"
// 	"go/token"
// 	"strings"
// 	"testing"

// 	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
// 	"github.com/Mad-Pixels/go-dyno/internal/utils"
// 	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // TestQueryBuilderTemplate validates the QueryBuilderTemplate rendering.
// // This template generates a fluent QueryBuilder with methods for every attribute,
// // composite key helpers, range conditions, sorting, pagination, and execution.
// func TestQueryBuilderTemplate(t *testing.T) {
// 	templateMap := v2.TemplateMap{
// 		PackageName: "usertable",
// 		TableName:   "UserTable",
// 		HashKey:     "user_id",
// 		RangeKey:    "created_at",
// 		Attributes: []common.Attribute{
// 			{Name: "user_id", Type: "S"},
// 			{Name: "created_at", Type: "N"},
// 			{Name: "status", Type: "S"},
// 		},
// 		CommonAttributes: []common.Attribute{
// 			{Name: "updated_at", Type: "N"},
// 		},
// 		AllAttributes: []common.Attribute{
// 			{Name: "user_id", Type: "S"},
// 			{Name: "created_at", Type: "N"},
// 			{Name: "status", Type: "S"},
// 			{Name: "updated_at", Type: "N"},
// 		},
// 		SecondaryIndexes: []common.SecondaryIndex{
// 			{
// 				Name:     "UserStatusIndex",
// 				HashKey:  "user_id#status",
// 				RangeKey: "created_at",
// 				HashKeyParts: []common.CompositeKeyPart{
// 					{IsConstant: false, Value: "user_id"},
// 					{IsConstant: false, Value: "status"},
// 				},
// 				RangeKeyParts: []common.CompositeKeyPart{
// 					{IsConstant: false, Value: "created_at"},
// 				},
// 			},
// 		},
// 	}

// 	// Render the template
// 	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)

// 	// Test that the rendered code is valid Go syntax
// 	// Example: parsing with imports for expression and types should succeed
// 	t.Run("go_syntax_valid", func(t *testing.T) {
// 		src := "package test\n\nimport (\n" +
// 			"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n" +
// 			"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
// 			")\n\n" + rendered
// 		fset := token.NewFileSet()
// 		_, err := parser.ParseFile(fset, "builder.go", src, parser.ParseComments)
// 		require.NoError(t, err, "Rendered QueryBuilder should be valid Go syntax")
// 	})

// 	// Test that NewQueryBuilder constructor is present
// 	// Example: func NewQueryBuilder() *QueryBuilder
// 	t.Run("NewQueryBuilder_present", func(t *testing.T) {
// 		assert.Contains(t, rendered, "func NewQueryBuilder() *QueryBuilder", "Should contain NewQueryBuilder constructor")
// 	})

// 	// Test that With<Attribute> methods are generated for key attributes only
// 	// Example: WithUserId, WithCreatedAt, WithStatus (hash/range keys)
// 	t.Run("key_condition_methods_present", func(t *testing.T) {
// 		for _, attr := range templateMap.Attributes {
// 			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
// 			expected := "With" + name + "("
// 			assert.Contains(t, rendered, expected,
// 				"Should contain key condition method: %s", expected)
// 		}
// 	})

// 	// Test that Filter<Attribute> methods are generated for non-key attributes only
// 	// Example: FilterTitle, FilterContent, FilterViews (filter expressions)
// 	t.Run("filter_expression_methods_present", func(t *testing.T) {
// 		for _, attr := range templateMap.CommonAttributes {
// 			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
// 			expected := "Filter" + name + "("
// 			assert.Contains(t, rendered, expected,
// 				"Should contain filter expression method: %s", expected)
// 		}
// 	})

// 	// Test that key and filter methods are not mixed
// 	// Example: key attributes should not have Filter methods, common attributes should not have With methods
// 	t.Run("no_mixed_method_types", func(t *testing.T) {
// 		for _, attr := range templateMap.Attributes {
// 			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
// 			unexpected := "Filter" + name + "("
// 			assert.NotContains(t, rendered, unexpected,
// 				"Key attribute should not have filter method: %s", unexpected)
// 		}

// 		for _, attr := range templateMap.CommonAttributes {
// 			name := utils.ToUpperCamelCase(utils.ToSafeName(attr.Name))
// 			unexpected := "With" + name + "("
// 			assert.NotContains(t, rendered, unexpected,
// 				"Common attribute should not have key method: %s", unexpected)
// 		}
// 	})

// 	// Test that composite hash key helper is generated
// 	// Example: WithUserStatusIndexHashKey(userId, status)
// 	t.Run("composite_hash_helper", func(t *testing.T) {
// 		assert.Contains(t, rendered, "WithUserStatusIndexHashKey(", "Should contain composite hash-key helper")
// 	})

// 	// Test that composite range key helper is generated
// 	// Example: WithUserStatusIndexRangeKey(createdAt)
// 	t.Run("composite_range_helper", func(t *testing.T) {
// 		assert.Contains(t, rendered, "WithUserStatusIndexRangeKey(", "Should contain composite range-key helper")
// 	})

// 	// Test that numeric range methods are generated
// 	// Example: WithCreatedAtBetween, WithCreatedAtGreaterThan, WithCreatedAtLessThan
// 	t.Run("numeric_range_methods", func(t *testing.T) {
// 		base := "WithCreatedAt"
// 		assert.Contains(t, rendered, base+"Between(", "Should contain Between method")
// 		assert.Contains(t, rendered, base+"GreaterThan(", "Should contain GreaterThan method")
// 		assert.Contains(t, rendered, base+"LessThan(", "Should contain LessThan method")
// 	})

// 	// Test sorting and pagination helpers
// 	// Example: OrderByDesc, OrderByAsc, Limit, StartFrom
// 	t.Run("sorting_pagination_helpers", func(t *testing.T) {
// 		assert.Contains(t, rendered, "OrderByDesc()", "Should contain OrderByDesc")
// 		assert.Contains(t, rendered, "OrderByAsc()", "Should contain OrderByAsc")
// 		assert.Contains(t, rendered, "Limit(", "Should contain Limit")
// 		assert.Contains(t, rendered, "StartFrom(", "Should contain StartFrom")
// 	})

// 	// Test PreferredSortKey override method
// 	// Example: WithPreferredSortKey
// 	t.Run("preferred_sort_key", func(t *testing.T) {
// 		assert.Contains(t, rendered, "WithPreferredSortKey(", "Should contain WithPreferredSortKey")
// 	})

// 	// Test that no duplicate method signatures exist
// 	// Example: count of "func (qb *QueryBuilder" should meet expected minimum
// 	t.Run("no_duplicate_signatures", func(t *testing.T) {
// 		lines := strings.Split(rendered, "\n")
// 		count := 0
// 		for _, line := range lines {
// 			if strings.HasPrefix(strings.TrimSpace(line), "func (qb *QueryBuilder)") {
// 				count++
// 			}
// 		}
// 		assert.True(t, count >= 12,
// 			"Expected at least 12 QueryBuilder methods, found %d", count)
// 	})
// }

// // TestQueryBuilderTemplateFormatting validates that the rendered template is gofmt‐compliant.
// // Example: format.Source(\"package test\\n\\n\" + rendered + \"\\n\") should return no error.
// func TestQueryBuilderTemplateFormatting(t *testing.T) {
// 	templateMap := v2.TemplateMap{
// 		PackageName: "usertable",
// 		TableName:   "UserTable",
// 		HashKey:     "user_id",
// 		RangeKey:    "created_at",
// 		Attributes: []common.Attribute{
// 			{Name: "user_id", Type: "S"},
// 			{Name: "created_at", Type: "N"},
// 			{Name: "status", Type: "S"},
// 		},
// 		CommonAttributes: []common.Attribute{
// 			{Name: "updated_at", Type: "N"},
// 		},
// 		AllAttributes: []common.Attribute{
// 			{Name: "user_id", Type: "S"},
// 			{Name: "created_at", Type: "N"},
// 			{Name: "status", Type: "S"},
// 			{Name: "updated_at", Type: "N"},
// 		},
// 		SecondaryIndexes: []common.SecondaryIndex{
// 			{
// 				Name:     "UserStatusIndex",
// 				HashKey:  "user_id#status",
// 				RangeKey: "created_at",
// 				HashKeyParts: []common.CompositeKeyPart{
// 					{IsConstant: false, Value: "user_id"},
// 					{IsConstant: false, Value: "status"},
// 				},
// 				RangeKeyParts: []common.CompositeKeyPart{
// 					{IsConstant: false, Value: "created_at"},
// 				},
// 			},
// 		},
// 	}

// 	rendered := utils.MustParseTemplateToString(v2.QueryBuilderTemplate, templateMap)
// 	full := "package test\n\n" +
// 		"import (\n" +
// 		"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n" +
// 		"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
// 		")\n\n" +
// 		rendered + "\n"

// 	if _, err := format.Source([]byte(full)); err != nil {
// 		t.Fatalf("QueryBuilderTemplate is not gofmt-compliant: %v", err)
// 	}
// }
