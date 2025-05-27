package parts

import (
	"go/parser"
	"go/token"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	"github.com/Mad-Pixels/go-dyno/templates/test"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQueryBuilderUtilsTemplate validates the generation of QueryBuilder utility functions.
func TestQueryBuilderUtilsTemplate(t *testing.T) {
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

		rendered := utils.MustParseTemplateFormattedToString(v2.QueryBuilderUtilsTemplate, templateMap)
		testQueryBuilderUtilsContent(t, rendered, templateMap)
	})
}

func testQueryBuilderUtilsContent(t *testing.T, rendered string, templateMap v2.TemplateMap) {
	t.Helper()

	// Test that generated code has valid Go syntax
	t.Run("go_syntax_valid", func(t *testing.T) {
		testCode := "package test\n\nimport (\n\t\"fmt\"\n\t\"strconv\"\n\t\"strings\"\n\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression\"\n)\n\ntype QueryBuilder struct{\n\tUsedKeys map[string]bool\n\tAttributes map[string]interface{}\n}\ntype CompositeKeyPart struct {\n\tIsConstant bool\n\tValue string\n}\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Generated QueryBuilder utils should be valid Go syntax")
	})

	// Test that all required utility methods are generated
	t.Run("required_methods_present", func(t *testing.T) {
		requiredMethods := []string{
			"func (qb *QueryBuilder) hasAllKeys(",
			"func (qb *QueryBuilder) buildCompositeKeyCondition(",
			"func (qb *QueryBuilder) getCompositeKeyName(",
			"func (qb *QueryBuilder) buildCompositeKeyValue(",
			"func (qb *QueryBuilder) formatAttributeValue(",
		}

		for _, method := range requiredMethods {
			assert.Contains(t, rendered, method,
				"Should contain required utility method: %s", method)
		}
	})

	// Test hasAllKeys method logic
	t.Run("hasAllKeys_method_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "for _, part := range parts")
		assert.Contains(t, rendered, "!part.IsConstant")
		assert.Contains(t, rendered, "qb.UsedKeys[part.Value]")
		assert.Contains(t, rendered, "return false")
		assert.Contains(t, rendered, "return true")
	})

	// Test buildCompositeKeyCondition handles different types
	t.Run("buildCompositeKeyCondition_type_handling", func(t *testing.T) {
		assert.Contains(t, rendered, "switch v := value.(type)")
		assert.Contains(t, rendered, "case string:")
		assert.Contains(t, rendered, "case int:")
		assert.Contains(t, rendered, "case bool:")
		assert.Contains(t, rendered, "strconv.Itoa(v)")
		assert.Contains(t, rendered, "expression.Key(")
	})

	// Test getCompositeKeyName optimization for different cases
	t.Run("getCompositeKeyName_optimization", func(t *testing.T) {
		assert.Contains(t, rendered, "switch len(parts)")
		assert.Contains(t, rendered, "case 0:")
		assert.Contains(t, rendered, "case 1:")
		assert.Contains(t, rendered, "strings.Join(")
		assert.Contains(t, rendered, "parts[0].Value")
	})

	// Test buildCompositeKeyValue has optimization cases
	t.Run("buildCompositeKeyValue_optimizations", func(t *testing.T) {
		assert.Contains(t, rendered, "switch len(parts)")
		assert.Contains(t, rendered, "case 1:")
		assert.Contains(t, rendered, "case 2:")
		assert.Contains(t, rendered, "case 3:")
		assert.Contains(t, rendered, "formatAttributeValue(")
		assert.Contains(t, rendered, "strings.Builder")
	})

	// Test formatAttributeValue boolean handling
	t.Run("formatAttributeValue_boolean_conversion", func(t *testing.T) {
		assert.Contains(t, rendered, "case bool:")
		assert.Contains(t, rendered, "return \"1\"")
		assert.Contains(t, rendered, "return \"0\"")
	})

	// Test error handling and edge cases
	t.Run("error_handling_edge_cases", func(t *testing.T) {
		assert.Contains(t, rendered, "case 0:")
		assert.Contains(t, rendered, "return \"\"")
		assert.Contains(t, rendered, "default:")
		assert.Contains(t, rendered, "fmt.Sprintf(")
	})

	// Test DynamoDB integration
	t.Run("dynamodb_integration", func(t *testing.T) {
		assert.Contains(t, rendered, "expression.Key(")
		assert.Contains(t, rendered, "expression.Value(")
		assert.Contains(t, rendered, "compositeKeyName")
	})
}

// TestQueryBuilderUtilsTemplateFormatting validates Go formatting standards
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

	rendered := utils.MustParseTemplateFormattedToString(v2.QueryBuilderUtilsTemplate, templateMap)

	// Create test file with proper import order (standard libs first, then external)
	testCode := `package test

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
)

type QueryBuilder struct {
	UsedKeys   map[string]bool
	Attributes map[string]interface{}
}

type CompositeKeyPart struct {
	IsConstant bool
	Value      string
}
` + rendered + `
`

	test.TestAllFormattersUnchanged(t, testCode)
}
