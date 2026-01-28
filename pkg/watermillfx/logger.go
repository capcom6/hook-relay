package watermillfx

import (
	"github.com/ThreeDotsLabs/watermill"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

// Debug implements watermill.LoggerAdapter.
func (l *Logger) Debug(msg string, fields watermill.LogFields) {
	l.logger.Debug(msg, l.convertFields(fields)...)
}

// Error implements watermill.LoggerAdapter.
func (l *Logger) Error(msg string, err error, fields watermill.LogFields) {
	l.logger.Error(msg, append(l.convertFields(fields), zap.Error(err))...)
}

// Info implements watermill.LoggerAdapter.
func (l *Logger) Info(msg string, fields watermill.LogFields) {
	l.logger.Info(msg, l.convertFields(fields)...)
}

// Trace implements watermill.LoggerAdapter.
func (l *Logger) Trace(msg string, fields watermill.LogFields) {
	l.logger.Debug(msg, l.convertFields(fields)...) // Using Debug level for Trace as zap doesn't have Trace level
}

// With implements watermill.LoggerAdapter.
func (l *Logger) With(fields watermill.LogFields) watermill.LoggerAdapter {
	newLogger := l.logger.With(l.convertFields(fields)...)
	return &Logger{
		logger: newLogger,
	}
}

func (l *Logger) convertFields(fields watermill.LogFields) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}
	return zapFields
}

var _ watermill.LoggerAdapter = (*Logger)(nil)
