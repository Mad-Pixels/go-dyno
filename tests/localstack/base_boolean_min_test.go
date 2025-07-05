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

	basebooleanmin "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/basebooleanmin"
)

// TestBaseBooleanMIN focuses on Boolean (BOOL) type operations with minimal generated code.
// This test validates boolean functionality using only basic methods (no sugar methods).
//
// Test Coverage:
// - Basic Boolean CRUD operations using universal methods
// - Boolean marshaling/unmarshaling
// - Core Query and Scan operations using .With() and .Filter()
// - Schema validation
//
// Schema: base-boolean__min.json
// - Table: "base-boolean-min"
// - Hash Key: id (S)
// - Range Key: version (N)
// - Common: is_active (BOOL), is_published (BOOL)
//
// Note: MIN mode only includes universal methods like .With() and .Filter()
// No convenience methods like .WithEQ(), .FilterEQ() etc.
func TestBaseBooleanMIN(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing Boolean MIN operations on: %s", basebooleanmin.TableName)

	t.Run("Boolean_MIN_BasicCRUD", func(t *testing.T) {
		testBooleanMINBasicCRUD(t, client, ctx)
	})

	t.Run("Boolean_MIN_QueryBuilder", func(t *testing.T) {
		testBooleanMINQueryBuilder(t, client, ctx)
	})

	t.Run("Boolean_MIN_ScanBuilder", func(t *testing.T) {
		testBooleanMINScanBuilder(t, client, ctx)
	})

	t.Run("Boolean_MIN_Schema", func(t *testing.T) {
		t.Parallel()
		testBooleanMINSchema(t)
	})
}

// ==================== Boolean MIN Basic CRUD ====================

func testBooleanMINBasicCRUD(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_create_and_read", func(t *testing.T) {
		item := basebooleanmin.SchemaItem{
			Id:          "min-test-001",
			Version:     1,
			IsActive:    true,
			IsPublished: false,
		}

		av, err := basebooleanmin.ItemInput(item)
		require.NoError(t, err, "Should marshal boolean item in MIN mode")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[basebooleanmin.ColumnId])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[basebooleanmin.ColumnVersion])
		assert.IsType(t, &types.AttributeValueMemberBOOL{}, av[basebooleanmin.ColumnIsActive])
		assert.IsType(t, &types.AttributeValueMemberBOOL{}, av[basebooleanmin.ColumnIsPublished])

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basebooleanmin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store boolean item in DynamoDB")

		key, err := basebooleanmin.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(basebooleanmin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve boolean item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, true, getResult.Item[basebooleanmin.ColumnIsActive].(*types.AttributeValueMemberBOOL).Value)
		assert.Equal(t, false, getResult.Item[basebooleanmin.ColumnIsPublished].(*types.AttributeValueMemberBOOL).Value)

		t.Logf("✅ MIN mode basic CRUD operations work correctly")
	})

	t.Run("min_raw_operations", func(t *testing.T) {
		key, err := basebooleanmin.KeyInputFromRaw("min-raw-001", 5)
		require.NoError(t, err, "Should create key from raw values")
		assert.NotEmpty(t, key, "Raw key should not be empty")

		updates := map[string]any{
			"is_active":    false,
			"is_published": true,
		}

		updateInput, err := basebooleanmin.UpdateItemInputFromRaw("min-raw-001", 5, updates)
		require.NoError(t, err, "Should create update input from raw values")
		assert.NotNil(t, updateInput, "Update input should be created")

		t.Logf("✅ MIN mode raw operations work correctly")
	})
}

// ==================== Boolean MIN QueryBuilder Tests ====================

func testBooleanMINQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupBooleanMINTestData(t, client, ctx)

	t.Run("min_query_universal_methods", func(t *testing.T) {
		qb := basebooleanmin.NewQueryBuilder().
			With("id", basebooleanmin.EQ, "min-query-test")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query using universal .With() method")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ MIN mode universal .With() method works")
	})

	t.Run("min_query_range_conditions", func(t *testing.T) {
		qb := basebooleanmin.NewQueryBuilder().
			With("id", basebooleanmin.EQ, "min-query-test").
			With("version", basebooleanmin.BETWEEN, 1, 3)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build range query using universal operators")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ MIN mode range conditions work with universal operators")
	})

	t.Run("min_query_with_filters", func(t *testing.T) {
		qb := basebooleanmin.NewQueryBuilder().
			With("id", basebooleanmin.EQ, "min-query-test").
			Filter("is_active", basebooleanmin.EQ, true)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with universal .Filter() method")
		assert.NotNil(t, queryInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ MIN mode universal .Filter() method works")
	})

	t.Run("min_query_execution", func(t *testing.T) {
		qb := basebooleanmin.NewQueryBuilder().
			With("id", basebooleanmin.EQ, "min-query-test")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN mode query")
		assert.NotEmpty(t, items, "Should return items")

		for _, item := range items {
			assert.Equal(t, "min-query-test", item.Id, "All items should have correct hash key")
		}

		t.Logf("✅ MIN mode query execution returned %d items", len(items))
	})
}

