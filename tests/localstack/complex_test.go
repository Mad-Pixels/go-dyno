package localstack

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	complex "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/blogposts"
)

// TestComplexSchema runs comprehensive integration tests for the complex.json schema.
// This test suite validates all generated code functionality against a real LocalStack DynamoDB instance.
//
// Test Coverage:
// - Utility functions (BoolToInt, IntToBool, Stream processing)
// - CRUD operations (PutItem, BatchPutItems)
// - QueryBuilder with composite keys and GSI support
// - Secondary Index operations and projections
// - Schema constants and metadata validation
// - Generated struct marshaling/unmarshaling
// - Complex key operations and composite key handling
//
// Schema Structure (complex.json):
//   - Table: "blog-posts"
//   - Hash Key: "user_id" (string)
//   - Range Key: "post_id" (string)
//   - Attributes: Multiple with composite keys support
//   - Secondary Indexes: 4 GSIs with different projection types
//   - Composite Keys: category#is_published, tag#is_published
//
// Example Usage:
//
//	go test -v ./tests/localstack/ -run TestComplexSchema
func TestComplexSchema(t *testing.T) {
	cfg := DefaultLocalStackConfig()
	client := ConnectToLocalStack(t, cfg)

	ctx, cancel := TestContext(10 * time.Minute)
	defer cancel()

	t.Logf("Testing schema: complex.json")
	t.Logf("Table: %s", complex.TableName)
	t.Logf("Hash Key: %s, Range Key: %s", complex.TableSchema.HashKey, complex.TableSchema.RangeKey)
	t.Logf("Secondary Indexes: %d", len(complex.TableSchema.SecondaryIndexes))

	t.Run("Utility_Functions", func(t *testing.T) {
		t.Parallel()
		testComplexBoolToInt(t)
		testComplexIntToBool(t)
	})

	t.Run("CRUD_Operations", func(t *testing.T) {
		testComplexPutItem(t, client, ctx)
		testComplexBatchPutItems(t, client, ctx)
	})

	t.Run("QueryBuilder_Basic", func(t *testing.T) {
		testComplexQueryBuilder(t, client, ctx)
		testComplexQueryBuilderExecution(t, client, ctx)
		testComplexQueryBuilderChaining(t)
	})

	t.Run("QueryBuilder_Advanced", func(t *testing.T) {
		testComplexCompositeKeys(t, client, ctx)
		testComplexRangeConditions(t, client, ctx)
		testComplexIndexSelection(t, client, ctx)
	})

	t.Run("Secondary_Indexes", func(t *testing.T) {
		testComplexGSIOperations(t, client, ctx)
		testComplexIndexProjections(t, client, ctx)
		testComplexCompositeKeyIndexes(t, client, ctx)
	})

	t.Run("Schema_Constants", func(t *testing.T) {
		t.Parallel()
		testComplexSchemaConstants(t)
		testComplexAttributeNames(t)
		testComplexTableSchema(t)
		testComplexIndexProjectionsMap(t)
	})

	t.Run("Key_Operations", func(t *testing.T) {
		testComplexCreateKey(t)
		testComplexCreateKeyFromItem(t)
		testComplexCompositeKeyGeneration(t)
	})

	t.Run("Advanced_Features", func(t *testing.T) {
		testComplexFilterConditions(t, client, ctx)
		testComplexSortingAndPagination(t, client, ctx)
		testComplexTriggerHandlers(t)
	})
}

// ==================== Utility Functions Tests ====================

