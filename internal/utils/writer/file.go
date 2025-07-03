package writer

import "github.com/Mad-Pixels/go-dyno/internal/utils/fs"

// FileWriter writes to a file.
type FileWriter struct {
	path string
}

// NewFileWriter creates a new file writer.
func NewFileWriter(path string) *FileWriter {
	return &FileWriter{path: path}
}

// Write implements Writer interface for file output.
func (fw *FileWriter) Write(data []byte) error {
	if err := fs.IsFileOrCreate(fw.path); err != nil {
		return err
	}
	return fs.WriteToFile(fw.path, data)
}

// Type return writer type.
func (fw *FileWriter) Type() string {
	return "file: " + fw.path
}
