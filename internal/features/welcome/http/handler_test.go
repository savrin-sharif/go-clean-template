package http_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/apperr"
	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/domain"
	welcomehttp "github.com/savrin-sharif/go-clean-template/internal/features/welcome/http"
)

type stubRepo struct {
	info domain.Info
	err  error
}

func (s stubRepo) GetInfo(context.Context) (domain.Info, error) {
	return s.info, s.err
}

func TestWelcomeHandlerSuccess(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uc := domain.NewGetWelcomeInfo(stubRepo{
		info: domain.Info{
			AppName:     "Demo",
			Environment: "local",
			Message:     "Welcome to Demo",
		},
	})
	h := welcomehttp.NewHandler(uc)

	engine := gin.New()
	h.Register(engine.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/welcome", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}

	var body struct {
		Data struct {
			AppName     string `json:"app_name"`
			Environment string `json:"environment"`
			Message     string `json:"message"`
		} `json:"data"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Data.AppName != "Demo" || body.Data.Message != "Welcome to Demo" {
		t.Fatalf("unexpected body: %#v", body.Data)
	}
}

func TestWelcomeHandlerMapsNotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	uc := domain.NewGetWelcomeInfo(stubRepo{err: apperr.NotFound("missing welcome")})
	h := welcomehttp.NewHandler(uc)

	engine := gin.New()
	h.Register(engine.Group("/api/v1"))

	req := httptest.NewRequest(http.MethodGet, "/api/v1/welcome", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
}
