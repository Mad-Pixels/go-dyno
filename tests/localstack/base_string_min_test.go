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

	basestringmin "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basestringmin"
)

// TestBaseStringMIN focuses on String (S) type operations with minimal generated code.
// This test validates string functionality using only basic methods (no sugar methods).
//
// Test Coverage:
// - Basic String CRUD operations using universal methods
// - String marshaling/unmarshaling
// - Core Query and Scan operations using .With() and .Filter()
// - Schema validation
//
// Schema: base-string__min.json
// - Table: "base-string-min"
// - Hash Key: id (S)
// - Range Key: category (S)
// - Common: title (S), description (S)
//
// Note: MIN mode only includes universal methods like .With() and .Filter()
// No convenience methods like .WithEQ(), .FilterEQ(), .FilterContains() etc.
func TestBaseStringMin(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing MIN mode String operations on: %s", basestringmin.TableName)

	t.Run("StringMIN_Input", func(t *testing.T) {
		testStringMINInput(t, client, ctx)
	})

	t.Run("StringMIN_Input_Raw", func(t *testing.T) {
		testStringMINInputRaw(t, client, ctx)
	})

	t.Run("StringMIN_QueryBuilder", func(t *testing.T) {
		testStringMINQueryBuilder(t, client, ctx)
	})

	t.Run("StringMIN_ScanBuilder", func(t *testing.T) {
		testStringMINScanBuilder(t, client, ctx)
	})

	t.Run("StringMIN_Schema", func(t *testing.T) {
		t.Parallel()
		testStringMINSchema(t)
	})
}

// ==================== String MIN Object Input ====================

func testStringMINInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_basic_crud", func(t *testing.T) {
		item := basestringmin.SchemaItem{
			Id:          "min-string-001",
			Category:    "min-docs",
			Title:       "MIN String Guide",
			Description: "Guide for MIN mode string operations",
		}
		av, err := basestringmin.ItemInput(item)
		require.NoError(t, err, "Should marshal MIN string item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.Contains(t, av, "id", "Should contain id field")
		assert.Contains(t, av, "category", "Should contain category field")
		assert.Contains(t, av, "title", "Should contain title field")
		assert.Contains(t, av, "description", "Should contain description field")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestringmin.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestringmin.ColumnCategory])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestringmin.ColumnTitle])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[basestringmin.ColumnDescription])

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basestringmin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store MIN string item in DynamoDB")
		key, err := basestringmin.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestringmin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve MIN string item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "min-string-001", getResult.Item[basestringmin.ColumnId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "min-docs", getResult.Item[basestringmin.ColumnCategory].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "MIN String Guide", getResult.Item[basestringmin.ColumnTitle].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Guide for MIN mode string operations", getResult.Item[basestringmin.ColumnDescription].(*types.AttributeValueMemberS).Value)

		item.Title = "Updated MIN String Guide"
		item.Description = "Updated guide for MIN mode string operations"

		updateInput, err := basestringmin.UpdateItemInput(item)
		require.NoError(t, err, "Should create update input from item")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update MIN string item")

		getResult, err = client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestringmin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, "Updated MIN String Guide", getResult.Item[basestringmin.ColumnTitle].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Updated guide for MIN mode string operations", getResult.Item[basestringmin.ColumnDescription].(*types.AttributeValueMemberS).Value)

		deleteInput, err := basestringmin.DeleteItemInput(item)
		require.NoError(t, err, "Should create delete input from item")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete MIN string item")

		getResult, err = client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basestringmin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "MIN string item should be deleted")
		t.Logf("✅ MIN mode basic CRUD operations work correctly")
	})

	t.Run("min_raw_operations", func(t *testing.T) {
		key, err := basestringmin.KeyInputFromRaw("min-raw-001", "min-category")
		require.NoError(t, err, "Should create key from raw values")
		assert.NotEmpty(t, key, "Raw key should not be empty")

		updates := map[string]any{
			"title":       "Raw MIN String Test",
			"description": "Testing raw operations in MIN mode",
		}
		updateInput, err := basestringmin.UpdateItemInputFromRaw("min-raw-001", "min-category", updates)
		require.NoError(t, err, "Should create update input from raw values")
		assert.NotNil(t, updateInput, "Update input should be created")
		t.Logf("✅ MIN mode raw string operations work correctly")
	})

	t.Run("min_string_edge_cases", func(t *testing.T) {
		edgeCases := []basestringmin.SchemaItem{
			{Id: "min-edge-1", Category: "empty", Title: "", Description: "Empty title test"},
			{Id: "min-edge-2", Category: "special", Title: "Special: !@#$%^&*()", Description: "Unicode: 🚀✨"},
			{Id: "min-edge-3", Category: "long", Title: "Very " + string(make([]byte, 50)), Description: "Long string test"},
		}

		for _, item := range edgeCases {
			av, err := basestringmin.ItemInput(item)
			require.NoError(t, err, "Should handle MIN string edge case: %s", item.Id)
			assert.NotEmpty(t, av, "Marshaled edge case should not be empty")
		}
		t.Logf("✅ MIN mode string edge cases handled successfully")
	})
}

