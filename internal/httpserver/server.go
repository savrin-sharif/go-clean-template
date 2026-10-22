package httpserver

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/config"
)

// Server wraps http.Server with graceful shutdown.
type Server struct {
	httpServer *http.Server
	log        *slog.Logger
}

// New constructs an HTTP server from config and a configured Gin engine.
func New(cfg config.Config, log *slog.Logger, engine *gin.Engine) *Server {
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	return &Server{
		log: log,
		httpServer: &http.Server{
			Addr:         addr,
			Handler:      engine,
			ReadTimeout:  cfg.Server.ReadTimeout,
			WriteTimeout: cfg.Server.WriteTimeout,
			IdleTimeout:  cfg.Server.IdleTimeout,
		},
	}
}

// NewEngine builds a Gin engine with shared middleware.
func NewEngine(cfg config.Config, log *slog.Logger) *gin.Engine {
	if cfg.Env == config.EnvProd {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}

	engine := gin.New()
	engine.Use(
		RequestID(),
		Recover(log),
		AccessLog(log),
		cors.New(cors.Config{
			AllowOrigins:     cfg.Server.AllowedOrigins,
			AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
			AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", requestIDHeader},
			ExposeHeaders:    []string{requestIDHeader},
			AllowCredentials: false,
			MaxAge:           12 * time.Hour,
		}),
		AuthHook(),
	)
	return engine
}

// Start begins listening. It blocks until the server stops.
func (s *Server) Start() error {
	s.log.Info("http server listening", "addr", s.httpServer.Addr)
	err := s.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

// Shutdown gracefully stops the server.
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

// Addr returns the listen address.
func (s *Server) Addr() string {
	return s.httpServer.Addr
}
