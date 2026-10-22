package data

import (
	"context"
	"database/sql"

	"github.com/savrin-sharif/go-clean-template/internal/apperr"
	"github.com/savrin-sharif/go-clean-template/internal/config"
	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/domain"
)

// PostgresRepository is an optional welcome store backed by Postgres.
// The default template uses LocalRepository; wire this when DATABASE_URL is set
// and you want a DB-backed welcome row.
type PostgresRepository struct {
	db      *sql.DB
	appName string
	env     config.Environment
}

// NewPostgresRepository constructs a Postgres-backed welcome repository.
// It does not require a migrations table for the sample query; if no row exists
// it falls back to config-derived values.
func NewPostgresRepository(db *sql.DB, cfg config.Config) *PostgresRepository {
	return &PostgresRepository{
		db:      db,
		appName: cfg.AppName,
		env:     cfg.Env,
	}
}

// GetInfo tries to read a welcome message from the database, falling back to config.
func (r *PostgresRepository) GetInfo(ctx context.Context) (domain.Info, error) {
	var message string
	err := r.db.QueryRowContext(ctx, `
		SELECT message FROM welcome_messages
		WHERE active = true
		ORDER BY id ASC
		LIMIT 1
	`).Scan(&message)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.Info{
				AppName:     r.appName,
				Environment: string(r.env),
				Message:     "Welcome to " + r.appName,
			}, nil
		}
		// Missing relation or other DB errors: treat as data/infrastructure failure.
		return domain.Info{}, apperr.WrapInternal(err, "welcome query failed")
	}
	return domain.Info{
		AppName:     r.appName,
		Environment: string(r.env),
		Message:     message,
	}, nil
}
