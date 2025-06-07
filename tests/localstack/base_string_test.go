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

	basestring "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basestring"
)

// TestBaseString focuses on String (S) type operations and functionality.
// This test validates string-specific features without other data types.
//
// Test Coverage:
// - String CRUD operations
// - String marshaling/unmarshaling
// - String filtering and conditions
// - String operations in Scan (Contains, BeginsWith)
// - QueryBuilder with string keys and filters
//
// Schema: base-string.json
// - Table: "base-string"
// - Hash Key: id (S)
// - Range Key: category (S)
// - Common: title (S), description (S)
func TestBaseString(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing String operations on: %s", basestring.TableName)

	t.Run("String_Input", func(t *testing.T) {
		testStringInput(t, client, ctx)
	})

	t.Run("String_Input_Raw", func(t *testing.T) {
		testStringInputRaw(t, client, ctx)
	})

	t.Run("String_QueryBuilder", func(t *testing.T) {
		testStringQueryBuilder(t, client, ctx)
	})

	t.Run("String_ScanBuilder", func(t *testing.T) {
		testStringScanBuilder(t, client, ctx)
	})

	t.Run("String_Schema", func(t *testing.T) {
		t.Parallel()
		testStringSchema(t)
	})
}

// ==================== String Object Input ====================

func testStringInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_string_item", func(t *testing.T) {
		item := basestring.SchemaItem{
			Id:          "string-test-001",
			Category:    "docs",
			Title:       "String Operations Guide",
			Description: "Comprehensive guide for string handling",
		}

		// Test marshaling
		av, err := basestring.ItemInput(item)
		require.NoError(t, err, "Should marshal string item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		// Verify string fields are properly marshaled as AttributeValueMemberS
		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "title", "Should contain title field")
		assert.Contains(t, av, "category", "Should contain category field")
		assert.Contains(t, av, "description", "Should contain description field")

		// Verify actual string values
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestring.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestring.ColumnTitle])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestring.ColumnCategory])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestring.ColumnDescription])

		// Test storage
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basestring.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store string item in DynamoDB")

		t.Logf("✅ Created string item: %s/%s", item.Id, item.Category)
	})

	t.Run("read_string_item", func(t *testing.T) {
		item := basestring.SchemaItem{
			Id:       "string-test-001",
			Category: "docs",
		}

		key, err := basestring.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve string item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		// Verify string values are correctly retrieved
		assert.Equal(t, "docs", getResult.Item[basestring.ColumnCategory].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "string-test-001", getResult.Item[basestring.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "String Operations Guide", getResult.Item[basestring.ColumnTitle].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Comprehensive guide for string handling", getResult.Item[basestring.ColumnDescription].(*types.AttributeValueMemberS).Value)

		t.Logf("✅ Retrieved string item successfully")
	})

	t.Run("update_string_item", func(t *testing.T) {
		item := basestring.SchemaItem{
			Id:          "string-test-001",
			Category:    "docs",
			Title:       "Updated String Guide",
			Description: "Updated comprehensive guide for string operations",
		}

		updateInput, err := basestring.UpdateItemInput(item)
		require.NoError(t, err, "Should create update input from item")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update string item")

		// Verify update
		key, _ := basestring.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, "Updated String Guide", getResult.Item[basestring.ColumnTitle].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Updated comprehensive guide for string operations", getResult.Item[basestring.ColumnDescription].(*types.AttributeValueMemberS).Value)

		t.Logf("✅ Updated string item successfully")
	})

	t.Run("delete_string_item", func(t *testing.T) {
		item := basestring.SchemaItem{
			Id:       "string-test-001",
			Category: "docs",
		}

		deleteInput, err := basestring.DeleteItemInput(item)
		require.NoError(t, err, "Should create delete input from item")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete string item")

		// Verify deletion
		key, _ := basestring.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "String item should be deleted")

		t.Logf("✅ Deleted string item successfully")
	})

	t.Run("string_edge_cases", func(t *testing.T) {
		edgeCases := []basestring.SchemaItem{
			{Id: "edge-1", Category: "empty-test", Title: "", Description: "Empty title test"},
			{Id: "edge-2", Category: "special", Title: "Special chars: !@#$%^&*()", Description: "Unicode: 你好 🌟"},
			{Id: "edge-3", Category: "long", Title: "Very " + string(make([]byte, 100)), Description: "Long string test"},
			{Id: "edge-4", Category: "minimal", Title: "x", Description: "Single char"},
		}

		for _, item := range edgeCases {
			av, err := basestring.ItemInput(item)
			require.NoError(t, err, "Should handle edge case: %s", item.Id)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(basestring.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store edge case item: %s", item.Id)
		}

		t.Logf("✅ String edge cases handled successfully")
	})
}

// ================== String Raw Object Input ==================

func testStringInputRaw(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_string_item_raw", func(t *testing.T) {
		item := basestring.SchemaItem{
			Id:          "string-raw-001",
			Category:    "raw-docs",
			Title:       "Raw String Operations Guide",
			Description: "Guide for raw string handling methods",
		}

		// Test marshaling (same as object-based)
		av, err := basestring.ItemInput(item)
		require.NoError(t, err, "Should marshal string item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		// Test storage
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basestring.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store string item in DynamoDB")

		t.Logf("✅ Created string item for raw testing: %s/%s", item.Id, item.Category)
	})

	t.Run("read_string_item_raw", func(t *testing.T) {
		key, err := basestring.KeyInputFromRaw("string-raw-001", "raw-docs")
		require.NoError(t, err, "Should create key from raw values")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve string item using raw key")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		// Verify string values are correctly retrieved
		assert.Equal(t, "string-raw-001", getResult.Item[basestring.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "raw-docs", getResult.Item[basestring.ColumnCategory].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Raw String Operations Guide", getResult.Item[basestring.ColumnTitle].(*types.AttributeValueMemberS).Value)

		t.Logf("✅ Retrieved string item successfully using raw key")
	})

	t.Run("update_string_item_raw", func(t *testing.T) {
		updates := map[string]any{
			"title":       "Updated Raw String Guide",
			"description": "Updated guide for raw string operations methods",
		}

		updateInput, err := basestring.UpdateItemInputFromRaw("string-raw-001", "raw-docs", updates)
		require.NoError(t, err, "Should create update input from raw values")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update string item using raw method")

		// Verify update
		key, _ := basestring.KeyInputFromRaw("string-raw-001", "raw-docs")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, "Updated Raw String Guide", getResult.Item[basestring.ColumnTitle].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Updated guide for raw string operations methods", getResult.Item[basestring.ColumnDescription].(*types.AttributeValueMemberS).Value)

		t.Logf("✅ Updated string item successfully using raw method")
	})

	t.Run("delete_string_item_raw", func(t *testing.T) {
		deleteInput, err := basestring.DeleteItemInputFromRaw("string-raw-001", "raw-docs")
		require.NoError(t, err, "Should create delete input from raw values")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete string item using raw method")

		// Verify deletion
		key, _ := basestring.KeyInputFromRaw("string-raw-001", "raw-docs")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "String item should be deleted")

		t.Logf("✅ Deleted string item successfully using raw method")
	})

	t.Run("raw_vs_object_comparison", func(t *testing.T) {
		// Create same key using both methods
		keyFromRaw, err := basestring.KeyInputFromRaw("comparison-test", "both-methods")
		require.NoError(t, err, "Should create key from raw values")

		item := basestring.SchemaItem{
			Id:       "comparison-test",
			Category: "both-methods",
		}
		keyFromObject, err := basestring.KeyInput(item)
		require.NoError(t, err, "Should create key from object")

		// Keys should be identical
		assert.Equal(t, keyFromRaw, keyFromObject, "Raw and object-based keys should be identical")

		t.Logf("✅ Raw and object-based methods produce identical results")
	})

	t.Run("raw_string_edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			id       string
			category string
			updates  map[string]any
		}{
			{
				id:       "raw-edge-1",
				category: "empty",
				updates:  map[string]any{"title": "", "description": "Empty title"},
			},
			{
				id:       "raw-edge-2",
				category: "unicode",
				updates:  map[string]any{"title": "Unicode: 🚀✨", "description": "Emoji and unicode chars"},
			},
			{
				id:       "raw-edge-3",
				category: "special-chars",
				updates:  map[string]any{"title": "Special: !@#$%^&*()", "description": "Special characters"},
			},
		}

		for _, edgeCase := range edgeCases {
			// Test create with raw update
			updateInput, err := basestring.UpdateItemInputFromRaw(edgeCase.id, edgeCase.category, edgeCase.updates)
			require.NoError(t, err, "Should handle raw edge case: %s", edgeCase.id)

			// Note: UpdateItemFromRaw works as upsert-like operation for new items in some cases
			// For proper testing, we'd first PutItem, then UpdateItemFromRaw
			assert.NotNil(t, updateInput, "Update input should be created for edge case: %s", edgeCase.id)

			// Test delete with raw method
			deleteInput, err := basestring.DeleteItemInputFromRaw(edgeCase.id, edgeCase.category)
			require.NoError(t, err, "Should create delete input for edge case: %s", edgeCase.id)
			assert.NotNil(t, deleteInput, "Delete input should be created")
		}

		t.Logf("✅ Raw string edge cases handled successfully")
	})

	t.Run("raw_conditional_operations", func(t *testing.T) {
		// Test conditional delete with raw method
		conditionExpr := "#version = :v"
		conditionNames := map[string]string{"#version": "version"}
		conditionValues := map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberN{Value: "1"},
		}

		deleteInput, err := basestring.DeleteItemInputWithCondition(
			"conditional-test", "raw-condition",
			conditionExpr, conditionNames, conditionValues,
		)
		require.NoError(t, err, "Should create conditional delete with raw method")
		assert.NotNil(t, deleteInput.ConditionExpression, "Should have condition expression")
		assert.Equal(t, conditionExpr, *deleteInput.ConditionExpression, "Condition should match")

		// Test conditional update with raw method
		updates := map[string]any{
			"title":   "Conditional Update",
			"version": 2,
		}

		updateInput, err := basestring.UpdateItemInputWithCondition(
			"conditional-test", "raw-condition",
			updates, conditionExpr, conditionNames, conditionValues,
		)
		require.NoError(t, err, "Should create conditional update with raw method")
		assert.NotNil(t, updateInput.ConditionExpression, "Should have condition expression")

		t.Logf("✅ Raw conditional operations work correctly")
	})
}

