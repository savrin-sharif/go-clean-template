// Package logging configures slog and redacts sensitive values.
package logging

import (
	"log/slog"
	"os"
	"strings"

	"github.com/savrin-sharif/go-clean-template/internal/config"
)

// New builds a slog logger from config.
func New(cfg config.LogConfig) *slog.Logger {
	level := parseLevel(cfg.Level)
	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: cfg.AddSource,
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if shouldRedact(a.Key) {
				return slog.String(a.Key, "[REDACTED]")
			}
			return a
		},
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	return slog.New(handler)
}

func parseLevel(raw string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func shouldRedact(key string) bool {
	k := strings.ToLower(key)
	switch {
	case strings.Contains(k, "password"),
		strings.Contains(k, "secret"),
		strings.Contains(k, "token"),
		strings.Contains(k, "apikey"),
		strings.Contains(k, "api_key"),
		strings.Contains(k, "authorization"),
		strings.Contains(k, "cookie"):
		return true
	default:
		return false
	}
}

// RedactMap returns a copy of m with sensitive keys redacted.
func RedactMap(m map[string]string) map[string]string {
	if m == nil {
		return nil
	}
	out := make(map[string]string, len(m))
	for k, v := range m {
		if shouldRedact(k) {
			out[k] = "[REDACTED]"
			continue
		}
		out[k] = v
	}
	return out
}
