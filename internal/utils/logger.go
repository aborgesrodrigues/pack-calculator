package utils

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	logHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	return slog.New(logHandler)
}
