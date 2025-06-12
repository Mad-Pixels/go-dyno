package fs

import (
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
