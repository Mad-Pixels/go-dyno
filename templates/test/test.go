// Package test provides common testing utilities for validating template generation and code formatting.
//
// This package contains reusable functions for testing Go code templates, including:
//   - Go formatter validation (go fmt, goimports, gofumpt)
//   - Temporary file creation for testing
//   - Command execution helpers with proper error handling
//   - Binary availability checks with helpful skip messages
//
// The utilities are designed to ensure generated code follows Go standards and best practices.
//
// Example usage:
//
//	func TestMyTemplate(t *testing.T) {
//	    generatedCode := "package main\n\nfunc main() {}\n"
//	    test.TestAllFormattersUnchanged(t, generatedCode)
//	}
package test
