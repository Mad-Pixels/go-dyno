package localstack

// import (
// 	"context"
// 	"math/big"
// 	"testing"
// 	"time"

// 	"github.com/aws/aws-sdk-go-v2/aws"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
// 	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
// 	"github.com/google/uuid"
// 	"github.com/shopspring/decimal"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	numbers "github.com/Mad-Pixels/go-dyno/tests/localstack/generated/numberstest"
// )

// // TestNumbersSchema runs comprehensive integration tests for the numbers.json schema.
// // This test suite validates AttributeSubtype functionality against a real LocalStack DynamoDB instance.
// //
// // Test Coverage:
// // - All supported numeric subtypes (int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64)
// // - Floating point subtypes (float32, float64)
// // - Arbitrary precision types (*big.Int, *decimal.Decimal)
// // - Special types (uuid.UUID, time.Time, []byte)
// // - Default types and fallback behavior
// // - CRUD operations with subtype marshaling/unmarshaling
// // - QueryBuilder with subtype parameters
// // - Range conditions on numeric fields
// // - Schema constants and metadata validation
// //
// // Schema Structure (numbers.json):
// //   - Table: "numbers-test"
// //   - Hash Key: "id" (string)
// //   - Range Key: "timestamp" (uint64)
// //   - Attributes: 10 different numeric subtypes
// //   - Common Attributes: 11 different subtypes including special types
// //   - Secondary Indexes: 3 GSIs with different projection types
// //
// // Example Usage:
// //
// //	go test -v ./tests/localstack/ -run TestNumbersSchema
// func TestNumbersSchema(t *testing.T) {
// 	cfg := DefaultLocalStackConfig()
// 	client := ConnectToLocalStack(t, cfg)

// 	ctx, cancel := TestContext(10 * time.Minute)
// 	defer cancel()

// 	t.Logf("Testing schema: numbers.json")
// 	t.Logf("Table: %s", numbers.TableName)
// 	t.Logf("Hash Key: %s, Range Key: %s", numbers.TableSchema.HashKey, numbers.TableSchema.RangeKey)
// 	t.Logf("Total Attributes: %d", len(numbers.TableSchema.AllAttributes))

// 	t.Run("Utility_Functions", func(t *testing.T) {
// 		t.Parallel()
// 		testNumbersBoolToInt(t)
// 		testNumbersIntToBool(t)
// 	})

// 	t.Run("CRUD_Operations", func(t *testing.T) {
// 		testNumbersPutItem(t, client, ctx)
// 		testNumbersBatchPutItems(t, client, ctx)
// 	})

// 	t.Run("Subtype_Validation", func(t *testing.T) {
// 		testNumbersSubtypeMarshaling(t, client, ctx)
// 		testNumbersSpecialTypes(t, client, ctx)
// 		testNumbersArbitraryPrecision(t, client, ctx)
// 	})

// 	t.Run("QueryBuilder_Subtypes", func(t *testing.T) {
// 		testNumbersQueryBuilderSubtypes(t, client, ctx)
// 		testNumbersRangeConditionsSubtypes(t, client, ctx)
// 		testNumbersFluentAPISubtypes(t)
// 	})

// 	t.Run("Numeric_Edge_Cases", func(t *testing.T) {
// 		testNumbersEdgeCases(t, client, ctx)
// 		testNumbersOverflowHandling(t, client, ctx)
// 		testNumbersZeroValues(t, client, ctx)
// 	})

// 	t.Run("Schema_Constants", func(t *testing.T) {
// 		t.Parallel()
// 		testNumbersSchemaConstants(t)
// 		testNumbersAttributeNames(t)
// 		testNumbersTableSchema(t)
// 	})

// 	t.Run("Key_Operations", func(t *testing.T) {
// 		testNumbersCreateKey(t)
// 		testNumbersCreateKeyFromItem(t)
// 	})

// 	t.Run("Secondary_Indexes", func(t *testing.T) {
// 		testNumbersGSIOperations(t, client, ctx)
// 		testNumbersIndexProjections(t, client, ctx)
// 	})
// }

// // ==================== Utility Functions Tests ====================

// // testNumbersBoolToInt validates the BoolToInt utility function
// func testNumbersBoolToInt(t *testing.T) {
// 	t.Run("boolean_to_numeric_conversion", func(t *testing.T) {
// 		testCases := []struct {
// 			input    bool
// 			expected int
// 			desc     string
// 		}{
// 			{true, 1, "active status should convert to 1"},
// 			{false, 0, "inactive status should convert to 0"},
// 		}

// 		for _, tc := range testCases {
// 			result := numbers.BoolToInt(tc.input)
// 			assert.Equal(t, tc.expected, result, tc.desc)
// 		}

// 		t.Logf("✅ BoolToInt utility function works correctly for numbers schema")
// 	})
// }

// // testNumbersIntToBool validates the IntToBool utility function
// func testNumbersIntToBool(t *testing.T) {
// 	t.Run("numeric_to_boolean_conversion", func(t *testing.T) {
// 		testCases := []struct {
// 			input    int
// 			expected bool
// 			desc     string
// 		}{
// 			{1, true, "active status (1) should be true"},
// 			{0, false, "inactive status (0) should be false"},
// 			{42, true, "any positive number should be true"},
// 			{-1, true, "negative number should be true"},
// 		}

// 		for _, tc := range testCases {
// 			result := numbers.IntToBool(tc.input)
// 			assert.Equal(t, tc.expected, result, tc.desc)
// 		}

// 		t.Logf("✅ IntToBool utility function handles all numeric cases correctly")
// 	})
// }

// // ==================== CRUD Operations Tests ====================

// // testNumbersPutItem validates item creation with all numeric subtypes
// func testNumbersPutItem(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("create_item_with_all_subtypes", func(t *testing.T) {
// 		// Generate test UUID and time
// 		testUUID := uuid.New()
// 		testTime := time.Now().UTC()
// 		testBigInt := big.NewInt(999999999999999999)
// 		testDecimal := decimal.NewFromFloat(123.456789)
// 		testData := []byte("test data for bytes")

// 		item := numbers.SchemaItem{
// 			// Primary keys
// 			Id:        "numbers-test-001",
// 			Timestamp: 1640995200,

