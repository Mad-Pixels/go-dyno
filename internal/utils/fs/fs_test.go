package fs

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDirOrCreate_ExistingDir(t *testing.T) {
	tmpDir := t.TempDir()
	err := IsDirOrCreate(tmpDir)
	assert.NoError(t, err)
}

func TestIsDirOrCreate_CreateNewDir(t *testing.T) {
	tmpDir := t.TempDir()
	newDir := filepath.Join(tmpDir, "newdir")
	err := IsDirOrCreate(newDir)
	assert.NoError(t, err)

	info, err := os.Stat(newDir)
	assert.NoError(t, err)
	assert.True(t, info.IsDir())
}

func TestIsDirOrCreate_PathIsFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "testfile")
	os.WriteFile(filePath, []byte("test"), 0644)

	err := IsDirOrCreate(filePath)
	assert.Error(t, err)
}

func TestRemovePath_ExistingPath(t *testing.T) {
	tmpDir := t.TempDir()
	testPath := filepath.Join(tmpDir, "testdir")
	os.MkdirAll(testPath, 0755)

	err := RemovePath(testPath)
	assert.NoError(t, err)

	_, err = os.Stat(testPath)
	assert.True(t, os.IsNotExist(err))
}

func TestRemovePath_NonExistentPath(t *testing.T) {
	err := RemovePath("/nonexistent/path")
	assert.NoError(t, err)
}

func TestAddFileExt_NoExtension(t *testing.T) {
	result := AddFileExt("model", ".go")
	assert.Equal(t, "model.go", result)
}

func TestAddFileExt_ReplaceExtension(t *testing.T) {
	result := AddFileExt("config.json", ".yml")
	assert.Equal(t, "config.yml", result)
}

func TestAddFileExt_ExtensionWithoutDot(t *testing.T) {
	result := AddFileExt("model", "go")
	assert.Equal(t, "model.go", result)
}

func TestAddFileExt_EmptyExtension(t *testing.T) {
	result := AddFileExt("model.go", "")
	assert.Equal(t, "model.", result)
}

func TestReadFile_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	testData := []byte("hello world")
	os.WriteFile(filePath, testData, 0644)

	data, err := ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
}

func TestReadFile_NonExistentFile(t *testing.T) {
	_, err := ReadFile("/nonexistent/file.txt")
	assert.Error(t, err)
}

func TestReadAndParseJSON_ValidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.json")
	testData := map[string]string{"key": "value"}
	jsonData, _ := json.Marshal(testData)
	os.WriteFile(filePath, jsonData, 0644)

	var result map[string]string
	err := ReadAndParseJSON(filePath, &result)
	assert.NoError(t, err)
	assert.Equal(t, testData, result)
}

func TestReadAndParseJSON_InvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "invalid.json")
	os.WriteFile(filePath, []byte("invalid json"), 0644)

	var result map[string]string
	err := ReadAndParseJSON(filePath, &result)
	assert.Error(t, err)
}

func TestWriteToFile_NewFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "test.txt")
	testData := []byte("test content")

	err := WriteToFile(filePath, testData)
	assert.NoError(t, err)

	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, testData, data)
}

func TestWriteToFile_OverwriteExisting(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("old content"), 0644)

	newData := []byte("new content")
	err := WriteToFile(filePath, newData)
	assert.NoError(t, err)

	data, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, newData, data)
}

func TestIsFileOrError_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	err := IsFileOrError(filePath)
	assert.NoError(t, err)
}

func TestIsFileOrError_NonExistentFile(t *testing.T) {
	err := IsFileOrError("/nonexistent/file.txt")
	assert.Error(t, err)
}

func TestIsFileOrError_PathIsDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	err := IsFileOrError(tmpDir)
	assert.Error(t, err)
}

func TestIsFileOrCreate_ExistingFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	err := IsFileOrCreate(filePath)
	assert.NoError(t, err)
}

func TestIsFileOrCreate_CreateNewFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "newfile.txt")

	err := IsFileOrCreate(filePath)
	assert.NoError(t, err)

	info, err := os.Stat(filePath)
	assert.NoError(t, err)
	assert.False(t, info.IsDir())
}

func TestIsFileOrCreate_PathIsDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	err := IsFileOrCreate(tmpDir)
	assert.Error(t, err)
}
