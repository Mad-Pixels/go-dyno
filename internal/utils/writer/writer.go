// Package writer provides output abstraction.
package writer

// Writer defines the interface for writing generated code.
// Implementations: FileWriter (files), StdoutWriter (stdout).
type Writer interface {
	// Write outputs the generated code data.
	Write(data []byte) error

	// Type returns writer description for logging.
	Type() string
}