// ==================== Boolean MIN ScanBuilder Tests ====================

func testBooleanMINScanBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_scan_universal_filter", func(t *testing.T) {
		sb := basebooleanmin.NewScanBuilder().
			Filter("is_active", basebooleanmin.EQ, true)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with universal .Filter() method")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ MIN mode scan universal .Filter() method works")
	})

	t.Run("min_scan_multiple_filters", func(t *testing.T) {
		sb := basebooleanmin.NewScanBuilder().
			Filter("is_active", basebooleanmin.EQ, true).
			Filter("is_published", basebooleanmin.EQ, false).
			Filter("version", basebooleanmin.GT, 0)

		scanInput, err := sb.BuildScan()
		require.NoError(t, err, "Should build scan with multiple universal filters")
		assert.NotNil(t, scanInput.FilterExpression, "Should have filter expression")

		t.Logf("✅ MIN mode multiple universal filters work")
	})

	t.Run("min_scan_execution", func(t *testing.T) {
		sb := basebooleanmin.NewScanBuilder().
			Filter("id", basebooleanmin.EQ, "min-query-test").
			Limit(5)

		items, err := sb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN mode scan")

		for _, item := range items {
			assert.Equal(t, "min-query-test", item.Id, "Items should match filter")
		}

		t.Logf("✅ MIN mode scan execution returned %d items", len(items))
	})
}

// ==================== Boolean MIN Schema Tests ====================

func testBooleanMINSchema(t *testing.T) {
	t.Run("min_schema_structure", func(t *testing.T) {
		schema := basebooleanmin.TableSchema

		assert.Equal(t, "base-boolean-min", schema.TableName, "Table name should match MIN schema")
		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
		assert.Equal(t, "version", schema.RangeKey, "Range key should be 'version'")
		assert.Len(t, schema.SecondaryIndexes, 0, "Should have no secondary indexes")

		t.Logf("✅ MIN mode schema structure validated")
	})

	t.Run("min_constants", func(t *testing.T) {
		assert.Equal(t, "base-boolean-min", basebooleanmin.TableName, "TableName constant should be correct")
		assert.Equal(t, "id", basebooleanmin.ColumnId, "ColumnId should be correct")
		assert.Equal(t, "version", basebooleanmin.ColumnVersion, "ColumnVersion should be correct")
		assert.Equal(t, "is_active", basebooleanmin.ColumnIsActive, "ColumnIsActive should be correct")
		assert.Equal(t, "is_published", basebooleanmin.ColumnIsPublished, "ColumnIsPublished should be correct")

		t.Logf("✅ MIN mode constants validated")
	})

	t.Run("min_operators_available", func(t *testing.T) {
		assert.NotNil(t, basebooleanmin.EQ, "EQ operator should be available")
		assert.NotNil(t, basebooleanmin.GT, "GT operator should be available")
		assert.NotNil(t, basebooleanmin.LT, "LT operator should be available")
		assert.NotNil(t, basebooleanmin.BETWEEN, "BETWEEN operator should be available")

		t.Logf("✅ MIN mode universal operators available")
	})

	t.Run("min_no_sugar_methods", func(t *testing.T) {
		qb := basebooleanmin.NewQueryBuilder()
		assert.NotNil(t, qb, "QueryBuilder should be available")

		sb := basebooleanmin.NewScanBuilder()
		assert.NotNil(t, sb, "ScanBuilder should be available")

		t.Logf("✅ MIN mode builders available (sugar methods should be absent)")
	})
}

// ==================== Helper Functions ====================

func setupBooleanMINTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []basebooleanmin.SchemaItem{
		{Id: "min-query-test", Version: 1, IsActive: true, IsPublished: false},
		{Id: "min-query-test", Version: 2, IsActive: false, IsPublished: true},
		{Id: "min-query-test", Version: 3, IsActive: true, IsPublished: true},
	}

	for _, item := range testItems {
		av, err := basebooleanmin.ItemInput(item)
		require.NoError(t, err, "Should marshal MIN test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(basebooleanmin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store MIN test item")
	}

	t.Logf("MIN setup complete: inserted %d boolean test items", len(testItems))
}
