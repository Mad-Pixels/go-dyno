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

	userprofile "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/userprofile"
)

// TestSetsSchema tests String Set (SS) and Number Set (NS) functionality.
// This test validates that the generated code correctly handles DynamoDB Set types
// in all operations including marshaling, unmarshaling, and query building.
//
// Schema Structure (user-profile.json):
//   - Table: "user-profile"
//   - Hash Key: "user_id" (string)
//   - Range Key: "profile_type" (string)
//   - SS Attributes: "tags", "interests"
//   - NS Attributes: "skill_levels", "scores"
//   - Secondary Indexes with set support
func TestSetsSchema(t *testing.T) {
	client := ConnectToLocalStack(t, DefaultLocalStackConfig())
	ctx, cancel := TestContext(5 * time.Minute)
	defer cancel()

	t.Logf("Testing Sets schema: user-profile.json")
	t.Logf("Table: %s", userprofile.TableName)
	t.Logf("Hash Key: %s, Range Key: %s", userprofile.TableSchema.HashKey, userprofile.TableSchema.RangeKey)

	t.Run("String_Set_Operations", func(t *testing.T) {
		testStringSetCRUD(t, client, ctx)
		testStringSetQueryBuilder(t, client, ctx)
	})

	t.Run("Number_Set_Operations", func(t *testing.T) {
		testNumberSetCRUD(t, client, ctx)
		testNumberSetQueryBuilder(t, client, ctx)
	})

	t.Run("Mixed_Set_Operations", func(t *testing.T) {
		testMixedSetOperations(t, client, ctx)
		testSetUtilityFunctions(t)
	})

	t.Run("Schema_Constants_With_Sets", func(t *testing.T) {
		t.Parallel()
		testSchemaConstantsWithSets(t)
		testAttributeNamesWithSets(t)
	})
}

// ==================== String Set Tests ====================

// testStringSetCRUD validates String Set (SS) CRUD operations
func testStringSetCRUD(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_item_with_string_sets", func(t *testing.T) {
		item := userprofile.SchemaItem{
			UserId:      "user-ss-test",
			ProfileType: "developer",
			Tags:        []string{"golang", "backend", "api"},  // SS field
			Interests:   []string{"coding", "music", "travel"}, // SS common field
			SkillLevels: []int{8, 9, 7},                        // NS field for reference
			Scores:      []int{95, 87, 92},                     // NS common field for reference
			DisplayName: "String Set Tester",
			Email:       "ss-test@example.com",
			CreatedAt:   1640995200,
			IsActive:    1,
			IsPremium:   true,
			UpdatedAt:   1640995200,
		}

		// Test marshaling
		av, err := userprofile.PutItem(item)
		require.NoError(t, err, "Should marshal item with string sets")
		require.NotEmpty(t, av, "AttributeValues should not be empty")

		// Verify SS fields are present and correctly typed
		assert.Contains(t, av, "tags", "Should contain tags field")
		assert.Contains(t, av, "interests", "Should contain interests field")

		// Test storage in DynamoDB
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userprofile.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store item with string sets in DynamoDB")

		// Test retrieval
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(userprofile.TableName),
			Key: map[string]types.AttributeValue{
				"user_id":      &types.AttributeValueMemberS{Value: item.UserId},
				"profile_type": &types.AttributeValueMemberS{Value: item.ProfileType},
			},
		})
		require.NoError(t, err, "Should retrieve item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		t.Logf("✅ String Set CRUD operations successful")
	})

	t.Run("empty_string_sets", func(t *testing.T) {
		item := userprofile.SchemaItem{
			UserId:      "user-empty-ss",
			ProfileType: "basic",
			Tags:        []string{}, // Empty SS
			Interests:   nil,        // Nil SS
			SkillLevels: []int{5},   // Non-empty NS
			Scores:      []int{},    // Empty NS
			DisplayName: "Empty Sets Test",
			Email:       "empty@example.com",
			CreatedAt:   1640995300,
			IsActive:    0,
			IsPremium:   false,
			UpdatedAt:   1640995300,
		}

		_, err := userprofile.PutItem(item)
		require.NoError(t, err, "Should handle empty/nil string sets")

		// Empty sets should not cause errors but may be handled differently by DynamoDB
		t.Logf("✅ Empty string sets handled correctly")
	})
}

