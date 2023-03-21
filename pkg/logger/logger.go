package logger

import "github.com/rs/zerolog"

// ILogger on case of changing Logger in the future
type ILogger interface {
	Info(msg string)
	Fatal(msg string)
	Debug(msg string)
	Warn(msg string)
}

type Logger struct {
	l zerolog.Logger
}

func New(logger zerolog.Logger) ILogger {
	return &Logger{l: logger}
}

func (l Logger) Info(msg string) {
	l.l.Info().Msg(msg)
}

func (l Logger) Fatal(msg string) {
	l.l.Fatal().Msg(msg)
}

func (l Logger) Debug(msg string) {
	l.l.Debug().Msg(msg)
}

func (l Logger) Warn(msg string) {
	l.l.Warn().Msg(msg)
}
