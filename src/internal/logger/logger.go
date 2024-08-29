package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Debug(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
}

type slogger struct {
	logger *slog.Logger
}

func NewLogger() Logger {
	newLogger := &slogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
	slog.SetDefault(newLogger.logger)

	return newLogger
}

func (s *slogger) Info(msg string, keysAndValues ...interface{}) {
	s.logger.Info(msg, keysAndValues...)
}

func (s *slogger) Error(msg string, keysAndValues ...interface{}) {
	s.logger.Error(msg, keysAndValues...)
}

func (s *slogger) Debug(msg string, keysAndValues ...interface{}) {
	s.logger.Debug(msg, keysAndValues...)
}

func (s *slogger) Warn(msg string, keysAndValues ...interface{}) {
	s.logger.Warn(msg, keysAndValues...)
}
