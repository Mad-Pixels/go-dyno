// Package validation provides utilities for validating generated Go code quality.
//
// This package tests complete DynamoDB code generation templates for:
//   - Format compliance (go fmt, goimports, gofumpt)
//   - Compilation validation
//   - Template rendering correctness
package validation

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// EXAMPLES define a path to JSON-schemas examples which will use for tests.
const EXAMPLES = "../../examples"

// ExecResult contains the result of executing a command with stdout, stderr and error information.
type ExecResult struct {
	Output string
	Error  error
	Stderr string
}

// AllFormattersUnchanged checks that Go formatters (go fmt, goimports, gofumpt) don't modify the provided code.
// This validates that the code is already properly formatted according to Go standards.
// Example: AllFormattersUnchanged(t, "package main\n\nfunc main() {}\n")
func AllFormattersUnchanged(t *testing.T, originalCode string) {
	if !strings.HasSuffix(originalCode, "\n") {
		originalCode += "\n"
	}
	testGoFormatterUnchanged(t, "goimports_unchanged", originalCode, execGoImports)
	testGoFormatterUnchanged(t, "gofumpt_unchanged", originalCode, execGoFumpt)
	testGoFormatterUnchanged(t, "go_fmt_unchanged", originalCode, execGoFmt)
}

// CodeCompiles checks that the provided Go code compiles successfully and passes go vet.
// Creates a temporary module with required dependencies and attempts compilation.
// Example: CodeCompiles(t, generatedCode, "mypackage")
func CodeCompiles(t *testing.T, code, packageName string) {
	if !strings.HasSuffix(code, "\n") {
		code += "\n"
	}

	tempDir := t.TempDir()
	if err := createGoMod(tempDir); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	goFileName := fmt.Sprintf("%s.go", packageName)
	goFilePath := filepath.Join(tempDir, goFileName)
	if err := os.WriteFile(goFilePath, []byte(code), 0644); err != nil {
		t.Fatalf("Failed to write Go file: %v", err)
	}
	tidyResult := execGoModTidy(t, tempDir)
	if tidyResult.Error != nil {
		t.Fatalf("Failed to run go mod tidy: %v\nStderr: %s", tidyResult.Error, tidyResult.Stderr)
	}

	buildResult := execGoBuild(t, tempDir)
	if buildResult.Error != nil {
		t.Errorf("Generated code failed to compile")
		t.Logf("Build error: %v", buildResult.Error)
		t.Logf("Build stderr: %s", buildResult.Stderr)
		t.Logf("Build output: %s", buildResult.Output)
		return
	}
	vetResult := execGoVet(t, tempDir)
	if vetResult.Error != nil {
		t.Errorf("Generated code failed go vet checks")
		t.Logf("Vet error: %v", vetResult.Error)
		t.Logf("Vet stderr: %s", vetResult.Stderr)
		t.Logf("Vet output: %s", vetResult.Output)
	}
}

func execGoFmt(t *testing.T, filePath string) (string, error) {
	t.Helper()

	result := execCommand(t, "go", "fmt", filePath)
	if result.Error != nil {
		t.Logf("go fmt failed: %v", result.Error)
		if result.Stderr != "" {
			t.Logf("go fmt stderr: %s", result.Stderr)
		}
		return "", fmt.Errorf("go fmt failed: %v", result.Error)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read formatted file: %v", err)
	}
	return string(content), nil
}

func execGoImports(t *testing.T, filePath string) (string, error) {
	t.Helper()
	if !checkBinaryExists("goimports") {
		t.Skip("goimports not found in PATH - install with: go install golang.org/x/tools/cmd/goimports@latest")
	}
	result := execCommand(t, "goimports", "-w", filePath)
	if result.Error != nil {
		t.Logf("goimports failed: %v", result.Error)
		if result.Stderr != "" {
			t.Logf("goimports stderr: %s", result.Stderr)
		}
		return "", fmt.Errorf("goimports failed: %v", result.Error)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read goimports output: %v", err)
	}
	return string(content), nil
}

