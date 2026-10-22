// Package data implements welcome repository contracts.
package data

import (
	"context"

	"github.com/savrin-sharif/go-clean-template/internal/config"
	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/domain"
)

// LocalRepository serves welcome info from process config (no database).
type LocalRepository struct {
	appName string
	env     config.Environment
}

// NewLocalRepository constructs a config-backed welcome repository.
func NewLocalRepository(cfg config.Config) *LocalRepository {
	return &LocalRepository{
		appName: cfg.AppName,
		env:     cfg.Env,
	}
}

// GetInfo returns welcome info derived from configuration.
func (r *LocalRepository) GetInfo(_ context.Context) (domain.Info, error) {
	return domain.Info{
		AppName:     r.appName,
		Environment: string(r.env),
		Message:     "Welcome to " + r.appName,
	}, nil
}