// 			// Index keys (moved from common_attributes)
// 			IsActive:   1, // Changed to int (was bool)
// 			StatusCode: 200,
// 			Count:      1000,
// 			CreatedAt:  testTime,

// 			// Common attributes with all different subtypes
// 			TinyNumber:     -128,                 // int8 min
// 			SmallNumber:    -32768,               // int16 min
// 			MediumNumber:   -2147483648,          // int32 min
// 			BigNumber:      -9223372036854775808, // int64 min
// 			UnsignedTiny:   255,                  // uint8 max
// 			UnsignedSmall:  65535,                // uint16 max
// 			UnsignedMedium: 4294967295,           // uint32 max

// 			// Floating point subtypes
// 			DiscountRate: 0.15,    // float32
// 			Score:        98.7654, // float64

// 			// Arbitrary precision
// 			PriceCents: testBigInt,
// 			Balance:    &testDecimal,

// 			// Special types
// 			RequestId: testUUID,
// 			Data:      testData,

// 			// Default types
// 			DefaultNumber: 42.0, // Should be float64
// 			Description:   "Test item with all subtypes",
// 		}

// 		av, err := numbers.PutItem(item)
// 		require.NoError(t, err, "PutItem marshaling should succeed")
// 		require.NotEmpty(t, av, "AttributeValues should not be empty")

// 		// Validate all fields are present
// 		expectedFields := []string{
// 			"id", "timestamp", "is_active", "status_code", "count", "created_at",
// 			"tiny_number", "small_number", "medium_number", "big_number",
// 			"unsigned_tiny", "unsigned_small", "unsigned_medium",
// 			"price_cents", "discount_rate", "score", "balance",
// 			"request_id", "data", "default_number", "description",
// 		}

// 		for _, field := range expectedFields {
// 			assert.Contains(t, av, field, "Must contain field '%s'", field)
// 		}

// 		t.Logf("✅ Successfully created and marshaled item with all subtypes")
// 	})

// 	t.Run("put_item_to_dynamodb", func(t *testing.T) {
// 		testUUID := uuid.New()
// 		testTime := time.Now().UTC()
// 		testBigInt := big.NewInt(123456789)
// 		testDecimal := decimal.NewFromFloat(99.99)

// 		item := numbers.SchemaItem{
// 			Id:             "put-test-002",
// 			Timestamp:      1640995300,
// 			IsActive:       1, // int instead of bool
// 			StatusCode:     201,
// 			Count:          500,
// 			CreatedAt:      testTime,
// 			TinyNumber:     42,
// 			SmallNumber:    1000,
// 			MediumNumber:   100000,
// 			BigNumber:      9999999999,
// 			UnsignedTiny:   200,
// 			UnsignedSmall:  50000,
// 			UnsignedMedium: 3000000000,
// 			PriceCents:     testBigInt,
// 			DiscountRate:   0.10,
// 			Score:          85.5,
// 			Balance:        &testDecimal,
// 			RequestId:      testUUID,
// 			Data:           []byte("put test data"),
// 			DefaultNumber:  77.7,
// 			Description:    "Put test item",
// 		}

// 		av, err := numbers.PutItem(item)
// 		require.NoError(t, err, "Item marshaling should succeed")

// 		input := &dynamodb.PutItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Item:      av,
// 		}

// 		_, err = client.PutItem(ctx, input)
// 		require.NoError(t, err, "DynamoDB PutItem operation should succeed")

// 		// Verify item was stored
// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key: map[string]types.AttributeValue{
// 				"id":        &types.AttributeValueMemberS{Value: item.Id},
// 				"timestamp": &types.AttributeValueMemberN{Value: "1640995300"},
// 			},
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err, "GetItem verification should succeed")
// 		assert.NotEmpty(t, result.Item, "Stored item should be retrievable")

// 		t.Logf("✅ Successfully stored and verified item with all subtypes in DynamoDB")
// 	})
// }

// // testNumbersBatchPutItems validates batch operations with numeric subtypes
// func testNumbersBatchPutItems(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("batch_put_multiple_numeric_items", func(t *testing.T) {
// 		testTime := time.Now().UTC()
// 		items := []numbers.SchemaItem{
// 			{
// 				Id:             "batch-1",
// 				Timestamp:      1640995400,
// 				IsActive:       1,
// 				StatusCode:     100,
// 				Count:          100,
// 				CreatedAt:      testTime,
// 				TinyNumber:     10,
// 				SmallNumber:    1000,
// 				MediumNumber:   100000,
// 				BigNumber:      1000000000,
// 				UnsignedTiny:   50,
// 				UnsignedSmall:  5000,
// 				UnsignedMedium: 500000,
// 				PriceCents:     big.NewInt(9999),
// 				DiscountRate:   0.05,
// 				Score:          95.5,
// 				Balance:        func() *decimal.Decimal { d := decimal.NewFromFloat(1000.50); return &d }(),
// 				RequestId:      uuid.New(),
// 				Data:           []byte("batch data 1"),
// 				DefaultNumber:  10.1,
// 				Description:    "Batch item 1",
// 			},
// 			{
// 				Id:             "batch-2",
// 				Timestamp:      1640995500,
// 				IsActive:       0,
// 				StatusCode:     200,
// 				Count:          200,
// 				CreatedAt:      testTime.Add(time.Minute),
// 				TinyNumber:     -50,
// 				SmallNumber:    -5000,
// 				MediumNumber:   -500000,
// 				BigNumber:      -5000000000,
// 				UnsignedTiny:   100,
// 				UnsignedSmall:  10000,
// 				UnsignedMedium: 1000000,
// 				PriceCents:     big.NewInt(19999),
// 				DiscountRate:   0.15,
// 				Score:          87.3,
// 				Balance:        func() *decimal.Decimal { d := decimal.NewFromFloat(2000.75); return &d }(),
// 				RequestId:      uuid.New(),
// 				Data:           []byte("batch data 2"),
// 				DefaultNumber:  20.2,
// 				Description:    "Batch item 2",
// 			},
// 		}

// 		batchItems, err := numbers.BatchPutItems(items)
// 		require.NoError(t, err, "BatchPutItems should succeed")
// 		require.Len(t, batchItems, 2, "Should return AttributeValues for both items")

