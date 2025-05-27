package test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// ExecResult contains the result of executing a command with stdout, stderr and error information.
type ExecResult struct {
	Output string
	Error  error
	Stderr string
}

// TestAllFormattersUnchanged tests that Go formatters (go fmt, goimports, gofumpt) don't modify the provided code.
// This validates that the code is already properly formatted according to Go standards.
// Example: TestAllFormattersUnchanged(t, "package main\n\nfunc main() {}\n")
func TestAllFormattersUnchanged(t *testing.T, originalCode string) {
	if !strings.HasSuffix(originalCode, "\n") {
		originalCode += "\n"
	}
	testGoFormatterUnchanged(t, "goimports_unchanged", originalCode, execGoImports)
	testGoFormatterUnchanged(t, "gofumpt_unchanged", originalCode, execGoFumpt)
	testGoFormatterUnchanged(t, "go_fmt_unchanged", originalCode, execGoFmt)
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
