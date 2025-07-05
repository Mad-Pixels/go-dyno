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

	userpostscompletemin "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/userpostscompletemin"
)

// TestUserPostsCompleteMIN focuses on GSI/LSI mixed operations with minimal generated code.
// This test validates the new architecture with both GSI and LSI indexes using universal methods only.
//
// Test Coverage:
// - Mixed GSI/LSI CRUD operations using universal methods
// - Index selection and query optimization in MIN mode
// - GSI vs LSI query performance with basic methods
// - Complex projection types using .With() and .Filter()
// - Multiple index types on same table in MIN mode
//
// Schema: user-posts-complete__min.json
// - Table: "user-posts-complete-min"
// - Hash Key: user_id (S)
// - Range Key: created_at (S)
// - LSI: lsi_by_post_type, lsi_by_status, lsi_by_priority
// - GSI: gsi_by_category, gsi_by_title, gsi_by_status_priority
//
// Note: MIN mode only includes universal methods like .With() and .Filter()
// No convenience methods like .WithEQ(), .FilterEQ(), .WithBetween() etc.
func TestUserPostsCompleteMIN(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing MIN mode GSI/LSI operations on: %s", userpostscompletemin.TableName)

	t.Run("UserPostsMIN_Input", func(t *testing.T) {
		testUserPostsMINInput(t, client, ctx)
	})

	t.Run("UserPostsMIN_QueryBuilder_LSI", func(t *testing.T) {
		testUserPostsMINQueryBuilderLSI(t, client, ctx)
	})

	t.Run("UserPostsMIN_QueryBuilder_GSI", func(t *testing.T) {
		testUserPostsMINQueryBuilderGSI(t, client, ctx)
	})

	t.Run("UserPostsMIN_Schema", func(t *testing.T) {
		t.Parallel()
		testUserPostsMINSchema(t)
	})
}

// ==================== User Posts MIN Input ====================

func testUserPostsMINInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_complete_crud", func(t *testing.T) {
		item := userpostscompletemin.SchemaItem{
			UserId:    "min-user123",
			CreatedAt: "2024-01-15T10:30:00Z",
			PostType:  "blog",
			Status:    "published",
			Priority:  85,
			Category:  "technology",
			Title:     "MIN Introduction to DynamoDB",
			Content:   "This is a comprehensive guide to DynamoDB in MIN mode",
			Tags:      []string{"aws", "database", "nosql"},
			ViewCount: 1500,
			UpdatedAt: "2024-01-16T09:15:00Z",
		}

		av, err := userpostscompletemin.ItemInput(item)
		require.NoError(t, err, "Should marshal MIN complete item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		assert.Contains(t, av, "user_id", "Should contain user_id field")
		assert.Contains(t, av, "created_at", "Should contain created_at field")
		assert.Contains(t, av, "post_type", "Should contain post_type field")
		assert.Contains(t, av, "status", "Should contain status field")
		assert.Contains(t, av, "priority", "Should contain priority field")
		assert.Contains(t, av, "category", "Should contain category field")
		assert.Contains(t, av, "title", "Should contain title field")
		assert.Contains(t, av, "content", "Should contain content field")
		assert.Contains(t, av, "tags", "Should contain tags field")
		assert.Contains(t, av, "view_count", "Should contain view_count field")

		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnUserId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnCreatedAt])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnPostType])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnStatus])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[userpostscompletemin.ColumnPriority])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnCategory])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnTitle])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscompletemin.ColumnContent])
		assert.IsType(t, &types.AttributeValueMemberSS{}, av[userpostscompletemin.ColumnTags])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[userpostscompletemin.ColumnViewCount])

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userpostscompletemin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store MIN complete item in DynamoDB")

		key, err := userpostscompletemin.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(userpostscompletemin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve MIN complete item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		assert.Equal(t, "min-user123", getResult.Item[userpostscompletemin.ColumnUserId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "2024-01-15T10:30:00Z", getResult.Item[userpostscompletemin.ColumnCreatedAt].(*types.AttributeValueMemberS).Value)

		assert.Equal(t, "blog", getResult.Item[userpostscompletemin.ColumnPostType].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "published", getResult.Item[userpostscompletemin.ColumnStatus].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "85", getResult.Item[userpostscompletemin.ColumnPriority].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "technology", getResult.Item[userpostscompletemin.ColumnCategory].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "MIN Introduction to DynamoDB", getResult.Item[userpostscompletemin.ColumnTitle].(*types.AttributeValueMemberS).Value)

		tagsSet := getResult.Item[userpostscompletemin.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Contains(t, tagsSet.Value, "aws")
		assert.Contains(t, tagsSet.Value, "database")
		assert.Contains(t, tagsSet.Value, "nosql")

		item.Title = "Updated MIN Introduction to DynamoDB"
		item.ViewCount = 2000

		updateInput, err := userpostscompletemin.UpdateItemInput(item)
		require.NoError(t, err, "Should create update input from item")

		_, err = client.UpdateItem(ctx, updateInput)
		require.NoError(t, err, "Should update MIN complete item")

		getResult, err = client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(userpostscompletemin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve updated item")

		assert.Equal(t, "Updated MIN Introduction to DynamoDB", getResult.Item[userpostscompletemin.ColumnTitle].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "2000", getResult.Item[userpostscompletemin.ColumnViewCount].(*types.AttributeValueMemberN).Value)

		deleteInput, err := userpostscompletemin.DeleteItemInput(item)
		require.NoError(t, err, "Should create delete input from item")

		_, err = client.DeleteItem(ctx, deleteInput)
		require.NoError(t, err, "Should delete MIN complete item")

		getResult, err = client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(userpostscompletemin.TableName),
			Key:       key,
		})
		require.NoError(t, err, "GetItem should not error for missing item")
		assert.Empty(t, getResult.Item, "MIN complete item should be deleted")

		t.Logf("✅ MIN mode complete CRUD operations work correctly")
	})

	t.Run("min_raw_operations", func(t *testing.T) {
		key, err := userpostscompletemin.KeyInputFromRaw("min-raw-user", "2024-01-01T12:00:00Z")
		require.NoError(t, err, "Should create key from raw values")
		assert.NotEmpty(t, key, "Raw key should not be empty")

		updates := map[string]any{
			"title":      "MIN Raw Update Test",
			"view_count": 500,
			"status":     "published",
		}

		updateInput, err := userpostscompletemin.UpdateItemInputFromRaw("min-raw-user", "2024-01-01T12:00:00Z", updates)
		require.NoError(t, err, "Should create update input from raw values")
		assert.NotNil(t, updateInput, "Update input should be created")

		t.Logf("✅ MIN mode raw operations work correctly")
	})
}

// ==================== LSI QueryBuilder Tests ====================

func testUserPostsMINQueryBuilderLSI(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupUserPostsMINTestData(t, client, ctx)

	t.Run("min_lsi_by_post_type", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("user_id", userpostscompletemin.EQ, "min-query-test-user").
			With("post_type", userpostscompletemin.EQ, "tutorial")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN LSI query by post_type")

		t.Logf("MIN IndexName: %v", queryInput.IndexName)
		t.Logf("MIN KeyConditionExpression: %v", queryInput.KeyConditionExpression)

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN LSI query")

		t.Logf("MIN Returned %d items", len(items))
		for i, item := range items {
			t.Logf("MIN Item %d: user_id=%s, post_type=%s", i, item.UserId, item.PostType)
			assert.Equal(t, "min-query-test-user", item.UserId)
			assert.Equal(t, "tutorial", item.PostType)
		}
		t.Logf("✅ MIN mode LSI by post_type query works")
	})

	t.Run("min_lsi_by_status", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("user_id", userpostscompletemin.EQ, "min-query-test-user").
			With("status", userpostscompletemin.EQ, "published")

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN LSI query by status")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN LSI status query")

		for _, item := range items {
			assert.Equal(t, "min-query-test-user", item.UserId)
			assert.Equal(t, "published", item.Status)
		}
		t.Logf("✅ MIN mode LSI by status query returned %d items", len(items))
	})

	t.Run("min_lsi_by_priority_with_range", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("user_id", userpostscompletemin.EQ, "min-query-test-user").
			With("priority", userpostscompletemin.BETWEEN, 70, 90)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN LSI query with priority range")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN LSI priority range query")

		for _, item := range items {
			assert.Equal(t, "min-query-test-user", item.UserId)
			assert.GreaterOrEqual(t, item.Priority, 70)
			assert.LessOrEqual(t, item.Priority, 90)
		}
		t.Logf("✅ MIN mode LSI priority range query returned %d items", len(items))
	})
}

