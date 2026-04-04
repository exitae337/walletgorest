package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

const (
	pathToConfig = "config/config.env" // deafult value
)

type Config struct {
	Env string `env:"ENV" env-default:"local"`

	HTTPAddress     string        `env:"HTTP_ADDRESS" env-default:":8080"`
	HTTPTimeout     time.Duration `env:"HTTP_TIMEOUT" env-default:"5s"`
	HTTPIdleTimeout time.Duration `env:"HTTP_IDLE_TIMEOUT" env-default:"30s"`

	DBHost            string        `env:"DB_HOST" env-required:"true"`
	DBPort            int           `env:"DB_PORT" env-required:"true"`
	DBUser            string        `env:"DB_USER" env-required:"true"`
	DBPassword        string        `env:"DB_PASSWORD" env-required:"true"`
	DBName            string        `env:"DB_NAME" env-required:"true"`
	DBMaxOpenConns    int           `env:"DB_MAX_OPEN_CONNS" env-default:"30"`
	DBMaxIdleConns    int           `env:"DB_MAX_IDLE_CONNS" env-default:"25"`
	DBConnMaxLifetime time.Duration `env:"DB_CONN_MAX_LIFETIME" env-default:"5m"`
}

func MustLoad() *Config {
	configPath := os.Getenv("CONFIG_PATH")

	var cfg Config

	if configPath == "" {
		configPath = pathToConfig
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("file %s does not exists", configPath)
	}

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Error with reading configuration file: %s", err)
	}

	if err := cfg.Validate(); err != nil {
		log.Fatalf("Error in validating configuration: %s", err)
	}

	return &cfg
}

func (cfg *Config) Validate() error {
	if cfg.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if cfg.DBUser == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if cfg.DBPassword == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if cfg.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if cfg.DBPort <= 0 || cfg.DBPort > 65535 {
		return fmt.Errorf("DB_PORT must be between 1 and 65535, got %d", cfg.DBPort)
	}

	if !strings.HasPrefix(cfg.HTTPAddress, ":") && !strings.Contains(cfg.HTTPAddress, ":") {
		return fmt.Errorf("HTTP_ADDRESS must be in format :port or host:port, got %s", cfg.HTTPAddress)
	}

	if cfg.HTTPTimeout <= 0 {
		return fmt.Errorf("HTTP_TIMEOUT must be positive, got %s", cfg.HTTPTimeout)
	}
	if cfg.HTTPIdleTimeout <= 0 {
		return fmt.Errorf("HTTP_IDLE_TIMEOUT must be positive, got %s", cfg.HTTPIdleTimeout)
	}

	if cfg.DBMaxOpenConns < 5 {
		return fmt.Errorf("DB_MAX_OPEN_CONNS must be at least 5, got %d", cfg.DBMaxOpenConns)
	}
	if cfg.DBMaxIdleConns < 1 {
		return fmt.Errorf("DB_MAX_IDLE_CONNS must be at least 1, got %d", cfg.DBMaxIdleConns)
	}
	if cfg.DBMaxIdleConns > cfg.DBMaxOpenConns {
		return fmt.Errorf("DB_MAX_IDLE_CONNS (%d) cannot exceed DB_MAX_OPEN_CONNS (%d)",
			cfg.DBMaxIdleConns, cfg.DBMaxOpenConns)
	}
	if cfg.DBConnMaxLifetime <= 0 {
		return fmt.Errorf("DB_CONN_MAX_LIFETIME must be positive, got %s", cfg.DBConnMaxLifetime)
	}

	validEnvs := map[string]bool{"local": true, "dev": true, "production": true}
	if !validEnvs[cfg.Env] {
		return fmt.Errorf("ENV must be local/dev/production, got %s", cfg.Env)
	}

	return nil
}
