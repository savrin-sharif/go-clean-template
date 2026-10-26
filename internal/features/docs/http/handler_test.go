package http_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/config"
	docshttp "github.com/savrin-sharif/go-clean-template/internal/features/docs/http"
)

func newTestEngine(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)

	h := docshttp.NewHandler(config.Config{
		AppName: "Demo API",
		Env:     config.EnvLocal,
	})
	engine := gin.New()
	h.Register(engine)
	return engine
}

func TestLandingReturnsHTML(t *testing.T) {
	engine := newTestEngine(t)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", ct)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "Demo API") {
		t.Fatalf("body missing app name: %s", body)
	}
	if !strings.Contains(body, "/docs") {
		t.Fatalf("body missing docs link: %s", body)
	}
	if !strings.Contains(body, "Built with Go") {
		t.Fatalf("body missing Built with Go caption: %s", body)
	}
}

func TestDocsReturnsHTML(t *testing.T) {
	engine := newTestEngine(t)

	req := httptest.NewRequest(http.MethodGet, "/docs", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "text/html") {
		t.Fatalf("Content-Type = %q, want text/html", ct)
	}
	if !strings.Contains(rec.Body.String(), "/openapi.yaml") {
		t.Fatalf("docs page missing openapi url: %s", rec.Body.String())
	}
}

func TestOpenAPIReturnsYAML(t *testing.T) {
	engine := newTestEngine(t)

	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body.String())
	}
	ct := rec.Header().Get("Content-Type")
	if !strings.Contains(ct, "yaml") {
		t.Fatalf("Content-Type = %q, want yaml", ct)
	}
	body := rec.Body.String()
	if !strings.HasPrefix(strings.TrimSpace(body), "openapi:") {
		t.Fatalf("body does not look like OpenAPI YAML: %s", body[:min(80, len(body))])
	}
}
