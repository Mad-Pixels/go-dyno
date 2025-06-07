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

	basenumber "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basenumber"
)

// TestBaseNumber focuses on Number (N) type operations and functionality.
// This test validates numeric-specific features without other data types.
//
// Test Coverage:
// - Number CRUD operations
// - Number marshaling/unmarshaling
// - Numeric range conditions (Between, GreaterThan, LessThan)
// - Numeric operations in Query and Scan
// - Increment operations
// - Edge cases (zero, negative, large numbers)
//
// Schema: base-number.json
// - Table: "base-number"
// - Hash Key: id (S)
// - Range Key: timestamp (N)
// - Common: count (N), price (N)
func TestBaseNumber(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Number operations on: %s", basenumber.TableName)

	t.Run("Number_CRUD", func(t *testing.T) {
		testNumberCRUD(t, client, ctx)
	})

	t.Run("Number_QueryBuilder", func(t *testing.T) {
		testNumberQueryBuilder(t, client, ctx)
	})

	t.Run("Number_ScanBuilder", func(t *testing.T) {
		testNumberScanBuilder(t, client, ctx)
	})

	t.Run("Number_RangeConditions", func(t *testing.T) {
		testNumberRangeConditions(t, client, ctx)
	})

	t.Run("Number_IncrementOperations", func(t *testing.T) {
		testNumberIncrementOperations(t, client, ctx)
	})

	t.Run("Number_Schema", func(t *testing.T) {
		t.Parallel()
		testNumberSchema(t)
	})
}

// ==================== Number CRUD Operations ====================

func testNumberCRUD(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_number_item", func(t *testing.T) {
		item := basenumber.SchemaItem{
			Id:        "number-test-001",
			Timestamp: 1640995200,
			Count:     42,
			Price:     1999,
		}

		// Test marshaling
		av, err := basenumber.PutItem(item)
		require.NoError(t, err, "Should marshal number item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		// Verify number fields are properly marshaled as AttributeValueMemberN
		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "timestamp", "Should contain timestamp field")
		assert.Contains(t, av, "count", "Should contain count field")
		assert.Contains(t, av, "price", "Should contain price field")

		// Verify actual number values
		assert.IsType(t, &types.AttributeValueMemberS{}, av["id"])        // id is string
		assert.IsType(t, &types.AttributeValueMemberN{}, av["timestamp"]) // timestamp is number
		assert.IsType(t, &types.AttributeValueMemberN{}, av["count"])     // count is number
		assert.IsType(t, &types.AttributeValueMemberN{}, av["price"])     // price is number

		// Test storage
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basenumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store number item in DynamoDB")

		t.Logf("✅ Created number item: %s/%d", item.Id, item.Timestamp)
	})

	t.Run("read_number_item", func(t *testing.T) {
		key, err := basenumber.CreateKey("number-test-001", 1640995200)
		require.NoError(t, err, "Should create key")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve number item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		// Verify number values are correctly retrieved
		assert.Equal(t, "number-test-001", getResult.Item["id"].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "1640995200", getResult.Item["timestamp"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "42", getResult.Item["count"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "1999", getResult.Item["price"].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ Retrieved number item successfully")
	})

	t.Run("update_number_item", func(t *testing.T) {
		updates := map[string]interface{}{
			"count": 100,
			"price": 2499,
		}

		updateInput, err := basenumber.UpdateItem("number-test-001", 1640995200, updates)
		require.NoError(t, err, "Should create update input")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update number item")

		// Verify update
		key, _ := basenumber.CreateKey("number-test-001", 1640995200)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, "100", getResult.Item["count"].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "2499", getResult.Item["price"].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ Updated number item successfully")
	})

	t.Run("delete_number_item", func(t *testing.T) {
		deleteInput, err := basenumber.DeleteItem("number-test-001", 1640995200)
		require.NoError(t, err, "Should create delete input")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete number item")

		// Verify deletion
		key, _ := basenumber.CreateKey("number-test-001", 1640995200)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "Number item should be deleted")

		t.Logf("✅ Deleted number item successfully")
	})

	t.Run("number_edge_cases", func(t *testing.T) {
		edgeCases := []basenumber.SchemaItem{
			{Id: "edge-1", Timestamp: 0, Count: 0, Price: 0},                           // Zero values
			{Id: "edge-2", Timestamp: 1, Count: -100, Price: -50},                      // Negative numbers
			{Id: "edge-3", Timestamp: 9999999999, Count: 2147483647, Price: 999999999}, // Large numbers
			{Id: "edge-4", Timestamp: 1640995100, Count: 1, Price: 1},                  // Minimal positive
		}

		for _, item := range edgeCases {
			av, err := basenumber.PutItem(item)
			require.NoError(t, err, "Should handle number edge case: %s", item.Id)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(basenumber.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store number edge case item: %s", item.Id)
		}

		t.Logf("✅ Number edge cases handled successfully")
	})
}

// ==================== Number QueryBuilder Tests ====================

func testNumberQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup test data
	setupNumberTestData(t, client, ctx)

	t.Run("number_hash_key_query", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().WithId("query-number-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build number hash key query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		assert.Equal(t, basenumber.TableName, *queryInput.TableName, "Should target correct table")

		t.Logf("✅ Number hash key query built successfully")
	})

	t.Run("number_hash_and_range_query", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithTimestamp(1640995300)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build number hash+range query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ Number hash+range query built successfully")
	})

	t.Run("number_filters", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			FilterCount(25).
			FilterPrice(1500)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with number filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ Number filters query built successfully")
	})

	t.Run("number_query_execution", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().WithId("query-number-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute number query")
		assert.NotEmpty(t, items, "Should return number items")

		// Verify returned items have correct number types
		for _, item := range items {
			assert.Equal(t, "query-number-test", item.Id, "All items should have correct hash key")
			assert.Greater(t, item.Timestamp, 0, "All items should have positive timestamp")
			assert.IsType(t, 0, item.Count, "Count should be int type")
			assert.IsType(t, 0, item.Price, "Price should be int type")
		}

		t.Logf("✅ Number query execution returned %d items", len(items))
	})

	t.Run("number_sorting", func(t *testing.T) {
		// Test sorting by number range key (timestamp)
		qbAsc := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			OrderByAsc()

		itemsAsc, err := qbAsc.Execute(ctx, client)
		require.NoError(t, err, "Should execute ascending number query")

		qbDesc := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			OrderByDesc()

		itemsDesc, err := qbDesc.Execute(ctx, client)
		require.NoError(t, err, "Should execute descending number query")

		// Verify sorting is different for numbers
		if len(itemsAsc) > 1 && len(itemsDesc) > 1 {
			assert.NotEqual(t, itemsAsc[0].Timestamp, itemsDesc[0].Timestamp, "Number sorting should produce different order")
			// Verify ascending order
			if len(itemsAsc) > 1 {
				assert.LessOrEqual(t, itemsAsc[0].Timestamp, itemsAsc[1].Timestamp, "Ascending should be in increasing order")
			}
		}

		t.Logf("✅ Number sorting works correctly")
	})
}

