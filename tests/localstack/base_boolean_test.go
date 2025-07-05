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

	baseboolean "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basebooleanall"
)

// TestBaseBoolean focuses on Boolean (BOOL) type operations and functionality.
// This test validates boolean-specific features without other data types.
//
// Test Coverage:
// - Boolean CRUD operations
// - Boolean marshaling/unmarshaling
// - Boolean operations in Query and Scan
// - Boolean state transitions and filtering
// - Edge cases (true/false consistency)
//
// Schema: base-boolean__all.json
// - Table: "base-boolean-all"
// - Hash Key: id (S)
// - Range Key: version (N)
// - Common: is_active (BOOL), is_published (BOOL)
func TestBaseBoolean(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Boolean operations on: %s", baseboolean.TableName)

	t.Run("Boolean_Input", func(t *testing.T) {
		testBooleanInput(t, client, ctx)
	})

	t.Run("Boolean_Input_Raw", func(t *testing.T) {
		testBooleanInputRaw(t, client, ctx)
	})

	t.Run("Boolean_QueryBuilder", func(t *testing.T) {
		testBooleanQueryBuilder(t, client, ctx)
	})

	t.Run("Boolean_ScanBuilder", func(t *testing.T) {
		testBooleanScanBuilder(t, client, ctx)
	})

	t.Run("Boolean_StateTransitions", func(t *testing.T) {
		testBooleanStateTransitions(t, client, ctx)
	})

	t.Run("Boolean_Schema", func(t *testing.T) {
		t.Parallel()
		testBooleanSchema(t)
	})
}

// ==================== Boolean Object Input ====================

func testBooleanInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_boolean_item", func(t *testing.T) {
		item := baseboolean.SchemaItem{
			Id:          "boolean-test-001",
			Version:     1,
			IsActive:    true,
			IsPublished: false,
		}

		av, err := baseboolean.ItemInput(item)
		require.NoError(t, err, "Should marshal boolean item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "version", "Should contain version field")
		assert.Contains(t, av, "is_active", "Should contain is_active field")
		assert.Contains(t, av, "is_published", "Should contain is_published field")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[baseboolean.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[baseboolean.ColumnVersion])
		assert.IsType(t, &types.AttributeValueMemberBOOL{}, av[baseboolean.ColumnIsActive])
		assert.IsType(t, &types.AttributeValueMemberBOOL{}, av[baseboolean.ColumnIsPublished])

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(baseboolean.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store boolean item in DynamoDB")

		t.Logf("✅ Created boolean item: %s/%d (active=%t, published=%t)",
			item.Id, item.Version, item.IsActive, item.IsPublished)
	})

	t.Run("read_boolean_item", func(t *testing.T) {
		item := baseboolean.SchemaItem{
			Id:      "boolean-test-001",
			Version: 1,
		}

		key, err := baseboolean.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve boolean item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "boolean-test-001", getResult.Item[baseboolean.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "1", getResult.Item[baseboolean.ColumnVersion].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, true, getResult.Item[baseboolean.ColumnIsActive].(*types.AttributeValueMemberBOOL).Value)
		assert.Equal(t, false, getResult.Item[baseboolean.ColumnIsPublished].(*types.AttributeValueMemberBOOL).Value)

		t.Logf("✅ Retrieved boolean item successfully")
	})

	t.Run("update_boolean_item", func(t *testing.T) {
		item := baseboolean.SchemaItem{
			Id:          "boolean-test-001",
			Version:     1,
			IsActive:    false,
			IsPublished: true,
		}

		updateInput, err := baseboolean.UpdateItemInput(item)
		require.NoError(t, err, "Should create update input from item")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update boolean item")

		key, _ := baseboolean.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, false, getResult.Item[baseboolean.ColumnIsActive].(*types.AttributeValueMemberBOOL).Value)
		assert.Equal(t, true, getResult.Item[baseboolean.ColumnIsPublished].(*types.AttributeValueMemberBOOL).Value)

		t.Logf("✅ Updated boolean item successfully (flipped both boolean values)")
	})

	t.Run("delete_boolean_item", func(t *testing.T) {
		item := baseboolean.SchemaItem{
			Id:      "boolean-test-001",
			Version: 1,
		}

		deleteInput, err := baseboolean.DeleteItemInput(item)
		require.NoError(t, err, "Should create delete input from item")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete boolean item")

		key, _ := baseboolean.KeyInput(item)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "Boolean item should be deleted")

		t.Logf("✅ Deleted boolean item successfully")
	})

	t.Run("boolean_all_combinations", func(t *testing.T) {
		combinations := []baseboolean.SchemaItem{
			{Id: "combo-1", Version: 1, IsActive: true, IsPublished: true},
			{Id: "combo-2", Version: 2, IsActive: true, IsPublished: false},
			{Id: "combo-3", Version: 3, IsActive: false, IsPublished: true},
			{Id: "combo-4", Version: 4, IsActive: false, IsPublished: false},
		}

		for _, item := range combinations {
			av, err := baseboolean.ItemInput(item)
			require.NoError(t, err, "Should handle boolean combination: %s", item.Id)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(baseboolean.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store boolean combination item: %s", item.Id)
		}

		t.Logf("✅ Boolean combinations handled successfully")
	})
}

