package godyno

import "github.com/sirupsen/logrus"

const (
	Name      = "godyno"
	Usage     = "generate dynamoDB objects"
	EnvPrefix = "GODYNO"
)

var (
	Logger *logrus.Logger = logrus.New()
)
