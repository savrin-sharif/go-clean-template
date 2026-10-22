package data_test

import (
	"context"
	"testing"

	"github.com/savrin-sharif/go-clean-template/internal/config"
	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/data"
)

func TestLocalRepositoryGetInfo(t *testing.T) {
	t.Parallel()

	repo := data.NewLocalRepository(config.Config{
		AppName: "Demo API",
		Env:     config.EnvLocal,
	})
	info, err := repo.GetInfo(context.Background())
	if err != nil {
		t.Fatalf("GetInfo() error = %v", err)
	}
	if info.AppName != "Demo API" {
		t.Fatalf("AppName = %q", info.AppName)
	}
	if info.Environment != "local" {
		t.Fatalf("Environment = %q", info.Environment)
	}
	if info.Message != "Welcome to Demo API" {
		t.Fatalf("Message = %q", info.Message)
	}
}
