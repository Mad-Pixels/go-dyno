package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

// Failure represents a structured error with an optional cause and contextual fields.
// It is designed to provide rich error logging and propagation across the application.
type Failure struct {
	// Readable error message.
	Message string

	// Underlying error.
	Err error

	// Additional context for structured logging.
	Fields map[string]any
}

// NewFailure creates a new Failure instance with a given message and optional error.
//
// Example:
//
//	err := NewFailure("failed to open file", io.EOF)
func NewFailure(msg string, err error) *Failure {
	return &Failure{
		Message: msg,
		Err:     err,
		Fields:  make(map[string]any),
	}
}

// Error implements the error interface.
func (f *Failure) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %v", f.Message, f.Err)
	}
	return f.Message
}

// Unwrap returns the wrapped error.
func (f *Failure) Unwrap() error {
	return f.Err
}

// Log emits the error to the global logger at the specified level, including all attached fields.
//
// Example:
//
//	NewFailure("request failed", err).
//		With("method", "POST").
//		With("url", "/api").
//		Log(zerolog.ErrorLevel)
func (f *Failure) Log(level zerolog.Level) {
	e := Log.WithLevel(level)
	for k, v := range f.Fields {
		e = e.Any(k, v)
	}
	e.Err(f.Err).Msg(f.Message)
}

// With adds a structured field to the Failure, useful for attaching context.
//
// Example:
//
//	err := NewFailure("validation error", nil).
//		With("field", "email").
//		With("reason", "missing '@'")
func (f *Failure) With(key string, value any) *Failure {
	f.Fields[key] = value
	return f
}