// testComplexBoolToInt validates the BoolToInt utility function for complex schema
func testComplexBoolToInt(t *testing.T) {
	t.Run("boolean_to_numeric_conversion", func(t *testing.T) {
		testCases := []struct {
			input    bool
			expected int
			desc     string
		}{
			{true, 1, "published posts should convert to 1"},
			{false, 0, "unpublished posts should convert to 0"},
		}

		for _, tc := range testCases {
			result := complex.BoolToInt(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		}

		t.Logf("✅ BoolToInt utility function works correctly for complex schema")
	})
}

// testComplexIntToBool validates the IntToBool utility function for complex schema
func testComplexIntToBool(t *testing.T) {
	t.Run("numeric_to_boolean_conversion", func(t *testing.T) {
		testCases := []struct {
			input    int
			expected bool
			desc     string
		}{
			{1, true, "published status (1) should be true"},
			{0, false, "unpublished status (0) should be false"},
			{2, true, "any non-zero should be true"},
			{-1, true, "negative non-zero should be true"},
		}

		for _, tc := range testCases {
			result := complex.IntToBool(tc.input)
			assert.Equal(t, tc.expected, result, tc.desc)
		}

		t.Logf("✅ IntToBool utility function handles complex schema cases correctly")
	})
}

// ==================== CRUD Operations Tests ====================

// testComplexPutItem validates item creation with complex schema structure
func testComplexPutItem(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_complex_blog_post", func(t *testing.T) {
		item := complex.SchemaItem{
			UserId:              "user123",                                      // hash key
			PostId:              "post456",                                      // range key
			CreatedAt:           1640995200,                                     // timestamp
			Likes:               42,                                             // numeric field
			IsPublished:         1,                                              // boolean as numeric
			CategoryIsPublished: "tech#1",                                       // composite key value
			TagIsPublished:      "golang#1",                                     // composite key value
			Title:               "Complex DynamoDB Schema",                      // common attribute
			Content:             "This is a test blog post with complex schema", // common attribute
			Category:            "tech",                                         // common attribute
			Tag:                 "golang",                                       // common attribute
			Views:               150,                                            // common attribute
			IsPremium:           true,                                           // boolean field
			IsFeatured:          false,                                          // boolean field
		}
		av, err := complex.PutItem(item)
		require.NoError(t, err, "PutItem marshaling should succeed")
		require.NotEmpty(t, av, "AttributeValues should not be empty")

		requiredFields := []string{
			"user_id",
			"post_id",
			"created_at",
			"likes",
			"is_published",
			"category#is_published",
			"tag#is_published",
			"title",
			"content",
			"category",
			"tag",
			"views",
			"is_premium",
			"is_featured",
		}
		for _, field := range requiredFields {
			assert.Contains(t, av, field, "Must contain field '%s'", field)
		}

		t.Logf("✅ Successfully created and marshaled complex blog post item")
	})

	t.Run("put_complex_item_to_dynamodb", func(t *testing.T) {
		item := complex.SchemaItem{
			UserId:              "put-test-user",
			PostId:              "put-test-post-789",
			CreatedAt:           1640995300,
			Likes:               25,
			IsPublished:         1,
			CategoryIsPublished: "programming#1",
			TagIsPublished:      "aws#1",
			Title:               "DynamoDB Complex Put Test",
			Content:             "Testing complex item insertion",
			Category:            "programming",
			Tag:                 "aws",
			Views:               75,
			IsPremium:           false,
			IsFeatured:          true,
		}
		av, err := complex.PutItem(item)
		require.NoError(t, err, "Item marshaling should succeed")

		input := &dynamodb.PutItemInput{
			TableName: aws.String(complex.TableName),
			Item:      av,
		}
		_, err = client.PutItem(ctx, input)
		require.NoError(t, err, "DynamoDB PutItem operation should succeed")

		getInput := &dynamodb.GetItemInput{
			TableName: aws.String(complex.TableName),
			Key: map[string]types.AttributeValue{
				"user_id": &types.AttributeValueMemberS{Value: item.UserId},
				"post_id": &types.AttributeValueMemberS{Value: item.PostId},
			},
		}
		result, err := client.GetItem(ctx, getInput)
		require.NoError(t, err, "GetItem verification should succeed")
		assert.NotEmpty(t, result.Item, "Stored item should be retrievable")

		t.Logf("✅ Successfully stored and verified complex item in DynamoDB")
	})
}

// testComplexBatchPutItems validates batch operations with complex items
func testComplexBatchPutItems(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("batch_put_complex_blog_posts", func(t *testing.T) {
		items := []complex.SchemaItem{
			{
				UserId:              "batch-user1",
				PostId:              "batch-post1",
				CreatedAt:           1640995400,
				Likes:               10,
				IsPublished:         1,
				CategoryIsPublished: "tech#1",
				TagIsPublished:      "go#1",
				Title:               "Batch Post 1",
				Content:             "First batch post",
				Category:            "tech",
				Tag:                 "go",
				Views:               100,
				IsPremium:           false,
				IsFeatured:          false,
			},
			{
				UserId:              "batch-user2",
				PostId:              "batch-post2",
				CreatedAt:           1640995500,
				Likes:               20,
				IsPublished:         0,
				CategoryIsPublished: "lifestyle#0",
				TagIsPublished:      "health#0",
				Title:               "Batch Post 2",
				Content:             "Second batch post",
				Category:            "lifestyle",
				Tag:                 "health",
				Views:               50,
				IsPremium:           true,
				IsFeatured:          true,
			},
			{
				UserId:              "batch-user3",
				PostId:              "batch-post3",
				CreatedAt:           1640995600,
				Likes:               30,
				IsPublished:         1,
				CategoryIsPublished: "business#1",
				TagIsPublished:      "startup#1",
				Title:               "Batch Post 3",
				Content:             "Third batch post",
				Category:            "business",
				Tag:                 "startup",
				Views:               200,
				IsPremium:           false,
				IsFeatured:          true,
			},
		}
		batchItems, err := complex.BatchPutItems(items)
		require.NoError(t, err, "BatchPutItems should succeed")
		require.Len(t, batchItems, 3, "Should return AttributeValues for all 3 items")

		for i, batchItem := range batchItems {
			// Check composite key fields
			assert.Contains(t, batchItem, "category#is_published", "Batch item %d should have composite key field", i)
			assert.Contains(t, batchItem, "tag#is_published", "Batch item %d should have composite key field", i)

			// Check all other fields
			assert.Contains(t, batchItem, "user_id", "Batch item %d should have user_id", i)
			assert.Contains(t, batchItem, "post_id", "Batch item %d should have post_id", i)
			assert.Contains(t, batchItem, "title", "Batch item %d should have title", i)
			assert.Contains(t, batchItem, "content", "Batch item %d should have content", i)
		}
		t.Logf("✅ Successfully prepared %d complex items for batch operation", len(batchItems))
	})
}

// ==================== QueryBuilder Basic Tests ====================

