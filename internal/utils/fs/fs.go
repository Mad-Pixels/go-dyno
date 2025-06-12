// Package fs provides utility functions for safe and structured file system operations.
//
// It includes helpers for:
//   - Reading and parsing files (e.g. JSON configs)
//   - Writing and overwriting files with proper error handling
//   - Creating or verifying existence of directories and files
//   - Removing directories recursively
//   - Manipulating file extensions
//
// All errors are enriched with context using the internal logger,
// which makes this package suitable for CLI tools and code generators.
package fs

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
)

// IsDirOrCreate checks if the given path exists and is a directory.
// If it does not exist, it creates the directory (and all parent directories).
//
// Example:
//
//	err := fs.IsDirOrCreate("generated/output")
func IsDirOrCreate(path string) error {
	exist, isDir, err := statPath(path)
	if err != nil {
		return logger.NewFailure("failed to stat path", err).
			With("path", path)
	}
	if exist && !isDir {
		return logger.NewFailure("path already exist and it's not a directory", nil).
			With("path", path)
	}
	if exist {
		return nil
	}
	if err := os.MkdirAll(path, 0o755); err != nil {
		return logger.NewFailure("failed to create a dictionary", err).With("path", path)
	}
	return nil
}

// RemovePath removes a file or directory at the specified path recursively.
//
// Example:
//
//	err := fs.RemovePath("tmp/build")
func RemovePath(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return logger.NewFailure("failed to remove path", err).
			With("path", path)
	}
	return nil
}

// AddFileExt ensures that the given file path has the specified extension.
// If the file already has an extension, it will be replaced with the new one.
// The extension should include the dot, e.g. ".go" or ".json".
//
// Example:
//
//	fs.AddFileExt("model", ".go")        → "model.go"
//	fs.AddFileExt("config.json", ".yml") → "config.yml"
func AddFileExt(path string, ext string) string {
	if ext == "" || ext[0] != '.' {
		ext = "." + ext
	}
	base := filepath.Base(path)
	dir := filepath.Dir(path)
	name := base
	if currentExt := filepath.Ext(base); currentExt != "" {
		name = base[:len(base)-len(currentExt)]
	}
	return filepath.Join(dir, name+ext)
}

// ReadFile reads the contents of a file at the given path.
// Returns the file contents or an error wrapped with context.
//
// Example:
//
//	data, err := fs.ReadFile("config.json")
func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, logger.NewFailure("failed read content from file", err).
			With("path", path)
	}
	return data, nil
}

// ReadAndParseJSON reads a JSON file and decodes its contents into the provided object.
//
// Example:
//
//	var cfg Config
//	err := fs.ReadAndParseJSON("config.json", &cfg)
func ReadAndParseJSON(path string, obj any) error {
	b, err := ReadFile(path)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, obj); err != nil {
		return logger.NewFailure("failed to parse JSON", err).
			With("path", path)
	}
	return nil
}

// WriteToFile writes data to a file, creating it if necessary.
// Overwrites any existing content.
//
// Example:
//
//	err := fs.WriteToFile("data.json", []byte(`{"key": "value"}`))
func WriteToFile(path string, data []byte) error {
	if err := IsFileOrCreate(path); err != nil {
		return err
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return logger.NewFailure("couldn't open a file", err).
			With("path", path)
	}
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return logger.NewFailure("couldn't write to file", err).
			With("path", path)
	}
	return nil
}

// IsFileOrError checks if the given path exists and is a file.
// Returns a descriptive error if the path does not exist or is a directory.
//
// Example:
//
//	err := fs.IsFileOrError("data.json")
func IsFileOrError(path string) error {
	exist, isDir, err := statPath(path)
	if err != nil {
		return logger.NewFailure("failed to stat path", err).
			With("path", path)
	}
	if !exist {
		return logger.NewFailure("file doesn't exist", nil).
			With("path", path)
	}
	if isDir {
		return logger.NewFailure("path is not a file", nil).
			With("path", path)
	}
	return nil
}

// IsFileOrCreate checks if a file exists at the given path, and creates it if it does not.
// Ensures the parent directory exists.
//
// Example:
//
//	err := fs.IsFileOrCreate("output/result.txt")
func IsFileOrCreate(path string) error {
	exist, isDir, err := statPath(path)
	if err != nil {
		return err
	}
	if exist && isDir {
		return logger.NewFailure("already exist and it's not a file", nil).
			With("path", path)
	}
	if exist {
		return nil
	}
	if err := IsDirOrCreate(filepath.Dir(path)); err != nil {
		return logger.NewFailure("failed to create file", err).
			With("path", path)
	}
	file, err := os.Create(path)
	if err != nil {
		return logger.NewFailure("failed to create  file", err).
			With("path", path)
	}
	defer file.Close()
	return nil
}

func statPath(path string) (exist bool, isDir bool, err error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false, false, nil
		}
		return false, false, err
	}
	return true, info.IsDir(), nil
}
