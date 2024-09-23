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
	Close() error
}

type slogger struct {
	logger *slog.Logger
	file   *os.File
}

func NewLogger(filePath string) (Logger, error) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	newLogger := &slogger{
		logger: slog.New(slog.NewTextHandler(file, nil)),
		file:   file,
	}
	//slog.SetDefault(newLogger.logger)

	return newLogger, nil
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

func (s *slogger) Close() error {
	if s.file != nil {
		return s.file.Close()
	}

	return nil
}
