package localstack

import (
	"context"
	"maps"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	basesetnumber "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basesetnumberall"
)

// TestBaseSetNumber focuses on Number Set (NS) type operations and functionality.
// This test validates number set specific features without other data types.
//
// Test Coverage:
// - Number Set CRUD operations
// - Number Set marshaling/unmarshaling
// - Set operations (AddToSet, RemoveFromSet)
// - Set-specific filters (Contains, NotContains)
// - Set operations in Query and Scan
// - Edge cases (empty sets, large numbers, negative numbers)
//
// Schema: base-set-number__all.json
// - Table: "base-set-number-all"
// - Hash Key: user_id (S)
// - Range Key: session_id (S)
// - Common: scores (NS), ratings (NS)
func TestBaseSetNumber(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Number Set operations on: %s", basesetnumber.TableName)

	t.Run("NumberSet_Input", func(t *testing.T) {
		testNumberSetInput(t, client, ctx)
	})

	t.Run("NumberSet_Input_Raw", func(t *testing.T) {
		testNumberSetInputRaw(t, client, ctx)
	})

	t.Run("NumberSet_QueryBuilder", func(t *testing.T) {
		testNumberSetQueryBuilder(t, client, ctx)
	})

	t.Run("NumberSet_ScanBuilder", func(t *testing.T) {
		testNumberSetScanBuilder(t, client, ctx)
	})

	t.Run("NumberSet_SetOperations", func(t *testing.T) {
		testNumberSetOperations(t, client, ctx)
	})

	t.Run("NumberSet_Schema", func(t *testing.T) {
		t.Parallel()
		testNumberSetSchema(t)
	})
}

// ==================== Number Set Object Input ====================

func testNumberSetInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_number_set_item", func(t *testing.T) {
		item := basesetnumber.SchemaItem{
			UserId:    "user-001",
			SessionId: "session-2024-001",
			Scores:    []int{85, 92, 78, 96, 88},
			Ratings:   []int{4, 5, 3, 5, 4},
		}
		av, err := basesetnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal number set item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.Contains(t, av, "user_id", "Should contain user_id field")
		assert.Contains(t, av, "session_id", "Should contain session_id field")
		assert.Contains(t, av, "scores", "Should contain scores field")
		assert.Contains(t, av, "ratings", "Should contain ratings field")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[basesetnumber.ColumnUserId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basesetnumber.ColumnSessionId])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[basesetnumber.ColumnScores])
		assert.IsType(t, &types.AttributeValueMemberNS{}, av[basesetnumber.ColumnRatings])

		// Verify number set values
		scoresSet := av[basesetnumber.ColumnScores].(*types.AttributeValueMemberNS)
		assert.Len(t, scoresSet.Value, 5, "Scores set should have 5 elements")
		assert.Contains(t, scoresSet.Value, "85", "Scores should contain 85")
		assert.Contains(t, scoresSet.Value, "92", "Scores should contain 92")
		assert.Contains(t, scoresSet.Value, "96", "Scores should contain 96")

		ratingsSet := av[basesetnumber.ColumnRatings].(*types.AttributeValueMemberNS)
		assert.Len(t, ratingsSet.Value, 5, "Ratings set should have 5 elements")
		assert.Contains(t, ratingsSet.Value, "4", "Ratings should contain 4")
		assert.Contains(t, ratingsSet.Value, "5", "Ratings should contain 5")
		assert.Contains(t, ratingsSet.Value, "3", "Ratings should contain 3")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Item:      av,
		})
		t.Logf("Saved item attributes: %+v", av)
		require.NoError(t, err, "Should store number set item in DynamoDB")
		t.Logf("✅ Created number set item: %s/%s", item.UserId, item.SessionId)
	})

	t.Run("read_number_set_item", func(t *testing.T) {
		item := basesetnumber.SchemaItem{
			UserId:    "user-001",
			SessionId: "session-2024-001",
		}

		key, err := basesetnumber.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		t.Logf("Retrieved item: %+v", getResult.Item)
		t.Logf("Available keys: %v", maps.Keys(getResult.Item))
		require.NoError(t, err, "Should retrieve number set item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "user-001", getResult.Item[basesetnumber.ColumnUserId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "session-2024-001", getResult.Item[basesetnumber.ColumnSessionId].(*types.AttributeValueMemberS).Value)

		scoresSet := getResult.Item[basesetnumber.ColumnScores].(*types.AttributeValueMemberNS)
		assert.Len(t, scoresSet.Value, 5, "Scores set should have 5 elements")
		assert.Contains(t, scoresSet.Value, "85", "Scores should contain 85")

		ratingsSet := getResult.Item[basesetnumber.ColumnRatings].(*types.AttributeValueMemberNS)
		assert.Len(t, ratingsSet.Value, 5, "Ratings set should have 5 elements")
		assert.Contains(t, ratingsSet.Value, "4", "Ratings should contain 4")
		t.Logf("✅ Retrieved number set item successfully")
	})

	t.Run("update_number_set_item", func(t *testing.T) {
		item := basesetnumber.SchemaItem{
			UserId:    "user-001",
			SessionId: "session-2024-001",
			Scores:    []int{90, 95, 88, 92},
			Ratings:   []int{5, 4, 5},
		}
		updateInput, err := basesetnumber.UpdateItemInput(item)
		require.NoError(t, err, "Should create update input from item")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update number set item")

		key, _ := basesetnumber.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		scoresSet := getResult.Item[basesetnumber.ColumnScores].(*types.AttributeValueMemberNS)
		assert.Len(t, scoresSet.Value, 4, "Scores set should have 4 elements after update")
		assert.Contains(t, scoresSet.Value, "90", "Scores should contain 90")
		assert.Contains(t, scoresSet.Value, "95", "Scores should contain 95")

		ratingsSet := getResult.Item[basesetnumber.ColumnRatings].(*types.AttributeValueMemberNS)
		assert.Len(t, ratingsSet.Value, 3, "Ratings set should have 3 elements after update")
		assert.Contains(t, ratingsSet.Value, "5", "Ratings should contain 5")
		assert.NotContains(t, ratingsSet.Value, "3", "Ratings should not contain 3")
		t.Logf("✅ Updated number set item successfully")
	})

	t.Run("delete_number_set_item", func(t *testing.T) {
		item := basesetnumber.SchemaItem{
			UserId:    "user-001",
			SessionId: "session-2024-001",
		}
		deleteInput, err := basesetnumber.DeleteItemInput(item)
		require.NoError(t, err, "Should create delete input from item")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete number set item")

		key, _ := basesetnumber.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "Number set item should be deleted")
		t.Logf("✅ Deleted number set item successfully")
	})

	t.Run("number_set_edge_cases", func(t *testing.T) {
		edgeCases := []basesetnumber.SchemaItem{
			{
				UserId:    "edge-1",
				SessionId: "zero-test",
				Scores:    []int{0},
				Ratings:   []int{0, 1, 2},
			},
			{
				UserId:    "edge-2",
				SessionId: "negative-test",
				Scores:    []int{-100, -50, 0, 50, 100},
				Ratings:   []int{1, 2, 3, 4, 5},
			},
			{
				UserId:    "edge-3",
				SessionId: "large-test",
				Scores:    []int{999999, 1000000, 2147483647},
				Ratings:   []int{10, 20, 30, 40, 50, 60, 70, 80, 90, 100},
			},
		}
		for _, item := range edgeCases {
			av, err := basesetnumber.ItemInput(item)
			require.NoError(t, err, "Should handle number set edge case: %s", item.UserId)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(basesetnumber.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store number set edge case item: %s", item.UserId)
		}
		t.Logf("✅ Number set edge cases handled successfully")
	})
}

// ==================== Number Set Raw Object Input ====================

