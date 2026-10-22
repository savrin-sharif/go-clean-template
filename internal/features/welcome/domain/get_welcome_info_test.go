package domain_test

import (
	"context"
	"errors"
	"testing"

	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/domain"
)

type stubRepo struct {
	info domain.Info
	err  error
}

func (s stubRepo) GetInfo(context.Context) (domain.Info, error) {
	return s.info, s.err
}

func TestGetWelcomeInfoSuccess(t *testing.T) {
	t.Parallel()

	uc := domain.NewGetWelcomeInfo(stubRepo{
		info: domain.Info{AppName: "Demo", Environment: "local", Message: "hi"},
	})
	got, err := uc.Execute(context.Background())
	if err != nil {
		t.Fatalf("Execute() error = %v", err)
	}
	if got.AppName != "Demo" || got.Message != "hi" {
		t.Fatalf("unexpected info: %#v", got)
	}
}

func TestGetWelcomeInfoPropagatesError(t *testing.T) {
	t.Parallel()

	want := errors.New("repo down")
	uc := domain.NewGetWelcomeInfo(stubRepo{err: want})
	_, err := uc.Execute(context.Background())
	if !errors.Is(err, want) {
		t.Fatalf("Execute() error = %v, want %v", err, want)
	}
}
