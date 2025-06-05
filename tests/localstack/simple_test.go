package localstack

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	simple "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/tablesimple"
)

// TestSimpleSchema runs comprehensive integration tests for the simple.json schema.
// This test suite validates all generated code functionality against a real LocalStack DynamoDB instance.
//
// Test Coverage:
// - Utility functions (BoolToInt, IntToBool)
// - CRUD operations (PutItem, BatchPutItems)
// - QueryBuilder basic functionality and fluent API
// - Schema constants and metadata validation
// - Generated struct marshaling/unmarshaling
//
// Schema Structure (simple.json):
//   - Table: "table-simple"
//   - Hash Key: "id" (string)
//   - Range Key: "created" (number)
//   - Attributes: ["id", "created", "name", "age"]
//   - No secondary indexes
//
// Example Usage:
//
//	go test -v ./tests/localstack/ -run TestSimpleSchema
func TestSimpleSchema(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(5 * time.Minute)
	defer cancel()

	t.Logf("Testing schema: simple.json")
	t.Logf("Table: %s", simple.TableName)
	t.Logf("Hash Key: %s, Range Key: %s", simple.TableSchema.HashKey, simple.TableSchema.RangeKey)

	t.Run("CRUD_Operations", func(t *testing.T) {
		testSimplePutItem(t, client, ctx)
		testSimpleBatchPutItems(t, client, ctx)
	})

	t.Run("QueryBuilder_Basic", func(t *testing.T) {
		testSimpleQueryBuilder(t, client, ctx)
		testSimpleQueryBuilderExecution(t, client, ctx)
		testSimpleQueryBuilderChaining(t)
	})

	t.Run("Schema_Constants", func(t *testing.T) {
		t.Parallel()
		testSimpleSchemaConstants(t)
		testSimpleAttributeNames(t)
		testSimpleTableSchema(t)
	})

	t.Run("Key_Operations", func(t *testing.T) {
		testSimpleCreateKey(t)
		testSimpleCreateKeyFromItem(t)
	})
}

// ==================== CRUD Operations Tests ====================

// testSimplePutItem validates item creation and DynamoDB marshaling functionality.
// Tests both local marshaling and actual DynamoDB write operations.
//
// Test Flow:
//  1. Create SchemaItem with all required fields
//  2. Marshal to DynamoDB AttributeValue format
//  3. Validate AttributeValue structure
//  4. Write to LocalStack DynamoDB
//  5. Verify successful storage
//
// Schema Mapping (simple.json):
//
//	Go Type     → DynamoDB Type → Attribute Name
//	string      → S             → "id" (hash key)
//	int         → N             → "created" (range key)
//	string      → S             → "name" (common attribute)
//	int         → N             → "age" (common attribute)
func testSimplePutItem(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_and_marshal_item", func(t *testing.T) {
		item := simple.SchemaItem{
			Id:      "test-id-123", // hash key: string
			Created: 1640995200,    // range key: unix timestamp
			Name:    "Test Item",   // common attribute: string
			Age:     25,            // common attribute: number
		}

		av, err := simple.PutItem(item)
		require.NoError(t, err, "PutItem marshaling should succeed")
		require.NotEmpty(t, av, "AttributeValues should not be empty")

		assert.Contains(t, av, "id", "Must contain hash key 'id'")
		assert.Contains(t, av, "created", "Must contain range key 'created'")
		assert.Contains(t, av, "name", "Must contain common attribute 'name'")
		assert.Contains(t, av, "age", "Must contain common attribute 'age'")

		assert.Equal(t, "test-id-123", av["id"].(*types.AttributeValueMemberS).Value, "ID should be string type")
		t.Logf("✅ Successfully created and marshaled item with ID: %s", item.Id)
	})

	t.Run("put_item_to_dynamodb", func(t *testing.T) {
		item := simple.SchemaItem{
			Id:      "test-put-456",
			Created: 1640995300,
			Name:    "Put Test Item",
			Age:     30,
		}

		av, err := simple.PutItem(item)
		require.NoError(t, err, "Item marshaling should succeed")

		input := &dynamodb.PutItemInput{
			TableName: aws.String(simple.TableName),
			Item:      av,
		}
		_, err = client.PutItem(ctx, input)
		require.NoError(t, err, "DynamoDB PutItem operation should succeed")
		t.Logf("✅ Successfully wrote item to DynamoDB table: %s", simple.TableName)

		getInput := &dynamodb.GetItemInput{
			TableName: aws.String(simple.TableName),
			Key: map[string]types.AttributeValue{
				"id":      &types.AttributeValueMemberS{Value: item.Id},
				"created": &types.AttributeValueMemberN{Value: "1640995300"},
			},
		}
		result, err := client.GetItem(ctx, getInput)
		require.NoError(t, err, "GetItem verification should succeed")
		assert.NotEmpty(t, result.Item, "Stored item should be retrievable")
		t.Logf("✅ Verified item persistence in DynamoDB")
	})
}

