package logger

import (
	"github.com/rs/zerolog"
)

type logWriter struct {
	stdout zerolog.ConsoleWriter
	stderr zerolog.ConsoleWriter
}

func (w logWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func (w logWriter) WriteLevel(l zerolog.Level, p []byte) (n int, err error) {
	if l >= zerolog.WarnLevel {
		return w.stderr.Write(p)
	}
	return w.stdout.Write(p)
}
