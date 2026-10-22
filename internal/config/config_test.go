package config_test

import (
	"testing"

	"github.com/savrin-sharif/go-clean-template/internal/config"
)

func TestLoadDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "Test Service")
	t.Setenv("APP_ENV", "dev")
	t.Setenv("HTTP_PORT", "9090")
	t.Setenv("HTTP_ALLOWED_ORIGINS", "http://localhost:3000, https://example.com")
	t.Setenv("DATABASE_URL", "")

	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.AppName != "Test Service" {
		t.Fatalf("AppName = %q", cfg.AppName)
	}
	if cfg.Env != config.EnvDev {
		t.Fatalf("Env = %q", cfg.Env)
	}
	if cfg.Server.Port != 9090 {
		t.Fatalf("Port = %d", cfg.Server.Port)
	}
	if len(cfg.Server.AllowedOrigins) != 2 {
		t.Fatalf("AllowedOrigins = %#v", cfg.Server.AllowedOrigins)
	}
	if cfg.HasDatabase() {
		t.Fatal("HasDatabase() = true, want false")
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	t.Setenv("HTTP_PORT", "0")
	_, err := config.Load()
	if err == nil {
		t.Fatal("expected error for invalid port")
	}
}

func TestParseProdEnvAliases(t *testing.T) {
	t.Setenv("APP_ENV", "production")
	t.Setenv("HTTP_PORT", "8080")
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Env != config.EnvProd {
		t.Fatalf("Env = %q, want prod", cfg.Env)
	}
}