// ==================== Boolean Raw Object Input ====================

func testBooleanInputRaw(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_boolean_item_raw", func(t *testing.T) {
		item := baseboolean.SchemaItem{
			Id:          "boolean-raw-001",
			Version:     10,
			IsActive:    true,
			IsPublished: false,
		}

		av, err := baseboolean.ItemInput(item)
		require.NoError(t, err, "Should marshal boolean item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(baseboolean.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store boolean item in DynamoDB")

		t.Logf("✅ Created boolean item for raw testing: %s/%d", item.Id, item.Version)
	})

	t.Run("read_boolean_item_raw", func(t *testing.T) {
		key, err := baseboolean.KeyInputFromRaw("boolean-raw-001", 10)
		require.NoError(t, err, "Should create key from raw values")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve boolean item using raw key")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "boolean-raw-001", getResult.Item[baseboolean.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "10", getResult.Item[baseboolean.ColumnVersion].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, true, getResult.Item[baseboolean.ColumnIsActive].(*types.AttributeValueMemberBOOL).Value)
		assert.Equal(t, false, getResult.Item[baseboolean.ColumnIsPublished].(*types.AttributeValueMemberBOOL).Value)

		t.Logf("✅ Retrieved boolean item successfully using raw key")
	})

	t.Run("update_boolean_item_raw", func(t *testing.T) {
		updates := map[string]any{
			"is_active":    false,
			"is_published": true,
		}

		updateInput, err := baseboolean.UpdateItemInputFromRaw("boolean-raw-001", 10, updates)
		require.NoError(t, err, "Should create update input from raw values")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update boolean item using raw method")

		key, _ := baseboolean.KeyInputFromRaw("boolean-raw-001", 10)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, false, getResult.Item[baseboolean.ColumnIsActive].(*types.AttributeValueMemberBOOL).Value)
		assert.Equal(t, true, getResult.Item[baseboolean.ColumnIsPublished].(*types.AttributeValueMemberBOOL).Value)

		t.Logf("✅ Updated boolean item successfully using raw method")
	})

	t.Run("delete_boolean_item_raw", func(t *testing.T) {
		deleteInput, err := baseboolean.DeleteItemInputFromRaw("boolean-raw-001", 10)
		require.NoError(t, err, "Should create delete input from raw values")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete boolean item using raw method")

		key, _ := baseboolean.KeyInputFromRaw("boolean-raw-001", 10)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "Boolean item should be deleted")

		t.Logf("✅ Deleted boolean item successfully using raw method")
	})

	t.Run("raw_vs_object_boolean_comparison", func(t *testing.T) {
		keyFromRaw, err := baseboolean.KeyInputFromRaw("comparison-test", 99)
		require.NoError(t, err, "Should create key from raw values")

		item := baseboolean.SchemaItem{
			Id:      "comparison-test",
			Version: 99,
		}
		keyFromObject, err := baseboolean.KeyInput(item)
		require.NoError(t, err, "Should create key from object")

		assert.Equal(t, keyFromRaw, keyFromObject, "Raw and object-based keys should be identical")

		t.Logf("✅ Raw and object-based boolean methods produce identical results")
	})

	t.Run("raw_boolean_edge_cases", func(t *testing.T) {
		edgeCases := []struct {
			id      string
			version int
			updates map[string]any
		}{
			{
				id:      "raw-edge-1",
				version: 1,
				updates: map[string]any{"is_active": true, "is_published": true},
			},
			{
				id:      "raw-edge-2",
				version: 2,
				updates: map[string]any{"is_active": false, "is_published": false},
			},
			{
				id:      "raw-edge-3",
				version: 3,
				updates: map[string]any{"is_active": true, "is_published": false},
			},
		}

		for _, edgeCase := range edgeCases {
			updateInput, err := baseboolean.UpdateItemInputFromRaw(edgeCase.id, edgeCase.version, edgeCase.updates)
			require.NoError(t, err, "Should handle raw boolean edge case: %s", edgeCase.id)
			assert.NotNil(t, updateInput, "Update input should be created for edge case: %s", edgeCase.id)

			deleteInput, err := baseboolean.DeleteItemInputFromRaw(edgeCase.id, edgeCase.version)
			require.NoError(t, err, "Should create delete input for edge case: %s", edgeCase.id)
			assert.NotNil(t, deleteInput, "Delete input should be created")
		}

		t.Logf("✅ Raw boolean edge cases handled successfully")
	})

	t.Run("raw_conditional_operations", func(t *testing.T) {
		conditionExpr := "#is_active = :active"
		conditionNames := map[string]string{"#is_active": "is_active"}
		conditionValues := map[string]types.AttributeValue{
			":active": &types.AttributeValueMemberBOOL{Value: true},
		}

		deleteInput, err := baseboolean.DeleteItemInputWithCondition(
			"conditional-test", 100,
			conditionExpr, conditionNames, conditionValues,
		)
		require.NoError(t, err, "Should create conditional delete with raw method")
		assert.NotNil(t, deleteInput.ConditionExpression, "Should have condition expression")
		assert.Equal(t, conditionExpr, *deleteInput.ConditionExpression, "Condition should match")

		updates := map[string]any{
			"is_active":    false,
			"is_published": true,
		}

		updateInput, err := baseboolean.UpdateItemInputWithCondition(
			"conditional-test", 100,
			updates, conditionExpr, conditionNames, conditionValues,
		)
		require.NoError(t, err, "Should create conditional update with raw method")
		assert.NotNil(t, updateInput.ConditionExpression, "Should have condition expression")

		t.Logf("✅ Raw conditional operations work correctly")
	})
}

