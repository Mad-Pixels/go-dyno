package utils

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
)

// ReadFile reads the contents of a file at the given path.
// Returns the file contents or an error wrapped with context.
//
// Example:
//
//	data, err := utils.ReadFile("config.json")
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
//	err := utils.ReadAndParseJSON("config.json", &cfg)
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

// IsFileOrError checks if the given path exists and is a file.
// Returns a descriptive error if the path does not exist or is a directory.
//
// Example:
//
//	err := utils.IsFileOrError("data.json")
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
//	err := utils.IsFileOrCreate("output/result.txt")
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

// WriteToFile writes data to a file, creating it if necessary.
// Overwrites any existing content.
//
// Example:
//
//	err := utils.WriteToFile("data.json", []byte(`{"key": "value"}`))
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

// IsDirOrCreate checks if the given path exists and is a directory.
// If it does not exist, it creates the directory (and all parent directories).
//
// Example:
//
//	err := utils.IsDirOrCreate("generated/output")
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
//	err := utils.RemovePath("tmp/build")
func RemovePath(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return logger.NewFailure("failed to remove path", err).
			With("path", path)
	}
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
