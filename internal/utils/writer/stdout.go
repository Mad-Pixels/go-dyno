package writer

import "os"

// StdoutWriter writes to stdout.
type StdoutWriter struct{}

// NewStdoutWriter creates a new stdout writer.
func NewStdoutWriter() *StdoutWriter {
	return &StdoutWriter{}
}

// Write implements Writer interface for stdout output
func (sw *StdoutWriter) Write(data []byte) error {
	_, err := os.Stdout.Write(data)
	return err
}

// Type return writer type.
func (sw *StdoutWriter) Type() string {
	return "stdout"
}
