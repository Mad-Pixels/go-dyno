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

	userpostscomplete "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/userpostscomplete"
)

// TestUserPostsComplete focuses on GSI/LSI mixed operations and functionality.
// This test validates the new architecture with both GSI and LSI indexes.
//
// Test Coverage:
// - Mixed GSI/LSI CRUD operations
// - Index selection and query optimization
// - GSI vs LSI query performance
// - Complex projection types (ALL, KEYS_ONLY, INCLUDE)
// - Multiple index types on same table
//
// Schema: user-posts-complete.json
// - Table: "user-posts-complete"
// - Hash Key: user_id (S)
// - Range Key: created_at (S)
// - LSI: lsi_by_post_type, lsi_by_status, lsi_by_priority
// - GSI: gsi_by_category, gsi_by_title, gsi_by_status_priority
func TestUserPostsComplete(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(3 * time.Minute)
	defer cancel()

	t.Logf("Testing GSI/LSI operations on: %s", userpostscomplete.TableName)

	t.Run("UserPosts_Input", func(t *testing.T) {
		testUserPostsInput(t, client, ctx)
	})

	t.Run("UserPosts_QueryBuilder_LSI", func(t *testing.T) {
		testUserPostsQueryBuilderLSI(t, client, ctx)
	})

	t.Run("UserPosts_QueryBuilder_GSI", func(t *testing.T) {
		testUserPostsQueryBuilderGSI(t, client, ctx)
	})

	t.Run("UserPosts_Schema", func(t *testing.T) {
		t.Parallel()
		testUserPostsSchema(t)
	})
}

// ==================== User Posts Input ====================

func testUserPostsInput(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_complete_item", func(t *testing.T) {
		item := userpostscomplete.SchemaItem{
			UserId:    "user123",
			CreatedAt: "2024-01-15T10:30:00Z",
			PostType:  "blog",
			Status:    "published",
			Priority:  85,
			Category:  "technology",
			Title:     "Introduction to DynamoDB",
			Content:   "This is a comprehensive guide to DynamoDB",
			Tags:      []string{"aws", "database", "nosql"},
			ViewCount: 1500,
			UpdatedAt: "2024-01-16T09:15:00Z",
		}

		av, err := userpostscomplete.ItemInput(item)
		require.NoError(t, err, "Should marshal complete item")
		assert.NotEmpty(t, av, "Marshaled item should not be empty")

		// Verify all required fields are present
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

		// Verify correct DynamoDB types
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnUserId])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnCreatedAt])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnPostType])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnStatus])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[userpostscomplete.ColumnPriority])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnCategory])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnTitle])
		assert.IsType(t, &types.AttributeValueMemberS{}, av[userpostscomplete.ColumnContent])
		assert.IsType(t, &types.AttributeValueMemberSS{}, av[userpostscomplete.ColumnTags])
		assert.IsType(t, &types.AttributeValueMemberN{}, av[userpostscomplete.ColumnViewCount])

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userpostscomplete.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store complete item in DynamoDB")

		t.Logf("✅ Created complete user post: %s/%s", item.UserId, item.CreatedAt)
	})

	t.Run("read_complete_item", func(t *testing.T) {
		item := userpostscomplete.SchemaItem{
			UserId:    "user123",
			CreatedAt: "2024-01-15T10:30:00Z",
		}

		key, err := userpostscomplete.KeyInput(item)
		require.NoError(t, err, "Should create key from item")

		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(userpostscomplete.TableName),
			Key:       key,
		})
		require.NoError(t, err, "Should retrieve complete item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		// Verify key values
		assert.Equal(t, "user123", getResult.Item[userpostscomplete.ColumnUserId].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "2024-01-15T10:30:00Z", getResult.Item[userpostscomplete.ColumnCreatedAt].(*types.AttributeValueMemberS).Value)

		// Verify complex attributes
		assert.Equal(t, "blog", getResult.Item[userpostscomplete.ColumnPostType].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "published", getResult.Item[userpostscomplete.ColumnStatus].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "85", getResult.Item[userpostscomplete.ColumnPriority].(*types.AttributeValueMemberN).Value)
		assert.Equal(t, "technology", getResult.Item[userpostscomplete.ColumnCategory].(*types.AttributeValueMemberS).Value)
		assert.Equal(t, "Introduction to DynamoDB", getResult.Item[userpostscomplete.ColumnTitle].(*types.AttributeValueMemberS).Value)

		// Verify string set
		tagsSet := getResult.Item[userpostscomplete.ColumnTags].(*types.AttributeValueMemberSS)
		assert.Contains(t, tagsSet.Value, "aws")
		assert.Contains(t, tagsSet.Value, "database")
		assert.Contains(t, tagsSet.Value, "nosql")

		t.Logf("✅ Retrieved complete user post successfully")
	})
}