// ==================== Boolean QueryBuilder Tests ====================

func testBooleanQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	// Setup test data
	setupBooleanTestData(t, client, ctx)

	t.Run("boolean_hash_key_query", func(t *testing.T) {
		qb := baseboolean.NewQueryBuilder().WithEQ("id", "query-boolean-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build boolean hash key query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		assert.Equal(t, baseboolean.TableName, *queryInput.TableName, "Should target correct table")

		t.Logf("✅ Boolean hash key query built successfully")
	})

	t.Run("boolean_hash_and_range_query", func(t *testing.T) {
		qb := baseboolean.NewQueryBuilder().
			WithEQ("id", "query-boolean-test").
			WithEQ("version", 1)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build boolean hash+range query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ Boolean hash+range query built successfully")
	})

	t.Run("boolean_filters", func(t *testing.T) {
		qb := baseboolean.NewQueryBuilder().
			WithEQ("id", "query-boolean-test").
			FilterEQ("is_active", true).
			FilterEQ("is_published", false)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with boolean filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ Boolean filters query built successfully")
	})

	t.Run("boolean_query_execution", func(t *testing.T) {
		qb := baseboolean.NewQueryBuilder().WithEQ("id", "query-boolean-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute boolean query")
		assert.NotEmpty(t, items, "Should return boolean items")

		for _, item := range items {
			assert.Equal(t, "query-boolean-test", item.Id, "All items should have correct hash key")
			assert.Greater(t, item.Version, 0, "All items should have positive version")
			assert.IsType(t, true, item.IsActive, "IsActive should be bool type")
			assert.IsType(t, false, item.IsPublished, "IsPublished should be bool type")
		}

		t.Logf("✅ Boolean query execution returned %d items", len(items))
	})

	t.Run("boolean_version_range_conditions", func(t *testing.T) {
		qb := baseboolean.NewQueryBuilder().
			WithEQ("id", "query-boolean-test").
			WithBetween("version", 1, 3)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build version between query")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute version between query")

		for _, item := range items {
			assert.GreaterOrEqual(t, item.Version, 1, "Version should be >= 1")
			assert.LessOrEqual(t, item.Version, 3, "Version should be <= 3")
		}

		t.Logf("✅ Boolean version range condition returned %d items", len(items))
	})

	t.Run("boolean_sorting", func(t *testing.T) {
		qbAsc := baseboolean.NewQueryBuilder().
			WithEQ("id", "query-boolean-test").
			OrderByAsc()

		itemsAsc, err := qbAsc.Execute(ctx, client)
		require.NoError(t, err, "Should execute ascending boolean query")

		qbDesc := baseboolean.NewQueryBuilder().
			WithEQ("id", "query-boolean-test").
			OrderByDesc()

		itemsDesc, err := qbDesc.Execute(ctx, client)
		require.NoError(t, err, "Should execute descending boolean query")

		if len(itemsAsc) > 1 && len(itemsDesc) > 1 {
			assert.NotEqual(t, itemsAsc[0].Version, itemsDesc[0].Version, "Boolean sorting should produce different order")
			if len(itemsAsc) > 1 {
				assert.LessOrEqual(t, itemsAsc[0].Version, itemsAsc[1].Version, "Ascending should be in increasing order")
			}
		}

		t.Logf("✅ Boolean sorting works correctly")
	})
}

