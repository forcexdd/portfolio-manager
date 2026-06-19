package logger

import (
	"context"
	"log/slog"
	"os"
)

type ctxKey string

const loggerKey ctxKey = "logger"

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	With(keysAndValues ...interface{}) Logger
}

type logger struct {
	l *slog.Logger
}

func New() Logger {
	return &logger{l: slog.New(slog.NewJSONHandler(os.Stdout, nil))}
}

func IntoContext(ctx context.Context, log Logger) context.Context {
	return context.WithValue(ctx, loggerKey, log)
}

func FromContext(ctx context.Context, defaultLog Logger) Logger {
	if log, ok := ctx.Value(loggerKey).(Logger); ok && log != nil {
		return log
	}
	return defaultLog
}

func (l *logger) Info(msg string, kv ...interface{})  { l.l.Info(msg, kv...) }
func (l *logger) Error(msg string, kv ...interface{}) { l.l.Error(msg, kv...) }
func (l *logger) Debug(msg string, kv ...interface{}) { l.l.Debug(msg, kv...) }
func (l *logger) Warn(msg string, kv ...interface{})  { l.l.Warn(msg, kv...) }
func (l *logger) Fatal(msg string, kv ...interface{}) {
	l.l.Error(msg, kv...)
	os.Exit(1)
}

func (l *logger) With(kv ...interface{}) Logger {
	return &logger{l: l.l.With(kv...)}
}
