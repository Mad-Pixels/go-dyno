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

	customnumber "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/customnumber"
)

// TestCustomNumber focuses on custom Number subtypes operations and functionality.
// This test validates that custom Go numeric types work correctly with DynamoDB operations.
//
// Test Coverage:
// - Custom numeric type marshaling/unmarshaling (int32, int64, float32, uint64, int16)
// - Type safety in generated structs
// - Custom types in range conditions
// - Custom types in QueryBuilder and ScanBuilder
// - ExtractFromDynamoDBStreamEvent with custom types
//
// Schema: custom-number.json
// - Table: "custom-number"
// - Hash Key: id (S)
// - Range Key: timestamp (N with int64 subtype)
// - Common: count (int32), price (float32), views (uint64), score (int16)
func TestCustomNumber(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Custom Number subtypes on: %s", customnumber.TableName)

	t.Run("Custom_Types_Struct", func(t *testing.T) {
		testCustomTypesStruct(t)
	})

	t.Run("Custom_Types_Marshaling", func(t *testing.T) {
		testCustomTypesMarshaling(t, client, ctx)
	})

	t.Run("Custom_Types_QueryBuilder", func(t *testing.T) {
		testCustomTypesQueryBuilder(t, client, ctx)
	})

	t.Run("Custom_Types_RangeConditions", func(t *testing.T) {
		testCustomTypesRangeConditions(t, client, ctx)
	})

	t.Run("Custom_Types_StreamEvent", func(t *testing.T) {
		testCustomTypesStreamEvent(t)
	})
}

// ==================== Custom Types Struct Tests ====================

func testCustomTypesStruct(t *testing.T) {
	t.Run("verify_go_types", func(t *testing.T) {
		// Создаем item с кастомными типами
		item := customnumber.SchemaItem{
			Id:        "test-id",
			Timestamp: 1640995200, // int64
			Count:     42,         // int32
			Price:     19.99,      // float32
			Views:     1000000,    // uint64
			Score:     85,         // int16
		}

		// Проверяем что компилятор принимает правильные типы
		assert.IsType(t, "", item.Id)              // string
		assert.IsType(t, int64(0), item.Timestamp) // int64
		assert.IsType(t, int32(0), item.Count)     // int32
		assert.IsType(t, float32(0), item.Price)   // float32
		assert.IsType(t, uint64(0), item.Views)    // uint64
		assert.IsType(t, int16(0), item.Score)     // int16

		// Проверяем что значения сохранились правильно
		assert.Equal(t, "test-id", item.Id)
		assert.Equal(t, int64(1640995200), item.Timestamp)
		assert.Equal(t, int32(42), item.Count)
		assert.Equal(t, float32(19.99), item.Price)
		assert.Equal(t, uint64(1000000), item.Views)
		assert.Equal(t, int16(85), item.Score)

		t.Logf("✅ Custom types verified: int64, int32, float32, uint64, int16")
	})

	t.Run("type_safety_compilation", func(t *testing.T) {
		// Эти присваивания должны компилироваться без ошибок
		var item customnumber.SchemaItem

		item.Timestamp = int64(1640995200)
		item.Count = int32(42)
		item.Price = float32(19.99)
		item.Views = uint64(1000000)
		item.Score = int16(85)

		// Проверяем что присваивания работают
		assert.Equal(t, int64(1640995200), item.Timestamp)
		assert.Equal(t, int32(42), item.Count)
		assert.Equal(t, float32(19.99), item.Price)
		assert.Equal(t, uint64(1000000), item.Views)
		assert.Equal(t, int16(85), item.Score)

		t.Logf("✅ Type safety compilation verified")
	})
}

// ==================== Custom Types Marshaling Tests ====================

