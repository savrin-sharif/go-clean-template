package httpserver

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/apperr"
)

// Envelope is the standard JSON success body.
type Envelope struct {
	Data any `json:"data"`
}

// ErrorBody is the standard JSON error body.
type ErrorBody struct {
	Error ErrorPayload `json:"error"`
}

// ErrorPayload carries a user-safe error message and optional field errors.
type ErrorPayload struct {
	Kind        string            `json:"kind"`
	Message     string            `json:"message"`
	FieldErrors map[string]string `json:"field_errors,omitempty"`
}

// WriteJSON writes a success envelope.
func WriteJSON(c *gin.Context, status int, data any) {
	c.JSON(status, Envelope{Data: data})
}

// WriteError maps err to an HTTP JSON error response.
func WriteError(c *gin.Context, err error) {
	appErr := apperr.From(err)
	status := appErr.HTTPStatus()
	payload := ErrorPayload{
		Kind:        string(appErr.Kind),
		Message:     appErr.Error(),
		FieldErrors: appErr.FieldErrors,
	}
	// Never leak internal causes to clients.
	if appErr.Kind == apperr.KindInternal {
		payload.Message = "Something went wrong. Please try again."
	}
	c.AbortWithStatusJSON(status, ErrorBody{Error: payload})
}

// WriteNoContent writes 204.
func WriteNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}