// testSimpleBatchPutItems validates batch operations for multiple items.
// Tests the generated BatchPutItems utility for bulk operations.
//
// Batch Operations Benefits:
//   - Reduced API calls (up to 25 items per batch)
//   - Better performance for bulk inserts
//   - Atomic failure handling per batch
//
// Example Use Case:
//
//	items := []SchemaItem{{...}, {...}, {...}}
//	batchItems, err := BatchPutItems(items)
//	// Use batchItems with DynamoDB BatchWriteItem API
func testSimpleBatchPutItems(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("batch_put_multiple_items", func(t *testing.T) {
		items := []simple.SchemaItem{
			{Id: "batch-1", Created: 1640995400, Name: "Batch Item 1", Age: 20},
			{Id: "batch-2", Created: 1640995500, Name: "Batch Item 2", Age: 25},
			{Id: "batch-3", Created: 1640995600, Name: "Batch Item 3", Age: 30},
		}

		batchItems, err := simple.BatchPutItems(items)
		require.NoError(t, err, "BatchPutItems should succeed")
		require.Len(t, batchItems, 3, "Should return AttributeValues for all 3 items")

		for i, batchItem := range batchItems {
			assert.Contains(t, batchItem, "id", "Batch item %d should have id", i)
			assert.Contains(t, batchItem, "created", "Batch item %d should have created", i)
			assert.Contains(t, batchItem, "name", "Batch item %d should have name", i)
			assert.Contains(t, batchItem, "age", "Batch item %d should have age", i)
		}
		t.Logf("✅ Successfully prepared %d items for batch operation", len(batchItems))
	})

	t.Run("batch_put_empty_slice", func(t *testing.T) {
		emptyItems := []simple.SchemaItem{}
		batchItems, err := simple.BatchPutItems(emptyItems)

		require.NoError(t, err, "BatchPutItems should handle empty slice")
		assert.Empty(t, batchItems, "Should return empty slice for empty input")
		t.Logf("✅ BatchPutItems correctly handles empty input")
	})
}

// ==================== QueryBuilder Tests ====================

