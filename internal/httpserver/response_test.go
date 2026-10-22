package httpserver_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/apperr"
	"github.com/savrin-sharif/go-clean-template/internal/httpserver"
)

func TestWriteErrorHidesInternalCause(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/boom", func(c *gin.Context) {
		httpserver.WriteError(c, apperr.WrapInternal(assertErr{}, "secret internals"))
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d", rec.Code)
	}

	var body httpserver.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error.Message == "secret internals" {
		t.Fatal("internal message leaked to client")
	}
	if body.Error.Kind != string(apperr.KindInternal) {
		t.Fatalf("kind = %s", body.Error.Kind)
	}
}

func TestWriteErrorValidationFields(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.GET("/bad", func(c *gin.Context) {
		httpserver.WriteError(c, apperr.Validation("invalid", map[string]string{"name": "required"}))
	})

	req := httptest.NewRequest(http.MethodGet, "/bad", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d", rec.Code)
	}

	var body httpserver.ErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if body.Error.FieldErrors["name"] != "required" {
		t.Fatalf("field_errors = %#v", body.Error.FieldErrors)
	}
}

type assertErr struct{}

func (assertErr) Error() string { return "cause" }
