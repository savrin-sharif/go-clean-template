package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/savrin-sharif/go-clean-template/internal/app"
	"github.com/savrin-sharif/go-clean-template/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	application, err := app.New(cfg)
	if err != nil {
		log.Fatalf("app: %v", err)
	}
	defer func() {
		if cerr := application.Close(); cerr != nil {
			application.Log.Error("close resources", "error", cerr)
		}
	}()

	errCh := make(chan error, 1)
	go func() {
		errCh <- application.Server.Start()
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-sigCh:
		application.Log.Info("shutdown signal received", "signal", sig.String())
	case err := <-errCh:
		if err != nil && !errors.Is(err, context.Canceled) {
			application.Log.Error("server stopped", "error", err)
			os.Exit(1)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()
	if err := application.Server.Shutdown(ctx); err != nil {
		application.Log.Error("graceful shutdown failed", "error", err)
		os.Exit(1)
	}
	application.Log.Info("server stopped", "shutdown_timeout", cfg.Server.ShutdownTimeout.String())

	// Drain a tiny moment so buffered logs flush in tests/CI.
	time.Sleep(10 * time.Millisecond)
}
