package parts

import (
	"go/format"
	"go/parser"
	"go/token"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema/common"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestUtilityFunctionsTemplate validates the UtilityFunctionsTemplate rendering.
// This template generates a set of helper functions for DynamoDB operations:
// - BatchPutItems, PutItem
// - BoolToInt, IntToBool
// - ExtractFromDynamoDBStreamEvent with per-type branches (including SS/NS)
// - IsFieldModified, GetBoolFieldChanged
// - CreateKey, CreateKeyFromItem
// - CreateTriggerHandler
// - ConvertMapToAttributeValues with Set support
func TestUtilityFunctionsTemplate(t *testing.T) {
	// Test with schema including all supported types including SS and NS
	t.Run("basic_template_rendering", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "usertable",
			TableName:   "UserTable",
			HashKey:     "user_id",
			RangeKey:    "created_at",
			AllAttributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "status", Type: "S"},
				{Name: "updated_at", Type: "N"},
				{Name: "is_active", Type: "B"},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)
		testUtilityFunctionsContent(t, rendered)
	})

	// Test with schema including SS and NS types
	t.Run("template_with_set_types", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "userprofile",
			TableName:   "UserProfile",
			HashKey:     "user_id",
			RangeKey:    "profile_type",
			AllAttributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "profile_type", Type: "S"},
				{Name: "tags", Type: "SS"},         // String Set
				{Name: "skill_levels", Type: "NS"}, // Number Set
				{Name: "interests", Type: "SS"},    // String Set
				{Name: "scores", Type: "NS"},       // Number Set
				{Name: "is_premium", Type: "B"},
				{Name: "created_at", Type: "N"},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)
		testUtilityFunctionsWithSets(t, rendered)
	})
}

// testUtilityFunctionsContent validates the basic content of rendered utility functions template
func testUtilityFunctionsContent(t *testing.T, rendered string) {
	t.Helper()

	// Test that the rendered code is valid Go syntax
	// Example: parsing with imports for all used packages should succeed
	t.Run("go_syntax_valid", func(t *testing.T) {
		src := "package test\n\nimport (\n" +
			"\t\"context\"\n" +
			"\t\"encoding/json\"\n" +
			"\t\"fmt\"\n" +
			"\t\"strconv\"\n" +
			"\t\"github.com/aws/aws-lambda-go/events\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
			")\n\n" + rendered
		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "utils.go", src, parser.ParseComments)
		require.NoError(t, err, "Rendered UtilityFunctionsTemplate should be valid Go syntax")
	})

	// Test that all top-level helper functions are present
	t.Run("functions_present", func(t *testing.T) {
		funcs := []string{
			"func BatchPutItems(",
			"func PutItem(",
			"func BoolToInt(",
			"func IntToBool(",
			"func ExtractFromDynamoDBStreamEvent(",
			"func IsFieldModified(",
			"func GetBoolFieldChanged(",
			"func CreateKey(",
			"func CreateKeyFromItem(",
			"func CreateTriggerHandler(",
			"func ConvertMapToAttributeValues(",
		}
		for _, fn := range funcs {
			assert.Contains(t, rendered, fn, "Should contain %s", fn)
		}
	})

	// Test that ExtractFromDynamoDBStreamEvent has per-type branches for each attribute type
	t.Run("extract_branches", func(t *testing.T) {
		// String branch
		assert.Contains(t, rendered, "val.String()", "Should handle S/String attributes")
		// Number branch
		assert.Contains(t, rendered, "strconv.Atoi(val.Number())", "Should handle N/Number attributes")
		// Boolean branch
		assert.Contains(t, rendered, "val.Boolean()", "Should handle B/Boolean attributes")
	})

	// Test that ConvertMapToAttributeValues covers all switch cases
	t.Run("convert_map_branches", func(t *testing.T) {
		branches := []string{
			"case string:",
			"case int:", // Added support for int
			"case float64:",
			"case bool:",
			"case nil:",
			"case map[string]interface{}:",
			"case []interface{}:",
			"default:",
		}
		for _, br := range branches {
			assert.Contains(t, rendered, br, "ConvertMapToAttributeValues should cover %s", br)
		}
	})
}

