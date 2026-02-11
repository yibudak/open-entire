package logging

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Setup configures the global slog logger based on config level and quiet mode.
func Setup(level string, quiet bool) {
	var w io.Writer = os.Stderr
	if quiet {
		w = io.Discard
	}

	opts := &slog.HandlerOptions{
		Level: parseLevel(level),
	}
	handler := slog.NewTextHandler(w, opts)
	slog.SetDefault(slog.New(handler))
}

func parseLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