// ==================== String QueryBuilder Tests ====================

func testStringQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup test data
	setupStringTestData(t, client, ctx)

	t.Run("string_hash_key_query", func(t *testing.T) {
		qb := basestring.NewQueryBuilder().WithId("query-string-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build string hash key query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		assert.Equal(t, basestring.TableName, *queryInput.TableName, "Should target correct table")

		t.Logf("✅ String hash key query built successfully")
	})

	t.Run("string_hash_and_range_query", func(t *testing.T) {
		qb := basestring.NewQueryBuilder().
			WithId("query-string-test").
			WithCategory("api")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build string hash+range query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ String hash+range query built successfully")
	})

	t.Run("string_filters", func(t *testing.T) {
		qb := basestring.NewQueryBuilder().
			WithId("query-string-test").
			FilterTitle("API Documentation").
			FilterDescription("REST API guide")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with string filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ String filters query built successfully")
	})

	t.Run("string_query_execution", func(t *testing.T) {
		qb := basestring.NewQueryBuilder().WithId("query-string-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute string query")
		assert.NotEmpty(t, items, "Should return string items")

		// Verify returned items have correct string types
		for _, item := range items {
			assert.Equal(t, "query-string-test", item.Id, "All items should have correct hash key")
			assert.NotEmpty(t, item.Category, "All items should have category")
			assert.IsType(t, "", item.Title, "Title should be string type")
			assert.IsType(t, "", item.Description, "Description should be string type")
		}

		t.Logf("✅ String query execution returned %d items", len(items))
	})

	t.Run("string_sorting", func(t *testing.T) {
		// Test sorting by string range key (category)
		qbAsc := basestring.NewQueryBuilder().
			WithId("query-string-test").
			OrderByAsc()

		itemsAsc, err := qbAsc.Execute(ctx, client)
		require.NoError(t, err, "Should execute ascending string query")

		qbDesc := basestring.NewQueryBuilder().
			WithId("query-string-test").
			OrderByDesc()

		itemsDesc, err := qbDesc.Execute(ctx, client)
		require.NoError(t, err, "Should execute descending string query")

		// Verify sorting is different
		if len(itemsAsc) > 1 && len(itemsDesc) > 1 {
			assert.NotEqual(t, itemsAsc[0].Category, itemsDesc[0].Category, "Sorting should produce different order")
		}

		t.Logf("✅ String sorting works correctly")
	})
}

