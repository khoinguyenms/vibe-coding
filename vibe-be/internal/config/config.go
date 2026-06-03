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
	Tracing  TracingConfig
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

type TracingConfig struct {
	Enabled      bool
	Endpoint     string
	ServiceName  string
	SamplerRatio float64
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

		Tracing: TracingConfig{
			Enabled:      getEnvAsBool("TRACING_ENABLED", false),
			Endpoint:     getEnvAsString("TRACING_ENDPOINT", "localhost:4317"),
			ServiceName:  getEnvAsString("TRACING_SERVICE_NAME", "vibe-be"),
			SamplerRatio: getEnvAsFloat("TRACING_SAMPLER_RATIO", 1.0),
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

func getEnvAsBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		value, err := strconv.ParseBool(v)
		if err == nil {
			return value
		}
	}
	return fallback
}

func getEnvAsFloat(key string, fallback float64) float64 {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		value, err := strconv.ParseFloat(v, 64)
		if err == nil {
			return value
		}
	}
	return fallback
}
