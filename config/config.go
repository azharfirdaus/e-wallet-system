package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLMode  string
}

func LoadPostgresConfig() (PostgresConfig, error) {
	config := PostgresConfig{
		Host:     envOrDefault("POSTGRES_HOST", "localhost"),
		Port:     envOrDefault("POSTGRES_PORT", "5432"),
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		Database: os.Getenv("POSTGRES_DB"),
		SSLMode:  envOrDefault("POSTGRES_SSLMODE", "disable"),
	}

	var missing []string
	if config.User == "" {
		missing = append(missing, "POSTGRES_USER")
	}
	if config.Password == "" {
		missing = append(missing, "POSTGRES_PASSWORD")
	}
	if config.Database == "" {
		missing = append(missing, "POSTGRES_DB")
	}
	if len(missing) > 0 {
		return PostgresConfig{}, errors.New("missing required environment variables: " + strings.Join(missing, ", "))
	}

	return config, nil
}

func (config PostgresConfig) DSN() string {
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(config.User, config.Password),
		Host:   fmt.Sprintf("%s:%s", config.Host, config.Port),
		Path:   config.Database,
	}

	query := dsn.Query()
	query.Set("sslmode", config.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

func envOrDefault(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