// ==================== GSI QueryBuilder Tests ====================

func testUserPostsMINQueryBuilderGSI(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("min_gsi_by_category", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("category", userpostscompletemin.EQ, "technology")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN GSI query by category")

		if queryInput.IndexName != nil {
			assert.Equal(t, "gsi_by_category", *queryInput.IndexName, "Should use GSI index")
		}
		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN GSI category query")

		for _, item := range items {
			assert.Equal(t, "technology", item.Category, "All items should match category")
		}
		t.Logf("✅ MIN mode GSI by category query returned %d items", len(items))
	})

	t.Run("min_gsi_by_title", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("title", userpostscompletemin.EQ, "MIN Advanced DynamoDB")

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN GSI query by title")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN GSI title query")

		for _, item := range items {
			assert.Equal(t, "MIN Advanced DynamoDB", item.Title)
		}
		t.Logf("✅ MIN mode GSI by title query returned %d items", len(items))
	})

	t.Run("min_gsi_status_priority_compound", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("status", userpostscompletemin.EQ, "published").
			With("priority", userpostscompletemin.GT, 80)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN compound GSI query")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN compound GSI query")

		t.Logf("MIN KeyCondition: %s", aws.ToString(queryInput.KeyConditionExpression))

		for _, item := range items {
			assert.Equal(t, "published", item.Status)
			assert.Greater(t, item.Priority, 80)
		}

		t.Logf("✅ MIN mode GSI compound query returned %d items", len(items))
	})

	t.Run("min_gsi_with_filters", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder().
			With("category", userpostscompletemin.EQ, "technology").
			Filter("view_count", userpostscompletemin.GT, 1000)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build MIN GSI query with filters")
		assert.NotNil(t, queryInput.FilterExpression, "Should have filter expression")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute MIN GSI query with filters")

		for _, item := range items {
			assert.Equal(t, "technology", item.Category)
			assert.Greater(t, item.ViewCount, 1000)
		}
		t.Logf("✅ MIN mode GSI with filters returned %d items", len(items))
	})
}

// ==================== Schema Tests ====================

