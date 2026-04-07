package log

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
}

// Info logs a message at INFO level.
func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

// Error logs a message at ERROR level.
func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

// Warn logs a message at WARN level.
func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

// Debug logs a message at DEBUG level.
func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

// With returns a Logger with preset fields.
func With(args ...any) *slog.Logger {
	return logger.With(args...)
}