func testNumberSetInputRaw(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_number_set_item_raw", func(t *testing.T) {
		item := basesetnumber.SchemaItem{
			UserId:    "user-raw-001",
			SessionId: "session-raw-001",
			Scores:    []int{75, 82, 89, 91},
			Ratings:   []int{3, 4, 4, 5},
		}
		av, err := basesetnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal number set item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store number set item in DynamoDB")
		t.Logf("✅ Created number set item for raw testing: %s/%s", item.UserId, item.SessionId)
	})

	t.Run("read_number_set_item_raw", func(t *testing.T) {
		key, err := basesetnumber.KeyInputFromRaw("user-raw-001", "session-raw-001")
		require.NoError(t, err, "Should create key from raw values")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve number set item using raw key")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "user-raw-001", getResult.Item[basesetnumber.ColumnUserId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "session-raw-001", getResult.Item[basesetnumber.ColumnSessionId].(*types.AttributeValueMemberS).Value)

		scoresSet := getResult.Item[basesetnumber.ColumnScores].(*types.AttributeValueMemberNS)
		assert.Contains(t, scoresSet.Value, "75", "Scores should contain 75")
		assert.Contains(t, scoresSet.Value, "91", "Scores should contain 91")

		t.Logf("✅ Retrieved number set item successfully using raw key")
	})

	t.Run("update_number_set_item_raw", func(t *testing.T) {
		extractNumberValues := func(attr types.AttributeValue) []string {
			switch v := attr.(type) {
			case *types.AttributeValueMemberNS:
				return v.Value
			case *types.AttributeValueMemberL:
				var values []string
				for _, item := range v.Value {
					if n, ok := item.(*types.AttributeValueMemberN); ok {
						values = append(values, n.Value)
					}
				}
				return values
			default:
				t.Fatalf("Unexpected type for number set/list: %T", attr)
				return nil
			}
		}
		updates := map[string]any{
			"scores":  []int{80, 85, 90, 95, 100},
			"ratings": []int{4, 5, 5},
		}

		updateInput, err := basesetnumber.UpdateItemInputFromRaw("user-raw-001", "session-raw-001", updates)
		require.NoError(t, err, "Should create update input from raw values")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update number set item using raw method")

		key, _ := basesetnumber.KeyInputFromRaw("user-raw-001", "session-raw-001")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		scoresValues := extractNumberValues(getResult.Item[basesetnumber.ColumnScores])
		assert.Contains(t, scoresValues, "80", "Scores should contain 80")
		assert.Contains(t, scoresValues, "100", "Scores should contain 100")

		ratingsValues := extractNumberValues(getResult.Item[basesetnumber.ColumnRatings])
		assert.Contains(t, ratingsValues, "4", "Ratings should contain 4")
		assert.Contains(t, ratingsValues, "5", "Ratings should contain 5")
		t.Logf("✅ Updated number set item successfully using raw method")
	})

	t.Run("delete_number_set_item_raw", func(t *testing.T) {
		deleteInput, err := basesetnumber.DeleteItemInputFromRaw("user-raw-001", "session-raw-001")
		require.NoError(t, err, "Should create delete input from raw values")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete number set item using raw method")

		key, _ := basesetnumber.KeyInputFromRaw("user-raw-001", "session-raw-001")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "Number set item should be deleted")
		t.Logf("✅ Deleted number set item successfully using raw method")
	})

	t.Run("raw_vs_object_number_set_comparison", func(t *testing.T) {
		keyFromRaw, err := basesetnumber.KeyInputFromRaw("comparison-test", "both-methods")
		require.NoError(t, err, "Should create key from raw values")

		item := basesetnumber.SchemaItem{
			UserId:    "comparison-test",
			SessionId: "both-methods",
		}
		keyFromObject, err := basesetnumber.KeyInput(item)
		require.NoError(t, err, "Should create key from object")

		assert.Equal(t, keyFromRaw, keyFromObject, "Raw and object-based keys should be identical")
		t.Logf("✅ Raw and object-based number set methods produce identical results")
	})

	t.Run("raw_number_set_edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			userId    string
			sessionId string
			updates   map[string]any
		}{
			{
				userId:    "raw-edge-1",
				sessionId: "zero-numbers",
				updates:   map[string]any{"scores": []int{0}, "ratings": []int{0, 1}},
			},
			{
				userId:    "raw-edge-2",
				sessionId: "negative-numbers",
				updates:   map[string]any{"scores": []int{-100, -1, 0, 1, 100}, "ratings": []int{1, 2, 3}},
			},
		}
		for _, edgeCase := range edgeCases {
			updateInput, err := basesetnumber.UpdateItemInputFromRaw(edgeCase.userId, edgeCase.sessionId, edgeCase.updates)
			require.NoError(t, err, "Should handle raw number set edge case: %s", edgeCase.userId)
			assert.NotNil(t, updateInput, "Update input should be created for edge case: %s", edgeCase.userId)

			deleteInput, err := basesetnumber.DeleteItemInputFromRaw(edgeCase.userId, edgeCase.sessionId)
			require.NoError(t, err, "Should create delete input for edge case: %s", edgeCase.userId)
			assert.NotNil(t, deleteInput, "Delete input should be created")
		}
		t.Logf("✅ Raw number set edge cases handled successfully")
	})
}

