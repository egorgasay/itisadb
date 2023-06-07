package logger

import (
	"go.uber.org/zap"
)

// ILogger on case of changing Logger in the future
type ILogger interface {
	Info(msg string)
	Fatal(msg string)
	Debug(msg string)
	Warn(msg string)
}

type Logger struct {
	l *zap.Logger
}

func New(logger *zap.Logger) ILogger {
	return &Logger{l: logger}
}

func (l Logger) Info(msg string) {
	l.l.Info(msg)
}

func (l Logger) Fatal(msg string) {
	l.l.Fatal(msg)
}

func (l Logger) Debug(msg string) {
	l.l.Debug(msg)
}

func (l Logger) Warn(msg string) {
	l.l.Warn(msg)
}
