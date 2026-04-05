package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/exitae337/walletgorest/internal/config"
	"github.com/exitae337/walletgorest/internal/http-server/handler"
	"github.com/exitae337/walletgorest/internal/service"
	"github.com/exitae337/walletgorest/internal/storage/postgres"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
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

	// Context for interrupt signal
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 3. Init database connection (Connection Pool)
	databaseConnCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Init and Ping storage
	storage, err := initStorage(databaseConnCtx, cfg, logger)
	if err != nil {
		logger.Error("failed to connect to connection pool", slog.Any("error", err))
		os.Exit(1)
	}

	// Connection pool monitoring
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logger.Info("stopping pool stats monitoring")
				return
			case <-ticker.C:
				stats := storage.GetPoolStats()
				logger.Info("pool stats",
					"total_conns", stats.TotalConns(),
					"idle_conns", stats.IdleConns(),
					"acquired_conns", stats.AcquiredConns(),
					"max_conns", stats.MaxConns(),
				)
			}
		}
	}()

	// Init repo, service for DB and Handler for Server
	walletRepo := postgres.NewWalletRepo(storage)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	router := setupRouter(walletHandler, logger)

	// 4. Init and Start server
	srv := &http.Server{
		Addr:         cfg.HTTPAddress,
		Handler:      router,
		ReadTimeout:  cfg.HTTPTimeout,
		WriteTimeout: cfg.HTTPTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}

	go func() {
		logger.Info("starting http server...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start http server", slog.Any("error", err))
			os.Exit(1)
		}
	}()

	// Waiting for interrupt signal
	<-ctx.Done()
	logger.Info("shutdown signal received")

	shutdownContext, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	logger.Info("shutting down server...")
	if err := srv.Shutdown(shutdownContext); err != nil {
		logger.Error("http server shutdown failed", slog.Any("error", err))
	} else {
		logger.Info("HTTP Server stopped")
	}

	logger.Info("closing database connection...")
	storage.Close(logger)
	logger.Info("database connection closed")

	logger.Info("service stopped gracefully")
}

func setupRouter(handler *handler.WalletHandler, logger *slog.Logger) *chi.Mux {
	logger.Info("setuping router -> chi.Router")

	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "ok"}`))
	})

	router.Route("/api/v1", func(r chi.Router) {
		handler.RegisterRoutes(r)
	})

	return router
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
	storage, err := postgres.New(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}

	if err := storage.HealthCheck(ctx); err != nil {
		storage.Close(logger)
		return nil, err
	}

	logger.Info("database connection established and healthy",
		slog.Int("max_open_conns", cfg.DBMaxOpenConns),
		slog.Int("max_idle_conns", cfg.DBMaxIdleConns),
	)

	return storage, nil
}