// 		// Validate subtype fields are properly marshaled
// 		for i, batchItem := range batchItems {
// 			subtypeFields := []string{
// 				"tiny_number", "small_number", "medium_number", "big_number",
// 				"unsigned_tiny", "unsigned_small", "unsigned_medium", "count",
// 				"price_cents", "discount_rate", "score", "balance",
// 				"request_id", "created_at", "data", "is_active", "status_code",
// 			}

// 			for _, field := range subtypeFields {
// 				assert.Contains(t, batchItem, field, "Batch item %d should have subtype field: %s", i, field)
// 			}
// 		}

// 		t.Logf("✅ Successfully prepared %d items with numeric subtypes for batch operation", len(batchItems))
// 	})
// }

// // ==================== Subtype Validation Tests ====================

// // testNumbersSubtypeMarshaling validates that different numeric subtypes marshal correctly
// func testNumbersSubtypeMarshaling(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("integer_subtypes_marshaling", func(t *testing.T) {
// 		item := numbers.SchemaItem{
// 			Id:             "marshal-test-int",
// 			Timestamp:      1640995600,
// 			TinyNumber:     127,                 // int8 max
// 			SmallNumber:    32767,               // int16 max
// 			MediumNumber:   2147483647,          // int32 max
// 			BigNumber:      9223372036854775807, // int64 max
// 			UnsignedTiny:   255,                 // uint8 max
// 			UnsignedSmall:  65535,               // uint16 max
// 			UnsignedMedium: 4294967295,          // uint32 max
// 			Count:          ^uint(0),            // uint max (platform dependent)
// 			StatusCode:     500,
// 			DefaultNumber:  123.456,
// 			Description:    "Integer marshaling test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		// Retrieve and validate
// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		// Validate that numeric values are stored correctly
// 		assert.Contains(t, result.Item, "tiny_number")
// 		assert.Contains(t, result.Item, "small_number")
// 		assert.Contains(t, result.Item, "medium_number")
// 		assert.Contains(t, result.Item, "big_number")

// 		t.Logf("✅ Integer subtypes marshaled and stored correctly")
// 	})

// 	t.Run("floating_point_subtypes_marshaling", func(t *testing.T) {
// 		item := numbers.SchemaItem{
// 			Id:            "marshal-test-float",
// 			Timestamp:     1640995700,
// 			DiscountRate:  0.123456789,      // float32 precision test
// 			Score:         123.456789012345, // float64 precision test
// 			StatusCode:    200,
// 			DefaultNumber: 999.999,
// 			Description:   "Float marshaling test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		// Validate floating point values
// 		assert.Contains(t, result.Item, "discount_rate")
// 		assert.Contains(t, result.Item, "score")

// 		t.Logf("✅ Floating point subtypes marshaled and stored correctly")
// 	})
// }

// // testNumbersSpecialTypes validates special subtypes (UUID, Time, Bytes)
// func testNumbersSpecialTypes(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("uuid_subtype_handling", func(t *testing.T) {
// 		testUUID := uuid.New()

// 		item := numbers.SchemaItem{
// 			Id:            "special-uuid-test",
// 			Timestamp:     1640995800,
// 			RequestId:     testUUID,
// 			StatusCode:    200,
// 			DefaultNumber: 1.0,
// 			Description:   "UUID test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		// Validate UUID is handled correctly
// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		assert.Contains(t, result.Item, "request_id")

// 		t.Logf("✅ UUID subtype handled correctly: %s", testUUID.String())
// 	})

// 	t.Run("time_subtype_handling", func(t *testing.T) {
// 		testTime := time.Date(2024, 12, 25, 15, 30, 45, 0, time.UTC)

// 		item := numbers.SchemaItem{
// 			Id:            "special-time-test",
// 			Timestamp:     1640995900,
// 			CreatedAt:     testTime,
// 			StatusCode:    200,
// 			DefaultNumber: 2.0,
// 			Description:   "Time test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		assert.Contains(t, result.Item, "created_at")

// 		t.Logf("✅ Time subtype handled correctly: %s", testTime.Format(time.RFC3339))
// 	})

// 	t.Run("bytes_subtype_handling", func(t *testing.T) {
// 		testData := []byte("Hello, this is test binary data! 🎉")

// 		item := numbers.SchemaItem{
// 			Id:            "special-bytes-test",
// 			Timestamp:     1640996000,
// 			Data:          testData,
// 			StatusCode:    200,
// 			DefaultNumber: 3.0,
// 			Description:   "Bytes test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		assert.Contains(t, result.Item, "data")

// 		t.Logf("✅ Bytes subtype handled correctly: %d bytes", len(testData))
// 	})
// }

// // testNumbersArbitraryPrecision validates *big.Int and *decimal.Decimal subtypes
// func testNumbersArbitraryPrecision(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("big_int_subtype_handling", func(t *testing.T) {
// 		// Very large number that exceeds int64
// 		hugeBigInt := new(big.Int)
// 		hugeBigInt.SetString("123456789012345678901234567890", 10)

// 		item := numbers.SchemaItem{
// 			Id:            "precision-bigint-test",
// 			Timestamp:     1640996100,
// 			PriceCents:    hugeBigInt,
// 			StatusCode:    200,
// 			DefaultNumber: 4.0,
// 			Description:   "Big int test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		assert.Contains(t, result.Item, "price_cents")

// 		t.Logf("✅ Big int subtype handled correctly: %s", hugeBigInt.String())
// 	})

// 	t.Run("decimal_subtype_handling", func(t *testing.T) {
// 		// High precision decimal
// 		preciseDecimal, _ := decimal.NewFromString("123456.789012345678901234567890")

// 		item := numbers.SchemaItem{
// 			Id:            "precision-decimal-test",
// 			Timestamp:     1640996200,
// 			Balance:       &preciseDecimal,
// 			StatusCode:    200,
// 			DefaultNumber: 5.0,
// 			Description:   "Decimal test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		require.NotEmpty(t, result.Item)

// 		assert.Contains(t, result.Item, "balance")

// 		t.Logf("✅ Decimal subtype handled correctly: %s", preciseDecimal.String())
// 	})
// }

// // ==================== QueryBuilder Subtype Tests ====================

