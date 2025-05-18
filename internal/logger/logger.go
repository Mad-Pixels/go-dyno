package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "01-02 15:04",
	})

	Log.SetOutput(os.Stdout)
	Log.SetLevel(logrus.InfoLevel)

	if lvl, ok := os.LookupEnv("LOG_LEVEL"); ok {
		if parsed, err := logrus.ParseLevel(lvl); err == nil {
			Log.SetLevel(parsed)
		} else {
			Log.Warnf("invalid LOG_LEVEL '%s', using default InfoLevel", lvl)
		}
	}
}