// ==================== LSI QueryBuilder Tests ====================

func testUserPostsQueryBuilderLSI(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	setupUserPostsTestData(t, client, ctx)

	t.Run("lsi_by_post_type", func(t *testing.T) {
		qb := userpostscomplete.NewQueryBuilder().
			WithEQ("user_id", "query-test-user").
			WithEQ("post_type", "tutorial")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build LSI query by post_type")

		// *** ДОБАВИТЬ ОТЛАДКУ ***
		t.Logf("IndexName: %v", queryInput.IndexName)
		t.Logf("KeyConditionExpression: %v", queryInput.KeyConditionExpression)
		t.Logf("FilterExpression: %v", queryInput.FilterExpression)

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute LSI query")

		t.Logf("Returned %d items", len(items))
		for i, item := range items {
			t.Logf("Item %d: user_id=%s, post_type=%s", i, item.UserId, item.PostType)
			assert.Equal(t, "query-test-user", item.UserId)
			assert.Equal(t, "tutorial", item.PostType)
		}
	})

	t.Run("lsi_by_status", func(t *testing.T) {
		qb := userpostscomplete.NewQueryBuilder().
			WithEQ("user_id", "query-test-user").
			WithEQ("status", "published")

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build LSI query by status")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute LSI status query")

		for _, item := range items {
			assert.Equal(t, "query-test-user", item.UserId)
			assert.Equal(t, "published", item.Status)
		}

		t.Logf("✅ LSI by status query returned %d items", len(items))
	})

	t.Run("lsi_by_priority_with_range", func(t *testing.T) {
		qb := userpostscomplete.NewQueryBuilder().
			WithEQ("user_id", "query-test-user").
			WithBetween("priority", 70, 90)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build LSI query with priority range")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute LSI priority range query")

		for _, item := range items {
			assert.Equal(t, "query-test-user", item.UserId)
			assert.GreaterOrEqual(t, item.Priority, 70)
			assert.LessOrEqual(t, item.Priority, 90)
		}

		t.Logf("✅ LSI priority range query returned %d items", len(items))
	})
}

// ==================== GSI QueryBuilder Tests ====================

func testUserPostsQueryBuilderGSI(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("gsi_by_category", func(t *testing.T) {
		qb := userpostscomplete.NewQueryBuilder().
			WithEQ("category", "technology")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build GSI query by category")

		// Should use GSI
		if queryInput.IndexName != nil {
			assert.Equal(t, "gsi_by_category", *queryInput.IndexName, "Should use GSI index")
		}

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute GSI category query")

		for _, item := range items {
			assert.Equal(t, "technology", item.Category, "All items should match category")
		}

		t.Logf("✅ GSI by category query returned %d items", len(items))
	})

	t.Run("gsi_by_title", func(t *testing.T) {
		qb := userpostscomplete.NewQueryBuilder().
			WithEQ("title", "Advanced DynamoDB")

		t.Logf("UsedKeys: %+v", qb.UsedKeys)
		t.Logf("Attributes: %+v", qb.Attributes)

		_, err := qb.BuildQuery()
		require.NoError(t, err, "Should build GSI query by title")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute GSI title query")

		for _, item := range items {
			assert.Equal(t, "Advanced DynamoDB", item.Title)
		}

		t.Logf("✅ GSI by title query returned %d items", len(items))
	})

	t.Run("gsi_status_priority_compound", func(t *testing.T) {
		qb := userpostscomplete.NewQueryBuilder().
			WithEQ("status", "published").
			WithGT("priority", 80)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build compound GSI query")

		items, err := qb.Execute(ctx, client)
		require.NoError(t, err, "Should execute compound GSI query")

		t.Logf("KeyCondition: %s", aws.ToString(queryInput.KeyConditionExpression))

		for _, item := range items {
			assert.Equal(t, "published", item.Status)
			assert.Greater(t, item.Priority, 80)
		}

		t.Logf("✅ GSI compound query returned %d items", len(items))
	})
}

