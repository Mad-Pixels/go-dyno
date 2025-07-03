package validation

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/generator"
	"github.com/stretchr/testify/require"
)

// TestGeneratedCodeCompilation validates that complete DynamoDB code generation produces compilable Go code.
//
// Test process:
//  1. Reads JSON schema files from fixtures/ directory (skipping files with 'invalid-' prefix)
//  2. Generates complete Go code using DynamoDB templates via new Generator API
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

			g, err := generator.NewGenerator(schemaFile)
			require.NoError(t, err, "Failed to create generator: %s", schemaFile)

			err = g.Validate()
			require.NoError(t, err, "Failed to validate schema: %s", schemaFile)

			builder := g.NewRenderBuilder()
			generatedCode := builder.Build()
			require.NotEmpty(t, generatedCode, "Generated code is empty")

			CodeCompiles(t, generatedCode, g.PackageName())
		})
	}
}