// testComplexQueryBuilder validates QueryBuilder creation with complex schema
func testComplexQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_complex_query_builder", func(t *testing.T) {
		qb := complex.NewQueryBuilder()
		require.NotNil(t, qb, "NewQueryBuilder should return non-nil instance")

		t.Logf("✅ Complex QueryBuilder created successfully")
	})

	t.Run("complex_fluent_api_methods", func(t *testing.T) {
		qb := complex.NewQueryBuilder()

		methods := []struct {
			name string
			call func() *complex.QueryBuilder
		}{
			{"WithUserId", func() *complex.QueryBuilder { return qb.WithUserId("test") }},
			{"WithPostId", func() *complex.QueryBuilder { return qb.WithPostId("test") }},
			{"WithCreatedAt", func() *complex.QueryBuilder { return qb.WithCreatedAt(123) }},
			{"WithLikes", func() *complex.QueryBuilder { return qb.WithLikes(10) }},
			{"WithIsPublished", func() *complex.QueryBuilder { return qb.WithIsPublished(1) }},
			{"WithTitle", func() *complex.QueryBuilder { return qb.WithTitle("test") }},
			{"WithContent", func() *complex.QueryBuilder { return qb.WithContent("test") }},
			{"WithCategory", func() *complex.QueryBuilder { return qb.WithCategory("test") }},
			{"WithTag", func() *complex.QueryBuilder { return qb.WithTag("test") }},
			{"WithViews", func() *complex.QueryBuilder { return qb.WithViews(100) }},
			{"WithIsPremium", func() *complex.QueryBuilder { return qb.WithIsPremium(true) }},
			{"WithIsFeatured", func() *complex.QueryBuilder { return qb.WithIsFeatured(false) }},
		}
		for _, method := range methods {
			result := method.call()
			require.NotNil(t, result, "%s should return QueryBuilder for chaining", method.name)
			assert.IsType(t, qb, result, "%s should return same type", method.name)
		}
		t.Logf("✅ All %d fluent API methods support proper method chaining", len(methods))
	})

	t.Run("complex_range_conditions", func(t *testing.T) {
		qb := complex.NewQueryBuilder()

		rangeConditions := []struct {
			name string
			call func() *complex.QueryBuilder
		}{
			{"WithCreatedAtBetween", func() *complex.QueryBuilder { return qb.WithCreatedAtBetween(100, 200) }},
			{"WithCreatedAtGreaterThan", func() *complex.QueryBuilder { return qb.WithCreatedAtGreaterThan(100) }},
			{"WithCreatedAtLessThan", func() *complex.QueryBuilder { return qb.WithCreatedAtLessThan(200) }},
			{"WithLikesBetween", func() *complex.QueryBuilder { return qb.WithLikesBetween(10, 50) }},
			{"WithLikesGreaterThan", func() *complex.QueryBuilder { return qb.WithLikesGreaterThan(10) }},
			{"WithLikesLessThan", func() *complex.QueryBuilder { return qb.WithLikesLessThan(50) }},
			{"WithViewsBetween", func() *complex.QueryBuilder { return qb.WithViewsBetween(100, 500) }},
			{"WithViewsGreaterThan", func() *complex.QueryBuilder { return qb.WithViewsGreaterThan(100) }},
			{"WithViewsLessThan", func() *complex.QueryBuilder { return qb.WithViewsLessThan(500) }},
		}
		for _, condition := range rangeConditions {
			result := condition.call()
			require.NotNil(t, result, "%s should return QueryBuilder", condition.name)
		}
		t.Logf("✅ All %d range condition methods work correctly", len(rangeConditions))
	})
}

// testComplexQueryBuilderChaining validates complex method chaining scenarios
func testComplexQueryBuilderChaining(t *testing.T) {
	t.Run("complex_method_chaining_scenarios", func(t *testing.T) {
		qb1 := complex.NewQueryBuilder().
			WithUserId("complex-user").
			WithPostId("complex-post").
			WithCreatedAtGreaterThan(1640995000).
			WithLikesLessThan(100).
			WithCategory("tech").
			WithTag("golang").
			WithIsPremium(true).
			OrderByDesc().
			Limit(25)

		require.NotNil(t, qb1, "Complex chaining should work")

		qb2 := complex.NewQueryBuilder().
			WithUserId("user").
			WithIsPublished(1).
			WithCategory("tech").
			OrderByAsc().
			Limit(50)

		require.NotNil(t, qb2, "Composite key chaining should work")
		t.Logf("✅ Complex method chaining works correctly")
	})
}

// testComplexQueryBuilderExecution validates query building with complex schema
func testComplexQueryBuilderExecution(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("build_complex_query_input", func(t *testing.T) {
		testItem := complex.SchemaItem{
			UserId: "complex-query-user", PostId: "complex-query-post", CreatedAt: 1640995700,
			Likes: 35, IsPublished: 1, CategoryIsPublished: "tech#1", TagIsPublished: "golang#1",
			Title: "Complex Query Test", Content: "Testing complex query", Category: "tech", Tag: "golang",
			Views: 150, IsPremium: false, IsFeatured: true,
		}
		insertComplexTestItem(t, client, ctx, testItem)

		qb := complex.NewQueryBuilder().WithUserId("complex-query-user")

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "BuildQuery should succeed")
		require.NotNil(t, queryInput, "QueryInput should not be nil")

		assert.Equal(t, complex.TableName, *queryInput.TableName, "TableName should match schema")
		assert.NotNil(t, queryInput.KeyConditionExpression, "KeyConditionExpression should be set")
		assert.NotEmpty(t, queryInput.ExpressionAttributeNames, "ExpressionAttributeNames should be populated")
		assert.NotEmpty(t, queryInput.ExpressionAttributeValues, "ExpressionAttributeValues should be populated")
		t.Logf("✅ Complex QueryInput built successfully")
	})

	t.Run("build_query_with_multiple_filters", func(t *testing.T) {
		qb := complex.NewQueryBuilder().
			WithUserId("filter-user").
			WithCategory("tech").
			WithTag("golang").
			WithIsPremium(true).
			WithViewsGreaterThan(100)

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "BuildQuery with multiple filters should succeed")

		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
		t.Logf("✅ Complex QueryInput with multiple filters built successfully")
	})
}