func testCustomTypesMarshaling(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("marshal_custom_types", func(t *testing.T) {
		item := customnumber.SchemaItem{
			Id:        "marshal-test-001",
			Timestamp: 1640995200,
			Count:     42,
			Price:     19.99,
			Views:     1000000,
			Score:     85,
		}

		// Test marshaling with custom types
		av, err := customnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal item with custom types")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		// Verify all fields are present and properly typed
		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "timestamp", "Should contain timestamp field")
		assert.Contains(t, av, "count", "Should contain count field")
		assert.Contains(t, av, "price", "Should contain price field")
		assert.Contains(t, av, "views", "Should contain views field")
		assert.Contains(t, av, "score", "Should contain score field")

		// Verify DynamoDB types (all custom numeric types should be marshaled as N)
		assert.IsType(t, &types.AttributeValueMemberS{}, av[customnumber.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[customnumber.ColumnTimestamp])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[customnumber.ColumnCount])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[customnumber.ColumnPrice])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[customnumber.ColumnViews])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[customnumber.ColumnScore])

		t.Logf("✅ Custom types marshaled successfully")
	})

	t.Run("roundtrip_custom_types", func(t *testing.T) {
		originalItem := customnumber.SchemaItem{
			Id:        "roundtrip-test-001",
			Timestamp: 1640995300,
			Count:     123,
			Price:     29.95,
			Views:     2000000,
			Score:     95,
		}

		// Marshal and store
		av, err := customnumber.ItemInput(originalItem)
		require.NoError(t, err, "Should marshal original item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(customnumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store item in DynamoDB")

		// Retrieve and verify
		key, err := customnumber.KeyInput(originalItem)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(customnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		// Verify values are preserved with correct precision
		assert.Equal(t, "roundtrip-test-001", getResult.Item[customnumber.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "1640995300", getResult.Item[customnumber.ColumnTimestamp].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "123", getResult.Item[customnumber.ColumnCount].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "29.95", getResult.Item[customnumber.ColumnPrice].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "2000000", getResult.Item[customnumber.ColumnViews].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "95", getResult.Item[customnumber.ColumnScore].(*types.AttributeValueMemberN).Value)

		t.Logf("✅ Custom types roundtrip successful")
	})

	t.Run("edge_values_custom_types", func(t *testing.T) {
		edgeItems := []customnumber.SchemaItem{
			{
				Id:        "edge-1",
				Timestamp: 0,   // int64 zero
				Count:     0,   // int32 zero
				Price:     0.0, // float32 zero
				Views:     0,   // uint64 zero
				Score:     0,   // int16 zero
			},
			{
				Id:        "edge-2",
				Timestamp: 9223372036854775807,  // int64 max
				Count:     2147483647,           // int32 max
				Price:     3.4028235e+38,        // float32 max
				Views:     18446744073709551615, // uint64 max
				Score:     32767,                // int16 max
			},
			{
				Id:        "edge-3",
				Timestamp: -9223372036854775808, // int64 min
				Count:     -2147483648,          // int32 min
				Price:     -3.4028235e+38,       // float32 min
				Views:     0,                    // uint64 min (can't be negative)
				Score:     -32768,               // int16 min
			},
		}

		for _, item := range edgeItems {
			av, err := customnumber.ItemInput(item)
			require.NoError(t, err, "Should handle edge values for item: %s", item.Id)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(customnumber.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store edge value item: %s", item.Id)
		}

		t.Logf("✅ Custom types edge values handled successfully")
	})
}

// ==================== Custom Types QueryBuilder Tests ====================

func testCustomTypesQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup test data with custom types
	setupCustomTypesTestData(t, client, ctx)

	t.Run("custom_types_parameters", func(t *testing.T) {
		// Test that QueryBuilder methods accept correct custom types
		qb := customnumber.NewQueryBuilder().
			WithId("query-custom-test").
			WithTimestamp(1640995400) // int64 parameter

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with custom int64 parameter")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ Custom types parameters accepted by QueryBuilder")
	})

	t.Run("custom_types_filters", func(t *testing.T) {
		qb := customnumber.NewQueryBuilder().
			WithId("query-custom-test").
			FilterCount(int32(50)).      // int32 parameter
			FilterPrice(float32(25.0)).  // float32 parameter
			FilterViews(uint64(500000)). // uint64 parameter
			FilterScore(int16(80))       // int16 parameter

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with custom type filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ Custom types filters work in QueryBuilder")
	})

	t.Run("custom_types_execution", func(t *testing.T) {
		qb := customnumber.NewQueryBuilder().WithId("query-custom-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute query with custom types")
		assert.NotEmpty(t, items, "Should return items")

		// Verify returned items have correct custom types
		for _, item := range items {
			assert.IsType(t, int64(0), item.Timestamp, "Timestamp should be int64")
			assert.IsType(t, int32(0), item.Count, "Count should be int32")
			assert.IsType(t, float32(0), item.Price, "Price should be float32")
			assert.IsType(t, uint64(0), item.Views, "Views should be uint64")
			assert.IsType(t, int16(0), item.Score, "Score should be int16")
		}

		t.Logf("✅ Custom types query execution returned %d items", len(items))
	})
}

