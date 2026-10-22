// Package app is the composition root for the HTTP service.
package app

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"github.com/savrin-sharif/go-clean-template/internal/config"
	healthhttp "github.com/savrin-sharif/go-clean-template/internal/features/health/http"
	welcomedata "github.com/savrin-sharif/go-clean-template/internal/features/welcome/data"
	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/domain"
	welcomehttp "github.com/savrin-sharif/go-clean-template/internal/features/welcome/http"
	"github.com/savrin-sharif/go-clean-template/internal/httpserver"
	"github.com/savrin-sharif/go-clean-template/internal/logging"
)

// App holds running process dependencies.
type App struct {
	Config config.Config
	Log    *slog.Logger
	Server *httpserver.Server
	DB     *sql.DB
}

// New wires configuration, infrastructure, features, and the HTTP server.
func New(cfg config.Config) (*App, error) {
	log := logging.New(cfg.Log)

	var db *sql.DB
	if cfg.HasDatabase() {
		opened, err := openDatabase(cfg)
		if err != nil {
			return nil, err
		}
		db = opened
	}

	welcomeRepo := buildWelcomeRepository(cfg, db)
	getWelcome := domain.NewGetWelcomeInfo(welcomeRepo)

	engine := httpserver.NewEngine(cfg, log)
	registerRoutes(engine, getWelcome)

	server := httpserver.New(cfg, log, engine)

	return &App{
		Config: cfg,
		Log:    log,
		Server: server,
		DB:     db,
	}, nil
}

func buildWelcomeRepository(cfg config.Config, db *sql.DB) domain.Repository {
	if db != nil {
		return welcomedata.NewPostgresRepository(db, cfg)
	}
	return welcomedata.NewLocalRepository(cfg)
}

func registerRoutes(engine *gin.Engine, getWelcome *domain.GetWelcomeInfo) {
	healthhttp.NewHandler().Register(engine)

	api := engine.Group("/api/v1")
	welcomehttp.NewHandler(getWelcome).Register(api)
}

func openDatabase(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.Database.URL)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}
	db.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	db.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.Database.ConnMaxLife)

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ReadTimeout)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping database: %w", err)
	}
	return db, nil
}

// Close releases infrastructure resources.
func (a *App) Close() error {
	if a.DB != nil {
		return a.DB.Close()
	}
	return nil
}