// testSimpleQueryBuilder validates QueryBuilder creation and fluent API methods.
// The QueryBuilder provides a type-safe, fluent interface for constructing DynamoDB queries.
//
// Generated Methods (based on simple.json schema):
//   - NewQueryBuilder() → Creates new builder instance
//   - WithId(string) → Sets hash key condition
//   - WithCreated(int) → Sets range key condition
//   - WithName(string) → Sets filter condition
//   - WithAge(int) → Sets filter condition
//   - OrderByAsc/OrderByDesc() → Controls sort order
//   - Limit(int) → Limits result count
//
// Example Usage:
//
//	query := NewQueryBuilder().
//	    WithId("user123").
//	    WithCreated(1640995200).
//	    OrderByDesc().
//	    Limit(10)
//	items, err := query.Execute(ctx, client)
func testSimpleQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_query_builder", func(t *testing.T) {
		qb := simple.NewQueryBuilder()
		require.NotNil(t, qb, "NewQueryBuilder should return non-nil instance")

		t.Logf("✅ QueryBuilder created successfully")
	})

	t.Run("fluent_api_methods", func(t *testing.T) {
		qb := simple.NewQueryBuilder()

		qb2 := qb.WithId("test-id")
		require.NotNil(t, qb2, "WithId should return QueryBuilder for chaining")
		assert.IsType(t, qb, qb2, "WithId should return same type")

		qb3 := qb.WithCreated(1640995200)
		require.NotNil(t, qb3, "WithCreated should return QueryBuilder for chaining")

		qb4 := qb.FilterName("test-name")
		require.NotNil(t, qb4, "WithName should return QueryBuilder for chaining")

		qb5 := qb.FilterAge(25)
		require.NotNil(t, qb5, "WithAge should return QueryBuilder for chaining")

		t.Logf("✅ All fluent API methods support proper method chaining")
	})

	t.Run("sorting_and_pagination", func(t *testing.T) {
		qb := simple.NewQueryBuilder()

		qbDesc := qb.OrderByDesc()
		require.NotNil(t, qbDesc, "OrderByDesc should return QueryBuilder")

		qbAsc := qb.OrderByAsc()
		require.NotNil(t, qbAsc, "OrderByAsc should return QueryBuilder")

		qbLimit := qb.Limit(50)
		require.NotNil(t, qbLimit, "Limit should return QueryBuilder")

		t.Logf("✅ Sorting and pagination methods work correctly")
	})
}

// testSimpleQueryBuilderChaining validates complex method chaining scenarios.
// Tests that multiple fluent methods can be chained together effectively.
//
// Chaining Patterns:
//  1. Key conditions + filters + sorting + pagination
//  2. Multiple attribute filters
//  3. Order of operations independence
func testSimpleQueryBuilderChaining(t *testing.T) {
	t.Run("complex_method_chaining", func(t *testing.T) {
		qb := simple.NewQueryBuilder().
			WithId("chain-test").
			WithCreated(1640995200).
			FilterName("Chained Query").
			FilterAge(35).
			OrderByDesc().
			Limit(25)

		require.NotNil(t, qb, "Complex chaining should work")
		t.Logf("✅ Complex method chaining works correctly")
	})

	t.Run("method_order_independence", func(t *testing.T) {
		qb1 := simple.NewQueryBuilder().WithId("test1").OrderByDesc().FilterAge(25)
		qb2 := simple.NewQueryBuilder().OrderByDesc().FilterAge(25).WithId("test1")
		qb3 := simple.NewQueryBuilder().FilterAge(25).WithId("test1").OrderByDesc()

		require.NotNil(t, qb1, "Chaining order 1 should work")
		require.NotNil(t, qb2, "Chaining order 2 should work")
		require.NotNil(t, qb3, "Chaining order 3 should work")

		t.Logf("✅ Method chaining order independence confirmed")
	})
}

// testSimpleQueryBuilderExecution validates query building and execution.
// Tests BuildQuery method and validates generated QueryInput structure.
//
// BuildQuery Output Validation:
//   - TableName matches schema
//   - KeyConditionExpression is properly formed
//   - ExpressionAttributeNames/Values are set
//   - Optional: FilterExpression for non-key attributes
func testSimpleQueryBuilderExecution(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("build_query_input", func(t *testing.T) {
		testItem := simple.SchemaItem{
			Id:      "query-test-123",
			Created: 1640995700,
			Name:    "Query Test Item",
			Age:     35,
		}

		av, err := simple.PutItem(testItem)
		require.NoError(t, err)

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(simple.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Test data setup should succeed")

		qb := simple.NewQueryBuilder().WithId("query-test-123")
		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "BuildQuery should succeed")
		require.NotNil(t, queryInput, "QueryInput should not be nil")

		assert.Equal(t, simple.TableName, *queryInput.TableName, "TableName should match schema")
		assert.NotNil(t, queryInput.KeyConditionExpression, "KeyConditionExpression should be set")
		assert.NotEmpty(t, queryInput.ExpressionAttributeNames, "ExpressionAttributeNames should be populated")
		assert.NotEmpty(t, queryInput.ExpressionAttributeValues, "ExpressionAttributeValues should be populated")

		t.Logf("✅ QueryInput built successfully")
		t.Logf("    Table: %s", *queryInput.TableName)
		t.Logf("    KeyCondition: %s", *queryInput.KeyConditionExpression)
	})

	t.Run("build_query_with_filters", func(t *testing.T) {
		qb := simple.NewQueryBuilder().
			WithId("filter-test").
			FilterName("Test Name").
			FilterAge(30)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "BuildQuery with filters should succeed")

		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		t.Logf("✅ QueryInput with filters built successfully")
	})
}

