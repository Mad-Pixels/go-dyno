package logger

import (
	"fmt"

	"github.com/rs/zerolog"
)

type Failure struct {
	Message string
	Err     error
	Fields  map[string]any
}

func NewFailure(msg string, err error) *Failure {
	return &Failure{
		Message: msg,
		Err:     err,
		Fields:  make(map[string]any),
	}
}

func (f *Failure) Error() string {
	if f.Err != nil {
		return fmt.Sprintf("%s: %v", f.Message, f.Err)
	}
	return f.Message
}

func (f *Failure) Unwrap() error {
	return f.Err
}

func (f *Failure) Log(level zerolog.Level) {
	e := Log.WithLevel(level)
	for k, v := range f.Fields {
		e = e.Any(k, v)
	}
	e.Err(f.Err).Msg(f.Message)
}

func (f *Failure) With(key string, value any) *Failure {
	f.Fields[key] = value
	return f
}