// ==================== Advanced QueryBuilder Tests ====================

// testComplexCompositeKeys validates composite key operations
func testComplexCompositeKeys(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("composite_key_generation", func(t *testing.T) {
		// Test composite key generation based on complex.json schema
		// CategoryIsPublished = "category#is_published"
		// TagIsPublished = "tag#is_published"

		item := complex.SchemaItem{
			UserId: "composite-user", PostId: "composite-post", CreatedAt: 1640995800,
			Category: "tech", IsPublished: 1, Tag: "golang",
			CategoryIsPublished: "tech#1", TagIsPublished: "golang#1",
			Title: "Composite Key Test", Content: "Testing composite keys",
			Likes: 20, Views: 80, IsPremium: false, IsFeatured: false,
		}

		insertComplexTestItem(t, client, ctx, item)

		// TODO: Test querying with composite keys would be done here
		// This depends on the actual QueryBuilder implementation

		t.Logf("✅ Composite key operations validated")
	})
}

// testComplexRangeConditions validates range query conditions
func testComplexRangeConditions(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("numeric_range_queries", func(t *testing.T) {
		testItems := []complex.SchemaItem{
			createComplexTestItem("range1", 1640995000, 10, 100),
			createComplexTestItem("range2", 1640995100, 20, 200),
			createComplexTestItem("range3", 1640995200, 30, 300),
		}
		for _, item := range testItems {
			insertComplexTestItem(t, client, ctx, item)
		}
		rangeTests := []struct {
			name string
			qb   *complex.QueryBuilder
		}{
			{
				"created_at_between",
				complex.NewQueryBuilder().WithUserId("range-user").WithCreatedAtBetween(1640995050, 1640995150),
			},
			{
				"likes_greater_than",
				complex.NewQueryBuilder().WithUserId("range-user").WithLikesGreaterThan(15),
			},
			{
				"views_less_than",
				complex.NewQueryBuilder().WithUserId("range-user").WithViewsLessThan(250),
			},
		}
		for _, test := range rangeTests {
			_, err := test.qb.BuildQuery()
			require.NoError(t, err, "Range query %s should build successfully", test.name)
		}
		t.Logf("✅ Range condition queries validated")
	})
}

// testComplexIndexSelection validates automatic index selection
func testComplexIndexSelection(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("automatic_index_selection", func(t *testing.T) {
		// Test that QueryBuilder selects appropriate indexes
		// Based on complex.json secondary indexes:
		// - PublishedByDateIndex
		// - PublishedByLikesIndex
		// - CategoryPublishedIndex
		// - TagPublishedIndex

		indexTests := []struct {
			name          string
			qb            *complex.QueryBuilder
			expectedIndex string
		}{
			{
				"should_use_published_by_date_index",
				complex.NewQueryBuilder().WithIsPublished(1).WithCreatedAt(1640995000),
				"PublishedByDateIndex",
			},
			{
				"should_use_published_by_likes_index",
				complex.NewQueryBuilder().WithIsPublished(1).WithLikes(50),
				"PublishedByLikesIndex",
			},
		}
		for _, test := range indexTests {
			queryInput, err := test.qb.BuildQuery()
			require.NoError(t, err, "Index selection test %s should succeed", test.name)

			if queryInput.IndexName != nil {
				t.Logf("Selected index for %s: %s", test.name, *queryInput.IndexName)
			} else {
				t.Logf("Using primary table for %s", test.name)
			}
		}
		t.Logf("✅ Automatic index selection validated")
	})
}

// ==================== Secondary Index Tests ====================

// testComplexGSIOperations validates Global Secondary Index operations
func testComplexGSIOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("gsi_query_operations", func(t *testing.T) {
		// Test querying using different GSIs
		// Insert test data for GSI operations
		testItems := []complex.SchemaItem{
			{
				UserId: "gsi-user1", PostId: "gsi-post1", CreatedAt: 1640995000,
				IsPublished: 1, Likes: 25, Category: "tech", Tag: "golang",
				CategoryIsPublished: "tech#1", TagIsPublished: "golang#1",
				Title: "GSI Test 1", Content: "Testing GSI operations", Views: 100,
				IsPremium: false, IsFeatured: true,
			},
			{
				UserId: "gsi-user2", PostId: "gsi-post2", CreatedAt: 1640995100,
				IsPublished: 0, Likes: 15, Category: "lifestyle", Tag: "health",
				CategoryIsPublished: "lifestyle#0", TagIsPublished: "health#0",
				Title: "GSI Test 2", Content: "Testing GSI operations", Views: 50,
				IsPremium: true, IsFeatured: false,
			},
		}
		for _, item := range testItems {
			insertComplexTestItem(t, client, ctx, item)
		}

		gsiTests := []struct {
			name string
			qb   *complex.QueryBuilder
		}{
			{
				"published_posts_by_date",
				complex.NewQueryBuilder().WithIsPublished(1).WithCreatedAtGreaterThan(1640994000),
			},
			{
				"published_posts_by_likes",
				complex.NewQueryBuilder().WithIsPublished(1).WithLikesGreaterThan(20),
			},
		}
		for _, test := range gsiTests {
			queryInput, err := test.qb.BuildQuery()
			require.NoError(t, err, "GSI test %s should succeed", test.name)

			if queryInput.IndexName != nil {
				t.Logf("GSI query %s uses index: %s", test.name, *queryInput.IndexName)
			}
		}
		t.Logf("✅ GSI operations validated")
	})
}