// testUtilityFunctionsWithSets validates utility functions content with Set types
func testUtilityFunctionsWithSets(t *testing.T, rendered string) {
	t.Helper()

	// Test that Set type branches are present in ExtractFromDynamoDBStreamEvent
	t.Run("extract_set_branches", func(t *testing.T) {
		// String Set branch
		assert.Contains(t, rendered, "val.StringSet()", "Should handle SS/StringSet attributes")
		assert.Contains(t, rendered, "if ss := val.StringSet(); ss != nil", "Should check StringSet for nil")

		// Number Set branch
		assert.Contains(t, rendered, "val.NumberSet()", "Should handle NS/NumberSet attributes")
		assert.Contains(t, rendered, "if ns := val.NumberSet(); ns != nil", "Should check NumberSet for nil")
		assert.Contains(t, rendered, "for _, numStr := range ns", "Should iterate over NumberSet")
		assert.Contains(t, rendered, "strconv.Atoi(numStr)", "Should convert number strings to int")
	})

	// Test that ConvertMapToAttributeValues supports Set types
	t.Run("convert_map_set_branches", func(t *testing.T) {
		// String Set support
		assert.Contains(t, rendered, "case []string:", "Should handle []string for SS type")
		assert.Contains(t, rendered, "&types.AttributeValueMemberSS{Value: v}", "Should create SS AttributeValue")

		// Number Set support
		assert.Contains(t, rendered, "case []int:", "Should handle []int for NS type")
		assert.Contains(t, rendered, "&types.AttributeValueMemberNS{Value: numbers}", "Should create NS AttributeValue")
		assert.Contains(t, rendered, "numbers[i] = fmt.Sprintf(\"%d\", num)", "Should convert int to string for NS")

		// Integer support (for single numbers)
		assert.Contains(t, rendered, "case int:", "Should handle int type")
		assert.Contains(t, rendered, "&types.AttributeValueMemberN{Value: fmt.Sprintf(\"%d\", v)}", "Should create N AttributeValue from int")
	})

	// Test that Set type conditions are properly generated for all SS/NS attributes
	t.Run("set_attribute_conditions", func(t *testing.T) {
		setAttributes := []string{
			"tags",         // SS type
			"skill_levels", // NS type
			"interests",    // SS type
			"scores",       // NS type
		}

		for _, attr := range setAttributes {
			// Should have conditional blocks for each set attribute
			expectedFieldName := utils.ToUpperCamelCase(utils.ToSafeName(attr))
			assert.Contains(t, rendered, expectedFieldName,
				"Should contain field name %s for attribute %s", expectedFieldName, attr)
		}
	})

	// Test formatting and syntax for Set type handling
	t.Run("set_syntax_valid", func(t *testing.T) {
		src := "package test\n\nimport (\n" +
			"\t\"context\"\n" +
			"\t\"encoding/json\"\n" +
			"\t\"fmt\"\n" +
			"\t\"strconv\"\n" +
			"\t\"github.com/aws/aws-lambda-go/events\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue\"\n" +
			"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
			")\n\n" + rendered

		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "utils_sets.go", src, parser.ParseComments)
		require.NoError(t, err, "Rendered template with Sets should be valid Go syntax")
	})
}

