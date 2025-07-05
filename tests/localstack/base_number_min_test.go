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

	basenumbermin "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basenumbermin"
)

// TestBaseNumberMIN focuses on Number (N) type operations with minimal generated code.
// This test validates numeric functionality using only basic methods (no sugar methods).
//
// Test Coverage:
// - Basic Number CRUD operations using universal methods
// - Number marshaling/unmarshaling
// - Core Query and Scan operations using .With() and .Filter()
// - Basic increment operations
// - Schema validation
//
// Schema: base-number__min.json
// - Table: "base-number-min"
// - Hash Key: id (S)
// - Range Key: timestamp (N)
// - Common: count (N), price (N)
//
// Note: MIN mode only includes universal methods like .With() and .Filter()
// No convenience methods like .WithEQ(), .FilterEQ(), .WithBetween() etc.
func TestBaseNumberMIN(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Number MIN operations on: %s", basenumbermin.TableName)

	t.Run("Number_MIN_BasicCRUD", func(t *testing.T) {
		testNumberMINBasicCRUD(t, client, ctx)
	})

	t.Run("Number_MIN_QueryBuilder", func(t *testing.T) {
		testNumberMINQueryBuilder(t, client, ctx)
	})

	t.Run("Number_MIN_ScanBuilder", func(t *testing.T) {
		testNumberMINScanBuilder(t, client, ctx)
	})

	t.Run("Number_MIN_IncrementOperations", func(t *testing.T) {
		testNumberMINIncrementOperations(t, client, ctx)
	})

	t.Run("Number_MIN_Schema", func(t *testing.T) {
		t.Parallel()
		testNumberMINSchema(t)
	})
}

// ==================== Number MIN Basic CRUD ====================

func testNumberMINBasicCRUD(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_create_and_read", func(t *testing.T) {
		item := basenumbermin.SchemaItem{
			Id:        "min-number-001",
			Timestamp: 1640995200,
			Count:     42,
			Price:     1999,
		}

		av, err := basenumbermin.ItemInput(item)
		require.NoError(t, err, "Should marshal number item in MIN mode")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[basenumbermin.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[basenumbermin.ColumnTimestamp])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[basenumbermin.ColumnCount])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[basenumbermin.ColumnPrice])

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basenumbermin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store number item in DynamoDB")

		key, err := basenumbermin.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumbermin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve number item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "42", getResult.Item[basenumbermin.ColumnCount].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "1999", getResult.Item[basenumbermin.ColumnPrice].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "1640995200", getResult.Item[basenumbermin.ColumnTimestamp].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ MIN mode basic number CRUD operations work correctly")
	})

	t.Run("min_raw_operations", func(t *testing.T) {
		key, err := basenumbermin.KeyInputFromRaw("min-raw-001", 1640995300)
		require.NoError(t, err, "Should create key from raw values")
		assert.NotEmpty(t, key, "Raw key should not be empty")

		updates := map[string]any{
			"count": 100,
			"price": 2500,
		}

		updateInput, err := basenumbermin.UpdateItemInputFromRaw("min-raw-001", 1640995300, updates)
		require.NoError(t, err, "Should create update input from raw values")
		assert.NotNil(t, updateInput, "Update input should be created")

		t.Logf("✅ MIN mode raw number operations work correctly")
	})

	t.Run("min_number_edge_cases", func(t *testing.T) {
		edgeCases := []basenumbermin.SchemaItem{
			{Id: "min-edge-1", Timestamp: 0, Count: 0, Price: 0},
			{Id: "min-edge-2", Timestamp: 1, Count: -100, Price: -50},
			{Id: "min-edge-3", Timestamp: 9999999999, Count: 2147483647, Price: 999999999},
		}

		for _, item := range edgeCases {
			av, err := basenumbermin.ItemInput(item)
			require.NoError(t, err, "Should handle number edge case: %s", item.Id)
			assert.NotEmpty(t, av, "Marshaled edge case should not be empty")
		}

		t.Logf("✅ MIN mode number edge cases handled successfully")
	})
}

// ==================== Number MIN QueryBuilder Tests ====================

func testNumberMINQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupNumberMINTestData(t, client, ctx)

	t.Run("min_query_universal_methods", func(t *testing.T) {
		qb := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query using universal .With() method")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ MIN mode universal .With() method works")
	})

	t.Run("min_query_range_conditions", func(t *testing.T) {
		qb := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test").
			With("timestamp", basenumbermin.BETWEEN, 1640995200, 1640995400)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build range query using universal operators")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ MIN mode range conditions work with universal operators")
	})

	t.Run("min_query_numeric_comparisons", func(t *testing.T) {
		qbGT := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test").
			With("timestamp", basenumbermin.GT, 1640995300)

		queryInput, err := qbGT.BuildQuery()
		require.NoError(t, err, "Should build GT query using universal operators")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		qbLT := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test").
			With("timestamp", basenumbermin.LT, 1640995350)

		queryInput, err = qbLT.BuildQuery()
		require.NoError(t, err, "Should build LT query using universal operators")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ MIN mode numeric comparisons work with universal operators")
	})

	t.Run("min_query_with_filters", func(t *testing.T) {
		qb := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build basic query first")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		qbWithFilters := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test").
			Filter("count", basenumbermin.GT, 20)

		queryInputWithFilters, err := qbWithFilters.BuildQuery()
		require.NoError(t, err, "Should build query with universal .Filter() method")
		assert.NotNil(t, queryInputWithFilters.FilterExpression, "Should have filter expression")

		t.Logf("✅ MIN mode universal .Filter() method works for numbers")
	})

	t.Run("min_query_execution", func(t *testing.T) {
		qb := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN mode query")
		assert.NotEmpty(t, items, "Should return items")

		for _, item := range items {
			assert.Equal(t, "min-query-test", item.Id, "All items should have correct hash key")
			assert.Greater(t, item.Timestamp, 0, "All items should have positive timestamp")
		}
		t.Logf("✅ MIN mode query execution returned %d items", len(items))
	})

	t.Run("min_query_sorting", func(t *testing.T) {
		qbAsc := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test").
			OrderByAsc()

		itemsAsc, err := qbAsc.Execute(ctx, client)
		require.NoError(t, err, "Should execute ascending query")

		qbDesc := basenumbermin.NewQueryBuilder().
			With("id", basenumbermin.EQ, "min-query-test").
			OrderByDesc()

		itemsDesc, err := qbDesc.Execute(ctx, client)
		require.NoError(t, err, "Should execute descending query")

		if len(itemsAsc) > 1 && len(itemsDesc) > 1 {
			assert.NotEqual(t, itemsAsc[0].Timestamp, itemsDesc[0].Timestamp, "Sorting should produce different order")
		}
		t.Logf("✅ MIN mode sorting works correctly")
	})
}

// ==================== Number MIN ScanBuilder Tests ====================

func testNumberMINScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_scan_universal_filter", func(t *testing.T) {
		sb := basenumbermin.NewScanBuilder().
			Filter("count", basenumbermin.GT, 20)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with universal .Filter() method")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ MIN mode scan universal .Filter() method works")
	})

	t.Run("min_scan_multiple_numeric_filters", func(t *testing.T) {
		sb := basenumbermin.NewScanBuilder().
			Filter("count", basenumbermin.GT, 20).
			Filter("price", basenumbermin.LT, 3000).
			Filter("timestamp", basenumbermin.BETWEEN, 1640995200, 1640995500)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with multiple universal numeric filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")
		t.Logf("✅ MIN mode multiple universal numeric filters work")
	})

	t.Run("min_scan_execution", func(t *testing.T) {
		sb := basenumbermin.NewScanBuilder().
			Filter("id", basenumbermin.EQ, "min-query-test").
			Limit(5)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN mode scan")

		for _, item := range items {
			assert.Equal(t, "min-query-test", item.Id, "Items should match filter")
		}
		t.Logf("✅ MIN mode scan execution returned %d items", len(items))
	})
}

// ==================== Number MIN Increment Operations Tests ====================

func testNumberMINIncrementOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup item for increment testing
	testItem := basenumbermin.SchemaItem{
		Id:        "min-increment-test",
		Timestamp: 1640995888,
		Count:     10,
		Price:     100,
	}

	av, err := basenumbermin.ItemInput(testItem)
	require.NoError(t, err, "Should create test item for increment")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basenumbermin.TableName),
		Item:      av,
	})
	require.NoError(t, err, "Should store test item")

	t.Run("min_increment_basic", func(t *testing.T) {
		incrementInput, err := basenumbermin.IncrementAttribute("min-increment-test", 1640995888, "count", 5)
		require.NoError(t, err, "Should create increment input")
		assert.NotNil(t, incrementInput.UpdateExpression, "Should have update expression")
		assert.Contains(t, *incrementInput.UpdateExpression, "ADD", "Should use ADD operation")

		_, err = client.UpdateItem(ctx, incrementInput)
		require.NoError(t, err, "Should increment count")

		key, _ := basenumbermin.KeyInputFromRaw("min-increment-test", 1640995888)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumbermin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve incremented item")

		assert.Equal(t, "15", getResult.Item[basenumbermin.ColumnCount].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ MIN mode increment: 10 + 5 = 15")
	})

	t.Run("min_decrement_basic", func(t *testing.T) {
		decrementInput, err := basenumbermin.IncrementAttribute("min-increment-test", 1640995888, "price", -25)
		require.NoError(t, err, "Should create decrement input")

		_, err = client.UpdateItem(ctx, decrementInput)
		require.NoError(t, err, "Should decrement price")

		key, _ := basenumbermin.KeyInputFromRaw("min-increment-test", 1640995888)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basenumbermin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve decremented item")
		assert.Equal(t, "75", getResult.Item[basenumbermin.ColumnPrice].(*types.AttributeValueMemberN).Value)
		t.Logf("✅ MIN mode decrement: 100 - 25 = 75")
	})
}

// ==================== Number MIN Schema Tests ====================

func testNumberMINSchema(t *testing.T) {
	t.Run("min_schema_structure", func(t *testing.T) {
		schema := basenumbermin.TableSchema

		assert.Equal(t, "base-number-min", schema.TableName, "Table name should match MIN schema")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "timestamp", schema.RangeKey, "Range key should be 'timestamp'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")

		t.Logf("✅ MIN mode schema structure validated")
	})

	t.Run("min_constants", func(t *testing.T) {
		assert.Equal(t, "base-number-min", basenumbermin.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", basenumbermin.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "timestamp", basenumbermin.ColumnTimestamp, "ColumnTimestamp should be correct")
		assert.Equal(t, "count", basenumbermin.ColumnCount, "ColumnCount should be correct")
		assert.Equal(t, "price", basenumbermin.ColumnPrice, "ColumnPrice should be correct")

		t.Logf("✅ MIN mode constants validated")
	})

	t.Run("min_numeric_operators_available", func(t *testing.T) {
		assert.NotNil(t, basenumbermin.EQ, "EQ operator should be available")
		assert.NotNil(t, basenumbermin.GT, "GT operator should be available")
		assert.NotNil(t, basenumbermin.LT, "LT operator should be available")
		assert.NotNil(t, basenumbermin.GTE, "GTE operator should be available")
		assert.NotNil(t, basenumbermin.LTE, "LTE operator should be available")
		assert.NotNil(t, basenumbermin.BETWEEN, "BETWEEN operator should be available")
		t.Logf("✅ MIN mode universal numeric operators available")
	})

	t.Run("min_number_attributes", func(t *testing.T) {
		expectedPrimary := map[string]string{
			"id":        "S",
			"timestamp": "N",
		}
		for _, attr := range basenumbermin.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should have correct type", attr.Name)
		}
		expectedCommon := map[string]string{
			"count": "N",
			"price": "N",
		}
		for _, attr := range basenumbermin.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be number type", attr.Name)
		}
		t.Logf("✅ MIN mode number attributes validated")
	})

	t.Run("min_no_sugar_methods", func(t *testing.T) {
		qb := basenumbermin.NewQueryBuilder()
		assert.NotNil(t, qb, "QueryBuilder should be available")

		sb := basenumbermin.NewScanBuilder()
		assert.NotNil(t, sb, "ScanBuilder should be available")
		t.Logf("✅ MIN mode builders available (sugar methods should be absent)")
	})
}

// ==================== Helper Functions ====================

func setupNumberMINTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basenumbermin.SchemaItem{
		{Id: "min-query-test", Timestamp: 1640995300, Count: 25, Price: 1500},
		{Id: "min-query-test", Timestamp: 1640995400, Count: 35, Price: 2000},
		{Id: "min-query-test", Timestamp: 1640995500, Count: 45, Price: 2500},
	}
	for _, item := range testItems {
		av, err := basenumbermin.ItemInput(item)
		require.NoError(t, err, "Should marshal MIN test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basenumbermin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store MIN test item")
	}
	t.Logf("MIN setup complete: inserted %d number test items", len(testItems))
}
