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

	basesetstring "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basesetstring"
)

// TestBaseSetString focuses on String Set (SS) type operations and functionality.
// This test validates string set specific features without other data types.
//
// Test Coverage:
// - String Set CRUD operations
// - String Set marshaling/unmarshaling
// - Set operations (AddToSet, RemoveFromSet)
// - Set-specific filters (Contains, NotContains)
// - Set operations in Query and Scan
// - Edge cases (empty sets, large sets, duplicates)
//
// Schema: base-set-string.json
// - Table: "base-set-string"
// - Hash Key: id (S)
// - Range Key: group_id (S)
// - Common: tags (SS), categories (SS)
func TestBaseSetString(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing String Set operations on: %s", basesetstring.TableName)

	t.Run("StringSet_Input", func(t *testing.T) {
		testStringSetInput(t, client, ctx)
	})

	t.Run("StringSet_Input_Raw", func(t *testing.T) {
		testStringSetInputRaw(t, client, ctx)
	})

	t.Run("StringSet_QueryBuilder", func(t *testing.T) {
		testStringSetQueryBuilder(t, client, ctx)
	})

	t.Run("StringSet_ScanBuilder", func(t *testing.T) {
		testStringSetScanBuilder(t, client, ctx)
	})

	t.Run("StringSet_SetOperations", func(t *testing.T) {
		testStringSetOperations(t, client, ctx)
	})

	t.Run("StringSet_Schema", func(t *testing.T) {
		t.Parallel()
		testStringSetSchema(t)
	})
}

// ==================== String Set Object Input ====================

func testStringSetInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_string_set_item", func(t *testing.T) {
		item := basesetstring.SchemaItem{
			Id:         "set-test-001",
			GroupId:    "web-development",
			Tags:       []string{"javascript", "react", "nodejs"},
			Categories: []string{"frontend", "backend", "fullstack"},
		}

		av, err := basesetstring.ItemInput(item)
		require.NoError(t, err, "Should marshal string set item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "group_id", "Should contain group_id field")
		assert.Contains(t, av, "tags", "Should contain tags field")
		assert.Contains(t, av, "categories", "Should contain categories field")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[basesetstring.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basesetstring.ColumnGroupId])
		assert.IsType(t, &types.AttributeValueMemberSS{}, av[basesetstring.ColumnTags])
		assert.IsType(t, &types.AttributeValueMemberSS{}, av[basesetstring.ColumnCategories])

		// Verify string set values
		tagsSet := av[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Len(t, tagsSet.Value, 3, "Tags set should have 3 elements")
		assert.Contains(t, tagsSet.Value, "javascript", "Tags should contain javascript")
		assert.Contains(t, tagsSet.Value, "react", "Tags should contain react")
		assert.Contains(t, tagsSet.Value, "nodejs", "Tags should contain nodejs")

		categoriesSet := av[basesetstring.ColumnCategories].(*types.AttributeValueMemberSS)
		assert.Len(t, categoriesSet.Value, 3, "Categories set should have 3 elements")
		assert.Contains(t, categoriesSet.Value, "frontend", "Categories should contain frontend")
		assert.Contains(t, categoriesSet.Value, "backend", "Categories should contain backend")
		assert.Contains(t, categoriesSet.Value, "fullstack", "Categories should contain fullstack")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basesetstring.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store string set item in DynamoDB")

		t.Logf("✅ Created string set item: %s/%s", item.Id, item.GroupId)
	})

	t.Run("read_string_set_item", func(t *testing.T) {
		item := basesetstring.SchemaItem{
			Id:      "set-test-001",
			GroupId: "web-development",
		}

		key, err := basesetstring.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve string set item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "set-test-001", getResult.Item[basesetstring.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "web-development", getResult.Item[basesetstring.ColumnGroupId].(*types.AttributeValueMemberS).Value)

		// Verify string sets
		tagsSet := getResult.Item[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Len(t, tagsSet.Value, 3, "Tags set should have 3 elements")
		assert.Contains(t, tagsSet.Value, "javascript", "Tags should contain javascript")

		categoriesSet := getResult.Item[basesetstring.ColumnCategories].(*types.AttributeValueMemberSS)
		assert.Len(t, categoriesSet.Value, 3, "Categories set should have 3 elements")
		assert.Contains(t, categoriesSet.Value, "frontend", "Categories should contain frontend")

		t.Logf("✅ Retrieved string set item successfully")
	})

	t.Run("update_string_set_item", func(t *testing.T) {
		item := basesetstring.SchemaItem{
			Id:         "set-test-001",
			GroupId:    "web-development",
			Tags:       []string{"javascript", "react", "typescript", "vue"},
			Categories: []string{"frontend", "backend"},
		}

		updateInput, err := basesetstring.UpdateItemInput(item)
		require.NoError(t, err, "Should create update input from item")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update string set item")

		key, _ := basesetstring.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		// Verify updated string sets
		tagsSet := getResult.Item[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Len(t, tagsSet.Value, 4, "Tags set should have 4 elements after update")
		assert.Contains(t, tagsSet.Value, "typescript", "Tags should contain typescript")
		assert.Contains(t, tagsSet.Value, "vue", "Tags should contain vue")

		categoriesSet := getResult.Item[basesetstring.ColumnCategories].(*types.AttributeValueMemberSS)
		assert.Len(t, categoriesSet.Value, 2, "Categories set should have 2 elements after update")
		assert.NotContains(t, categoriesSet.Value, "fullstack", "Categories should not contain fullstack")

		t.Logf("✅ Updated string set item successfully")
	})

	t.Run("delete_string_set_item", func(t *testing.T) {
		item := basesetstring.SchemaItem{
			Id:      "set-test-001",
			GroupId: "web-development",
		}

		deleteInput, err := basesetstring.DeleteItemInput(item)
		require.NoError(t, err, "Should create delete input from item")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete string set item")

		key, _ := basesetstring.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "String set item should be deleted")

		t.Logf("✅ Deleted string set item successfully")
	})

	t.Run("string_set_edge_cases", func(t *testing.T) {
		edgeCases := []basesetstring.SchemaItem{
			{
				Id:         "edge-1",
				GroupId:    "empty-test",
				Tags:       []string{}, // Empty set
				Categories: []string{"single"},
			},
			{
				Id:         "edge-2",
				GroupId:    "special",
				Tags:       []string{"tag with spaces", "special-chars!@#", "unicode-тест"},
				Categories: []string{"UPPERCASE", "lowercase", "MiXeD"},
			},
			{
				Id:         "edge-3",
				GroupId:    "large",
				Tags:       []string{"tag1", "tag2", "tag3", "tag4", "tag5", "tag6", "tag7", "tag8", "tag9", "tag10"},
				Categories: []string{"category-with-very-long-name-that-tests-string-length-limits"},
			},
		}

		for _, item := range edgeCases {
			av, err := basesetstring.ItemInput(item)
			require.NoError(t, err, "Should handle string set edge case: %s", item.Id)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(basesetstring.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store string set edge case item: %s", item.Id)
		}

		t.Logf("✅ String set edge cases handled successfully")
	})
}

// ==================== String Set Raw Object Input ====================

func testStringSetInputRaw(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_string_set_item_raw", func(t *testing.T) {
		item := basesetstring.SchemaItem{
			Id:         "set-raw-001",
			GroupId:    "mobile-development",
			Tags:       []string{"swift", "kotlin", "react-native"},
			Categories: []string{"ios", "android", "cross-platform"},
		}

		av, err := basesetstring.ItemInput(item)
		require.NoError(t, err, "Should marshal string set item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basesetstring.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store string set item in DynamoDB")

		t.Logf("✅ Created string set item for raw testing: %s/%s", item.Id, item.GroupId)
	})

	t.Run("read_string_set_item_raw", func(t *testing.T) {
		key, err := basesetstring.KeyInputFromRaw("set-raw-001", "mobile-development")
		require.NoError(t, err, "Should create key from raw values")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve string set item using raw key")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "set-raw-001", getResult.Item[basesetstring.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "mobile-development", getResult.Item[basesetstring.ColumnGroupId].(*types.AttributeValueMemberS).Value)

		// Verify string sets
		tagsSet := getResult.Item[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Contains(t, tagsSet.Value, "swift", "Tags should contain swift")
		assert.Contains(t, tagsSet.Value, "kotlin", "Tags should contain kotlin")

		t.Logf("✅ Retrieved string set item successfully using raw key")
	})

	t.Run("update_string_set_item_raw", func(t *testing.T) {
		updates := map[string]any{
			"tags":       []string{"swift", "kotlin", "flutter", "xamarin"},
			"categories": []string{"native", "hybrid"},
		}

		updateInput, err := basesetstring.UpdateItemInputFromRaw("set-raw-001", "mobile-development", updates)
		require.NoError(t, err, "Should create update input from raw values")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update string set item using raw method")

		key, _ := basesetstring.KeyInputFromRaw("set-raw-001", "mobile-development")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		// Verify updated sets
		tagsSet := getResult.Item[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Contains(t, tagsSet.Value, "flutter", "Tags should contain flutter")
		assert.Contains(t, tagsSet.Value, "xamarin", "Tags should contain xamarin")

		categoriesSet := getResult.Item[basesetstring.ColumnCategories].(*types.AttributeValueMemberSS)
		assert.Contains(t, categoriesSet.Value, "native", "Categories should contain native")
		assert.Contains(t, categoriesSet.Value, "hybrid", "Categories should contain hybrid")

		t.Logf("✅ Updated string set item successfully using raw method")
	})

	t.Run("delete_string_set_item_raw", func(t *testing.T) {
		deleteInput, err := basesetstring.DeleteItemInputFromRaw("set-raw-001", "mobile-development")
		require.NoError(t, err, "Should create delete input from raw values")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete string set item using raw method")

		key, _ := basesetstring.KeyInputFromRaw("set-raw-001", "mobile-development")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "String set item should be deleted")

		t.Logf("✅ Deleted string set item successfully using raw method")
	})

	t.Run("raw_vs_object_string_set_comparison", func(t *testing.T) {
		keyFromRaw, err := basesetstring.KeyInputFromRaw("comparison-test", "both-methods")
		require.NoError(t, err, "Should create key from raw values")

		item := basesetstring.SchemaItem{
			Id:      "comparison-test",
			GroupId: "both-methods",
		}
		keyFromObject, err := basesetstring.KeyInput(item)
		require.NoError(t, err, "Should create key from object")

		assert.Equal(t, keyFromRaw, keyFromObject, "Raw and object-based keys should be identical")

		t.Logf("✅ Raw and object-based string set methods produce identical results")
	})

	t.Run("raw_string_set_edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			id      string
			groupId string
			updates map[string]any
		}{
			{
				id:      "raw-edge-1",
				groupId: "empty-sets",
				updates: map[string]any{"tags": []string{}, "categories": []string{"solo"}},
			},
			{
				id:      "raw-edge-2",
				groupId: "special-chars",
				updates: map[string]any{"tags": []string{"tag!@#", "спец-символы"}, "categories": []string{"CAPS", "lower"}},
			},
		}

		for _, edgeCase := range edgeCases {
			updateInput, err := basesetstring.UpdateItemInputFromRaw(edgeCase.id, edgeCase.groupId, edgeCase.updates)
			require.NoError(t, err, "Should handle raw string set edge case: %s", edgeCase.id)
			assert.NotNil(t, updateInput, "Update input should be created for edge case: %s", edgeCase.id)

			deleteInput, err := basesetstring.DeleteItemInputFromRaw(edgeCase.id, edgeCase.groupId)
			require.NoError(t, err, "Should create delete input for edge case: %s", edgeCase.id)
			assert.NotNil(t, deleteInput, "Delete input should be created")
		}

		t.Logf("✅ Raw string set edge cases handled successfully")
	})
}