// // testNumbersQueryBuilderSubtypes validates QueryBuilder with subtype parameters
// func testNumbersQueryBuilderSubtypes(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("query_builder_subtype_parameters", func(t *testing.T) {
// 		// Insert test data with various subtypes
// 		testItems := []numbers.SchemaItem{
// 			createNumbersTestItem("qb-subtype-1", 1640996300, 100, 200),
// 			createNumbersTestItem("qb-subtype-2", 1640996400, 150, 250),
// 			createNumbersTestItem("qb-subtype-3", 1640996500, 200, 300),
// 		}

// 		for _, item := range testItems {
// 			insertNumbersTestItem(t, client, ctx, item)
// 		}

// 		// Test QueryBuilder with different subtype parameters
// 		qb := numbers.NewQueryBuilder()

// 		// Test various With methods for different subtypes
// 		subtypeTests := []struct {
// 			name   string
// 			method func() *numbers.QueryBuilder
// 		}{
// 			{"WithId", func() *numbers.QueryBuilder { return qb.WithId("qb-subtype-1") }},
// 			{"WithTimestamp", func() *numbers.QueryBuilder { return qb.WithTimestamp(1640996300) }},       // uint64
// 			{"WithTinyNumber", func() *numbers.QueryBuilder { return qb.WithTinyNumber(100) }},            // int8
// 			{"WithSmallNumber", func() *numbers.QueryBuilder { return qb.WithSmallNumber(1000) }},         // int16
// 			{"WithMediumNumber", func() *numbers.QueryBuilder { return qb.WithMediumNumber(100000) }},     // int32
// 			{"WithBigNumber", func() *numbers.QueryBuilder { return qb.WithBigNumber(1000000000) }},       // int64
// 			{"WithUnsignedTiny", func() *numbers.QueryBuilder { return qb.WithUnsignedTiny(50) }},         // uint8
// 			{"WithUnsignedSmall", func() *numbers.QueryBuilder { return qb.WithUnsignedSmall(5000) }},     // uint16
// 			{"WithUnsignedMedium", func() *numbers.QueryBuilder { return qb.WithUnsignedMedium(500000) }}, // uint32
// 			{"WithCount", func() *numbers.QueryBuilder { return qb.WithCount(200) }},                      // uint
// 			{"WithDiscountRate", func() *numbers.QueryBuilder { return qb.WithDiscountRate(0.15) }},       // float32
// 			{"WithScore", func() *numbers.QueryBuilder { return qb.WithScore(95.5) }},                     // float64
// 			{"WithIsActive", func() *numbers.QueryBuilder { return qb.WithIsActive(1) }},                  // int instead of bool
// 			{"WithStatusCode", func() *numbers.QueryBuilder { return qb.WithStatusCode(200) }},            // int
// 		}

// 		for _, test := range subtypeTests {
// 			result := test.method()
// 			require.NotNil(t, result, "%s should return non-nil QueryBuilder", test.name)
// 			assert.IsType(t, qb, result, "%s should return correct type", test.name)
// 		}

// 		t.Logf("✅ All %d QueryBuilder subtype methods work correctly", len(subtypeTests))
// 	})
// }

// // testNumbersRangeConditionsSubtypes validates range conditions with subtypes
// func testNumbersRangeConditionsSubtypes(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("numeric_range_conditions_with_subtypes", func(t *testing.T) {
// 		// Insert test data with range of values
// 		baseTime := uint64(1640996600)
// 		testItems := []numbers.SchemaItem{
// 			createNumbersTestItem("range-1", baseTime+100, 10, 100),
// 			createNumbersTestItem("range-2", baseTime+200, 20, 200),
// 			createNumbersTestItem("range-3", baseTime+300, 30, 300),
// 			createNumbersTestItem("range-4", baseTime+400, 40, 400),
// 			createNumbersTestItem("range-5", baseTime+500, 50, 500),
// 		}

// 		for _, item := range testItems {
// 			insertNumbersTestItem(t, client, ctx, item)
// 		}

// 		// Test range conditions with different numeric subtypes
// 		rangeTests := []struct {
// 			name string
// 			qb   *numbers.QueryBuilder
// 			desc string
// 		}{
// 			{
// 				"timestamp_between_uint64",
// 				numbers.NewQueryBuilder().WithId("range-user").WithTimestampBetween(baseTime+150, baseTime+350),
// 				"uint64 range condition on timestamp",
// 			},
// 			{
// 				"tiny_number_greater_than_int8",
// 				numbers.NewQueryBuilder().WithId("range-user").WithTinyNumberGreaterThan(25),
// 				"int8 greater than condition",
// 			},
// 			{
// 				"small_number_less_than_int16",
// 				numbers.NewQueryBuilder().WithId("range-user").WithSmallNumberLessThan(35),
// 				"int16 less than condition",
// 			},
// 			{
// 				"medium_number_between_int32",
// 				numbers.NewQueryBuilder().WithId("range-user").WithMediumNumberBetween(15, 35),
// 				"int32 between condition",
// 			},
// 			{
// 				"big_number_greater_than_int64",
// 				numbers.NewQueryBuilder().WithId("range-user").WithBigNumberGreaterThan(25),
// 				"int64 greater than condition",
// 			},
// 			{
// 				"unsigned_tiny_less_than_uint8",
// 				numbers.NewQueryBuilder().WithId("range-user").WithUnsignedTinyLessThan(45),
// 				"uint8 less than condition",
// 			},
// 			{
// 				"unsigned_small_between_uint16",
// 				numbers.NewQueryBuilder().WithId("range-user").WithUnsignedSmallBetween(150, 350),
// 				"uint16 between condition",
// 			},
// 			{
// 				"unsigned_medium_greater_than_uint32",
// 				numbers.NewQueryBuilder().WithId("range-user").WithUnsignedMediumGreaterThan(250),
// 				"uint32 greater than condition",
// 			},
// 			{
// 				"count_less_than_uint",
// 				numbers.NewQueryBuilder().WithId("range-user").WithCountLessThan(450),
// 				"uint less than condition",
// 			},
// 			{
// 				"discount_rate_between_float32",
// 				numbers.NewQueryBuilder().WithId("range-user").WithDiscountRateBetween(0.05, 0.25),
// 				"float32 between condition",
// 			},
// 			{
// 				"score_greater_than_float64",
// 				numbers.NewQueryBuilder().WithId("range-user").WithScoreGreaterThan(80.0),
// 				"float64 greater than condition",
// 			},
// 			{
// 				"status_code_between_int",
// 				numbers.NewQueryBuilder().WithId("range-user").WithStatusCodeBetween(150, 350),
// 				"int between condition",
// 			},
// 		}