// ==================== Number ScanBuilder Tests ====================

func testNumberScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("number_scan_filters", func(t *testing.T) {
		sb := basenumber.NewScanBuilder().
			FilterId("query-number-test").
			FilterCount(25)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with number filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ Number scan filters built successfully")
	})

	t.Run("number_scan_execution", func(t *testing.T) {
		sb := basenumber.NewScanBuilder().
			FilterCountGreaterThan(20).
			Limit(10)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute number scan")

		// Verify items match filter criteria
		for _, item := range items {
			assert.Greater(t, item.Count, 20, "Items should match greater than filter")
		}

		t.Logf("✅ Number scan execution returned %d items", len(items))
	})
}

// ==================== Number Range Conditions Tests ====================

func testNumberRangeConditions(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("timestamp_between", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithTimestampBetween(1640995200, 1640995400)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build timestamp between query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute between query")

		// Verify all returned items are within range
		for _, item := range items {
			assert.GreaterOrEqual(t, item.Timestamp, 1640995200, "Timestamp should be >= start")
			assert.LessOrEqual(t, item.Timestamp, 1640995400, "Timestamp should be <= end")
		}

		t.Logf("✅ Timestamp between condition returned %d items", len(items))
	})

	t.Run("timestamp_greater_than", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithTimestampGreaterThan(1640995300)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build timestamp greater than query")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute greater than query")

		// Verify all returned items are greater than threshold
		for _, item := range items {
			assert.Greater(t, item.Timestamp, 1640995300, "Timestamp should be > threshold")
		}

		t.Logf("✅ Timestamp greater than condition returned %d items", len(items))
	})

	t.Run("timestamp_less_than", func(t *testing.T) {
		qb := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithTimestampLessThan(1640995350)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build timestamp less than query")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute less than query")

		// Verify all returned items are less than threshold
		for _, item := range items {
			assert.Less(t, item.Timestamp, 1640995350, "Timestamp should be < threshold")
		}

		t.Logf("✅ Timestamp less than condition returned %d items", len(items))
	})

	t.Run("count_range_conditions", func(t *testing.T) {
		// Test range conditions on common attributes (count)
		qbBetween := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithCountBetween(20, 40)

		_, err := qbBetween.BuildQuery()
		require.NoError(t, err, "Should build count between query")

		qbGreater := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithCountGreaterThan(30)

		_, err = qbGreater.BuildQuery()
		require.NoError(t, err, "Should build count greater than query")

		qbLess := basenumber.NewQueryBuilder().
			WithId("query-number-test").
			WithCountLessThan(35)

		_, err = qbLess.BuildQuery()
		require.NoError(t, err, "Should build count less than query")

		t.Logf("✅ Count range conditions built successfully")
	})
}