// ==================== String Set QueryBuilder Tests ====================

func testStringSetQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupStringSetTestData(t, client, ctx)

	t.Run("string_set_hash_key_query", func(t *testing.T) {
		qb := basesetstring.NewQueryBuilder().WithEQ("id", "query-set-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build string set hash key query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		assert.Equal(t, basesetstring.TableName, *queryInput.TableName, "Should target correct table")

		t.Logf("✅ String set hash key query built successfully")
	})

	t.Run("string_set_contains_filters", func(t *testing.T) {
		qb := basesetstring.NewQueryBuilder().
			WithEQ("id", "query-set-test").
			FilterContains("tags", "react").
			FilterContains("categories", "frontend")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with string set contains filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ String set contains filters query built successfully")
	})

	t.Run("string_set_query_execution", func(t *testing.T) {
		qb := basesetstring.NewQueryBuilder().WithEQ("id", "query-set-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute string set query")
		assert.NotEmpty(t, items, "Should return string set items")

		for _, item := range items {
			assert.Equal(t, "query-set-test", item.Id, "All items should have correct hash key")
			assert.NotEmpty(t, item.GroupId, "All items should have group ID")
			assert.IsType(t, []string{}, item.Tags, "Tags should be string slice type")
			assert.IsType(t, []string{}, item.Categories, "Categories should be string slice type")
		}

		t.Logf("✅ String set query execution returned %d items", len(items))
	})

	t.Run("string_set_sorting", func(t *testing.T) {
		qbAsc := basesetstring.NewQueryBuilder().
			WithEQ("id", "query-set-test").
			OrderByAsc()

		itemsAsc, err := qbAsc.Execute(ctx, client)
		require.NoError(t, err, "Should execute ascending string set query")

		qbDesc := basesetstring.NewQueryBuilder().
			WithEQ("id", "query-set-test").
			OrderByDesc()

		itemsDesc, err := qbDesc.Execute(ctx, client)
		require.NoError(t, err, "Should execute descending string set query")

		if len(itemsAsc) > 1 && len(itemsDesc) > 1 {
			assert.NotEqual(t, itemsAsc[0].GroupId, itemsDesc[0].GroupId, "String set sorting should produce different order")
		}

		t.Logf("✅ String set sorting works correctly")
	})
}