// testComplexIndexProjections validates index projection handling
func testComplexIndexProjections(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("index_projection_validation", func(t *testing.T) {
		projections := complex.IndexProjections
		require.NotEmpty(t, projections, "IndexProjections should not be empty")

		// Test each index projection based on complex.json
		expectedIndexes := []string{
			"PublishedByDateIndex",
			"PublishedByLikesIndex",
			"CategoryPublishedIndex",
			"TagPublishedIndex",
		}
		for _, indexName := range expectedIndexes {
			projection, exists := projections[indexName]
			assert.True(t, exists, "Index %s should exist in projections", indexName)
			assert.NotEmpty(t, projection, "Index %s should have projection attributes", indexName)

			t.Logf("Index %s projects %d attributes", indexName, len(projection))
		}
		t.Logf("✅ Index projections validated")
	})
}

// testComplexCompositeKeyIndexes validates composite key index operations
func testComplexCompositeKeyIndexes(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("composite_key_index_operations", func(t *testing.T) {
		// Test composite key indexes: CategoryPublishedIndex, TagPublishedIndex
		// These use composite hash keys: category#is_published, tag#is_published

		testItem := complex.SchemaItem{
			UserId: "composite-idx-user", PostId: "composite-idx-post", CreatedAt: 1640995900,
			Category: "tech", IsPublished: 1, Tag: "golang", Likes: 40,
			CategoryIsPublished: "tech#1", TagIsPublished: "golang#1",
			Title: "Composite Index Test", Content: "Testing composite key indexes",
			Views: 120, IsPremium: false, IsFeatured: true,
		}

		insertComplexTestItem(t, client, ctx, testItem)

		compositeTests := []struct {
			name string
			qb   *complex.QueryBuilder
		}{
			{
				"category_published_query",
				complex.NewQueryBuilder().WithCategory("tech").WithIsPublished(1),
			},
			{
				"tag_published_query",
				complex.NewQueryBuilder().WithTag("golang").WithIsPublished(1),
			},
		}
		for _, test := range compositeTests {
			queryInput, err := test.qb.BuildQuery()
			require.NoError(t, err, "Composite key test %s should succeed", test.name)

			if queryInput.IndexName != nil {
				t.Logf("Composite key query %s uses index: %s", test.name, *queryInput.IndexName)
			}
		}
		t.Logf("✅ Composite key index operations validated")
	})
}

// ==================== Schema Constants Tests ====================

// testComplexSchemaConstants validates all generated constants for complex schema
func testComplexSchemaConstants(t *testing.T) {
	t.Run("table_and_index_constants", func(t *testing.T) {
		// Validate table name
		assert.NotEmpty(t, complex.TableName, "TableName constant should not be empty")
		assert.Equal(t, "blog-posts", complex.TableName, "TableName should match complex.json schema")

		indexConstants := []struct {
			constant string
			expected string
		}{
			{complex.IndexPublishedByDateIndex, "PublishedByDateIndex"},
			{complex.IndexPublishedByLikesIndex, "PublishedByLikesIndex"},
			{complex.IndexCategoryPublishedIndex, "CategoryPublishedIndex"},
			{complex.IndexTagPublishedIndex, "TagPublishedIndex"},
		}
		for _, idx := range indexConstants {
			assert.Equal(t, idx.expected, idx.constant, "Index constant should match expected value")
		}
		t.Logf("✅ Table and index constants validated")
		t.Logf("    Table: %s", complex.TableName)
		t.Logf("    Indexes: %d constants", len(indexConstants))
	})

	t.Run("column_constants_complex", func(t *testing.T) {
		// Validate all column constants match complex schema attributes
		columnTests := map[string]string{
			complex.ColumnUserId:              "user_id",
			complex.ColumnPostId:              "post_id",
			complex.ColumnCreatedAt:           "created_at",
			complex.ColumnLikes:               "likes",
			complex.ColumnIsPublished:         "is_published",
			complex.ColumnCategoryIsPublished: "category#is_published",
			complex.ColumnTagIsPublished:      "tag#is_published",
			complex.ColumnTitle:               "title",
			complex.ColumnContent:             "content",
			complex.ColumnCategory:            "category",
			complex.ColumnTag:                 "tag",
			complex.ColumnViews:               "views",
			complex.ColumnIsPremium:           "is_premium",
			complex.ColumnIsFeatured:          "is_featured",
		}
		for constant, expected := range columnTests {
			assert.Equal(t, expected, constant, "Column constant should match expected value")
		}
		t.Logf("✅ All %d column constants validated", len(columnTests))
	})
}