// 		for _, test := range rangeTests {
// 			queryInput, err := test.qb.BuildQuery()
// 			require.NoError(t, err, "Range test %s should build successfully: %s", test.name, test.desc)

// 			// Validate that query input is properly constructed
// 			assert.NotNil(t, queryInput.KeyConditionExpression, "Query %s should have key condition", test.name)

// 			t.Logf("Range test %s: %s - built successfully", test.name, test.desc)
// 		}

// 		t.Logf("✅ All %d range condition tests with subtypes passed", len(rangeTests))
// 	})
// }

// // testNumbersFluentAPISubtypes validates fluent API method chaining with subtypes
// func testNumbersFluentAPISubtypes(t *testing.T) {
// 	t.Run("complex_fluent_chaining_with_subtypes", func(t *testing.T) {
// 		// Test complex method chaining with various subtypes
// 		qb1 := numbers.NewQueryBuilder().
// 			WithId("fluent-test").
// 			WithTimestamp(1640996700). // uint64
// 			WithTinyNumber(42).        // int8
// 			WithBigNumber(9999999999). // int64
// 			WithDiscountRate(0.15).    // float32
// 			WithScore(95.5).           // float64
// 			WithIsActive(true).        // bool
// 			WithStatusCode(200).       // int
// 			OrderByDesc().
// 			Limit(10)

// 		require.NotNil(t, qb1, "Complex fluent chaining with subtypes should work")

// 		// Test method order independence
// 		qb2 := numbers.NewQueryBuilder().
// 			WithDiscountRate(0.10).
// 			WithId("order-test").
// 			WithIsActive(false).
// 			WithTimestamp(1640996800).
// 			WithScore(88.8).
// 			OrderByAsc().
// 			Limit(5)

// 		require.NotNil(t, qb2, "Method order independence should work with subtypes")

// 		// Test range conditions in fluent chain
// 		qb3 := numbers.NewQueryBuilder().
// 			WithId("range-fluent-test").
// 			WithTimestampBetween(1640996000, 1640997000).
// 			WithTinyNumberGreaterThan(10).
// 			WithScoreLessThan(100.0).
// 			WithCountBetween(100, 1000).
// 			OrderByDesc()

// 		require.NotNil(t, qb3, "Range conditions in fluent chain should work")

// 		t.Logf("✅ Complex fluent API chaining works correctly with all subtypes")
// 	})
// }

// // ==================== Edge Cases Tests ====================

// // testNumbersEdgeCases validates edge cases and boundary values
// func testNumbersEdgeCases(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("boundary_values_handling", func(t *testing.T) {
// 		// Test with boundary values for each subtype
// 		item := numbers.SchemaItem{
// 			Id:             "edge-case-boundaries",
// 			Timestamp:      ^uint64(0),              // uint64 max
// 			TinyNumber:     127,                     // int8 max
// 			SmallNumber:    32767,                   // int16 max
// 			MediumNumber:   2147483647,              // int32 max
// 			BigNumber:      9223372036854775807,     // int64 max
// 			UnsignedTiny:   255,                     // uint8 max
// 			UnsignedSmall:  65535,                   // uint16 max
// 			UnsignedMedium: 4294967295,              // uint32 max
// 			Count:          ^uint(0),                // uint max
// 			DiscountRate:   3.4028235e+38,           // float32 max (approximately)
// 			Score:          1.7976931348623157e+308, // float64 max (approximately)
// 			PriceCents:     new(big.Int).SetUint64(^uint64(0)),
// 			Balance:        func() *decimal.Decimal { d := decimal.NewFromFloat(999999999.999999999); return &d }(),
// 			IsActive:       true,
// 			StatusCode:     2147483647, // int max
// 			DefaultNumber:  1e308,
// 			Description:    "Boundary values test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		// Validate retrieval
// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		assert.NotEmpty(t, result.Item)

// 		t.Logf("✅ Boundary values handled correctly")
// 	})

// 	t.Run("negative_values_handling", func(t *testing.T) {
// 		// Test with negative values for signed types
// 		item := numbers.SchemaItem{
// 			Id:           "edge-case-negative",
// 			Timestamp:    1640996900,
// 			TinyNumber:   -128,                 // int8 min
// 			SmallNumber:  -32768,               // int16 min
// 			MediumNumber: -2147483648,          // int32 min
// 			BigNumber:    -9223372036854775808, // int64 min
// 			// Unsigned types can't be negative, use small positive values
// 			UnsignedTiny:   1,
// 			UnsignedSmall:  1,
// 			UnsignedMedium: 1,
// 			Count:          1,
// 			DiscountRate:   -0.5,   // Negative discount (surcharge)
// 			Score:          -100.5, // Negative score
// 			PriceCents:     big.NewInt(-999999),
// 			Balance:        func() *decimal.Decimal { d := decimal.NewFromFloat(-1000.50); return &d }(),
// 			IsActive:       false,
// 			StatusCode:     -1,
// 			DefaultNumber:  -999.999,
// 			Description:    "Negative values test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		assert.NotEmpty(t, result.Item)

// 		t.Logf("✅ Negative values handled correctly")
// 	})
// }

// // testNumbersOverflowHandling validates overflow handling for numeric types
// func testNumbersOverflowHandling(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("safe_overflow_prevention", func(t *testing.T) {
// 		// Test that our types handle expected ranges properly
// 		// Note: In real applications, overflow would be caught at compile time
// 		// or during value assignment, not during DynamoDB operations

// 		safeItem := numbers.SchemaItem{
// 			Id:             "overflow-safe",
// 			Timestamp:      1640997000,
// 			TinyNumber:     100,           // Safe int8 value
// 			SmallNumber:    20000,         // Safe int16 value
// 			MediumNumber:   1000000,       // Safe int32 value
// 			BigNumber:      1000000000000, // Safe int64 value
// 			UnsignedTiny:   200,           // Safe uint8 value
// 			UnsignedSmall:  50000,         // Safe uint16 value
// 			UnsignedMedium: 3000000000,    // Safe uint32 value
// 			Count:          1000000,       // Safe uint value
// 			DiscountRate:   0.99,          // Safe float32 value
// 			Score:          999.99,        // Safe float64 value
// 			PriceCents:     big.NewInt(999999999999),
// 			Balance:        func() *decimal.Decimal { d := decimal.NewFromFloat(99999.99); return &d }(),
// 			IsActive:       true,
// 			StatusCode:     200,
// 			DefaultNumber:  123.456,
// 			Description:    "Safe values test",
// 		}