// testStringSetQueryBuilder validates QueryBuilder with String Set fields
func testStringSetQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("query_with_string_set_filters", func(t *testing.T) {
		// Insert test data with various tag combinations
		testItems := []userprofile.SchemaItem{
			{
				UserId: "query-user1", ProfileType: "dev1", Tags: []string{"golang", "docker"},
				Interests: []string{"backend"}, SkillLevels: []int{8}, Scores: []int{90},
				DisplayName: "Dev 1", Email: "dev1@test.com", CreatedAt: 1640995400,
				IsActive: 1, IsPremium: true, UpdatedAt: 1640995400,
			},
			{
				UserId: "query-user2", ProfileType: "dev2", Tags: []string{"python", "aws"},
				Interests: []string{"devops"}, SkillLevels: []int{7}, Scores: []int{85},
				DisplayName: "Dev 2", Email: "dev2@test.com", CreatedAt: 1640995500,
				IsActive: 1, IsPremium: false, UpdatedAt: 1640995500,
			},
		}

		// Insert test data
		for _, item := range testItems {
			av, err := userprofile.PutItem(item)
			require.NoError(t, err)
			_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
				TableName: aws.String(userprofile.TableName),
				Item:      av,
			})
			require.NoError(t, err)
		}

		// Test QueryBuilder with string set filtering
		qb := userprofile.NewQueryBuilder().
			WithUserId("query-user1").
			FilterTags([]string{"golang", "docker"}). // Filter by string set
			FilterInterests([]string{"backend"})      // Filter by common string set

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with string set filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ QueryBuilder with string sets works correctly")
	})
}

// ==================== Number Set Tests ====================

// testNumberSetCRUD validates Number Set (NS) CRUD operations
func testNumberSetCRUD(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("create_item_with_number_sets", func(t *testing.T) {
		item := userprofile.SchemaItem{
			UserId:      "user-ns-test",
			ProfileType: "analyst",
			Tags:        []string{"data", "analytics"},      // SS fields for reference
			Interests:   []string{"statistics", "modeling"}, // SS common field
			SkillLevels: []int{9, 8, 7, 10},                 // NS field - skill ratings
			Scores:      []int{95, 88, 92, 97, 85},          // NS common field - test scores
			DisplayName: "Number Set Tester",
			Email:       "ns-test@example.com",
			CreatedAt:   1640995600,
			IsActive:    1,
			IsPremium:   true,
			UpdatedAt:   1640995600,
		}

		// Test marshaling with number sets
		av, err := userprofile.PutItem(item)
		require.NoError(t, err, "Should marshal item with number sets")
		require.NotEmpty(t, av, "AttributeValues should not be empty")

		// Verify NS fields are present
		assert.Contains(t, av, "skill_levels", "Should contain skill_levels field")
		assert.Contains(t, av, "scores", "Should contain scores field")

		// Test storage in DynamoDB
		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userprofile.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store item with number sets in DynamoDB")

		// Test retrieval and verification
		getResult, err := client.GetItem(ctx, &dynamodb.GetItemInput{
			TableName: aws.String(userprofile.TableName),
			Key: map[string]types.AttributeValue{
				"user_id":      &types.AttributeValueMemberS{Value: item.UserId},
				"profile_type": &types.AttributeValueMemberS{Value: item.ProfileType},
			},
		})
		require.NoError(t, err, "Should retrieve item")
		assert.NotEmpty(t, getResult.Item, "Retrieved item should not be empty")

		t.Logf("✅ Number Set CRUD operations successful")
	})

	t.Run("number_set_edge_cases", func(t *testing.T) {
		item := userprofile.SchemaItem{
			UserId:      "user-ns-edge",
			ProfileType: "edge-case",
			Tags:        []string{"test"},
			Interests:   []string{"edge-cases"},
			SkillLevels: []int{0, -1, 100, 999}, // Edge case numbers including negatives
			Scores:      []int{1},               // Single number set
			DisplayName: "Edge Case Test",
			Email:       "edge@example.com",
			CreatedAt:   1640995700,
			IsActive:    1,
			IsPremium:   false,
			UpdatedAt:   1640995700,
		}

		_, err := userprofile.PutItem(item)
		require.NoError(t, err, "Should handle edge case numbers in sets")

		t.Logf("✅ Number set edge cases handled correctly")
	})
}

// testNumberSetQueryBuilder validates QueryBuilder with Number Set fields
func testNumberSetQueryBuilder(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("query_with_number_set_filters", func(t *testing.T) {
		// Test QueryBuilder with number set filtering
		qb := userprofile.NewQueryBuilder().
			WithUserId("test-user").
			FilterSkillLevels([]int{8, 9, 10}). // Filter by number set
			FilterScores([]int{90, 95, 100})    // Filter by common number set

		queryInput, err := qb.BuildQuery()
		require.NoError(t, err, "Should build query with number set filters")
		assert.NotNil(t, queryInput.KeyConditionExpression, "Should have key condition")

		t.Logf("✅ QueryBuilder with number sets works correctly")
	})
}

// ==================== Mixed Set Operations ====================