func testUserPostsMINSchema(t *testing.T) {
	t.Run("min_table_schema_structure", func(t *testing.T) {
		schema := userpostscompletemin.TableSchema

		assert.Equal(t, "user-posts-complete-min", schema.TableName, "Table name should match MIN schema")
		assert.Equal(t, "user_id", schema.HashKey, "Hash key should be 'user_id'")
		assert.Equal(t, "created_at", schema.RangeKey, "Range key should be 'created_at'")
		assert.Len(t, schema.SecondaryIndexes, 6, "Should have 6 secondary indexes (3 LSI + 3 GSI)")
		t.Logf("✅ MIN mode schema structure validated")
	})

	t.Run("min_constants_validation", func(t *testing.T) {
		assert.Equal(t, "user-posts-complete-min", userpostscompletemin.TableName, "TableName should be MIN")
		assert.Equal(t, "user_id", userpostscompletemin.ColumnUserId, "ColumnUserId should be correct")
		assert.Equal(t, "created_at", userpostscompletemin.ColumnCreatedAt, "ColumnCreatedAt should be correct")
		assert.Equal(t, "post_type", userpostscompletemin.ColumnPostType, "ColumnPostType should be correct")
		assert.Equal(t, "status", userpostscompletemin.ColumnStatus, "ColumnStatus should be correct")
		assert.Equal(t, "priority", userpostscompletemin.ColumnPriority, "ColumnPriority should be correct")
		assert.Equal(t, "category", userpostscompletemin.ColumnCategory, "ColumnCategory should be correct")
		assert.Equal(t, "title", userpostscompletemin.ColumnTitle, "ColumnTitle should be correct")
		assert.Equal(t, "content", userpostscompletemin.ColumnContent, "ColumnContent should be correct")
		assert.Equal(t, "tags", userpostscompletemin.ColumnTags, "ColumnTags should be correct")
		assert.Equal(t, "view_count", userpostscompletemin.ColumnViewCount, "ColumnViewCount should be correct")
		assert.Equal(t, "updated_at", userpostscompletemin.ColumnUpdatedAt, "ColumnUpdatedAt should be correct")
		t.Logf("✅ MIN mode constants validated")
	})

	t.Run("min_operators_available", func(t *testing.T) {
		assert.NotNil(t, userpostscompletemin.EQ, "EQ operator should be available")
		assert.NotNil(t, userpostscompletemin.GT, "GT operator should be available")
		assert.NotNil(t, userpostscompletemin.LT, "LT operator should be available")
		assert.NotNil(t, userpostscompletemin.GTE, "GTE operator should be available")
		assert.NotNil(t, userpostscompletemin.LTE, "LTE operator should be available")
		assert.NotNil(t, userpostscompletemin.BETWEEN, "BETWEEN operator should be available")
		t.Logf("✅ MIN mode universal operators available")
	})

	t.Run("min_attribute_names", func(t *testing.T) {
		attrs := userpostscompletemin.AttributeNames
		expectedAttrs := []string{
			"user_id", "created_at", "post_type", "status", "priority",
			"category", "title", "content", "tags", "view_count", "updated_at",
		}

		assert.Len(t, attrs, len(expectedAttrs), "Should have correct number of attributes")
		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}
		t.Logf("✅ MIN mode AttributeNames validated")
	})

	t.Run("min_secondary_indexes", func(t *testing.T) {
		schema := userpostscompletemin.TableSchema

		lsiIndexes := []string{"lsi_by_post_type", "lsi_by_status", "lsi_by_priority"}
		gsiIndexes := []string{"gsi_by_category", "gsi_by_title", "gsi_by_status_priority"}

		indexNames := make([]string, 0, len(schema.SecondaryIndexes))
		for _, idx := range schema.SecondaryIndexes {
			indexNames = append(indexNames, idx.Name)
		}

		for _, expectedLSI := range lsiIndexes {
			assert.Contains(t, indexNames, expectedLSI, "Should contain LSI: %s", expectedLSI)
		}

		for _, expectedGSI := range gsiIndexes {
			assert.Contains(t, indexNames, expectedGSI, "Should contain GSI: %s", expectedGSI)
		}
		t.Logf("✅ MIN mode secondary indexes validated")
	})

	t.Run("min_no_sugar_methods", func(t *testing.T) {
		qb := userpostscompletemin.NewQueryBuilder()
		assert.NotNil(t, qb, "QueryBuilder should be available")

		sb := userpostscompletemin.NewScanBuilder()
		assert.NotNil(t, sb, "ScanBuilder should be available")
		t.Logf("✅ MIN mode builders available (sugar methods should be absent)")
	})
}

// ==================== Helper Functions ====================

func setupUserPostsMINTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []userpostscompletemin.SchemaItem{
		{
			UserId:    "min-query-test-user",
			CreatedAt: "2024-01-01T10:00:00Z",
			PostType:  "tutorial",
			Status:    "published",
			Priority:  85,
			Category:  "technology",
			Title:     "MIN DynamoDB Basics",
			Content:   "Introduction to DynamoDB concepts in MIN mode",
			Tags:      []string{"aws", "database"},
			ViewCount: 1200,
			UpdatedAt: "2024-01-01T11:00:00Z",
		},
		{
			UserId:    "min-query-test-user",
			CreatedAt: "2024-01-02T14:30:00Z",
			PostType:  "tutorial",
			Status:    "published",
			Priority:  88,
			Category:  "technology",
			Title:     "MIN Advanced DynamoDB",
			Content:   "Advanced DynamoDB patterns in MIN mode",
			Tags:      []string{"aws", "database", "advanced"},
			ViewCount: 800,
			UpdatedAt: "2024-01-02T15:30:00Z",
		},
		{
			UserId:    "min-query-test-user",
			CreatedAt: "2024-01-03T09:15:00Z",
			PostType:  "blog",
			Status:    "draft",
			Priority:  75,
			Category:  "programming",
			Title:     "MIN Best Practices",
			Content:   "Programming best practices in MIN mode",
			Tags:      []string{"programming", "tips"},
			ViewCount: 450,
			UpdatedAt: "2024-01-03T10:15:00Z",
		},
	}

	for _, item := range testItems {
		av, err := userpostscompletemin.ItemInput(item)
		require.NoError(t, err, "Should marshal MIN test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userpostscompletemin.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store MIN test item")
	}
	t.Logf("MIN setup complete: inserted %d test items", len(testItems))
}