// ==================== String Set ScanBuilder Tests ====================

func testStringSetScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("string_set_scan_contains", func(t *testing.T) {
		sb := basesetstring.NewScanBuilder().
			FilterContains("tags", "javascript")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with string set contains filter")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ String set contains scan filter built successfully")
	})

	t.Run("string_set_multiple_contains", func(t *testing.T) {
		sb := basesetstring.NewScanBuilder().
			FilterContains("tags", "react").
			FilterContains("categories", "frontend")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with multiple string set contains filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ Multiple string set contains filters built successfully")
	})

	t.Run("string_set_scan_execution", func(t *testing.T) {
		sb := basesetstring.NewScanBuilder().
			FilterContains("tags", "javascript").
			Limit(5)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute string set scan")

		for _, item := range items {
			assert.Contains(t, item.Tags, "javascript", "Items should match contains filter")
		}

		t.Logf("✅ String set scan execution returned %d items", len(items))
	})
}

// ==================== String Set Operations Tests ====================

func testStringSetOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup item for set operations testing
	testItem := basesetstring.SchemaItem{
		Id:         "set-ops-test",
		GroupId:    "operations",
		Tags:       []string{"initial", "test"},
		Categories: []string{"category1"},
	}

	av, err := basesetstring.ItemInput(testItem)
	require.NoError(t, err, "Should create test item for set operations")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(basesetstring.TableName),
		Item:      av,
	})
	require.NoError(t, err, "Should store test item")

	t.Run("add_to_string_set", func(t *testing.T) {
		addInput, err := basesetstring.AddToSet("set-ops-test", "operations", "tags", []string{"added1", "added2"})
		require.NoError(t, err, "Should create add to set input")

		_, err = client.UpdateItem(ctx, addInput)
		require.NoError(t, err, "Should add to string set")

		key, _ := basesetstring.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after add to set")

		tagsSet := getResult.Item[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Contains(t, tagsSet.Value, "initial", "Should still contain initial value")
		assert.Contains(t, tagsSet.Value, "test", "Should still contain test value")
		assert.Contains(t, tagsSet.Value, "added1", "Should contain added1")
		assert.Contains(t, tagsSet.Value, "added2", "Should contain added2")

		t.Logf("✅ Added to string set successfully")
	})

	t.Run("remove_from_string_set", func(t *testing.T) {
		removeInput, err := basesetstring.RemoveFromSet("set-ops-test", "operations", "tags", []string{"test", "added1"})
		require.NoError(t, err, "Should create remove from set input")

		_, err = client.UpdateItem(ctx, removeInput)
		require.NoError(t, err, "Should remove from string set")

		key, _ := basesetstring.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after remove from set")

		tagsSet := getResult.Item[basesetstring.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Contains(t, tagsSet.Value, "initial", "Should still contain initial value")
		assert.Contains(t, tagsSet.Value, "added2", "Should still contain added2")
		assert.NotContains(t, tagsSet.Value, "test", "Should not contain removed test")
		assert.NotContains(t, tagsSet.Value, "added1", "Should not contain removed added1")

		t.Logf("✅ Removed from string set successfully")
	})

	t.Run("add_to_categories_set", func(t *testing.T) {
		addInput, err := basesetstring.AddToSet("set-ops-test", "operations", "categories", []string{"category2", "category3"})
		require.NoError(t, err, "Should create add to categories set input")

		_, err = client.UpdateItem(ctx, addInput)
		require.NoError(t, err, "Should add to categories set")

		key, _ := basesetstring.KeyInputFromRaw("set-ops-test", "operations")
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basesetstring.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve item after add to categories set")

		categoriesSet := getResult.Item[basesetstring.ColumnCategories].(*types.AttributeValueMemberSS)
		assert.Contains(t, categoriesSet.Value, "category1", "Should still contain category1")
		assert.Contains(t, categoriesSet.Value, "category2", "Should contain added category2")
		assert.Contains(t, categoriesSet.Value, "category3", "Should contain added category3")

		t.Logf("✅ Added to categories set successfully")
	})
}

// ==================== String Set Schema Tests ====================

func testStringSetSchema(t *testing.T) {
	t.Run("string_set_table_schema", func(t *testing.T) {
		schema := basesetstring.TableSchema

		assert.Equal(t, "base-set-string", schema.TableName, "Table name should match")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "group_id", schema.RangeKey, "Range key should be 'group_id'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")

		t.Logf("✅ String set schema structure validated")
	})

	t.Run("string_set_attributes", func(t *testing.T) {
		expectedPrimary := map[string]string{
			"id":       "S",
			"group_id": "S",
		}

		for _, attr := range basesetstring.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should have correct type", attr.Name)
		}

		expectedCommon := map[string]string{
			"tags":       "SS",
			"categories": "SS",
		}

		for _, attr := range basesetstring.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be string set type", attr.Name)
		}

		t.Logf("✅ String set attributes validated")
	})

	t.Run("string_set_constants", func(t *testing.T) {
		assert.Equal(t, "base-set-string", basesetstring.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", basesetstring.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "group_id", basesetstring.ColumnGroupId, "ColumnGroupId should be correct")
		assert.Equal(t, "tags", basesetstring.ColumnTags, "ColumnTags should be correct")
		assert.Equal(t, "categories", basesetstring.ColumnCategories, "ColumnCategories should be correct")

		t.Logf("✅ String set constants validated")
	})

	t.Run("string_set_attribute_names", func(t *testing.T) {
		attrs := basesetstring.AttributeNames
		expectedAttrs := []string{"id", "group_id", "tags", "categories"}

		assert.Len(t, attrs, len(expectedAttrs), "Should have correct number of attributes")

		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}

		t.Logf("✅ String set AttributeNames validated")
	})

	t.Run("string_set_go_types", func(t *testing.T) {
		// Verify that generated struct has correct Go types for string sets
		item := basesetstring.SchemaItem{}

		// These should compile without type errors
		item.Tags = []string{"test1", "test2"}
		item.Categories = []string{"cat1", "cat2"}

		assert.IsType(t, []string{}, item.Tags, "Tags should be []string type")
		assert.IsType(t, []string{}, item.Categories, "Categories should be []string type")
		assert.IsType(t, "", item.Id, "Id should be string type")
		assert.IsType(t, "", item.GroupId, "GroupId should be string type")

		t.Logf("✅ String set Go types validated")
	})
}

// ==================== Helper Functions ====================

func setupStringSetTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basesetstring.SchemaItem{
		{
			Id:         "query-set-test",
			GroupId:    "frontend",
			Tags:       []string{"javascript", "react", "css"},
			Categories: []string{"frontend", "ui", "responsive"},
		},
		{
			Id:         "query-set-test",
			GroupId:    "backend",
			Tags:       []string{"nodejs", "express", "mongodb"},
			Categories: []string{"backend", "api", "database"},
		},
		{
			Id:         "query-set-test",
			GroupId:    "fullstack",
			Tags:       []string{"javascript", "react", "nodejs", "mongodb"},
			Categories: []string{"frontend", "backend", "fullstack"},
		},
	}

	for _, item := range testItems {
		av, err := basesetstring.ItemInput(item)
		require.NoError(t, err, "Should marshal string set test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basesetstring.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store string set test item")
	}

	t.Logf("Setup complete: inserted %d string set test items", len(testItems))
}
