package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds application configuration loaded from the environment.
type Config struct {
	DatabaseURL      string
	JWTSecret        string
	Port             string
	CORSAllowOrigins string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Load reads configuration from the environment. It attempts to load a local .env file
// (ignored in production) and never fails if the file is missing.
func Load() (*Config, error) {
	_ = godotenv.Load()

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	cors := os.Getenv("CORS_ALLOW_ORIGINS")
	if cors == "" {
		cors = "*"
	}

	maxOpen := 25
	if v := os.Getenv("DB_MAX_OPEN_CONNS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("DB_MAX_OPEN_CONNS: %w", err)
		}
		maxOpen = n
	}

	maxIdle := 5
	if v := os.Getenv("DB_MAX_IDLE_CONNS"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return nil, fmt.Errorf("DB_MAX_IDLE_CONNS: %w", err)
		}
		maxIdle = n
	}

	connLife := 15 * time.Minute
	if v := os.Getenv("DB_CONN_MAX_LIFETIME"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			return nil, fmt.Errorf("DB_CONN_MAX_LIFETIME: %w", err)
		}
		connLife = d
	}

	return &Config{
		DatabaseURL:      dbURL,
		JWTSecret:        jwtSecret,
		Port:             port,
		CORSAllowOrigins: cors,
		MaxOpenConns:     maxOpen,
		MaxIdleConns:     maxIdle,
		ConnMaxLifetime:  connLife,
	}, nil
}
