package parts

// import (
// 	"go/parser"
// 	"go/token"
// 	"strings"
// 	"testing"

// 	v2 "github.com/Mad-Pixels/go-dyno/templates/v2"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"
// )

// // TestImportsTemplate validates the basic structure and content of the imports template.
// // The imports template is static and should contain all required packages for DynamoDB operations.
// func TestImportsTemplate(t *testing.T) {
// 	// Test that the template is properly structured with import block
// 	// Example: should contain "import (" and ")" to form a valid import block
// 	t.Run("static_template", func(t *testing.T) {
// 		rendered := v2.ImportsTemplate

// 		assert.NotEmpty(t, rendered, "Imports template should not be empty")
// 		assert.Contains(t, rendered, "import (", "Should contain import block")
// 		assert.Contains(t, rendered, ")", "Should close import block")
// 	})

// 	// Test that generated imports produce valid Go syntax when parsed
// 	// Example: "package test\nimport (...)" should be parseable by go/parser
// 	t.Run("go_syntax_valid", func(t *testing.T) {
// 		testCode := "package test\n\n" + v2.ImportsTemplate

// 		fset := token.NewFileSet()
// 		_, err := parser.ParseFile(fset, "test.go", testCode, parser.ParseComments)
// 		require.NoError(t, err, "Generated imports should be valid Go syntax")
// 	})

// 	// Test that no import appears twice in the template
// 	// Example: "context" should only appear once, not duplicated
// 	t.Run("no_duplicate_imports", func(t *testing.T) {
// 		lines := strings.Split(v2.ImportsTemplate, "\n")
// 		importLines := make([]string, 0)

// 		inImportBlock := false
// 		for _, line := range lines {
// 			trimmed := strings.TrimSpace(line)
// 			if trimmed == "import (" {
// 				inImportBlock = true
// 				continue
// 			}
// 			if trimmed == ")" && inImportBlock {
// 				break
// 			}
// 			if inImportBlock && trimmed != "" {
// 				cleanImport := strings.Trim(trimmed, `"`)
// 				importLines = append(importLines, cleanImport)
// 			}
// 		}

// 		seen := make(map[string]bool)
// 		for _, imp := range importLines {
// 			assert.False(t, seen[imp], "Import should not be duplicated: %s", imp)
// 			seen[imp] = true
// 		}
// 	})
// }