// ==================== Boolean ScanBuilder Tests ====================

func testBooleanScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("boolean_scan_filters", func(t *testing.T) {
		sb := baseboolean.NewScanBuilder().
			FilterEQ("id", "query-boolean-test").
			FilterEQ("is_active", true)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with boolean filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ Boolean scan filters built successfully")
	})

	t.Run("boolean_advanced_filters", func(t *testing.T) {
		sb := baseboolean.NewScanBuilder().
			FilterEQ("is_active", true).
			FilterEQ("is_published", false).
			FilterBetween("version", 1, 5)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with advanced boolean filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ Advanced boolean filters built successfully")
	})

	t.Run("boolean_scan_execution", func(t *testing.T) {
		sb := baseboolean.NewScanBuilder().
			FilterEQ("is_active", true).
			Limit(10)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute boolean scan")

		for _, item := range items {
			assert.Equal(t, true, item.IsActive, "Items should match boolean filter")
		}

		t.Logf("✅ Boolean scan execution returned %d items", len(items))
	})

	t.Run("boolean_not_equal_filter", func(t *testing.T) {
		sb := baseboolean.NewScanBuilder().FilterNE("is_published", true)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with boolean not equal filter")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ Boolean not equal filter built successfully")
	})
}

// ==================== Boolean State Transitions Tests ====================