func execGoFumpt(t *testing.T, filePath string) (string, error) {
	t.Helper()
	if !checkBinaryExists("gofumpt") {
		t.Skip("gofumpt not found in PATH - install with: go install mvdan.cc/gofumpt@latest")
	}
	result := execCommand(t, "gofumpt", "-w", filePath)
	if result.Error != nil {
		t.Logf("gofumpt failed: %v", result.Error)
		if result.Stderr != "" {
			t.Logf("gofumpt stderr: %s", result.Stderr)
		}
		return "", fmt.Errorf("gofumpt failed: %v", result.Error)
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read gofumpt output: %v", err)
	}
	return string(content), nil
}

// execGoModTidy runs "go mod tidy" in the specified directory
func execGoModTidy(t *testing.T, dir string) ExecResult {
	t.Helper()
	return execCommand(t, "go", "mod", "tidy", "-C", dir)
}

// execGoBuild runs "go build" in the specified directory
func execGoBuild(t *testing.T, dir string) ExecResult {
	t.Helper()
	return execCommand(t, "go", "build", "-C", dir, "./...")
}

// execGoVet runs "go vet" in the specified directory
func execGoVet(t *testing.T, dir string) ExecResult {
	t.Helper()
	return execCommand(t, "go", "vet", "-C", dir, "./...")
}

func execCommand(t *testing.T, name string, args ...string) ExecResult {
	t.Helper()

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	return ExecResult{
		Output: stdout.String(),
		Error:  err,
		Stderr: stderr.String(),
	}
}

func checkBinaryExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func createTempGoFile(t *testing.T, content string) string {
	t.Helper()

	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.go")

	err := os.WriteFile(testFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	return testFile
}

// createGoMod creates a go.mod file with required DynamoDB dependencies
// Uses current Go version and latest package versions
func createGoMod(dir string) error {
	goVersion, err := getCurrentGoVersion()
	if err != nil {
		return err
	}

	goModTemplate := `module testmodule

go %s

require (
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go-v2 v1.24.0
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue v1.12.14
	github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression v1.6.14
	github.com/aws/aws-sdk-go-v2/service/dynamodb v1.26.7
)
`

	goModContent := fmt.Sprintf(goModTemplate, goVersion)
	return os.WriteFile(filepath.Join(dir, "go.mod"), []byte(goModContent), 0644)
}

func getCurrentGoVersion() (string, error) {
	cmd := exec.Command("go", "version")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	versionStr := string(output)
	parts := strings.Fields(versionStr)
	if len(parts) < 3 {
		return "", fmt.Errorf("unexpected go version output: %s", versionStr)
	}
	fullVersion := parts[2]
	if !strings.HasPrefix(fullVersion, "go") {
		return "", fmt.Errorf("unexpected version format: %s", fullVersion)
	}

	version := strings.TrimPrefix(fullVersion, "go")
	versionParts := strings.Split(version, ".")
	if len(versionParts) < 2 {
		return "", fmt.Errorf("unexpected version format: %s", version)
	}
	return fmt.Sprintf("%s.%s", versionParts[0], versionParts[1]), nil
}

func testGoFormatterUnchanged(t *testing.T, testName string, originalCode string, formatterFunc func(*testing.T, string) (string, error)) {
	t.Helper()

	t.Run(testName, func(t *testing.T) {
		filePath := createTempGoFile(t, originalCode)

		formattedCode, err := formatterFunc(t, filePath)
		if err != nil {
			t.Fatalf("Formatter failed: %v", err)
		}

		if originalCode != formattedCode {
			t.Errorf("Formatter changed the code - it should already be properly formatted")
			t.Logf("Original:\n%s", originalCode)
			t.Logf("Formatted:\n%s", formattedCode)
		}
	})
}
