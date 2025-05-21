package utils

import (
	"encoding/json"
	"os"

	"github.com/Mad-Pixels/go-dyno/internal/logger"
)

func ReadFile(path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, logger.NewFailure("failed read content from file", err).
			With("path", path)
	}
	return data, nil
}

func ReadAndParseJsonFile(path string, obj any) error {
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
	if err := os.MkdirAll(path, 0755); err != nil {
		return logger.NewFailure("failed to create a dictionary", err).With("path", path)
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