func testBooleanStateTransitions(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("activation_workflow", func(t *testing.T) {
		// Create inactive, unpublished item
		item := baseboolean.SchemaItem{
			Id:          "workflow-test",
			Version:     1,
			IsActive:    false,
			IsPublished: false,
		}

		av, err := baseboolean.ItemInput(item)
		require.NoError(t, err, "Should create workflow item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(baseboolean.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store workflow item")

		// Step 1: Activate
		updates := map[string]any{"is_active": true}
		updateInput, err := baseboolean.UpdateItemInputFromRaw("workflow-test", 1, updates)
		require.NoError(t, err, "Should create activation update")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should activate item")

		// Step 2: Publish (only if active)
		updates = map[string]any{"is_published": true}
		conditionExpr := "#is_active = :true"
		conditionNames := map[string]string{"#is_active": "is_active"}
		conditionValues := map[string]types.AttributeValue{
			":true": &types.AttributeValueMemberBOOL{Value: true},
		}

		updateInput, err = baseboolean.UpdateItemInputWithCondition(
			"workflow-test", 1, updates, conditionExpr, conditionNames, conditionValues,
		)
		require.NoError(t, err, "Should create conditional publish update")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should publish item conditionally")

		// Verify final state
		key, _ := baseboolean.KeyInputFromRaw("workflow-test", 1)
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(baseboolean.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve final state")

		assert.Equal(t, true, getResult.Item[baseboolean.ColumnIsActive].(*types.AttributeValueMemberBOOL).Value)
		assert.Equal(t, true, getResult.Item[baseboolean.ColumnIsPublished].(*types.AttributeValueMemberBOOL).Value)

		t.Logf("✅ Activation workflow completed successfully")
	})

	t.Run("boolean_state_combinations", func(t *testing.T) {
		testCases := []struct {
			name        string
			isActive    bool
			isPublished bool
			description string
		}{
			{"draft", false, false, "Inactive and unpublished"},
			{"ready", true, false, "Active but not published"},
			{"live", true, true, "Active and published"},
			{"disabled", false, true, "Inactive but published (edge case)"},
		}

		for i, tc := range testCases {
			item := baseboolean.SchemaItem{
				Id:          "state-test",
				Version:     i + 1,
				IsActive:    tc.isActive,
				IsPublished: tc.isPublished,
			}

			av, err := baseboolean.ItemInput(item)
			require.NoError(t, err, "Should create %s state item", tc.name)

			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(baseboolean.TableName),
				Item:      av,
			})
			require.NoError(t, err, "Should store %s state item", tc.name)

			t.Logf("✅ Created %s state: active=%t, published=%t", tc.name, tc.isActive, tc.isPublished)
		}
	})
}

// ==================== Boolean Schema Tests ====================

func testBooleanSchema(t *testing.T) {
	t.Run("boolean_table_schema", func(t *testing.T) {
		schema := baseboolean.TableSchema

		assert.Equal(t, "base-boolean-all", schema.TableName, "Table name should match")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "version", schema.RangeKey, "Range key should be 'version'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")

		t.Logf("✅ Boolean schema structure validated")
	})

	t.Run("boolean_attributes", func(t *testing.T) {
		// Check primary attributes
		expectedPrimary := map[string]string{
			"id":      "S", // hash key is string
			"version": "N", // range key is number
		}

		for _, attr := range baseboolean.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should have correct type", attr.Name)
		}

		// Check common attributes (all boolean type)
		expectedCommon := map[string]string{
			"is_active":    "BOOL",
			"is_published": "BOOL",
		}

		for _, attr := range baseboolean.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be boolean type", attr.Name)
		}

		t.Logf("✅ Boolean attributes validated")
	})

	t.Run("boolean_constants", func(t *testing.T) {
		assert.Equal(t, "base-boolean-all", baseboolean.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", baseboolean.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "version", baseboolean.ColumnVersion, "ColumnVersion should be correct")
		assert.Equal(t, "is_active", baseboolean.ColumnIsActive, "ColumnIsActive should be correct")
		assert.Equal(t, "is_published", baseboolean.ColumnIsPublished, "ColumnIsPublished should be correct")

		t.Logf("✅ Boolean constants validated")
	})

	t.Run("boolean_attribute_names", func(t *testing.T) {
		attrs := baseboolean.AttributeNames
		expectedAttrs := []string{"id", "version", "is_active", "is_published"}

		assert.Len(t, attrs, len(expectedAttrs), "Should have correct number of attributes")

		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}

		t.Logf("✅ Boolean AttributeNames validated")
	})
}

// ==================== Helper Functions ====================

func setupBooleanTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []baseboolean.SchemaItem{
		{Id: "query-boolean-test", Version: 1, IsActive: true, IsPublished: false},
		{Id: "query-boolean-test", Version: 2, IsActive: false, IsPublished: true},
		{Id: "query-boolean-test", Version: 3, IsActive: true, IsPublished: true},
		{Id: "query-boolean-test", Version: 4, IsActive: false, IsPublished: false},
	}

	for _, item := range testItems {
		av, err := baseboolean.ItemInput(item)
		require.NoError(t, err, "Should marshal boolean test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(baseboolean.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store boolean test item")
	}

	t.Logf("Setup complete: inserted %d boolean test items", len(testItems))
}