// ==================== String MIN Raw Object Input ====================

func testStringMINInputRaw(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("raw_vs_object_min_comparison", func(t *testing.T) {
		keyFromRaw, err := basestringmin.KeyInputFromRaw("comparison-min-test", "both-methods")
		require.NoError(t, err, "Should create key from raw values")

		item := basestringmin.SchemaItem{
			Id:       "comparison-min-test",
			Category: "both-methods",
		}
		keyFromObject, err := basestringmin.KeyInput(item)
		require.NoError(t, err, "Should create key from object")

		assert.Equal(t, keyFromRaw, keyFromObject, "Raw and object-based keys should be identical")

		t.Logf("✅ Raw and object-based MIN methods produce identical results")
	})

	t.Run("raw_conditional_operations_min", func(t *testing.T) {
		conditionExpr := "#version = :v"
		conditionNames := map[string]string{"#version": "version"}
		conditionValues := map[string]types.AttributeValue{
			":v": &types.AttributeValueMemberN{Value: "1"},
		}

		deleteInput, err := basestringmin.DeleteItemInputWithCondition(
			"conditional-min-test", "min-condition",
			conditionExpr, conditionNames, conditionValues,
		)
		require.NoError(t, err, "Should create conditional delete with raw method")
		assert.NotNil(t, deleteInput.ConditionExpression, "Should have condition expression")
		assert.Equal(t, conditionExpr, *deleteInput.ConditionExpression, "Condition should match")
		t.Logf("✅ Raw conditional operations work in MIN mode")
	})
}

// ==================== String MIN QueryBuilder Tests ====================

func testStringMINQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupStringMINTestData(t, client, ctx)

	t.Run("min_query_universal_methods", func(t *testing.T) {
		qb := basestringmin.NewQueryBuilder().
			With("id", basestringmin.EQ, "min-query-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query using universal .With() method")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		t.Logf("✅ MIN mode universal .With() method works")
	})

	t.Run("min_query_range_conditions", func(t *testing.T) {
		qb := basestringmin.NewQueryBuilder().
			With("id", basestringmin.EQ, "min-query-test").
			With("category", basestringmin.BETWEEN, "api", "tutorial")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build range query using universal operators")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		t.Logf("✅ MIN mode range conditions work with universal operators")
	})

	t.Run("min_query_with_filters", func(t *testing.T) {
		qb := basestringmin.NewQueryBuilder().
			With("id", basestringmin.EQ, "min-query-test").
			Filter("title", basestringmin.EQ, "MIN API Documentation")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with universal .Filter() method")
		assert.NotNil(t, queryInput.FilterExpression, "Should have filter expression")
		t.Logf("✅ MIN mode universal .Filter() method works")
	})

	t.Run("min_query_execution", func(t *testing.T) {
		qb := basestringmin.NewQueryBuilder().
			With("id", basestringmin.EQ, "min-query-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN mode query")
		assert.NotEmpty(t, items, "Should return items")

		for _, item := range items {
			assert.Equal(t, "min-query-test", item.Id, "All items should have correct hash key")
			assert.NotEmpty(t, item.Category, "All items should have category")
			assert.IsType(t, "", item.Title, "Title should be string type")
			assert.IsType(t, "", item.Description, "Description should be string type")
		}
		t.Logf("✅ MIN mode query execution returned %d items", len(items))
	})
}

// ==================== String MIN ScanBuilder Tests ====================

func testStringMINScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_scan_universal_filter", func(t *testing.T) {
		sb := basestringmin.NewScanBuilder().
			Filter("id", basestringmin.EQ, "min-query-test")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with universal .Filter() method")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")
		t.Logf("✅ MIN mode scan universal .Filter() method works")
	})

	t.Run("min_scan_multiple_filters", func(t *testing.T) {
		sb := basestringmin.NewScanBuilder().
			Filter("id", basestringmin.EQ, "min-query-test").
			Filter("category", basestringmin.EQ, "api").
			Filter("title", basestringmin.GT, "A")

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with multiple universal filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")
		t.Logf("✅ MIN mode multiple universal filters work")
	})

	t.Run("min_scan_execution", func(t *testing.T) {
		sb := basestringmin.NewScanBuilder().
			Filter("id", basestringmin.EQ, "min-query-test").
			Limit(5)
		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN mode scan")

		for _, item := range items {
			assert.Equal(t, "min-query-test", item.Id, "Items should match filter")
		}
		t.Logf("✅ MIN mode scan execution returned %d items", len(items))
	})
}

// ==================== String MIN Schema Tests ====================

func testStringMINSchema(t *testing.T) {
	t.Run("min_schema_structure", func(t *testing.T) {
		schema := basestringmin.TableSchema
		assert.Equal(t, "base-string-min", schema.TableName, "Table name should match MIN schema")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "category", schema.RangeKey, "Range key should be 'category'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")
		t.Logf("✅ MIN mode schema structure validated")
	})

	t.Run("min_constants", func(t *testing.T) {
		assert.Equal(t, "base-string-min", basestringmin.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", basestringmin.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "category", basestringmin.ColumnCategory, "ColumnCategory should be correct")
		assert.Equal(t, "title", basestringmin.ColumnTitle, "ColumnTitle should be correct")
		assert.Equal(t, "description", basestringmin.ColumnDescription, "ColumnDescription should be correct")
		t.Logf("✅ MIN mode constants validated")
	})

	t.Run("min_operators_available", func(t *testing.T) {
		assert.NotNil(t, basestringmin.EQ, "EQ operator should be available")
		assert.NotNil(t, basestringmin.GT, "GT operator should be available")
		assert.NotNil(t, basestringmin.LT, "LT operator should be available")
		assert.NotNil(t, basestringmin.GTE, "GTE operator should be available")
		assert.NotNil(t, basestringmin.LTE, "LTE operator should be available")
		assert.NotNil(t, basestringmin.BETWEEN, "BETWEEN operator should be available")
		t.Logf("✅ MIN mode universal operators available")
	})

	t.Run("min_string_attributes", func(t *testing.T) {
		expectedPrimary := map[string]string{
			"id":       "S",
			"category": "S",
		}
		for _, attr := range basestringmin.TableSchema.Attributes {
			expectedType, exists := expectedPrimary[attr.Name]
			assert.True(t, exists, "Primary attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be string type", attr.Name)
		}
		expectedCommon := map[string]string{
			"title":       "S",
			"description": "S",
		}
		for _, attr := range basestringmin.TableSchema.CommonAttributes {
			expectedType, exists := expectedCommon[attr.Name]
			assert.True(t, exists, "Common attribute %s should be expected", attr.Name)
			assert.Equal(t, expectedType, attr.Type, "Attribute %s should be string type", attr.Name)
		}
		t.Logf("✅ MIN mode string attributes validated")
	})

	t.Run("min_no_sugar_methods", func(t *testing.T) {
		qb := basestringmin.NewQueryBuilder()
		assert.NotNil(t, qb, "QueryBuilder should be available")

		sb := basestringmin.NewScanBuilder()
		assert.NotNil(t, sb, "ScanBuilder should be available")
		t.Logf("✅ MIN mode builders available (sugar methods should be absent)")
	})
}

// ==================== Helper Functions ====================

func setupStringMINTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basestringmin.SchemaItem{
		{Id: "min-query-test", Category: "api", Title: "MIN API Documentation", Description: "REST API guide for MIN mode"},
		{Id: "min-query-test", Category: "sdk", Title: "MIN SDK Reference", Description: "Complete SDK documentation for MIN"},
		{Id: "min-query-test", Category: "tutorial", Title: "MIN Getting Started", Description: "Quick start tutorial for MIN mode"},
	}
	for _, item := range testItems {
		av, err := basestringmin.ItemInput(item)
		require.NoError(t, err, "Should marshal MIN test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basestringmin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store MIN test item")
	}
	t.Logf("MIN setup complete: inserted %d string test items", len(testItems))
}