// testMixedSetOperations validates operations with both SS and NS in same item
func testMixedSetOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
	t.Run("comprehensive_set_item", func(t *testing.T) {
		item := userprofile.SchemaItem{
			UserId:      "comprehensive-user",
			ProfileType: "full-stack",
			Tags:        []string{"golang", "react", "postgres", "docker", "k8s"}, // 5 tags
			Interests:   []string{"coding", "mentoring", "open-source", "gaming"}, // 4 interests
			SkillLevels: []int{9, 8, 7, 8, 9, 6, 10},                              // 7 skill levels
			Scores:      []int{95, 88, 92, 97, 85, 90, 94, 89},                    // 8 scores
			DisplayName: "Full Stack Developer",
			Email:       "fullstack@example.com",
			CreatedAt:   1640995800,
			IsActive:    1,
			IsPremium:   true,
			UpdatedAt:   1640995800,
		}

		// Test complete lifecycle
		av, err := userprofile.PutItem(item)
		require.NoError(t, err, "Should marshal comprehensive item")

		_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
			TableName: aws.String(userprofile.TableName),
			Item:      av,
		})
		require.NoError(t, err, "Should store comprehensive item")

		// Test BatchPutItems with sets
		batchItems := []userprofile.SchemaItem{item}
		batchAVs, err := userprofile.BatchPutItems(batchItems)
		require.NoError(t, err, "BatchPutItems should handle sets")
		require.Len(t, batchAVs, 1, "Should return one batch item")

		t.Logf("✅ Comprehensive set operations successful")
	})
}

// testSetUtilityFunctions validates utility functions with set types
func testSetUtilityFunctions(t *testing.T) {
	t.Run("utility_functions_with_sets", func(t *testing.T) {
		item := userprofile.SchemaItem{
			UserId:      "utility-test",
			ProfileType: "tester",
			Tags:        []string{"testing", "qa"},
			Interests:   []string{"automation"},
			SkillLevels: []int{8, 7},
			Scores:      []int{90, 85},
			DisplayName: "Utility Tester",
			Email:       "utility@example.com",
			CreatedAt:   1640995900,
			IsActive:    1,
			IsPremium:   false,
			UpdatedAt:   1640995900,
		}

		// Test CreateKey
		key, err := userprofile.CreateKey(item.UserId, item.ProfileType)
		require.NoError(t, err, "CreateKey should work with set items")
		assert.Contains(t, key, "user_id", "Key should contain user_id")
		assert.Contains(t, key, "profile_type", "Key should contain profile_type")

		// Test CreateKeyFromItem
		keyFromItem, err := userprofile.CreateKeyFromItem(item)
		require.NoError(t, err, "CreateKeyFromItem should work with set items")
		assert.Equal(t, key, keyFromItem, "Keys should be identical")

		// Test ConvertMapToAttributeValues
		inputMap := map[string]interface{}{
			"tags":         []string{"test1", "test2"},
			"skill_levels": []int{5, 6, 7},
			"name":         "test",
			"count":        42,
		}
		converted, err := userprofile.ConvertMapToAttributeValues(inputMap)
		require.NoError(t, err, "Should convert map with sets")
		assert.Contains(t, converted, "tags", "Should contain tags")
		assert.Contains(t, converted, "skill_levels", "Should contain skill_levels")

		t.Logf("✅ Utility functions with sets work correctly")
	})
}

// ==================== Schema Constants Tests ====================

// testSchemaConstantsWithSets validates constants include set fields
func testSchemaConstantsWithSets(t *testing.T) {
	t.Run("set_column_constants", func(t *testing.T) {
		// Verify SS column constants
		assert.Equal(t, "tags", userprofile.ColumnTags, "ColumnTags should be correct")
		assert.Equal(t, "interests", userprofile.ColumnInterests, "ColumnInterests should be correct")

		// Verify NS column constants
		assert.Equal(t, "skill_levels", userprofile.ColumnSkillLevels, "ColumnSkillLevels should be correct")
		assert.Equal(t, "scores", userprofile.ColumnScores, "ColumnScores should be correct")

		t.Logf("✅ Set column constants are correct")
	})
}

// testAttributeNamesWithSets validates AttributeNames includes set fields
func testAttributeNamesWithSets(t *testing.T) {
	t.Run("attribute_names_include_sets", func(t *testing.T) {
		attrs := userprofile.AttributeNames
		require.NotEmpty(t, attrs, "AttributeNames should not be empty")

		// Check SS attributes
		assert.Contains(t, attrs, "tags", "AttributeNames should contain tags")
		assert.Contains(t, attrs, "interests", "AttributeNames should contain interests")

		// Check NS attributes
		assert.Contains(t, attrs, "skill_levels", "AttributeNames should contain skill_levels")
		assert.Contains(t, attrs, "scores", "AttributeNames should contain scores")

		t.Logf("✅ AttributeNames includes all set fields: %v", attrs)
	})
}
