package spn

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	level := slog.LevelInfo
	if os.Getenv("GOSPN_DEBUG") != "" {
		level = slog.LevelDebug
	}
	logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level}))
}