// ==================== Number Set QueryBuilder Tests ====================

func testNumberSetQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupNumberSetTestData(t, client, ctx)

	t.Run("number_set_hash_key_query", func(t *testing.T) {
		qb := basesetnumber.NewQueryBuilder().WithEQ("user_id", "query-set-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build number set hash key query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		assert.Equal(t, basesetnumber.TableName, *queryInput.TableName, "Should target correct table")
		t.Logf("✅ Number set hash key query built successfully")
	})

	t.Run("number_set_contains_filters", func(t *testing.T) {
		qb := basesetnumber.NewQueryBuilder().
			WithEQ("user_id", "query-set-test").
			FilterContains("scores", 85).
			FilterContains("ratings", 5)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with number set contains filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		t.Logf("✅ Number set contains filters query built successfully")
	})

	t.Run("number_set_query_execution", func(t *testing.T) {
		qb := basesetnumber.NewQueryBuilder().WithEQ("user_id", "query-set-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute number set query")
		assert.NotEmpty(t, items, "Should return number set items")

		for _, item := range items {
			assert.Equal(t, "query-set-test", item.UserId, "All items should have correct hash key")
			assert.NotEmpty(t, item.SessionId, "All items should have session ID")
			assert.IsType(t, []int{}, item.Scores, "Scores should be int slice type")
			assert.IsType(t, []int{}, item.Ratings, "Ratings should be int slice type")
		}
		t.Logf("✅ Number set query execution returned %d items", len(items))
	})

	t.Run("number_set_sorting", func(t *testing.T) {
		qbAsc := basesetnumber.NewQueryBuilder().
			WithEQ("user_id", "query-set-test").
			OrderByAsc()

		itemsAsc, err := qbAsc.Execute(ctx, client)
		require.NoError(t, err, "Should execute ascending number set query")

		qbDesc := basesetnumber.NewQueryBuilder().
			WithEQ("user_id", "query-set-test").
			OrderByDesc()

		itemsDesc, err := qbDesc.Execute(ctx, client)
		require.NoError(t, err, "Should execute descending number set query")

		if len(itemsAsc) > 1 && len(itemsDesc) > 1 {
			assert.NotEqual(t, itemsAsc[0].SessionId, itemsDesc[0].SessionId, "Number set sorting should produce different order")
		}
		t.Logf("✅ Number set sorting works correctly")
	})
}

// ==================== Number Set ScanBuilder Tests ====================

func testNumberSetScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("number_set_scan_contains", func(t *testing.T) {
		sb := basesetnumber.NewScanBuilder().
			FilterContains("scores", 85)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with number set contains filter")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")
		t.Logf("✅ Number set contains scan filter built successfully")
	})

	t.Run("number_set_multiple_contains", func(t *testing.T) {
		sb := basesetnumber.NewScanBuilder().
			FilterContains("scores", 90).
			FilterContains("ratings", 5)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with multiple number set contains filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")
		t.Logf("✅ Multiple number set contains filters built successfully")
	})

	t.Run("number_set_scan_execution", func(t *testing.T) {
		sb := basesetnumber.NewScanBuilder().
			FilterContains("scores", 85).
			Limit(5)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute number set scan")

		for _, item := range items {
			assert.Contains(t, item.Scores, 85, "Items should match contains filter")
		}
		t.Logf("✅ Number set scan execution returned %d items", len(items))
	})
}

// ==================== Number Set Operations Tests ====================

func testNumberSetOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup item for set operations testing
	testItem := basesetnumber.SchemaItem{
		UserId:    "set-ops-test",
		SessionId: "operations",
		Scores:    []int{10, 20},
		Ratings:   []int{1},
	}

	av, err := basesetnumber.ItemInput(testItem)
	require.NoError(t, err, "Should create test item for set operations")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basesetnumber.TableName),
		Item:      av,
	})
	require.NoError(t, err, "Should store test item")

	t.Run("add_to_number_set", func(t *testing.T) {
		addInput, err := basesetnumber.AddToSet("set-ops-test", "operations", "scores", []int{30, 40})
		require.NoError(t, err, "Should create add to set input")

		_, err = client.UpdateItem(ctx, addInput)
		require.NoError(t, err, "Should add to number set")

		key, _ := basesetnumber.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after add to set")

		scoresSet := getResult.Item[basesetnumber.ColumnScores].(*types.AttributeValueMemberNS)
		assert.Contains(t, scoresSet.Value, "10", "Should still contain initial value")
		assert.Contains(t, scoresSet.Value, "20", "Should still contain initial value")
		assert.Contains(t, scoresSet.Value, "30", "Should contain added value")
		assert.Contains(t, scoresSet.Value, "40", "Should contain added value")
		t.Logf("✅ Added to number set successfully")
	})

	t.Run("remove_from_number_set", func(t *testing.T) {
		removeInput, err := basesetnumber.RemoveFromSet("set-ops-test", "operations", "scores", []int{20, 30})
		require.NoError(t, err, "Should create remove from set input")

		_, err = client.UpdateItem(ctx, removeInput)
		require.NoError(t, err, "Should remove from number set")

		key, _ := basesetnumber.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after remove from set")

		scoresSet := getResult.Item[basesetnumber.ColumnScores].(*types.AttributeValueMemberNS)
		assert.Contains(t, scoresSet.Value, "10", "Should still contain initial value")
		assert.Contains(t, scoresSet.Value, "40", "Should still contain added value")
		assert.NotContains(t, scoresSet.Value, "20", "Should not contain removed value")
		assert.NotContains(t, scoresSet.Value, "30", "Should not contain removed value")
		t.Logf("✅ Removed from number set successfully")
	})

	t.Run("add_to_ratings_set", func(t *testing.T) {
		addInput, err := basesetnumber.AddToSet("set-ops-test", "operations", "ratings", []int{2, 3, 4, 5})
		require.NoError(t, err, "Should create add to ratings set input")

		_, err = client.UpdateItem(ctx, addInput)
		require.NoError(t, err, "Should add to ratings set")

		key, _ := basesetnumber.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after add to ratings set")

		ratingsSet := getResult.Item[basesetnumber.ColumnRatings].(*types.AttributeValueMemberNS)
		assert.Contains(t, ratingsSet.Value, "1", "Should still contain rating 1")
		assert.Contains(t, ratingsSet.Value, "2", "Should contain added rating 2")
		assert.Contains(t, ratingsSet.Value, "5", "Should contain added rating 5")
		t.Logf("✅ Added to ratings set successfully")
	})
}

// ==================== Number Set Schema Tests ====================

