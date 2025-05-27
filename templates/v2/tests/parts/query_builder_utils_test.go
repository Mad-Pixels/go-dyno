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

// TestQueryBuilderUtilsTemplate validates the generation of QueryBuilder utility functions.
// This template produces methods for composite key handling and attribute formatting.
func TestQueryBuilderUtilsTemplate(t *testing.T) {
	// Test simple composite keys scenario
	// Example: single hash and range attribute with no composite parts
	t.Run("simple_composite_keys", func(t *testing.T) {
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
					Name:     "StatusIndex",
					HashKey:  "status",
					RangeKey: "created_at",
				},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.QueryBuilderUtilsTemplate, templateMap)
		testQueryBuilderUtilsContent(t, rendered, templateMap)
	})

	// Test complex composite keys scenario
	// Example: composite hash key with multiple parts
	t.Run("complex_composite_keys", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "ordertable",
			TableName:   "OrderTable",
			HashKey:     "order_id",
			RangeKey:    "created_at",
			Attributes: []common.Attribute{
				{Name: "order_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "user_id", Type: "S"},
				{Name: "status", Type: "S"},
				{Name: "level", Type: "S"},
				{Name: "is_public", Type: "N"},
			},
			CommonAttributes: []common.Attribute{
				{Name: "updated_at", Type: "N"},
				{Name: "total_amount", Type: "N"},
			},
			AllAttributes: []common.Attribute{
				{Name: "order_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "user_id", Type: "S"},
				{Name: "status", Type: "S"},
				{Name: "level", Type: "S"},
				{Name: "is_public", Type: "N"},
				{Name: "updated_at", Type: "N"},
				{Name: "total_amount", Type: "N"},
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
				},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.QueryBuilderUtilsTemplate, templateMap)
		testQueryBuilderUtilsContent(t, rendered, templateMap)
	})
}

// testQueryBuilderUtilsContent validates the content of rendered QueryBuilder utils template.
func testQueryBuilderUtilsContent(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	// Test that generated code has valid Go syntax
	// Example: parsing "package test\n\n" + rendered with go/parser should succeed
	t.Run("go_syntax_valid", func(t *testing.T) {
		testCode := "package test\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Generated QueryBuilder utils should be valid Go syntax")
	})

	// Test that hasAllKeys method is generated
	// Example: func (qb *QueryBuilder) hasAllKeys(parts []CompositeKeyPart) bool
	t.Run("hasAllKeys_method_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) hasAllKeys(")
	})

	// Test hasAllKeys method logic
	// Example: code iterates parts and checks qb.UsedKeys[part.Value]
	t.Run("hasAllKeys_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "for _, part := range parts")
		assert.Contains(t, rendered, "qb.UsedKeys[part.Value]")
		assert.Contains(t, rendered, "return false")
		assert.Contains(t, rendered, "return true")
	})

	// Test that buildCompositeKeyCondition method is generated
	// Example: func (qb *QueryBuilder) buildCompositeKeyCondition(key string, value interface{}) expression.ConditionBuilder
	t.Run("buildCompositeKeyCondition_method_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) buildCompositeKeyCondition(")
	})

	// Test buildCompositeKeyCondition type handling
	// Example: switch v := value.(type); case string: case int: case bool:
	t.Run("buildCompositeKeyCondition_type_handling", func(t *testing.T) {
		assert.Contains(t, rendered, "switch v := value.(type)")
		assert.Contains(t, rendered, "case string:")
		assert.Contains(t, rendered, "case int:")
		assert.Contains(t, rendered, "case bool:")
	})

	// Test buildCompositeKeyValue optimizations
	// Example: switch len(parts); case 1:, case 2:, case 3:
	t.Run("buildCompositeKeyValue_optimizations", func(t *testing.T) {
		assert.Contains(t, rendered, "switch len(parts)")
		assert.Contains(t, rendered, "case 1:")
		assert.Contains(t, rendered, "case 2:")
		assert.Contains(t, rendered, "case 3:")
	})

	// Test formatAttributeValue boolean conversion
	// Example: return "1" for true, "0" for false
	t.Run("formatAttributeValue_boolean_conversion", func(t *testing.T) {
		assert.Contains(t, rendered, "case bool:")
		assert.Contains(t, rendered, "return \"1\"")
		assert.Contains(t, rendered, "return \"0\"")
	})

	// Test DynamoDB integration usage
	// Example: methods call expression.Key(...) and expression.Value(...)
	t.Run("dynamodb_integration", func(t *testing.T) {
		assert.Contains(t, rendered, "expression.Key(")
		assert.Contains(t, rendered, "expression.Value(")
	})
}

// TestQueryBuilderUtilsTemplateFormatting validates Go formatting standards.
// Example: format.Source("package test\n\n" + rendered + "\n") should return no error.
func TestQueryBuilderUtilsTemplateFormatting(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "testtable",
		TableName:   "TestTable",
		HashKey:     "id",
		RangeKey:    "sort_key",
		Attributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "user_id", Type: "S"},
			{Name: "status", Type: "S"},
		},
		CommonAttributes: []common.Attribute{
			{Name: "created_at", Type: "N"},
			{Name: "is_active", Type: "N"},
		},
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "sort_key", Type: "S"},
			{Name: "user_id", Type: "S"},
			{Name: "status", Type: "S"},
			{Name: "created_at", Type: "N"},
			{Name: "is_active", Type: "N"},
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
			},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.QueryBuilderUtilsTemplate, templateMap)

	fullCode := "package test\n\n" + rendered + "\n"
	if _, err := format.Source([]byte(fullCode)); err != nil {
		t.Fatalf("QueryBuilderUtilsTemplate is not gofmt-compliant: %v", err)
	}
}
