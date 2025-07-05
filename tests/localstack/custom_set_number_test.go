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

	customsetnumber "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/customsetnumberall"
)

// TestCustomSetNumber focuses on custom Number Set subtypes operations and functionality.
// This test validates that custom Go numeric set types work correctly with DynamoDB operations.
//
// Test Coverage:
// - Custom numeric set type marshaling/unmarshaling ([]int32, []int64, []float32, []uint64, []int16)
// - Type safety in generated structs
// - Custom types in set operations (AddToSet, RemoveFromSet)
// - Custom types in QueryBuilder and ScanBuilder with Contains filters
// - Set operations with different numeric types
//
// Schema: custom-set-number__all.json
//   - Table: "custom-set-number-all"
//   - Hash Key: id (S)
//   - Range Key: group_id (S)
//   - Common: int32_scores (NS with int32 subtype), int64_timestamps (NS with int64 subtype),
//     float32_rates (NS with float32 subtype), uint64_counters (NS with uint64 subtype),
//     int16_values (NS with int16 subtype)
func TestCustomSetNumber(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Custom Number Set subtypes on: %s", customsetnumber.TableName)

	t.Run("CustomSet_Types_Struct", func(t *testing.T) {
		testCustomSetTypesStruct(t)
	})

	t.Run("CustomSet_Types_Marshaling", func(t *testing.T) {
		testCustomSetTypesMarshaling(t, client, ctx)
	})

	t.Run("CustomSet_Types_QueryBuilder", func(t *testing.T) {
		testCustomSetTypesQueryBuilder(t, client, ctx)
	})

	t.Run("CustomSet_SetOperations", func(t *testing.T) {
		testCustomSetOperations(t, client, ctx)
	})

	t.Run("CustomSet_EdgeValues", func(t *testing.T) {
		testCustomSetEdgeValues(t, client, ctx)
	})
}

// ==================== Custom Set Types Struct Tests ====================

func testCustomSetTypesStruct(t *testing.T) {
	t.Run("verify_go_set_types", func(t *testing.T) {
		item := customsetnumber.SchemaItem{
			Id:              "test-id",
			GroupId:         "test-group",
			Int32Scores:     []int32{85, 92, 78},
			Int64Timestamps: []int64{1640995200, 1640995300, 1640995400},
			Float32Rates:    []float32{19.99, 25.50, 30.75},
			Uint64Counters:  []uint64{1000000, 2000000, 3000000},
			Int16Values:     []int16{100, 200, 300},
		}

		assert.IsType(t, "", item.Id)                     // string
		assert.IsType(t, "", item.GroupId)                // string
		assert.IsType(t, []int32{}, item.Int32Scores)     // []int32
		assert.IsType(t, []int64{}, item.Int64Timestamps) // []int64
		assert.IsType(t, []float32{}, item.Float32Rates)  // []float32
		assert.IsType(t, []uint64{}, item.Uint64Counters) // []uint64
		assert.IsType(t, []int16{}, item.Int16Values)     // []int16

		assert.Equal(t, "test-id", item.Id)
		assert.Equal(t, "test-group", item.GroupId)
		assert.Equal(t, []int32{85, 92, 78}, item.Int32Scores)
		assert.Equal(t, []int64{1640995200, 1640995300, 1640995400}, item.Int64Timestamps)
		assert.Equal(t, []float32{19.99, 25.50, 30.75}, item.Float32Rates)
		assert.Equal(t, []uint64{1000000, 2000000, 3000000}, item.Uint64Counters)
		assert.Equal(t, []int16{100, 200, 300}, item.Int16Values)
		t.Logf("✅ Custom set types verified: []int32, []int64, []float32, []uint64, []int16")
	})

	t.Run("type_safety_compilation", func(t *testing.T) {
		var item customsetnumber.SchemaItem

		item.Int32Scores = []int32{42, 100}
		item.Int64Timestamps = []int64{1640995200}
		item.Float32Rates = []float32{19.99, 29.99}
		item.Uint64Counters = []uint64{1000000}
		item.Int16Values = []int16{85, 95}

		assert.Equal(t, []int32{42, 100}, item.Int32Scores)
		assert.Equal(t, []int64{1640995200}, item.Int64Timestamps)
		assert.Equal(t, []float32{19.99, 29.99}, item.Float32Rates)
		assert.Equal(t, []uint64{1000000}, item.Uint64Counters)
		assert.Equal(t, []int16{85, 95}, item.Int16Values)
		t.Logf("✅ Type safety compilation verified")
	})
}

