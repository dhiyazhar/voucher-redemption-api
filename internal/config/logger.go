package config

import (
	"log/slog"
	"os"

	"github.com/spf13/viper"
)

func NewLogger(viper *viper.Viper) *slog.Logger {
	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})

	return slog.New(handler)
}