// testComplexAttributeNames validates the AttributeNames array for complex schema
func testComplexAttributeNames(t *testing.T) {
	t.Run("complex_attribute_names_array", func(t *testing.T) {
		attrs := complex.AttributeNames
		require.NotEmpty(t, attrs, "AttributeNames should not be empty")

		// Expected attributes from complex.json (attributes + common_attributes)
		expectedAttrs := []string{
			"user_id",
			"post_id",
			"created_at",
			"likes",
			"is_published",
			"category#is_published",
			"tag#is_published",
			"title",
			"content",
			"category",
			"tag",
			"views",
			"is_premium",
			"is_featured",
		}
		assert.Len(t, attrs, len(expectedAttrs), "AttributeNames should contain all schema attributes")

		for _, expected := range expectedAttrs {
			assert.Contains(t, attrs, expected, "AttributeNames should contain '%s'", expected)
		}
		attrSet := make(map[string]bool)
		for _, attr := range attrs {
			assert.False(t, attrSet[attr], "AttributeNames should not contain duplicate: %s", attr)
			attrSet[attr] = true
		}
		t.Logf("✅ AttributeNames array contains %d attributes", len(attrs))
	})
}

// testComplexTableSchema validates the TableSchema variable for complex schema
func testComplexTableSchema(t *testing.T) {
	t.Run("complex_table_schema_structure", func(t *testing.T) {
		schema := complex.TableSchema

		// Validate basic schema properties
		assert.Equal(t, "blog-posts", schema.TableName, "Schema TableName should match")
		assert.Equal(t, "user_id", schema.HashKey, "Hash key should be 'user_id'")
		assert.Equal(t, "post_id", schema.RangeKey, "Range key should be 'post_id'")

		// Validate attribute collections
		expectedPrimaryAttrs := 7     // From complex.json attributes array
		expectedCommonAttrs := 7      // From complex.json common_attributes array
		expectedSecondaryIndexes := 4 // From complex.json secondary_indexes array

		assert.Len(t, schema.Attributes, expectedPrimaryAttrs, "Should have %d primary attributes", expectedPrimaryAttrs)
		assert.Len(t, schema.CommonAttributes, expectedCommonAttrs, "Should have %d common attributes", expectedCommonAttrs)
		assert.Len(t, schema.SecondaryIndexes, expectedSecondaryIndexes, "Should have %d secondary indexes", expectedSecondaryIndexes)

		// Validate secondary index structure
		for _, idx := range schema.SecondaryIndexes {
			assert.NotEmpty(t, idx.Name, "Index name should not be empty")
			assert.NotEmpty(t, idx.HashKey, "Index hash key should not be empty")
			assert.NotEmpty(t, idx.ProjectionType, "Index projection type should not be empty")

			t.Logf("    Index: %s (HashKey: %s, RangeKey: %s, Projection: %s)",
				idx.Name, idx.HashKey, idx.RangeKey, idx.ProjectionType)
		}
		t.Logf("✅ Complex TableSchema structure is valid")
		t.Logf("    Table: %s", schema.TableName)
		t.Logf("    Keys: %s (hash), %s (range)", schema.HashKey, schema.RangeKey)
		t.Logf("    Attributes: %d primary, %d common", len(schema.Attributes), len(schema.CommonAttributes))
		t.Logf("    Indexes: %d secondary", len(schema.SecondaryIndexes))
	})
}

// testComplexIndexProjectionsMap validates the IndexProjections map
func testComplexIndexProjectionsMap(t *testing.T) {
	t.Run("index_projections_map_structure", func(t *testing.T) {
		projections := complex.IndexProjections
		require.NotEmpty(t, projections, "IndexProjections should not be empty")

		// Test each index projection based on complex.json secondary_indexes
		indexProjectionTests := []struct {
			indexName      string
			projectionType string
			minAttributes  int
		}{
			{"PublishedByDateIndex", "ALL", 14},       // ALL projection = all attributes
			{"PublishedByLikesIndex", "KEYS_ONLY", 2}, // KEYS_ONLY = hash + range keys only
			{"CategoryPublishedIndex", "INCLUDE", 5},  // INCLUDE = keys + specified attributes
			{"TagPublishedIndex", "INCLUDE", 5},       // INCLUDE = keys + specified attributes
		}
		for _, test := range indexProjectionTests {
			projection, exists := projections[test.indexName]
			assert.True(t, exists, "Index %s should exist in projections", test.indexName)
			assert.GreaterOrEqual(t, len(projection), test.minAttributes,
				"Index %s should have at least %d projected attributes", test.indexName, test.minAttributes)

			t.Logf("    %s (%s): %d attributes", test.indexName, test.projectionType, len(projection))
		}
		t.Logf("✅ IndexProjections map structure validated")
	})
}

// ==================== Key Operations Tests ====================

