package validation

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Mad-Pixels/go-dyno/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// findProjectRoot finds the project root by looking for go.mod file
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", os.ErrNotExist
}

// getSchemaPath returns absolute path to schema file in test data directory
func getSchemaPath(t *testing.T, filename string) string {
	t.Helper()

	projectRoot, err := findProjectRoot()
	require.NoError(t, err, "Should find project root")

	schemaPath := filepath.Join(projectRoot, "tests", "fixtures", filename)

	_, err = os.Stat(schemaPath)
	require.NoError(t, err, "Schema file should exist: %s", schemaPath)
	return schemaPath
}

// TestSchemaValidation tests that invalid JSON schemas are properly rejected
// during the LoadSchema phase, before any code generation occurs.
//
// Test Coverage:
// - Invalid subtype combinations (string with float, number with string, etc.)
// - Empty attribute names
// - Unknown DynamoDB types
// - Incompatible type/subtype pairs
//
// These tests ensure that users get clear error messages for invalid schemas
// rather than generating broken code or runtime errors.
func TestSchemaValidation(t *testing.T) {
	testCases := []struct {
		name          string
		schemaFile    string
		expectError   bool
		errorContains string
		description   string
	}{
		{
			name:        "valid_schema_should_pass_base-string-all",
			schemaFile:  "base-string__all.json",
			expectError: false,
			description: "Valid schema should load without errors",
		},
		{
			name:        "valid_schema_should_pass_base-boolean-all",
			schemaFile:  "base-boolean__all.json",
			expectError: false,
			description: "Valid schema should load without errors",
		},
		{
			name:        "valid_schema_should_pass_user-posts-compile-all",
			schemaFile:  "user-posts-complete__all.json",
			expectError: false,
			description: "Valid schema should load without errors",
		},
		{
			name:        "valid_schema_should_pass_user-posts-compile-min",
			schemaFile:  "user-posts-complete__min.json",
			expectError: false,
			description: "Valid schema should load without errors",
		},
		{
			name:          "invalid_string_with_float_subtype",
			schemaFile:    "invalid-string-with-float.json",
			expectError:   true,
			errorContains: "incompatible subtype",
			description:   "String attribute cannot have float32 subtype",
		},
		{
			name:          "invalid_number_with_string_subtype",
			schemaFile:    "invalid-number-with-string.json",
			expectError:   true,
			errorContains: "incompatible subtype",
			description:   "Number attribute cannot have string subtype",
		},
		{
			name:          "invalid_empty_attribute_name",
			schemaFile:    "invalid-empty-name.json",
			expectError:   true,
			errorContains: "attribute name cannot be empty",
			description:   "Attribute name cannot be empty",
		},
		{
			name:          "invalid_unknown_dynamodb_type",
			schemaFile:    "invalid-unknown-type.json",
			expectError:   true,
			errorContains: "invalid attribute type",
			description:   "Unknown DynamoDB types should be rejected",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			schemaPath := getSchemaPath(t, tc.schemaFile)
			g, err := generator.NewGenerator(schemaPath)
			if tc.expectError {
				if err != nil {
					assert.Error(t, err, "Expected validation error for %s", tc.name)
					assert.Contains(t, err.Error(), tc.errorContains,
						"Error should contain expected message for %s", tc.name)
					t.Logf("✅ Correctly rejected invalid schema: %s", err.Error())
					return
				}

				require.NotNil(t, g, "Generator should be created")
				err = g.Validate()
				assert.Error(t, err, "Expected validation error for %s", tc.name)
				assert.Contains(t, err.Error(), tc.errorContains,
					"Error should contain expected message for %s", tc.name)
				t.Logf("✅ Correctly rejected invalid schema: %s", err.Error())
			} else {
				assert.NoError(t, err, "Valid schema should load without error")
				assert.NotNil(t, g, "Valid schema should not be nil")

				err = g.Validate()
				assert.NoError(t, err, "Valid schema should validate without error")
				t.Logf("✅ Valid schema loaded successfully")
			}
		})
	}
}

// TestValidationErrorMessages ensures error messages are helpful for users
func TestValidationErrorMessages(t *testing.T) {
	testCases := []struct {
		name         string
		schemaFile   string
		expectedMsgs []string
	}{
		{
			name:       "subtype_compatibility_error",
			schemaFile: "invalid-string-with-float.json",
			expectedMsgs: []string{
				"incompatible subtype",
			},
		},
		{
			name:       "empty_name_error",
			schemaFile: "invalid-empty-name.json",
			expectedMsgs: []string{
				"attribute name cannot be empty",
			},
		},
		{
			name:       "unknown_type_error",
			schemaFile: "invalid-unknown-type.json",
			expectedMsgs: []string{
				"invalid attribute type",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			schemaPath := getSchemaPath(t, tc.schemaFile)
			g, err := generator.NewGenerator(schemaPath)
			if err == nil {
				err = g.Validate()
			}
			require.Error(t, err, "Should get validation error")
			for _, expectedMsg := range tc.expectedMsgs {
				assert.Contains(t, err.Error(), expectedMsg,
					"Error message should contain: %s", expectedMsg)
			}
			t.Logf("✅ Error message contains expected content: %s", err.Error())
		})
	}
}