// ==================== Schema Tests ====================

func testUserPostsSchema(t *testing.T) {
	t.Run("table_schema_structure", func(t *testing.T) {
		schema := userpostscomplete.TableSchema

		assert.Equal(t, "user-posts-complete", schema.TableName)
		assert.Equal(t, "user_id", schema.HashKey)
		assert.Equal(t, "created_at", schema.RangeKey)
		assert.Len(t, schema.SecondaryIndexes, 6, "Should have 6 secondary indexes (3 LSI + 3 GSI)")
	})

	t.Run("constants_validation", func(t *testing.T) {
		assert.Equal(t, "user-posts-complete", userpostscomplete.TableName)
		assert.Equal(t, "user_id", userpostscomplete.ColumnUserId)
		assert.Equal(t, "created_at", userpostscomplete.ColumnCreatedAt)
		assert.Equal(t, "post_type", userpostscomplete.ColumnPostType)
		assert.Equal(t, "status", userpostscomplete.ColumnStatus)
		assert.Equal(t, "priority", userpostscomplete.ColumnPriority)
		assert.Equal(t, "category", userpostscomplete.ColumnCategory)
		assert.Equal(t, "title", userpostscomplete.ColumnTitle)
		assert.Equal(t, "content", userpostscomplete.ColumnContent)
		assert.Equal(t, "tags", userpostscomplete.ColumnTags)
		assert.Equal(t, "view_count", userpostscomplete.ColumnViewCount)
		assert.Equal(t, "updated_at", userpostscomplete.ColumnUpdatedAt)
	})

	t.Run("attribute_names", func(t *testing.T) {
		attrs := userpostscomplete.AttributeNames
		expectedAttrs := []string{
			"user_id", "created_at", "post_type", "status", "priority",
			"category", "title", "content", "tags", "view_count", "updated_at",
		}

		assert.Len(t, attrs, len(expectedAttrs))
		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain %s", expected)
		}
	})
}

// ==================== Helper Functions ====================

func setupUserPostsTestData(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Helper()

	testItems := []userpostscomplete.SchemaItem{
		{
			UserId:    "query-test-user",
			CreatedAt: "2024-01-01T10:00:00Z",
			PostType:  "tutorial",
			Status:    "published",
			Priority:  85,
			Category:  "technology",
			Title:     "DynamoDB Basics",
			Content:   "Introduction to DynamoDB concepts",
			Tags:      []string{"aws", "database"},
			ViewCount: 1200,
			UpdatedAt: "2024-01-01T11:00:00Z",
		},
		{
			UserId:    "query-test-user",
			CreatedAt: "2024-01-02T14:30:00Z",
			PostType:  "tutorial",
			Status:    "published",
			Priority:  88,
			Category:  "technology",
			Title:     "Advanced DynamoDB",
			Content:   "Advanced DynamoDB patterns",
			Tags:      []string{"aws", "database", "advanced"},
			ViewCount: 800,
			UpdatedAt: "2024-01-02T15:30:00Z",
		},
		{
			UserId:    "query-test-user",
			CreatedAt: "2024-01-03T09:15:00Z",
			PostType:  "blog",
			Status:    "draft",
			Priority:  75,
			Category:  "programming",
			Title:     "Best Practices",
			Content:   "Programming best practices",
			Tags:      []string{"programming", "tips"},
			ViewCount: 450,
			UpdatedAt: "2024-01-03T10:15:00Z",
		},
	}

	for _, item := range testItems {
		av, err := userpostscomplete.ItemInput(item)
		require.NoError(t, err, "Should marshal test item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userpostscomplete.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store test item")
	}

	t.Logf("Setup complete: inserted %d test items", len(testItems))
}
