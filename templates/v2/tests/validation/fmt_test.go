package validation

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/schema"
	"github.com/Mad-Pixels/go-dyno/internal/utils"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/require"
)

// TestGeneratedCodeFormatting validates that complete DynamoDB code generation produces properly formatted Go code.
//
// Test process:
//  1. Reads JSON schema files from .tmpl/ directory
//  2. Generates complete Go code using DynamoDB templates
//  3. Runs formatting validation (go fmt, goimports, gofumpt)
//
// This ensures generated code is production-ready and passes all standard Go formatting tools.
func TestGeneratedCodeFormatting(t *testing.T) {
	templatesDir := filepath.Join(".", ".tmpl")
	schemaFiles, err := filepath.Glob(filepath.Join(templatesDir, "*.json"))
	require.NoError(t, err, "Failed to read template files")
	require.NotEmpty(t, schemaFiles, "No JSON files found in %s", templatesDir)

	for _, schemaFile := range schemaFiles {
		schemaName := strings.TrimSuffix(filepath.Base(schemaFile), ".json")

		t.Run(schemaName, func(t *testing.T) {
			dynamoSchema, err := schema.LoadSchema(schemaFile)
			require.NoError(t, err, "Failed to load schema: %s", schemaFile)

			templateMap := v2.TemplateMap{
				PackageName:      dynamoSchema.PackageName(),
				TableName:        dynamoSchema.TableName(),
				HashKey:          dynamoSchema.HashKey(),
				RangeKey:         dynamoSchema.RangeKey(),
				Attributes:       dynamoSchema.Attributes(),
				CommonAttributes: dynamoSchema.CommonAttributes(),
				AllAttributes:    dynamoSchema.AllAttributes(),
				SecondaryIndexes: dynamoSchema.SecondaryIndexes(),
			}

			generatedCode := utils.MustParseTemplateFormattedToString(v2.CodeTemplate, templateMap)
			require.NotEmpty(t, generatedCode, "Generated code is empty")
			AllFormattersUnchanged(t, generatedCode)
		})
	}
}
