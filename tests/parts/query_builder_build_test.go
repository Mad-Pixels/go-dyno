package parts

import (
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestQueryBuilderBuildTemplate validates the QueryBuilderBuildTemplate content.
// This template implements the core Build, BuildQuery and Execute logic for QueryBuilder.
func TestQueryBuilderBuildTemplate(t *testing.T) {
	rendered := v2.QueryBuilderBuildTemplate

	// Test that the template parses as valid Go syntax
	// Example: parsing "package test\n\n<rendered>" should succeed without errors
	t.Run("go_syntax_valid", func(t *testing.T) {
		src := "package test\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "build.go", src, parser.ParseComments)
		require.NoError(t, err, "QueryBuilderBuildTemplate should be valid Go syntax")
	})

	// Test that Build method is present
	// Example: func (qb *QueryBuilder) Build() (...)
	t.Run("Build_method_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) Build()", "Should contain Build method signature")
	})

	// Test that Build uses sort.Slice and calculateIndexParts
	// Example: sort.Slice(sortedIndexes, ..., qb.calculateIndexParts(...)
	t.Run("Build_sort_and_priority", func(t *testing.T) {
		assert.Contains(t, rendered, "sort.Slice(sortedIndexes", "Build must use sort.Slice")
		assert.Contains(t, rendered, "qb.calculateIndexParts", "Build must call calculateIndexParts")
	})

	// Test that calculateIndexParts method is present and counts parts
	// Example: func (qb *QueryBuilder) calculateIndexParts(idx SecondaryIndex) int
	t.Run("calculateIndexParts_method_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) calculateIndexParts(", "Should contain calculateIndexParts method")
		assert.Contains(t, rendered, "len(idx.HashKeyParts)", "Should count len(idx.HashKeyParts)")
		assert.Contains(t, rendered, "len(idx.RangeKeyParts)", "Should count len(idx.RangeKeyParts)")
	})

	// Test that buildHashKeyCondition handles composite and simple keys
	// Example: qb.hasAllKeys(idx.HashKeyParts); expression.Key(idx.HashKey)
	t.Run("buildHashKeyCondition_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) buildHashKeyCondition(", "Should contain buildHashKeyCondition signature")
		assert.Contains(t, rendered, "qb.hasAllKeys(idx.HashKeyParts)", "Should check composite parts via hasAllKeys")
		assert.Contains(t, rendered, "expression.Key(idx.HashKey)", "Should build simple key condition")
	})

	// Test that buildRangeKeyCondition supports optional and pre-built conditions
	// Example: idx.RangeKeyParts; qb.KeyConditions[idx.RangeKey]
	t.Run("buildRangeKeyCondition_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) buildRangeKeyCondition(", "Should contain buildRangeKeyCondition signature")
		assert.Contains(t, rendered, "idx.RangeKeyParts", "Should handle composite range keys")
		assert.Contains(t, rendered, "qb.KeyConditions[idx.RangeKey]", "Should use pre-built range conditions")
	})

	// Test that buildFilterCondition skips key attributes and combines filters
	// Example: for attrName, value := range qb.Attributes; expression.Name(attrName).Equal(...)
	t.Run("buildFilterCondition_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) buildFilterCondition(", "Should contain buildFilterCondition signature")
		assert.Contains(t, rendered, "for attrName, value := range qb.Attributes", "Should iterate over Attributes")
		assert.Contains(t, rendered, "expression.Name(attrName).Equal(expression.Value(value))", "Should build equality filters")
	})

	// Test that isPartOfIndexKey correctly identifies key attributes
	// Example: for _, part := range idx.HashKeyParts; attrName == idx.HashKey
	t.Run("isPartOfIndexKey_logic", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) isPartOfIndexKey(", "Should contain isPartOfIndexKey signature")
		assert.Contains(t, rendered, "for _, part := range idx.HashKeyParts", "Should check composite hash parts")
		assert.Contains(t, rendered, "attrName == idx.HashKey", "Should check simple hash key")
	})

	// Test that BuildQuery wraps Build and constructs AWS SDK call
	// Example: indexName, keyCond, filterCond, exclusiveStartKey, err := qb.Build()
	t.Run("BuildQuery_method_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) BuildQuery()", "Should contain BuildQuery signature")
		assert.Contains(t, rendered, "qb.Build()", "BuildQuery must call Build")
		assert.NotContains(t, rendered, "format.Source", "BuildQuery should not call format.Source")
	})

	// Test that Execute method is present and handles DynamoDB client
	// Example: func (qb *QueryBuilder) Execute(ctx context.Context, client *dynamodb.Client)
	t.Run("Execute_method_present", func(t *testing.T) {
		assert.Contains(t, rendered, "func (qb *QueryBuilder) Execute(ctx context.Context, client *dynamodb.Client)", "Should contain Execute signature")
		assert.Contains(t, rendered, "client.Query(ctx, input)", "Execute must call client.Query")
		assert.Contains(t, rendered, "attributevalue.UnmarshalListOfMaps", "Execute must unmarshal items")
	})

	// Test that no duplicate function signatures exist
	// Example: count occurrences of "func (qb *QueryBuilder)"—we expect one per method
	t.Run("no_duplicate_signatures", func(t *testing.T) {
		lines := strings.Split(rendered, "\n")
		count := 0
		for _, line := range lines {
			if strings.HasPrefix(strings.TrimSpace(line), "func (qb *QueryBuilder)") {
				count++
			}
		}
		assert.True(t, count >= 8, "Expected at least 8 QueryBuilder methods, found %d", count)
	})
}

// TestQueryBuilderBuildTemplateFormatting validates that the template is gofmt‐compliant.
// Example: format.Source("package test\n\n" + rendered + "\n") should return no error.
func TestQueryBuilderBuildTemplateFormatting(t *testing.T) {
	rendered := v2.QueryBuilderBuildTemplate
	full := "package test\n\n" + rendered + "\n"
	if _, err := format.Source([]byte(full)); err != nil {
		t.Fatalf("QueryBuilderBuildTemplate is not gofmt-compliant: %v", err)
	}
}
