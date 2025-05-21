package utils

import (
	"fmt"
	"os"
)

func IsFileOrError(path string) error {
	exist, isDir, err := statPath(path)
	if err != nil {
		return err
	}
	if !exist {
		return fmt.Errorf("ERROR: file '%s' doesn't exist", path)
	}
	if isDir {
		return fmt.Errorf("ERROR: '%s' is not a file", path)
	}
	return nil
}

func IsDirOrCreate(path string) error {
	exist, isDir, err := statPath(path)
	if err != nil {
		return err
	}
	if exist && !isDir {
		return fmt.Errorf("ERROR: '%s' already exist and it's not a directory", path)
	}
	if exist {
		return nil
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("ERROR: failed to create '%s' directory, %w", path, err)
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