// 		insertNumbersTestItem(t, client, ctx, safeItem)

// 		key, err := numbers.CreateKeyFromItem(safeItem)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		assert.NotEmpty(t, result.Item)

// 		t.Logf("✅ Safe value ranges handled correctly")
// 	})
// }

// // testNumbersZeroValues validates zero value handling
// func testNumbersZeroValues(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("zero_values_handling", func(t *testing.T) {
// 		// Test with zero values for all numeric types
// 		item := numbers.SchemaItem{
// 			Id:             "zero-values-test",
// 			Timestamp:      0,   // uint64 zero
// 			TinyNumber:     0,   // int8 zero
// 			SmallNumber:    0,   // int16 zero
// 			MediumNumber:   0,   // int32 zero
// 			BigNumber:      0,   // int64 zero
// 			UnsignedTiny:   0,   // uint8 zero
// 			UnsignedSmall:  0,   // uint16 zero
// 			UnsignedMedium: 0,   // uint32 zero
// 			Count:          0,   // uint zero
// 			DiscountRate:   0.0, // float32 zero
// 			Score:          0.0, // float64 zero
// 			PriceCents:     big.NewInt(0),
// 			Balance:        func() *decimal.Decimal { d := decimal.Zero; return &d }(),
// 			RequestId:      uuid.Nil,    // Zero UUID
// 			CreatedAt:      time.Time{}, // Zero time
// 			Data:           []byte{},    // Empty bytes
// 			IsActive:       false,       // bool zero (false)
// 			StatusCode:     0,           // int zero
// 			DefaultNumber:  0.0,         // float64 zero
// 			Description:    "",          // string zero (empty)
// 		}

// 		insertNumbersTestItem(t, client, ctx, item)

// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err)

// 		getInput := &dynamodb.GetItemInput{
// 			TableName: aws.String(numbers.TableName),
// 			Key:       key,
// 		}

// 		result, err := client.GetItem(ctx, getInput)
// 		require.NoError(t, err)
// 		assert.NotEmpty(t, result.Item)

// 		t.Logf("✅ Zero values handled correctly for all subtypes")
// 	})
// }

// // ==================== Schema Constants Tests ====================

// // testNumbersSchemaConstants validates generated constants
// func testNumbersSchemaConstants(t *testing.T) {
// 	t.Run("table_and_index_constants", func(t *testing.T) {
// 		// Validate table name
// 		assert.NotEmpty(t, numbers.TableName, "TableName constant should not be empty")
// 		assert.Equal(t, "numbers-test", numbers.TableName, "TableName should match numbers.json schema")

// 		// Validate index constants
// 		indexConstants := []struct {
// 			constant string
// 			expected string
// 		}{
// 			{numbers.IndexStatusIndex, "StatusIndex"},
// 			{numbers.IndexTypeCountIndex, "TypeCountIndex"},
// 			{numbers.IndexTimeRangeIndex, "TimeRangeIndex"},
// 		}

// 		for _, idx := range indexConstants {
// 			assert.Equal(t, idx.expected, idx.constant, "Index constant should match expected value")
// 		}

// 		t.Logf("✅ Table and index constants validated")
// 		t.Logf("    Table: %s", numbers.TableName)
// 		t.Logf("    Indexes: %d constants", len(indexConstants))
// 	})

// 	t.Run("column_constants_numbers", func(t *testing.T) {
// 		// Validate all column constants for numbers schema
// 		columnTests := map[string]string{
// 			numbers.ColumnId:             "id",
// 			numbers.ColumnTimestamp:      "timestamp",
// 			numbers.ColumnTinyNumber:     "tiny_number",
// 			numbers.ColumnSmallNumber:    "small_number",
// 			numbers.ColumnMediumNumber:   "medium_number",
// 			numbers.ColumnBigNumber:      "big_number",
// 			numbers.ColumnUnsignedTiny:   "unsigned_tiny",
// 			numbers.ColumnUnsignedSmall:  "unsigned_small",
// 			numbers.ColumnUnsignedMedium: "unsigned_medium",
// 			numbers.ColumnCount:          "count",
// 			numbers.ColumnPriceCents:     "price_cents",
// 			numbers.ColumnDiscountRate:   "discount_rate",
// 			numbers.ColumnScore:          "score",
// 			numbers.ColumnBalance:        "balance",
// 			numbers.ColumnRequestId:      "request_id",
// 			numbers.ColumnCreatedAt:      "created_at",
// 			numbers.ColumnData:           "data",
// 			numbers.ColumnIsActive:       "is_active",
// 			numbers.ColumnStatusCode:     "status_code",
// 			numbers.ColumnDefaultNumber:  "default_number",
// 			numbers.ColumnDescription:    "description",
// 		}

// 		for constant, expected := range columnTests {
// 			assert.Equal(t, expected, constant, "Column constant should match expected value")
// 		}

// 		t.Logf("✅ All %d column constants validated for numbers schema", len(columnTests))
// 	})
// }

// // testNumbersAttributeNames validates the AttributeNames array
// func testNumbersAttributeNames(t *testing.T) {
// 	t.Run("numbers_attribute_names_array", func(t *testing.T) {
// 		attrs := numbers.AttributeNames
// 		require.NotEmpty(t, attrs, "AttributeNames should not be empty")

// 		// Expected attributes from numbers.json
// 		expectedAttrs := []string{
// 			"id", "timestamp", "is_active", "status_code", "count", "created_at",
// 			"tiny_number", "small_number", "medium_number", "big_number",
// 			"unsigned_tiny", "unsigned_small", "unsigned_medium",
// 			"price_cents", "discount_rate", "score", "balance",
// 			"request_id", "data", "default_number", "description",
// 		}

// 		assert.Len(t, attrs, len(expectedAttrs), "AttributeNames should contain all schema attributes")

// 		for _, expected := range expectedAttrs {
// 			assert.Contains(t, attrs, expected, "AttributeNames should contain '%s'", expected)
// 		}