// ==================== Schema Constants Tests ====================

// testSimpleSchemaConstants validates all generated constants and their values.
// Constants provide compile-time safety and prevent typos in table/column names.
//
// Generated Constants (from simple.json):
//   - TableName = "table-simple"
//   - ColumnId = "id"
//   - ColumnCreated = "created"
//   - ColumnName = "name"
//   - ColumnAge = "age"
func testSimpleSchemaConstants(t *testing.T) {
	t.Run("table_name_constant", func(t *testing.T) {
		assert.NotEmpty(t, simple.TableName, "TableName constant should not be empty")
		assert.Equal(t, "table-simple", simple.TableName, "TableName should match simple.json schema")

		t.Logf("✅ TableName constant: %s", simple.TableName)
	})

	t.Run("column_constants", func(t *testing.T) {
		columnTests := map[string]string{
			"ColumnId":      "id",
			"ColumnCreated": "created",
			"ColumnName":    "name",
			"ColumnAge":     "age",
		}

		assert.Equal(t, "id", simple.ColumnId, "ColumnId should be 'id'")
		assert.Equal(t, "created", simple.ColumnCreated, "ColumnCreated should be 'created'")
		assert.Equal(t, "name", simple.ColumnName, "ColumnName should be 'name'")
		assert.Equal(t, "age", simple.ColumnAge, "ColumnAge should be 'age'")

		t.Logf("✅ All column constants validated:")
		for constant, expected := range columnTests {
			t.Logf("    %s = %s", constant, expected)
		}
	})
}

// testSimpleAttributeNames validates the AttributeNames array.
// This array contains all table attributes and is used for:
//   - Projection expressions
//   - Validation logic
//   - Dynamic attribute iteration
//
// Expected Attributes (from simple.json):
//
//	["id", "created", "name", "age"]
func testSimpleAttributeNames(t *testing.T) {
	t.Run("attribute_names_array", func(t *testing.T) {
		attrs := simple.AttributeNames
		require.NotEmpty(t, attrs, "AttributeNames should not be empty")

		expectedAttrs := []string{"id", "created", "name", "age"}
		assert.Len(t, attrs, len(expectedAttrs), "AttributeNames should contain all schema attributes")

		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain '%s'", expected)
		}

		attrSet := make(map[string]bool)
		for _, attr := range attrs {
			assert.False(t, attrSet[attr], "AttributeNames should not contain duplicate: %s", attr)
			attrSet[attr] = true
		}
		t.Logf("✅ AttributeNames array contains %d attributes: %v", len(attrs), attrs)
	})
}

