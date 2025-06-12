package validation

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/generator"
	"github.com/Mad-Pixels/go-dyno/internal/utils/tmplkit"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/require"
)

// TestGeneratedCodeCompilation validates that complete DynamoDB code generation produces compilable Go code.
//
// Test process:
//  1. Reads JSON schema files from .tmpl/ directory (skipping files with 'invalid-' prefix)
//  2. Generates complete Go code using DynamoDB templates
//  3. Creates temporary module with proper dependencies
//  4. Runs "go build" to ensure compilation succeeds
//
// This ensures generated code is syntactically correct and all dependencies resolve properly.
func TestGeneratedCodeCompilation(t *testing.T) {
	schemaFiles, err := filepath.Glob(filepath.Join(EXAMPLES, "*.json"))
	require.NoError(t, err, "Failed to read template files")
	require.NotEmpty(t, schemaFiles, "No JSON files found in %s", EXAMPLES)

	for _, schemaFile := range schemaFiles {
		schemaFile := schemaFile
		schemaName := strings.TrimSuffix(filepath.Base(schemaFile), ".json")

		if strings.HasPrefix(schemaName, "invalid-") {
			t.Logf("Skipping invalid schema for compilation test: %s", schemaName)
			continue
		}

		t.Run(schemaName, func(t *testing.T) {
			t.Parallel()

			dynamoSchema, err := generator.Load(schemaFile)
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

			generatedCode := tmplkit.MustParseTemplateFormattedToString(v2.CodeTemplate, templateMap)
			require.NotEmpty(t, generatedCode, "Generated code is empty")

			CodeCompiles(t, generatedCode, dynamoSchema.PackageName())
		})
	}
}