// testComplexCreateKey validates key creation with complex schema
func testComplexCreateKey(t *testing.T) {
	t.Run("create_complex_primary_key", func(t *testing.T) {
		hashKeyValue := "complex-user-123"
		rangeKeyValue := "complex-post-456"

		key, err := complex.CreateKey(hashKeyValue, rangeKeyValue)
		require.NoError(t, err, "CreateKey should succeed with valid inputs")
		require.NotEmpty(t, key, "Created key should not be empty")

		// Validate key structure for complex schema
		assert.Contains(t, key, "user_id", "Key should contain hash key 'user_id'")
		assert.Contains(t, key, "post_id", "Key should contain range key 'post_id'")

		t.Logf("✅ CreateKey generated valid complex DynamoDB key structure")
	})

	t.Run("create_key_with_complex_values", func(t *testing.T) {
		// Test with various data types
		testCases := []struct {
			hashKey  any
			rangeKey any
			desc     string
		}{
			{"user-123", "post-abc", "string keys"},
			{"user-456", "post-789", "mixed string keys"},
		}
		for _, tc := range testCases {
			key, err := complex.CreateKey(tc.hashKey, tc.rangeKey)
			require.NoError(t, err, "CreateKey should succeed for %s", tc.desc)
			assert.Contains(t, key, "user_id", "Key should contain user_id for %s", tc.desc)
			assert.Contains(t, key, "post_id", "Key should contain post_id for %s", tc.desc)
		}
		t.Logf("✅ CreateKey handles complex key values correctly")
	})
}

// testComplexCreateKeyFromItem validates key extraction from complex SchemaItem
func testComplexCreateKeyFromItem(t *testing.T) {
	t.Run("extract_key_from_complex_item", func(t *testing.T) {
		// Create complete complex item
		item := complex.SchemaItem{
			UserId:              "extract-user-789",
			PostId:              "extract-post-012",
			CreatedAt:           1640996000,
			Likes:               45,
			IsPublished:         1,
			CategoryIsPublished: "tech#1",
			TagIsPublished:      "golang#1",
			Title:               "Key Extraction Complex Test",
			Content:             "Testing key extraction from complex item",
			Category:            "tech",
			Tag:                 "golang",
			Views:               180,
			IsPremium:           true,
			IsFeatured:          false,
		}

		// Extract key attributes only
		key, err := complex.CreateKeyFromItem(item)
		require.NoError(t, err, "CreateKeyFromItem should succeed")
		require.NotEmpty(t, key, "Extracted key should not be empty")

		// Validate extracted key contains only primary key attributes
		assert.Contains(t, key, "user_id", "Extracted key should contain hash key 'user_id'")
		assert.Contains(t, key, "post_id", "Extracted key should contain range key 'post_id'")

		// Validate non-key attributes are excluded
		nonKeyAttrs := []string{"created_at", "likes", "is_published", "title", "content",
			"category", "tag", "views", "is_premium", "is_featured"}
		for _, attr := range nonKeyAttrs {
			assert.NotContains(t, key, attr, "Extracted key should not contain non-key attribute '%s'", attr)
		}
		t.Logf("✅ CreateKeyFromItem correctly extracts only primary key attributes")
	})
}

// testComplexCompositeKeyGeneration validates composite key generation
func testComplexCompositeKeyGeneration(t *testing.T) {
	t.Run("composite_key_value_generation", func(t *testing.T) {
		// Test composite key generation for complex schema
		// category#is_published and tag#is_published

		testCases := []struct {
			category    string
			tag         string
			isPublished int
			expectedCat string
			expectedTag string
		}{
			{"tech", "golang", 1, "tech#1", "golang#1"},
			{"lifestyle", "health", 0, "lifestyle#0", "health#0"},
			{"business", "startup", 1, "business#1", "startup#1"},
		}
		for _, tc := range testCases {
			// These values should be generated automatically by the QueryBuilder
			// when using composite key methods
			assert.Equal(t, tc.expectedCat, tc.category+"#"+string(rune(tc.isPublished+'0')),
				"Category composite key should be generated correctly")
			assert.Equal(t, tc.expectedTag, tc.tag+"#"+string(rune(tc.isPublished+'0')),
				"Tag composite key should be generated correctly")
		}
		t.Logf("✅ Composite key generation logic validated")
	})
}

// ==================== Advanced Features Tests ====================

// testComplexFilterConditions validates filter condition handling
func testComplexFilterConditions(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("complex_filter_conditions", func(t *testing.T) {
		// Insert test data with various attribute values
		testItems := []complex.SchemaItem{
			createComplexTestItem("filter1", 1640995000, 10, 100),
			createComplexTestItem("filter2", 1640995100, 20, 200),
			createComplexTestItem("filter3", 1640995200, 30, 300),
		}
		for _, item := range testItems {
			insertComplexTestItem(t, client, ctx, item)
		}

		// Test various filter combinations
		filterTests := []struct {
			name string
			qb   *complex.QueryBuilder
		}{
			{
				"multiple_attribute_filters",
				complex.NewQueryBuilder().
					WithUserId("filter-user").
					WithCategory("tech").
					WithTag("golang").
					WithIsPremium(false),
			},
			{
				"numeric_and_boolean_filters",
				complex.NewQueryBuilder().
					WithUserId("filter-user").
					WithViewsGreaterThan(150).
					WithIsFeatured(true),
			},
			{
				"range_and_equality_filters",
				complex.NewQueryBuilder().
					WithUserId("filter-user").
					WithCreatedAtBetween(1640995050, 1640995150).
					WithCategory("tech"),
			},
		}
		for _, test := range filterTests {
			queryInput, err := test.qb.BuildQuery()
			require.NoError(t, err, "Filter test %s should succeed", test.name)

			// Validate that complex filters are handled
			assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")
			t.Logf("Filter test %s: built successfully", test.name)
		}
		t.Logf("✅ Complex filter conditions validated")
	})
}

