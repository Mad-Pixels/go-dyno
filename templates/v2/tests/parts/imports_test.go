package parts

import (
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/Mad-Pixels/go-dyno/templates/test"
	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestImportsTemplate validates the basic structure and content of the imports template.
// The imports template is static and should contain all required packages for DynamoDB operations.
func TestImportsTemplate(t *testing.T) {
	// Test that the template is properly structured with import block
	// Example: should contain "import (" and ")" to form a valid import block
	t.Run("static_template", func(t *testing.T) {
		rendered := v2.ImportsTemplate

		assert.NotEmpty(t, rendered, "Imports template should not be empty")
		assert.Contains(t, rendered, "import (", "Should contain import block")
		assert.Contains(t, rendered, ")", "Should close import block")
	})

	// Test that generated imports produce valid Go syntax when parsed
	// Example: "package test\nimport (...)" should be parseable by go/parser
	t.Run("go_syntax_valid", func(t *testing.T) {
		testCode := "package test\n\n" + v2.ImportsTemplate

		fset := token.NewFileSet()
		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
		require.NoError(t, err, "Generated imports should be valid Go syntax")
	})

	// Test that all essential packages for DynamoDB operations are included
	// Example: AWS SDK packages, standard library packages like "context", "fmt", etc.
	t.Run("required_imports_present", func(t *testing.T) {
		requiredImports := []string{
			"fmt",
			"context",
			"encoding/json",
			"strings",
			"strconv",
			"sort",
			"github.com/aws/aws-sdk-go-v2/aws",
			"github.com/aws/aws-sdk-go-v2/service/dynamodb",
			"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue",
			"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression",
			"github.com/aws/aws-sdk-go-v2/service/dynamodb/types",
			"github.com/aws/aws-lambda-go/events",
		}
		for _, imp := range requiredImports {
			assert.Contains(t, v2.ImportsTemplate, imp,
				"Should contain required import: %s", imp)
		}
	})

	// Test that no import appears twice in the template
	// Example: "context" should only appear once, not duplicated
	t.Run("no_duplicate_imports", func(t *testing.T) {
		lines := strings.Split(v2.ImportsTemplate, "\n")
		importLines := make([]string, 0)

		inImportBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "import (" {
				inImportBlock = true
				continue
			}
			if trimmed == ")" && inImportBlock {
				break
			}
			if inImportBlock && trimmed != "" {
				cleanImport := strings.Trim(trimmed, `"`)
				importLines = append(importLines, cleanImport)
			}
		}

		seen := make(map[string]bool)
		for _, imp := range importLines {
			assert.False(t, seen[imp], "Import should not be duplicated: %s", imp)
			seen[imp] = true
		}
	})
}

// TestImportsTemplateFormatting validates that the imports template follows Go formatting standards.
// This ensures the generated code doesn't need additional formatting and is production-ready.
func TestImportsTemplateFormatting(t *testing.T) {
	// Test that Go formatters (go fmt, goimports, gofumpt) don't change the template
	// Example: the template should already be properly formatted with correct tabs/spaces
	testCode := "package test\n\n" + v2.ImportsTemplate + "\n\nfunc dummy() {}\n"
	test.TestAllFormattersUnchanged(t, testCode)

	// Test that imports are properly grouped: standard library first, then external packages
	// Example: "fmt", "context" should come before "github.com/aws/..." packages
	t.Run("proper_import_grouping", func(t *testing.T) {
		lines := strings.Split(v2.ImportsTemplate, "\n")
		importLines := make([]string, 0)

		inImportBlock := false
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed == "import (" {
				inImportBlock = true
				continue
			}
			if trimmed == ")" && inImportBlock {
				break
			}
			if inImportBlock && trimmed != "" {
				importLines = append(importLines, trimmed)
			}
		}
		stdImports := []string{"fmt", "context", "encoding/json", "strings", "strconv", "sort"}
		awsImports := []string{"github.com/aws"}
		foundStd := false
		foundAws := false

		for _, line := range importLines {
			for _, std := range stdImports {
				if strings.Contains(line, std) && !strings.Contains(line, "/") {
					foundStd = true
					assert.False(t, foundAws, "Standard imports should come before AWS imports")
				}
			}
			for _, aws := range awsImports {
				if strings.Contains(line, aws) {
					foundAws = true
				}
			}
		}
		assert.True(t, foundStd, "Should have standard imports")
		assert.True(t, foundAws, "Should have AWS imports")
	})
}
