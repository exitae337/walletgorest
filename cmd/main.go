package main

import (
	"context"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/exitae337/walletgorest/internal/config"
	"github.com/exitae337/walletgorest/internal/storage/postgres"
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

	// 3. Init database connection (Connection Pool)
	databaseConnCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storage, err := initStorage(databaseConnCtx, cfg, logger)
	if err != nil {
		logger.Error("failed to connect to connection pool", slog.Any("error", err))
		os.Exit(1)
	}
	defer storage.Close(logger)

	// Connection pool monitoring
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			stats := storage.GetPoolStats()
			slog.Info("pool stats",
				"total_conns", stats.TotalConns(),
				"idle_conns", stats.IdleConns(),
				"acquired_conns", stats.AcquiredConns(),
				"max_conns", stats.MaxConns(),
			)
		}
	}()

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

func initStorage(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*postgres.Storage, error) {
	logger.Info("init database connection",
		slog.String("host", cfg.DBHost),
		slog.String("name", cfg.DBName),
		slog.String("max open conns", strconv.Itoa(cfg.DBMaxOpenConns)),
	)
	storage, err := postgres.New(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	logger.Info("database init success")
	return storage, nil
}
