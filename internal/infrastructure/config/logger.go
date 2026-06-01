package config

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// NewLogger creates a structured JSON logger at the given level.
// If level is empty, defaults to "info".
func NewLogger(level string) *slog.Logger {
	var lvl slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lvl = slog.LevelDebug
	case "warn", "warning":
		lvl = slog.LevelWarn
	case "error":
		lvl = slog.LevelError
	default:
		lvl = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: lvl,
	}

	var writer io.Writer = os.Stdout

	logger := slog.New(slog.NewJSONHandler(writer, opts))
	return logger
}