// 		// Ensure no duplicates
// 		attrSet := make(map[string]bool)
// 		for _, attr := range attrs {
// 			assert.False(t, attrSet[attr], "AttributeNames should not contain duplicate: %s", attr)
// 			attrSet[attr] = true
// 		}

// 		t.Logf("✅ AttributeNames array contains %d attributes", len(attrs))
// 	})
// }

// // testNumbersTableSchema validates the TableSchema variable
// func testNumbersTableSchema(t *testing.T) {
// 	t.Run("numbers_table_schema_structure", func(t *testing.T) {
// 		schema := numbers.TableSchema

// 		// Validate basic schema properties
// 		assert.Equal(t, "numbers-test", schema.TableName, "Schema TableName should match")
// 		assert.Equal(t, "id", schema.HashKey, "Hash key should be 'id'")
// 		assert.Equal(t, "timestamp", schema.RangeKey, "Range key should be 'timestamp'")

// 		// Validate attribute collections
// 		expectedPrimaryAttrs := 6     // From numbers.json attributes array (id, timestamp, is_active, status_code, count, created_at)
// 		expectedCommonAttrs := 15     // From numbers.json common_attributes array
// 		expectedSecondaryIndexes := 3 // From numbers.json secondary_indexes array

// 		assert.Len(t, schema.Attributes, expectedPrimaryAttrs, "Should have %d primary attributes", expectedPrimaryAttrs)
// 		assert.Len(t, schema.CommonAttributes, expectedCommonAttrs, "Should have %d common attributes", expectedCommonAttrs)
// 		assert.Len(t, schema.SecondaryIndexes, expectedSecondaryIndexes, "Should have %d secondary indexes", expectedSecondaryIndexes)

// 		// Validate secondary index structure
// 		for _, idx := range schema.SecondaryIndexes {
// 			assert.NotEmpty(t, idx.Name, "Index name should not be empty")
// 			assert.NotEmpty(t, idx.HashKey, "Index hash key should not be empty")
// 			assert.NotEmpty(t, idx.ProjectionType, "Index projection type should not be empty")

// 			t.Logf("    Index: %s (HashKey: %s, RangeKey: %s, Projection: %s)",
// 				idx.Name, idx.HashKey, idx.RangeKey, idx.ProjectionType)
// 		}

// 		t.Logf("✅ Numbers TableSchema structure is valid")
// 		t.Logf("    Table: %s", schema.TableName)
// 		t.Logf("    Keys: %s (hash), %s (range)", schema.HashKey, schema.RangeKey)
// 		t.Logf("    Attributes: %d primary, %d common", len(schema.Attributes), len(schema.CommonAttributes))
// 		t.Logf("    Indexes: %d secondary", len(schema.SecondaryIndexes))
// 	})
// }

// // ==================== Key Operations Tests ====================

// // testNumbersCreateKey validates key creation with numeric subtypes
// func testNumbersCreateKey(t *testing.T) {
// 	t.Run("create_key_with_numeric_subtypes", func(t *testing.T) {
// 		hashKeyValue := "numbers-user-456"
// 		rangeKeyValue := uint64(1640997200) // uint64 subtype

// 		key, err := numbers.CreateKey(hashKeyValue, rangeKeyValue)
// 		require.NoError(t, err, "CreateKey should succeed with subtype values")
// 		require.NotEmpty(t, key, "Created key should not be empty")

// 		// Validate key structure
// 		assert.Contains(t, key, "id", "Key should contain hash key 'id'")
// 		assert.Contains(t, key, "timestamp", "Key should contain range key 'timestamp'")

// 		t.Logf("✅ CreateKey with numeric subtypes works correctly")
// 	})

// 	t.Run("create_key_with_various_numeric_types", func(t *testing.T) {
// 		// Test with different numeric types for range key
// 		testCases := []struct {
// 			hashKey  any
// 			rangeKey any
// 			desc     string
// 		}{
// 			{"user-1", uint64(1640997300), "uint64 range key"},
// 			{"user-2", int64(1640997400), "int64 range key"},
// 			{"user-3", uint32(1640997500), "uint32 range key"},
// 		}

// 		for _, tc := range testCases {
// 			key, err := numbers.CreateKey(tc.hashKey, tc.rangeKey)
// 			require.NoError(t, err, "CreateKey should succeed for %s", tc.desc)
// 			assert.Contains(t, key, "id", "Key should contain id for %s", tc.desc)
// 			assert.Contains(t, key, "timestamp", "Key should contain timestamp for %s", tc.desc)
// 		}

// 		t.Logf("✅ CreateKey handles various numeric types correctly")
// 	})
// }

// // testNumbersCreateKeyFromItem validates key extraction from numbers SchemaItem
// func testNumbersCreateKeyFromItem(t *testing.T) {
// 	t.Run("extract_key_from_numbers_item", func(t *testing.T) {
// 		// Create complete numbers item with all subtypes
// 		testUUID := uuid.New()
// 		testTime := time.Now().UTC()
// 		testBigInt := big.NewInt(123456789)
// 		testDecimal := decimal.NewFromFloat(999.99)

// 		item := numbers.SchemaItem{
// 			Id:             "extract-numbers-test",
// 			Timestamp:      1640997600,
// 			TinyNumber:     42,
// 			SmallNumber:    1000,
// 			MediumNumber:   100000,
// 			BigNumber:      1000000000,
// 			UnsignedTiny:   100,
// 			UnsignedSmall:  10000,
// 			UnsignedMedium: 1000000,
// 			Count:          500,
// 			PriceCents:     testBigInt,
// 			DiscountRate:   0.15,
// 			Score:          95.5,
// 			Balance:        &testDecimal,
// 			RequestId:      testUUID,
// 			CreatedAt:      testTime,
// 			Data:           []byte("test data"),
// 			IsActive:       true,
// 			StatusCode:     200,
// 			DefaultNumber:  123.45,
// 			Description:    "Key extraction test",
// 		}

// 		// Extract key attributes only
// 		key, err := numbers.CreateKeyFromItem(item)
// 		require.NoError(t, err, "CreateKeyFromItem should succeed")
// 		require.NotEmpty(t, key, "Extracted key should not be empty")