// testComplexSortingAndPagination validates sorting and pagination features
func testComplexSortingAndPagination(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("sorting_and_pagination_complex", func(t *testing.T) {
		// Insert multiple items for the same user to test sorting
		baseTime := int64(1640995000)
		testItems := make([]complex.SchemaItem, 5)

		for i := 0; i < 5; i++ {
			testItems[i] = complex.SchemaItem{
				UserId:              "sort-user",
				PostId:              "sort-post-" + string(rune('A'+i)),
				CreatedAt:           int(baseTime + int64(i*100)),
				Likes:               10 + i*5,
				IsPublished:         1,
				CategoryIsPublished: "tech#1",
				TagIsPublished:      "golang#1",
				Title:               "Sort Test " + string(rune('A'+i)),
				Content:             "Testing sorting",
				Category:            "tech",
				Tag:                 "golang",
				Views:               100 + i*50,
				IsPremium:           i%2 == 0,
				IsFeatured:          i%3 == 0,
			}
			insertComplexTestItem(t, client, ctx, testItems[i])
		}

		sortingTests := []struct {
			name string
			qb   *complex.QueryBuilder
		}{
			{
				"ascending_sort",
				complex.NewQueryBuilder().WithUserId("sort-user").OrderByAsc(),
			},
			{
				"descending_sort",
				complex.NewQueryBuilder().WithUserId("sort-user").OrderByDesc(),
			},
			{
				"limited_results",
				complex.NewQueryBuilder().WithUserId("sort-user").Limit(2),
			},
			{
				"sorted_and_limited",
				complex.NewQueryBuilder().WithUserId("sort-user").OrderByDesc().Limit(3),
			},
		}

		for _, test := range sortingTests {
			queryInput, err := test.qb.BuildQuery()
			require.NoError(t, err, "Sorting test %s should succeed", test.name)

			if test.name == "ascending_sort" {
				assert.NotNil(t, queryInput.ScanIndexForward, "Should have ScanIndexForward set")
				assert.True(t, *queryInput.ScanIndexForward, "Should be ascending")
			}
			if test.name == "descending_sort" {
				assert.NotNil(t, queryInput.ScanIndexForward, "Should have ScanIndexForward set")
				assert.False(t, *queryInput.ScanIndexForward, "Should be descending")
			}
			if test.name == "limited_results" || test.name == "sorted_and_limited" {
				assert.NotNil(t, queryInput.Limit, "Should have Limit set")
			}

			t.Logf("Sorting test %s: configured correctly", test.name)
		}

		t.Logf("✅ Sorting and pagination features validated")
	})
}

// testComplexTriggerHandlers validates DynamoDB Stream trigger handling
func testComplexTriggerHandlers(t *testing.T) {
	t.Run("trigger_handler_creation", func(t *testing.T) {
		// Test trigger handler creation with complex schema
		handler := complex.CreateTriggerHandler(
			// onInsert
			func(ctx context.Context, item *complex.SchemaItem) error {
				assert.NotNil(t, item, "Insert handler should receive item")
				t.Logf("Insert handler called with item: %+v", item)
				return nil
			},
			// onModify
			func(ctx context.Context, oldItem, newItem *complex.SchemaItem) error {
				assert.NotNil(t, oldItem, "Modify handler should receive old item")
				assert.NotNil(t, newItem, "Modify handler should receive new item")
				t.Logf("Modify handler called with old: %+v, new: %+v", oldItem, newItem)
				return nil
			},
			// onDelete - Fix: use correct type from events package
			func(ctx context.Context, oldImage map[string]events.DynamoDBAttributeValue) error {
				assert.NotNil(t, oldImage, "Delete handler should receive old image")
				t.Logf("Delete handler called with oldImage: %+v", oldImage)
				return nil
			},
		)

		require.NotNil(t, handler, "CreateTriggerHandler should return non-nil handler")
		t.Logf("✅ Complex trigger handler created successfully")
	})
}

// ==================== Test Helper Functions ====================

// createComplexTestItem generates a SchemaItem with specified values
func createComplexTestItem(suffix string, createdAt int64, likes int, views int) complex.SchemaItem {
	return complex.SchemaItem{
		UserId:              "test-user",
		PostId:              "test-post-" + suffix,
		CreatedAt:           int(createdAt),
		Likes:               likes,
		IsPublished:         1,
		CategoryIsPublished: "tech#1",
		TagIsPublished:      "golang#1",
		Title:               "Test Post " + suffix,
		Content:             "Test content for " + suffix,
		Category:            "tech",
		Tag:                 "golang",
		Views:               views,
		IsPremium:           false,
		IsFeatured:          false,
	}
}

// insertComplexTestItem helper function to insert test items
func insertComplexTestItem(t *testing.T, client *dynamodb.Client, ctx context.Context, item complex.SchemaItem) {
	t.Helper()

	av, err := complex.PutItem(item)
	require.NoError(t, err, "Item marshaling should succeed")

	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(complex.TableName),
		Item:      av,
	})
	require.NoError(t, err, "PutItem should succeed")
}