// ==================== Number Increment Operations Tests ====================

func testNumberIncrementOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup item for increment testing
	testItem := basenumber.SchemaItem{
		Id:        "increment-test",
		Timestamp: 1640995999,
		Count:     10,
		Price:     100,
	}

	av, err := basenumber.PutItem(testItem)
	require.NoError(t, err, "Should create test item for increment")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basenumber.TableName),
		Item:      av,
	})
	require.NoError(t, err, "Should store test item")

	t.Run("increment_count", func(t *testing.T) {
		// Increment count by 5
		incrementInput, err := basenumber.IncrementAttribute("increment-test", 1640995999, "count", 5)
		require.NoError(t, err, "Should create increment input")

		_, err = client.UpdateItem(ctx, incrementInput)
		require.NoError(t, err, "Should increment count")

		// Verify increment
		key, _ := basenumber.CreateKey("increment-test", 1640995999)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve incremented item")

		assert.Equal(t, "15", getResult.Item["count"].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ Count incremented successfully: 10 + 5 = 15")
	})

	t.Run("decrement_price", func(t *testing.T) {
		// Decrement price by 25 (negative increment)
		decrementInput, err := basenumber.IncrementAttribute("increment-test", 1640995999, "price", -25)
		require.NoError(t, err, "Should create decrement input")

		_, err = client.UpdateItem(ctx, decrementInput)
		require.NoError(t, err, "Should decrement price")

		// Verify decrement
		key, _ := basenumber.CreateKey("increment-test", 1640995999)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve decremented item")

		assert.Equal(t, "75", getResult.Item["price"].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ Price decremented successfully: 100 - 25 = 75")
	})
}

// ==================== Number Schema Tests ====================

func testNumberSchema(t *testing.T) {
	t.Run("number_table_schema", func(t *testing.T) {
		schema := basenumber.TableSchema

		assert.Equal(t, "base-number", schema.TableName, "Table name should match")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "timestamp", schema.RangeKey, "Range key should be 'timestamp'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")

		t.Logf("✅ Number schema structure validated")
	})

	t.Run("number_attributes", func(t *testing.T) {
		// Check primary attributes
		expectedPrimary := map[string]string{
			"id":        "S", // hash key is string
			"timestamp": "N", // range key is number
		}

		for _, attr := range basenumber.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should have correct type", attr.Name)
		}

		// Check common attributes (all number type)
		expectedCommon := map[string]string{
			"count": "N",
			"price": "N",
		}

		for _, attr := range basenumber.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be number type", attr.Name)
		}

		t.Logf("✅ Number attributes validated")
	})

	t.Run("number_constants", func(t *testing.T) {
		assert.Equal(t, "base-number", basenumber.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", basenumber.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "timestamp", basenumber.ColumnTimestamp, "ColumnTimestamp should be correct")
		assert.Equal(t, "count", basenumber.ColumnCount, "ColumnCount should be correct")
		assert.Equal(t, "price", basenumber.ColumnPrice, "ColumnPrice should be correct")

		t.Logf("✅ Number constants validated")
	})

	t.Run("number_attribute_names", func(t *testing.T) {
		attrs := basenumber.AttributeNames
		expectedAttrs := []string{"id", "timestamp", "count", "price"}

		assert.Len(t, attrs, len(expectedAttrs), "Should have correct number of attributes")

		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}

		t.Logf("✅ Number AttributeNames validated")
	})
}

// ==================== Helper Functions ====================

func setupNumberTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basenumber.SchemaItem{
		{Id: "query-number-test", Timestamp: 1640995300, Count: 25, Price: 1500},
		{Id: "query-number-test", Timestamp: 1640995400, Count: 35, Price: 2000},
		{Id: "query-number-test", Timestamp: 1640995500, Count: 45, Price: 2500},
	}

	for _, item := range testItems {
		av, err := basenumber.PutItem(item)
		require.NoError(t, err, "Should marshal number test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basenumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store number test item")
	}

	t.Logf("Setup complete: inserted %d number test items", len(testItems))
}
