package postgres

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/exitae337/walletgorest/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	maxConnIdleTime   = 2 * time.Minute
	healthCheckPeriod = 1 * time.Minute
	connectionTimeout = 5 * time.Second
	appName           = "walletgorest"
)

type Storage struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config, logger *slog.Logger) (*Storage, error) {
	const op = "storage.postgres.New"

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Connection pool settings
	poolConfig.MaxConns = int32(cfg.DBMaxOpenConns)
	poolConfig.MinConns = int32(cfg.DBMaxIdleConns)
	poolConfig.MaxConnLifetime = cfg.DBConnMaxLifetime
	poolConfig.MaxConnIdleTime = maxConnIdleTime
	poolConfig.HealthCheckPeriod = healthCheckPeriod

	poolConfig.ConnConfig.ConnectTimeout = connectionTimeout
	poolConfig.ConnConfig.RuntimeParams["app_name"] = appName

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, err
	}

	createTableSQL := `
    CREATE TABLE IF NOT EXISTS walletdb (
        walletid UUID PRIMARY KEY,
        amount INTEGER NOT NULL DEFAULT 0 CHECK (amount >= 0)
    );
    CREATE INDEX IF NOT EXISTS idx_walletid ON walletdb(walletid);
    `

	if _, err := pool.Exec(ctx, createTableSQL); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	logger.Info("database pool connected, table created",
		"max_conns", poolConfig.MaxConns,
		"min_conns", poolConfig.MinConns,
		"max_conn_lifetime", poolConfig.MaxConnLifetime,
		"max_conn_idle_time", poolConfig.MaxConnIdleTime)

	return &Storage{pool: pool}, nil
}

func (s *Storage) Close(logger *slog.Logger) error {
	const op = "storage.postgres.Close"
	logger.Info("closing database...")
	s.pool.Close()
	logger.Info("database connection closed successfully")
	return nil
}

func (s *Storage) HealthCheck(ctx context.Context) error {
	const op = "storage.postgres.HealthCheck"
	if err := s.pool.Ping(ctx); err != nil {
		return fmt.Errorf("%s: pool health check failed: %w", op, err)
	}
	return nil
}

func (s *Storage) Pool() *pgxpool.Pool {
	return s.pool
}

func (s *Storage) GetPoolStats() *pgxpool.Stat {
	return s.pool.Stat()
}
