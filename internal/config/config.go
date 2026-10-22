// Package config loads environment-backed application configuration.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Environment is the resolved runtime environment.
type Environment string

const (
	EnvLocal Environment = "local"
	EnvDev   Environment = "dev"
	EnvProd  Environment = "prod"
)

// Config holds application configuration.
type Config struct {
	AppName  string
	Env      Environment
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// ServerConfig contains HTTP server settings.
type ServerConfig struct {
	Port            int
	ShutdownTimeout time.Duration
	AllowedOrigins  []string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
}

// DatabaseConfig contains optional database settings.
type DatabaseConfig struct {
	URL          string
	MaxOpenConns int
	MaxIdleConns int
	ConnMaxLife  time.Duration
}

// LogConfig controls structured logging.
type LogConfig struct {
	Level     string
	AddSource bool
}

// Load reads .env when present, then environment variables.
func Load() (Config, error) {
	_ = godotenv.Load()

	env := parseEnv(getEnv("APP_ENV", "local"))
	cfg := Config{
		AppName: getEnv("APP_NAME", "Go Clean Template"),
		Env:     env,
		Server: ServerConfig{
			Port:            getEnvInt("HTTP_PORT", 8080),
			ShutdownTimeout: getEnvDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second),
			AllowedOrigins:  splitCSV(getEnv("HTTP_ALLOWED_ORIGINS", "*")),
			ReadTimeout:     getEnvDuration("HTTP_READ_TIMEOUT", 15*time.Second),
			WriteTimeout:    getEnvDuration("HTTP_WRITE_TIMEOUT", 15*time.Second),
			IdleTimeout:     getEnvDuration("HTTP_IDLE_TIMEOUT", 60*time.Second),
		},
		Database: DatabaseConfig{
			URL:          os.Getenv("DATABASE_URL"),
			MaxOpenConns: getEnvInt("DATABASE_MAX_OPEN_CONNS", 10),
			MaxIdleConns: getEnvInt("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLife:  getEnvDuration("DATABASE_CONN_MAX_LIFETIME", time.Hour),
		},
		Log: LogConfig{
			Level:     getEnv("LOG_LEVEL", defaultLogLevel(env)),
			AddSource: getEnvBool("LOG_ADD_SOURCE", false),
		},
	}

	if err := cfg.validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) validate() error {
	if c.AppName == "" {
		return fmt.Errorf("APP_NAME must not be empty")
	}
	if c.Server.Port < 1 || c.Server.Port > 65535 {
		return fmt.Errorf("HTTP_PORT must be between 1 and 65535")
	}
	if c.Env == EnvProd && strings.TrimSpace(c.Database.URL) == "" {
		// Postgres is optional for the template default (local store).
		// Prod still allows empty URL when DATABASE_REQUIRED is not set.
		if getEnvBool("DATABASE_REQUIRED", false) {
			return fmt.Errorf("DATABASE_URL is required when DATABASE_REQUIRED=true")
		}
	}
	return nil
}

// HasDatabase reports whether a database URL is configured.
func (c Config) HasDatabase() bool {
	return strings.TrimSpace(c.Database.URL) != ""
}

func parseEnv(raw string) Environment {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "prod", "production":
		return EnvProd
	case "dev", "development":
		return EnvDev
	default:
		return EnvLocal
	}
}

func defaultLogLevel(env Environment) string {
	if env == EnvProd {
		return "info"
	}
	return "debug"
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return n
}

func getEnvBool(key string, fallback bool) bool {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	v, err := strconv.ParseBool(raw)
	if err != nil {
		return fallback
	}
	return v
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	raw, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(raw) == "" {
		return fallback
	}
	d, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return d
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, p)
		}
	}
	if len(out) == 0 {
		return []string{"*"}
	}
	return out
}