// TestUtilityFunctionsTemplateFormatting validates that the rendered template is gofmt-compliant.
// Example: format.Source("package test\n\n" + rendered + "\n") should return no error.
func TestUtilityFunctionsTemplateFormatting(t *testing.T) {
	// Test formatting with basic types
	t.Run("basic_types_formatting", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "usertable",
			TableName:   "UserTable",
			HashKey:     "user_id",
			RangeKey:    "created_at",
			AllAttributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "created_at", Type: "N"},
				{Name: "status", Type: "S"},
				{Name: "updated_at", Type: "N"},
				{Name: "is_active", Type: "B"},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)
		testFormattingCompliance(t, rendered)
	})

	// Test formatting with Set types
	t.Run("set_types_formatting", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "userprofile",
			TableName:   "UserProfile",
			HashKey:     "user_id",
			RangeKey:    "profile_type",
			AllAttributes: []common.Attribute{
				{Name: "user_id", Type: "S"},
				{Name: "profile_type", Type: "S"},
				{Name: "tags", Type: "SS"},
				{Name: "skill_levels", Type: "NS"},
				{Name: "interests", Type: "SS"},
				{Name: "scores", Type: "NS"},
				{Name: "is_premium", Type: "B"},
				{Name: "created_at", Type: "N"},
			},
		}

		rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)
		testFormattingCompliance(t, rendered)
	})

	// Test with edge cases and complex attribute names
	t.Run("edge_cases_formatting", func(t *testing.T) {
		templateMap := v2.TemplateMap{
			PackageName: "complextest",
			TableName:   "ComplexTest",
			HashKey:     "complex_id",
			RangeKey:    "sort_key",
			AllAttributes: []common.Attribute{
				{Name: "complex_id", Type: "S"},
				{Name: "sort_key", Type: "S"},
				{Name: "multi_word_tags", Type: "SS"},     // Complex name
				{Name: "skill-levels", Type: "NS"},        // With hyphen
				{Name: "user_interests", Type: "SS"},      // Underscore
				{Name: "test_scores_2024", Type: "NS"},    // With numbers
				{Name: "is_active_user", Type: "B"},       // Boolean
				{Name: "created_at_timestamp", Type: "N"}, // Long name
			},
		}

		rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)
		testFormattingCompliance(t, rendered)
	})
}

// testFormattingCompliance checks if the rendered template is gofmt-compliant
func testFormattingCompliance(t *testing.T, rendered string) {
	t.Helper()

	full := "package test\n\nimport (\n" +
		"\t\"context\"\n" +
		"\t\"encoding/json\"\n" +
		"\t\"fmt\"\n" +
		"\t\"strconv\"\n" +
		"\t\"github.com/aws/aws-lambda-go/events\"\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue\"\n" +
		"\t\"github.com/aws/aws-sdk-go-v2/service/dynamodb/types\"\n" +
		")\n\n" +
		rendered + "\n"

	_, err := format.Source([]byte(full))
	require.NoError(t, err, "UtilityFunctionsTemplate should be gofmt-compliant")
}

// TestUtilityFunctionsTemplateExtractLogic validates the ExtractFromDynamoDBStreamEvent logic
func TestUtilityFunctionsTemplateExtractLogic(t *testing.T) {
	templateMap := v2.TemplateMap{
		PackageName: "testextract",
		TableName:   "TestExtract",
		HashKey:     "id",
		RangeKey:    "created",
		AllAttributes: []common.Attribute{
			{Name: "id", Type: "S"},
			{Name: "created", Type: "N"},
			{Name: "name", Type: "S"},
			{Name: "tags", Type: "SS"},
			{Name: "scores", Type: "NS"},
			{Name: "is_active", Type: "B"},
		},
	}

	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

	t.Run("extract_all_type_branches", func(t *testing.T) {
		// Check that each attribute type has proper extraction logic
		typeChecks := map[string][]string{
			"S":  {"val.String()", "item.Id =", "item.Name ="},
			"N":  {"strconv.Atoi(val.Number())", "item.Created ="},
			"B":  {"val.Boolean()", "item.IsActive ="},
			"SS": {"val.StringSet()", "item.Tags ="},
			"NS": {"val.NumberSet()", "item.Scores ="},
		}

		for typeCode, checks := range typeChecks {
			for _, check := range checks {
				assert.Contains(t, rendered, check,
					"Should contain %s logic for type %s", check, typeCode)
			}
		}
	})

	t.Run("extract_error_handling", func(t *testing.T) {
		// Check proper error handling for number conversions
		assert.Contains(t, rendered, "if n, err := strconv.Atoi(val.Number()); err == nil",
			"Should have error handling for number conversion")
		assert.Contains(t, rendered, "if num, err := strconv.Atoi(numStr); err == nil",
			"Should have error handling for number set conversion")
	})

	t.Run("extract_nil_checks", func(t *testing.T) {
		// Check nil checks for Set types
		assert.Contains(t, rendered, "if ss := val.StringSet(); ss != nil",
			"Should check StringSet for nil")
		assert.Contains(t, rendered, "if ns := val.NumberSet(); ns != nil",
			"Should check NumberSet for nil")
	})
}