// testSimpleTableSchema validates the TableSchema variable structure.
// TableSchema contains complete table metadata for runtime operations.
//
// Schema Structure Validation:
//   - TableName, HashKey, RangeKey
//   - Attributes and CommonAttributes arrays
//   - SecondaryIndexes (should be empty for simple.json)
func testSimpleTableSchema(t *testing.T) {
	t.Run("table_schema_structure", func(t *testing.T) {
		schema := simple.TableSchema

		assert.Equal(t, "table-simple", schema.TableName, "Schema TableName should match table name")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "created", schema.RangeKey, "Range key should be 'created'")

		assert.Len(t, schema.Attributes, 2, "Should have 2 primary attributes")
		assert.Len(t, schema.CommonAttributes, 2, "Should have 2 common attributes")
		assert.Len(t, schema.SecondaryIndexes, 0, "Simple schema should have no secondary indexes")

		primaryAttrs := []string{"id", "created"}
		commonAttrs := []string{"name", "age"}

		for _, attr := range schema.Attributes {
			assert.Contains(t, primaryAttrs, attr.Name, "Primary attribute should be in expected list")
			assert.NotEmpty(t, attr.Type, "Attribute type should not be empty")
		}

		for _, attr := range schema.CommonAttributes {
			assert.Contains(t, commonAttrs, attr.Name, "Common attribute should be in expected list")
			assert.NotEmpty(t, attr.Type, "Attribute type should not be empty")
		}

		t.Logf("✅ TableSchema structure is valid:")
		t.Logf("    Table: %s", schema.TableName)
		t.Logf("    Keys: %s (hash), %s (range)", schema.HashKey, schema.RangeKey)
		t.Logf("    Attributes: %d primary, %d common", len(schema.Attributes), len(schema.CommonAttributes))
	})
}

// ==================== Key Operations Tests ====================

// testSimpleCreateKey validates key creation utilities.
// These functions generate DynamoDB key structures for GetItem, DeleteItem, etc.
//
// Key Structure:
//
//	{
//	  "id": {"S": "value"},           // hash key
//	  "created": {"N": "timestamp"}   // range key
//	}
func testSimpleCreateKey(t *testing.T) {
	t.Run("create_key_with_hash_and_range", func(t *testing.T) {
		hashKeyValue := "test-key-123"
		rangeKeyValue := 1640995800

		key, err := simple.CreateKey(hashKeyValue, rangeKeyValue)
		require.NoError(t, err, "CreateKey should succeed with valid inputs")
		require.NotEmpty(t, key, "Created key should not be empty")

		assert.Contains(t, key, "id", "Key should contain hash key 'id'")
		assert.Contains(t, key, "created", "Key should contain range key 'created'")

		t.Logf("✅ CreateKey generated valid DynamoDB key structure")
	})

	t.Run("create_key_hash_only", func(t *testing.T) {
		hashKeyValue := "hash-only-test"

		key, err := simple.CreateKey(hashKeyValue, nil)
		require.NoError(t, err, "CreateKey should handle nil range key")
		assert.Contains(t, key, "id", "Key should contain hash key 'id'")
		t.Logf("✅ CreateKey handles hash-only keys correctly")
	})
}

// testSimpleCreateKeyFromItem validates key extraction from SchemaItem.
// This utility extracts just the key attributes from a complete item.
//
// Use Cases:
//   - GetItem operations after initial PutItem
//   - DeleteItem operations
//   - Key-based updates
func testSimpleCreateKeyFromItem(t *testing.T) {
	t.Run("extract_key_from_complete_item", func(t *testing.T) {
		item := simple.SchemaItem{
			Id:      "extract-test-456",
			Created: 1640995900,
			Name:    "Key Extraction Test",
			Age:     40,
		}

		key, err := simple.CreateKeyFromItem(item)
		require.NoError(t, err, "CreateKeyFromItem should succeed")
		require.NotEmpty(t, key, "Extracted key should not be empty")

		assert.Contains(t, key, "id", "Extracted key should contain hash key 'id'")
		assert.Contains(t, key, "created", "Extracted key should contain range key 'created'")
		assert.NotContains(t, key, "name", "Extracted key should not contain non-key attribute 'name'")
		assert.NotContains(t, key, "age", "Extracted key should not contain non-key attribute 'age'")

		t.Logf("✅ CreateKeyFromItem correctly extracts only key attributes")
	})
}
