// Package logger provides a structured logging utility built on top of zerolog,
// with support for centralized failure reporting, colored console output,
// log level control via environment variables, and separation of stdout/stderr streams.
//
// The logger is initialized with the Init() function, which reads environment
// variables prefixed by godyno.EnvPrefix (e.g., GODYNO_LOG_LEVEL, GODYNO_LOG_NO_COLOR)
// to configure the behavior dynamically.
//
// Usage:
//
//	logger.Init()
//	logger.Log.Info().Msg("Hello world")
//
//	err := someFunc()
//	if err != nil {
//	    logger.NewFailure("failed to run someFunc", err).
//	        With("context", "example").
//	        Log(zerolog.ErrorLevel)
//	}
package logger

import (
	"fmt"
	"os"
	"strings"

	godyno "github.com/Mad-Pixels/go-dyno"

	"github.com/rs/zerolog"
)

var (
	// Log is the global structured logger instance used throughout the application.
	Log zerolog.Logger

	logNoColor = false
	logLevel   = zerolog.InfoLevel
	logParts   = []string{"level", "message"}
	logFormat  = func(i any) string { return strings.ToUpper(i.(string)) }
)

// Init initializes the global logger (Log) with settings based on environment variables.
// It configures colored output, global log level, and separates stdout/stderr depending on severity.
//
// Environment Variables:
//
//	GODYNO_LOG_LEVEL     — one of "debug", "info", "warn", "error", etc.
//	GODYNO_LOG_NO_COLOR  — "true" disables colored output
//
// Example:
//
//	os.Setenv("GODYNO_LOG_LEVEL", "debug")
//	logger.Init()
//	logger.Log.Debug().Msg("debug message")
//	logger.Log.Error().Msg("error message")
func Init() {
	if lvlStr, ok := os.LookupEnv(fmt.Sprintf("%s_LOG_LEVEL", godyno.EnvPrefix)); ok {
		if lvl, err := zerolog.ParseLevel(strings.ToLower(lvlStr)); err == nil {
			logLevel = lvl
		}
	}
	if noColorStr, ok := os.LookupEnv(fmt.Sprintf("%s_LOG_NO_COLOR", godyno.EnvPrefix)); ok {
		if strings.ToLower(noColorStr) == "true" {
			logNoColor = true
		}
	}
	zerolog.SetGlobalLevel(logLevel)

	stdOutWriter := zerolog.ConsoleWriter{
		Out:         os.Stdout,
		NoColor:     logNoColor,
		PartsOrder:  logParts,
		FormatLevel: logFormat,
	}

	stdErrWriter := zerolog.ConsoleWriter{
		Out:         os.Stderr,
		NoColor:     logNoColor,
		PartsOrder:  logParts,
		FormatLevel: logFormat,
	}

	Log = zerolog.New(logWriter{
		stdout: stdOutWriter,
		stderr: stdErrWriter,
	}).
		Level(logLevel).
		With().
		Logger()
}
