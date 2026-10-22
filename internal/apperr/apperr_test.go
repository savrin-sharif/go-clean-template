package apperr_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/savrin-sharif/go-clean-template/internal/apperr"
)

func TestHTTPStatus(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		err    *apperr.Error
		status int
	}{
		{"validation", apperr.Validation("bad", map[string]string{"email": "required"}), http.StatusUnprocessableEntity},
		{"unauthorized", apperr.Unauthorized(""), http.StatusUnauthorized},
		{"forbidden", apperr.Forbidden(""), http.StatusForbidden},
		{"not_found", apperr.NotFound(""), http.StatusNotFound},
		{"conflict", apperr.Conflict(""), http.StatusConflict},
		{"data", apperr.Data(""), http.StatusBadRequest},
		{"internal", apperr.Internal(""), http.StatusInternalServerError},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := tc.err.HTTPStatus(); got != tc.status {
				t.Fatalf("HTTPStatus() = %d, want %d", got, tc.status)
			}
		})
	}
}

func TestFromWrapsUnknown(t *testing.T) {
	t.Parallel()

	err := errors.New("boom")
	appErr := apperr.From(err)
	if appErr.Kind != apperr.KindInternal {
		t.Fatalf("Kind = %s, want %s", appErr.Kind, apperr.KindInternal)
	}
	if !errors.Is(appErr, err) {
		t.Fatalf("expected cause to unwrap to original error")
	}
}

func TestAsExtractsTyped(t *testing.T) {
	t.Parallel()

	original := apperr.NotFound("missing")
	got, ok := apperr.As(original)
	if !ok {
		t.Fatal("As() returned false")
	}
	if got.Kind != apperr.KindNotFound {
		t.Fatalf("Kind = %s, want %s", got.Kind, apperr.KindNotFound)
	}
}