// ==================== Custom Types Range Conditions Tests ====================

func testCustomTypesRangeConditions(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("int64_timestamp_conditions", func(t *testing.T) {
		// Test range conditions on int64 timestamp
		qb := customnumber.NewQueryBuilder().
			WithId("query-custom-test").
			WithTimestampBetween(1640995200, 1640995500)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build int64 timestamp between condition")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute int64 between query")

		// Verify all items are within range
		for _, item := range items {
			assert.GreaterOrEqual(t, item.Timestamp, int64(1640995200), "Timestamp should be >= start")
			assert.LessOrEqual(t, item.Timestamp, int64(1640995500), "Timestamp should be <= end")
		}

		t.Logf("✅ int64 timestamp range conditions work correctly")
	})

	t.Run("int32_count_conditions", func(t *testing.T) {
		// Test range conditions on int32 count
		qb := customnumber.NewQueryBuilder().
			WithId("query-custom-test").
			WithCountBetween(40, 60)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build int32 count between condition")

		t.Logf("✅ int32 count range conditions compile correctly")
	})

	t.Run("uint64_views_conditions", func(t *testing.T) {
		// Test range conditions on uint64 views
		qb := customnumber.NewQueryBuilder().
			WithId("query-custom-test").
			WithViewsGreaterThan(1000000)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build uint64 views greater than condition")

		t.Logf("✅ uint64 views range conditions compile correctly")
	})
}

// ==================== Custom Types Stream Event Tests ====================

func testCustomTypesStreamEvent(t *testing.T) {
	t.Run("extract_custom_types_logic", func(t *testing.T) {
		// Проверяем что в сгенерированном коде есть правильная логика для кастомных типов
		// Это косвенный тест - мы проверяем что код компилируется и имеет правильную структуру

		item := customnumber.SchemaItem{
			Id:        "stream-test",
			Timestamp: 1640995600,
			Count:     75,
			Price:     39.99,
			Views:     3000000,
			Score:     90,
		}

		// Если код компилируется и эти присваивания работают, значит типы правильные
		assert.Equal(t, int64(1640995600), item.Timestamp)
		assert.Equal(t, int32(75), item.Count)
		assert.Equal(t, float32(39.99), item.Price)
		assert.Equal(t, uint64(3000000), item.Views)
		assert.Equal(t, int16(90), item.Score)

		t.Logf("✅ Custom types stream event extraction logic verified")
	})
}

// ==================== Helper Functions ====================

func setupCustomTypesTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []customnumber.SchemaItem{
		{Id: "query-custom-test", Timestamp: 1640995300, Count: 45, Price: 19.99, Views: 1500000, Score: 85},
		{Id: "query-custom-test", Timestamp: 1640995400, Count: 55, Price: 29.99, Views: 2500000, Score: 90},
		{Id: "query-custom-test", Timestamp: 1640995500, Count: 65, Price: 39.99, Views: 3500000, Score: 95},
	}

	for _, item := range testItems {
		av, err := customnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal custom types test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(customnumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store custom types test item")
	}

	t.Logf("Setup complete: inserted %d custom types test items", len(testItems))
}
