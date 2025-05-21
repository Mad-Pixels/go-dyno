package logger

import (
	"fmt"
	"os"
	"strings"

	godyno "github.com/Mad-Pixels/go-dyno"

	"github.com/rs/zerolog"
)

var (
	Log zerolog.Logger

	logNoColor = false
	logLevel   = zerolog.InfoLevel
	logParts   = []string{"level", "message"}
	logFormat  = func(i any) string { return strings.ToUpper(i.(string)) }
)

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