// ==================== Custom Set Types Marshaling Tests ====================

func testCustomSetTypesMarshaling(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("marshal_custom_set_types", func(t *testing.T) {
		item := customsetnumber.SchemaItem{
			Id:              "marshal-test-001",
			GroupId:         "marshal-group",
			Int32Scores:     []int32{85, 92, 78, 96, 88},
			Int64Timestamps: []int64{1640995200, 1640995300},
			Float32Rates:    []float32{19.99, 25.50, 30.75},
			Uint64Counters:  []uint64{1000000, 2000000},
			Int16Values:     []int16{100, 200, 300, 400},
		}
		av, err := customsetnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal item with custom set types")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "group_id", "Should contain group_id field")
		assert.Contains(t, av, "int32_scores", "Should contain int32_scores field")
		assert.Contains(t, av, "int64_timestamps", "Should contain int64_timestamps field")
		assert.Contains(t, av, "float32_rates", "Should contain float32_rates field")
		assert.Contains(t, av, "uint64_counters", "Should contain uint64_counters field")
		assert.Contains(t, av, "int16_values", "Should contain int16_values field")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[customsetnumber.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[customsetnumber.ColumnGroupId])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[customsetnumber.ColumnInt32Scores])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[customsetnumber.ColumnInt64Timestamps])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[customsetnumber.ColumnFloat32Rates])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[customsetnumber.ColumnUint64Counters])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[customsetnumber.ColumnInt16Values])
		t.Logf("✅ Custom set types marshaled successfully")
	})

	t.Run("roundtrip_custom_set_types", func(t *testing.T) {
		originalItem := customsetnumber.SchemaItem{
			Id:              "roundtrip-test-001",
			GroupId:         "roundtrip-group",
			Int32Scores:     []int32{123, 456, 789},
			Int64Timestamps: []int64{1640995300, 1640995400, 1640995500},
			Float32Rates:    []float32{29.95, 39.95},
			Uint64Counters:  []uint64{2000000, 3000000, 4000000},
			Int16Values:     []int16{95, 85, 75},
		}

		av, err := customsetnumber.ItemInput(originalItem)
		require.NoError(t, err, "Should marshal original item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(customsetnumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store item in DynamoDB")

		key, err := customsetnumber.KeyInput(originalItem)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(customsetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "roundtrip-test-001", getResult.Item[customsetnumber.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "roundtrip-group", getResult.Item[customsetnumber.ColumnGroupId].(*types.AttributeValueMemberS).Value)

		int32Set := getResult.Item[customsetnumber.ColumnInt32Scores].(*types.AttributeValueMemberNS)
		assert.Contains(t, int32Set.Value, "123")
		assert.Contains(t, int32Set.Value, "456")
		assert.Contains(t, int32Set.Value, "789")

		int64Set := getResult.Item[customsetnumber.ColumnInt64Timestamps].(*types.AttributeValueMemberNS)
		assert.Contains(t, int64Set.Value, "1640995300")
		assert.Contains(t, int64Set.Value, "1640995400")

		float32Set := getResult.Item[customsetnumber.ColumnFloat32Rates].(*types.AttributeValueMemberNS)
		assert.Contains(t, float32Set.Value, "29.95")
		assert.Contains(t, float32Set.Value, "39.95")

		uint64Set := getResult.Item[customsetnumber.ColumnUint64Counters].(*types.AttributeValueMemberNS)
		assert.Contains(t, uint64Set.Value, "2000000")
		assert.Contains(t, uint64Set.Value, "3000000")

		int16Set := getResult.Item[customsetnumber.ColumnInt16Values].(*types.AttributeValueMemberNS)
		assert.Contains(t, int16Set.Value, "95")
		assert.Contains(t, int16Set.Value, "85")
		t.Logf("✅ Custom set types roundtrip successful")
	})
}

