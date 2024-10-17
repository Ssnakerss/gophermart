package logger

import (
	"log/slog"
	"os"
)

func Setup(devOrProd string) {
	var logger *slog.Logger
	switch devOrProd {
	case "DEV":
		opts := &slog.HandlerOptions{Level: slog.LevelDebug}
		logger = slog.New(slog.NewTextHandler(os.Stdout, opts))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
	}
	slog.SetDefault(logger)
}
