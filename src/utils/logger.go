package utils

import (
	"os"

	"go.uber.org/zap"
)

type Logger struct {
	Logger *zap.Logger
}

func NewLogger() *Logger {
	env := os.Getenv("ENVIRONMENT")
	if env == "production" {
		logger, _ := zap.NewProduction()
		return &Logger{Logger: logger}
	}

	logger, _ := zap.NewDevelopment()
	return &Logger{Logger: logger}
}

func (l *Logger) Info(message string) {
	l.Logger.Info(message)
}

func (l *Logger) Error(message string) {
	l.Logger.Error(message)
}

func (l *Logger) Warn(message string) {
	l.Logger.Warn(message)
}

func (l *Logger) Debug(message string) {
	l.Logger.Debug(message)
}