// ==================== Custom Set Types QueryBuilder Tests ====================

func testCustomSetTypesQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupCustomSetTypesTestData(t, client, ctx)

	t.Run("custom_set_types_contains_filters", func(t *testing.T) {
		qb := customsetnumber.NewQueryBuilder().
			WithEQ("id", "query-custom-set-test").
			FilterContains("int32_scores", int32(85)).
			FilterContains("float32_rates", float32(25.0)).
			FilterContains("uint64_counters", uint64(1500000)).
			FilterContains("int16_values", int16(200))

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with custom set type contains filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		t.Logf("✅ Custom set types contains filters work in QueryBuilder")
	})

	t.Run("custom_set_types_execution", func(t *testing.T) {
		qb := customsetnumber.NewQueryBuilder().WithEQ("id", "query-custom-set-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute query with custom set types")
		assert.NotEmpty(t, items, "Should return items")

		for _, item := range items {
			assert.IsType(t, []int32{}, item.Int32Scores, "Int32Scores should be []int32")
			assert.IsType(t, []int64{}, item.Int64Timestamps, "Int64Timestamps should be []int64")
			assert.IsType(t, []float32{}, item.Float32Rates, "Float32Rates should be []float32")
			assert.IsType(t, []uint64{}, item.Uint64Counters, "Uint64Counters should be []uint64")
			assert.IsType(t, []int16{}, item.Int16Values, "Int16Values should be []int16")
		}
		t.Logf("✅ Custom set types query execution returned %d items", len(items))
	})
}

// ==================== Custom Set Operations Tests ====================

func testCustomSetOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	testItem := customsetnumber.SchemaItem{
		Id:              "set-ops-test",
		GroupId:         "operations",
		Int32Scores:     []int32{10, 20},
		Int64Timestamps: []int64{1000, 2000},
		Float32Rates:    []float32{1.5, 2.5},
		Uint64Counters:  []uint64{100, 200},
		Int16Values:     []int16{5, 10},
	}

	av, err := customsetnumber.ItemInput(testItem)
	require.NoError(t, err, "Should create test item for set operations")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(customsetnumber.TableName),
		Item:      av,
	})
	require.NoError(t, err, "Should store test item")

	t.Run("add_to_int32_set", func(t *testing.T) {
		addInput, err := customsetnumber.AddToSet("set-ops-test", "operations", "int32_scores", []int32{30, 40})
		require.NoError(t, err, "Should create add to int32 set input")

		_, err = client.UpdateItem(ctx, addInput)
		require.NoError(t, err, "Should add to int32 set")

		key, _ := customsetnumber.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(customsetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after add to int32 set")

		int32Set := getResult.Item[customsetnumber.ColumnInt32Scores].(*types.AttributeValueMemberNS)
		assert.Contains(t, int32Set.Value, "10", "Should still contain initial value")
		assert.Contains(t, int32Set.Value, "20", "Should still contain initial value")
		assert.Contains(t, int32Set.Value, "30", "Should contain added value")
		assert.Contains(t, int32Set.Value, "40", "Should contain added value")
		t.Logf("✅ Added to int32 set successfully")
	})

	t.Run("add_to_float32_set", func(t *testing.T) {
		addInput, err := customsetnumber.AddToSet("set-ops-test", "operations", "float32_rates", []float32{3.5, 4.5})
		require.NoError(t, err, "Should create add to float32 set input")

		_, err = client.UpdateItem(ctx, addInput)
		require.NoError(t, err, "Should add to float32 set")

		key, _ := customsetnumber.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(customsetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after add to float32 set")

		float32Set := getResult.Item[customsetnumber.ColumnFloat32Rates].(*types.AttributeValueMemberNS)
		assert.Contains(t, float32Set.Value, "1.5", "Should still contain initial value")
		assert.Contains(t, float32Set.Value, "2.5", "Should still contain initial value")
		assert.Contains(t, float32Set.Value, "3.5", "Should contain added value")
		assert.Contains(t, float32Set.Value, "4.5", "Should contain added value")
		t.Logf("✅ Added to float32 set successfully")
	})

	t.Run("remove_from_uint64_set", func(t *testing.T) {
		removeInput, err := customsetnumber.RemoveFromSet("set-ops-test", "operations", "uint64_counters", []uint64{100})
		require.NoError(t, err, "Should create remove from uint64 set input")

		_, err = client.UpdateItem(ctx, removeInput)
		require.NoError(t, err, "Should remove from uint64 set")

		key, _ := customsetnumber.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(customsetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after remove from uint64 set")

		uint64Set := getResult.Item[customsetnumber.ColumnUint64Counters].(*types.AttributeValueMemberNS)
		assert.Contains(t, uint64Set.Value, "200", "Should still contain remaining value")
		assert.NotContains(t, uint64Set.Value, "100", "Should not contain removed value")
		t.Logf("✅ Removed from uint64 set successfully")
	})
}

// ==================== Custom Set Edge Values Tests ====================

func testCustomSetEdgeValues(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("edge_values_custom_set_types", func(t *testing.T) {
		edgeItems := []customsetnumber.SchemaItem{
			{
				Id:              "edge-1",
				GroupId:         "zero-values",
				Int32Scores:     []int32{0},
				Int64Timestamps: []int64{0},
				Float32Rates:    []float32{0.0},
				Uint64Counters:  []uint64{0},
				Int16Values:     []int16{0},
			},
			{
				Id:              "edge-2",
				GroupId:         "max-values",
				Int32Scores:     []int32{2147483647},
				Int64Timestamps: []int64{9223372036854775807},
				Float32Rates:    []float32{3.4028235e+38},
				Uint64Counters:  []uint64{18446744073709551615},
				Int16Values:     []int16{32767},
			},
			{
				Id:              "edge-3",
				GroupId:         "min-values",
				Int32Scores:     []int32{-2147483648},
				Int64Timestamps: []int64{-9223372036854775808},
				Float32Rates:    []float32{-3.4028235e+38},
				Uint64Counters:  []uint64{0},
				Int16Values:     []int16{-32768},
			},
		}
		for _, item := range edgeItems {
			av, err := customsetnumber.ItemInput(item)
			require.NoError(t, err, "Should handle edge values for item: %s", item.Id)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(customsetnumber.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store edge value item: %s", item.Id)
		}
		t.Logf("✅ Custom set types edge values handled successfully")
	})
}

// ==================== Helper Functions ====================

func setupCustomSetTypesTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []customsetnumber.SchemaItem{
		{
			Id:              "query-custom-set-test",
			GroupId:         "group-A",
			Int32Scores:     []int32{85, 92, 78},
			Int64Timestamps: []int64{1640995400, 1640995500},
			Float32Rates:    []float32{25.0, 30.5},
			Uint64Counters:  []uint64{1500000, 2500000},
			Int16Values:     []int16{200, 300},
		},
		{
			Id:              "query-custom-set-test",
			GroupId:         "group-B",
			Int32Scores:     []int32{88, 95, 82},
			Int64Timestamps: []int64{1640995600, 1640995700},
			Float32Rates:    []float32{35.0, 40.5},
			Uint64Counters:  []uint64{3500000, 4500000},
			Int16Values:     []int16{400, 500},
		},
	}
	for _, item := range testItems {
		av, err := customsetnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal custom set types test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(customsetnumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store custom set types test item")
	}
	t.Logf("Setup complete: inserted %d custom set types test items", len(testItems))
}