func testNumberSetSchema(t *testing.T) {
	t.Run("number_set_table_schema", func(t *testing.T) {
		schema := basesetnumber.TableSchema

		assert.Equal(t, "base-set-number-all", schema.TableName, "Table name should match")
		assert.Equal(t, "user_id", schema.HashKey, "Hash key should be 'user_id'")
		assert.Equal(t, "session_id", schema.RangeKey, "Range key should be 'session_id'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")
		t.Logf("✅ Number set schema structure validated")
	})

	t.Run("number_set_attributes", func(t *testing.T) {
		expectedPrimary := map[string]string{
			"session_id": "S",
		}
		for _, attr := range basesetnumber.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should have correct type", attr.Name)
		}
		expectedCommon := map[string]string{
			"scores":  "NS",
			"ratings": "NS",
		}
		for _, attr := range basesetnumber.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be number set type", attr.Name)
		}
		t.Logf("✅ Number set attributes validated")
	})

	t.Run("number_set_constants", func(t *testing.T) {
		assert.Equal(t, "base-set-number-all", basesetnumber.TableName, "TableName constant should be correct")
		assert.Equal(t, "user_id", basesetnumber.ColumnUserId, "ColumnUserId should be correct")
		assert.Equal(t, "session_id", basesetnumber.ColumnSessionId, "ColumnSessionId should be correct")
		assert.Equal(t, "scores", basesetnumber.ColumnScores, "ColumnScores should be correct")
		assert.Equal(t, "ratings", basesetnumber.ColumnRatings, "ColumnRatings should be correct")
		t.Logf("✅ Number set constants validated")
	})

	t.Run("number_set_attribute_names", func(t *testing.T) {
		attrs := basesetnumber.AttributeNames
		expectedAttrs := []string{"user_id", "session_id", "scores", "ratings"}

		assert.Len(t, attrs, len(expectedAttrs), "Should have correct number of attributes")
		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}
		t.Logf("✅ Number set AttributeNames validated")
	})

	t.Run("number_set_go_types", func(t *testing.T) {
		item := basesetnumber.SchemaItem{}
		item.Scores = []int{1, 2, 3}
		item.Ratings = []int{4, 5}

		assert.IsType(t, []int{}, item.Scores, "Scores should be []int type")
		assert.IsType(t, []int{}, item.Ratings, "Ratings should be []int type")
		assert.IsType(t, "", item.UserId, "UserId should be string type")
		assert.IsType(t, "", item.SessionId, "SessionId should be string type")
		t.Logf("✅ Number set Go types validated")
	})

	t.Run("number_set_edge_values", func(t *testing.T) {
		item := basesetnumber.SchemaItem{
			UserId:    "edge-test",
			SessionId: "edge-session",
			Scores:    []int{-2147483648, 0, 2147483647},
			Ratings:   []int{-100, -1, 0, 1, 100},
		}
		assert.Len(t, item.Scores, 3, "Should handle edge score values")
		assert.Len(t, item.Ratings, 5, "Should handle edge rating values")
		assert.Contains(t, item.Scores, 0, "Should handle zero value")
		assert.Contains(t, item.Scores, -2147483648, "Should handle negative values")
		assert.Contains(t, item.Ratings, -100, "Should handle negative ratings")
		t.Logf("✅ Number set edge values validated")
	})
}

// ==================== Helper Functions ====================

func setupNumberSetTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basesetnumber.SchemaItem{
		{
			UserId:    "query-set-test",
			SessionId: "session-A",
			Scores:    []int{85, 90, 78, 92},
			Ratings:   []int{4, 5, 3, 5},
		},
		{
			UserId:    "query-set-test",
			SessionId: "session-B",
			Scores:    []int{88, 95, 82, 91},
			Ratings:   []int{5, 4, 4, 5},
		},
		{
			UserId:    "query-set-test",
			SessionId: "session-C",
			Scores:    []int{76, 89, 85, 93, 87},
			Ratings:   []int{3, 4, 5, 5, 4},
		},
	}
	for _, item := range testItems {
		av, err := basesetnumber.ItemInput(item)
		require.NoError(t, err, "Should marshal number set test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basesetnumber.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store number set test item")
	}
	t.Logf("Setup complete: inserted %d number set test items", len(testItems))
}