// 		// Validate extracted key contains only primary key attributes
// 		assert.Contains(t, key, "id", "Extracted key should contain hash key 'id'")
// 		assert.Contains(t, key, "timestamp", "Extracted key should contain range key 'timestamp'")

// 		// Validate non-key attributes are excluded
// 		nonKeyAttrs := []string{
// 			"tiny_number", "small_number", "medium_number", "big_number",
// 			"unsigned_tiny", "unsigned_small", "unsigned_medium", "count",
// 			"price_cents", "discount_rate", "score", "balance",
// 			"request_id", "created_at", "data", "is_active", "status_code",
// 			"default_number", "description",
// 		}

// 		for _, attr := range nonKeyAttrs {
// 			assert.NotContains(t, key, attr, "Extracted key should not contain non-key attribute '%s'", attr)
// 		}

// 		t.Logf("✅ CreateKeyFromItem correctly extracts only primary key attributes")
// 	})
// }

// // ==================== Secondary Index Tests ====================

// // testNumbersGSIOperations validates GSI operations with numeric subtypes
// func testNumbersGSIOperations(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("gsi_operations_with_subtypes", func(t *testing.T) {
// 		// Insert test data for GSI operations
// 		testItems := []numbers.SchemaItem{
// 			{
// 				Id:            "gsi-test-1",
// 				Timestamp:     1640997700,
// 				StatusCode:    200,
// 				Count:         100,
// 				IsActive:      true,
// 				DefaultNumber: 1.0,
// 				Description:   "GSI test 1",
// 			},
// 			{
// 				Id:            "gsi-test-2",
// 				Timestamp:     1640997800,
// 				StatusCode:    201,
// 				Count:         200,
// 				IsActive:      false,
// 				DefaultNumber: 2.0,
// 				Description:   "GSI test 2",
// 			},
// 			{
// 				Id:            "gsi-test-3",
// 				Timestamp:     1640997900,
// 				StatusCode:    200,
// 				Count:         300,
// 				IsActive:      true,
// 				DefaultNumber: 3.0,
// 				Description:   "GSI test 3",
// 			},
// 		}

// 		for _, item := range testItems {
// 			insertNumbersTestItem(t, client, ctx, item)
// 		}

// 		// Test queries using different GSIs
// 		gsiTests := []struct {
// 			name string
// 			qb   *numbers.QueryBuilder
// 		}{
// 			{
// 				"status_index_query",
// 				numbers.NewQueryBuilder().WithIsActive(1).WithTimestampGreaterThan(1640997650),
// 			},
// 			{
// 				"type_count_index_query",
// 				numbers.NewQueryBuilder().WithStatusCode(200).WithCountGreaterThan(150),
// 			},
// 			{
// 				"time_range_index_query",
// 				numbers.NewQueryBuilder().WithIsActive(1),
// 			},
// 		}

// 		for _, test := range gsiTests {
// 			queryInput, err := test.qb.BuildQuery()
// 			require.NoError(t, err, "GSI test %s should succeed", test.name)

// 			if queryInput.IndexName != nil {
// 				t.Logf("GSI query %s uses index: %s", test.name, *queryInput.IndexName)
// 			}
// 		}

// 		t.Logf("✅ GSI operations with subtypes validated")
// 	})
// }

// // testNumbersIndexProjections validates index projection handling
// func testNumbersIndexProjections(t *testing.T, client *dynamodb.Client, ctx context.Context) {
// 	t.Run("index_projections_with_subtypes", func(t *testing.T) {
// 		projections := numbers.IndexProjections
// 		require.NotEmpty(t, projections, "IndexProjections should not be empty")

// 		// Test each index projection from numbers.json
// 		expectedIndexes := []string{
// 			"StatusIndex",
// 			"TypeCountIndex",
// 			"TimeRangeIndex",
// 		}

// 		for _, indexName := range expectedIndexes {
// 			projection, exists := projections[indexName]
// 			assert.True(t, exists, "Index %s should exist in projections", indexName)
// 			assert.NotEmpty(t, projection, "Index %s should have projection attributes", indexName)

// 			t.Logf("Index %s projects %d attributes", indexName, len(projection))
// 		}

// 		t.Logf("✅ Index projections with subtypes validated")
// 	})
// }

// // ==================== Test Helper Functions ====================

// // createNumbersTestItem generates a SchemaItem with specified values and subtypes
// func createNumbersTestItem(suffix string, timestamp uint64, statusCode int, count uint) numbers.SchemaItem {
// 	testUUID := uuid.New()
// 	testTime := time.Now().UTC()
// 	testBigInt := big.NewInt(int64(statusCode * 1000))
// 	testDecimal := decimal.NewFromFloat(float64(count) / 10.0)

// 	return numbers.SchemaItem{
// 		Id:             "test-" + suffix,
// 		Timestamp:      timestamp,
// 		IsActive:       statusCode % 2, // 0 or 1 instead of bool
// 		StatusCode:     statusCode,
// 		Count:          count,
// 		CreatedAt:      testTime,
// 		TinyNumber:     int8(statusCode % 128),
// 		SmallNumber:    int16(statusCode * 10),
// 		MediumNumber:   int32(statusCode * 100),
// 		BigNumber:      int64(statusCode * 1000),
// 		UnsignedTiny:   uint8(statusCode % 256),
// 		UnsignedSmall:  uint16(statusCode * 2),
// 		UnsignedMedium: uint32(statusCode * 20),
// 		PriceCents:     testBigInt,
// 		DiscountRate:   float32(statusCode) / 1000.0,
// 		Score:          float64(statusCode) / 2.0,
// 		Balance:        &testDecimal,
// 		RequestId:      testUUID,
// 		Data:           []byte("test data " + suffix),
// 		DefaultNumber:  float64(statusCode),
// 		Description:    "Test item " + suffix,
// 	}
// }

// // insertNumbersTestItem helper function to insert test items with all subtypes
// func insertNumbersTestItem(t *testing.T, client *dynamodb.Client, ctx context.Context, item numbers.SchemaItem) {
// 	t.Helper()

// 	av, err := numbers.PutItem(item)
// 	require.NoError(t, err, "Item marshaling should succeed")

// 	_, err = client.PutItem(ctx, &dynamodb.PutItemInput{
// 		TableName: aws.String(numbers.TableName),
// 		Item:      av,
// 	})
// 	require.NoError(t, err, "PutItem should succeed")
// }