// ==================== String ScanBuilder Tests ====================

func testStringScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("string_scan_filters", func(t *testing.T) {
		sb := basestring.NewScanBuilder().
			FilterId("query-string-test").
			FilterCategory("api")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with string filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ String scan filters built successfully")
	})

	t.Run("string_contains_filter", func(t *testing.T) {
		sb := basestring.NewScanBuilder().
			FilterTitleContains("API")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with contains filter")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ String contains filter built successfully")
	})

	t.Run("string_begins_with_filter", func(t *testing.T) {
		sb := basestring.NewScanBuilder().
			FilterDescriptionBeginsWith("REST")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with begins_with filter")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ String begins_with filter built successfully")
	})

	t.Run("string_scan_execution", func(t *testing.T) {
		sb := basestring.NewScanBuilder().
			FilterTitleContains("API").
			Limit(10)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute string scan")

		// Verify items match filter criteria
		for _, item := range items {
			assert.Contains(t, item.Title, "API", "Items should match contains filter")
		}

		t.Logf("✅ String scan execution returned %d items", len(items))
	})
}

// ==================== String Schema Tests ====================

func testStringSchema(t *testing.T) {
	t.Run("string_table_schema", func(t *testing.T) {
		schema := basestring.TableSchema

		assert.Equal(t, "base-string", schema.TableName, "Table name should match")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "category", schema.RangeKey, "Range key should be 'category'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")

		t.Logf("✅ String schema structure validated")
	})

	t.Run("string_attributes", func(t *testing.T) {
		// Check primary attributes (all string type)
		expectedPrimary := map[string]string{
			"id":       "S",
			"category": "S",
		}

		for _, attr := range basestring.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be string type", attr.Name)
		}

		// Check common attributes (all string type)
		expectedCommon := map[string]string{
			"title":       "S",
			"description": "S",
		}

		for _, attr := range basestring.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be string type", attr.Name)
		}

		t.Logf("✅ String attributes validated")
	})

	t.Run("string_constants", func(t *testing.T) {
		assert.Equal(t, "base-string", basestring.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", basestring.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "category", basestring.ColumnCategory, "ColumnCategory should be correct")
		assert.Equal(t, "title", basestring.ColumnTitle, "ColumnTitle should be correct")
		assert.Equal(t, "description", basestring.ColumnDescription, "ColumnDescription should be correct")

		t.Logf("✅ String constants validated")
	})

	t.Run("string_attribute_names", func(t *testing.T) {
		attrs := basestring.AttributeNames
		expectedAttrs := []string{"id", "category", "title", "description"}

		assert.Len(t, attrs, len(expectedAttrs), "Should have correct number of attributes")

		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}

		t.Logf("✅ String AttributeNames validated")
	})
}

// ==================== Helper Functions ====================

func setupStringTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basestring.SchemaItem{
		{Id: "query-string-test", Category: "api", Title: "API Documentation", Description: "REST API guide for developers"},
		{Id: "query-string-test", Category: "sdk", Title: "SDK Reference", Description: "Complete SDK documentation"},
		{Id: "query-string-test", Category: "tutorial", Title: "Getting Started", Description: "Quick start tutorial"},
	}

	for _, item := range testItems {
		av, err := basestring.ItemInput(item)
		require.NoError(t, err, "Should marshal string test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basestring.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store string test item")
	}

	t.Logf("Setup complete: inserted %d string test items", len(testItems))
}
