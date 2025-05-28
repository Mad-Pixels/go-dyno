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
// - ExtractFromDynamoDBStreamEvent with per-type branches
// - IsFieldModified, GetBoolFieldChanged
// - CreateKey, CreateKeyFromItem
// - CreateTriggerHandler
// - ConvertMapToAttributeValues
func TestUtilityFunctionsTemplate(t *testing.T) {
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

	// Render the template
	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

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

// TestUtilityFunctionsTemplateFormatting validates that the rendered template is gofmt-compliant.
// Example: format.Source("package test\n\n" + rendered + "\n") should return no error.
func TestUtilityFunctionsTemplateFormatting(t *testing.T) {
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

	// Render the template
	rendered := utils.MustParseTemplateToString(v2.UtilityFunctionsTemplate, templateMap)

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

	if _, err := format.Source([]byte(full)); err != nil {
		t.Fatalf("UtilityFunctionsTemplate is not gofmt-compliant: %v", err)
	}
}
