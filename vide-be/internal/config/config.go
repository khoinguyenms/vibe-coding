package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	AppEnv   string
	AppPort  string
	LogLevel string
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSLMode)
}

func Load() *Config {
	return &Config{
		AppEnv:   getEnvAsString("APP_ENV", "development"),
		AppPort:  getEnvAsString("APP_PORT", "8080"),
		LogLevel: getEnvAsString("LOG_LEVEL", "info"),

		Database: DatabaseConfig{
			Host:     getEnvAsString("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnvAsString("DB_USER", "postgres"),
			Password: getEnvAsString("DB_PASSWORD", "postgres"),
			Name:     getEnvAsString("DB_NAME", "vide_be"),
			SSLMode:  getEnvAsString("DB_SSLMODE", "disable"),
		},
	}
}

func getEnvAsString(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		value, err := strconv.Atoi(v)
		if err == nil {
			return value
		}
	}
	return fallback
}
