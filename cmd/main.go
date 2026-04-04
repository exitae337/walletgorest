package main

import (
	"log/slog"
	"os"

	"github.com/exitae337/walletgorest/internal/config"
)

const (
	envLocal = "local"
	envProd  = "production"
	envDev   = "dev"
)

func main() {
	// 1. Read and Init config file
	cfg := config.MustLoad()
	// 2. Init logger for REST app
	logger := setupLogger(cfg.Env)
	logger.Info("starting wallet service", slog.String("env", cfg.Env))
	logger.Debug("debug messages are enabled")
	// 3. Init database connection
	// 4. Init Server and Handlers
	// 5. Start server
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger
	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return logger
}
