// Package http exposes health HTTP handlers.
package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/httpserver"
)

// Handler serves health endpoints.
type Handler struct{}

// NewHandler constructs the health handler.
func NewHandler() *Handler {
	return &Handler{}
}

type healthResponse struct {
	Status string `json:"status"`
}

// Live handles GET /health.
func (h *Handler) Live(c *gin.Context) {
	httpserver.WriteJSON(c, nethttp.StatusOK, healthResponse{Status: "ok"})
}

// Register mounts health routes on the engine root.
func (h *Handler) Register(engine *gin.Engine) {
	engine.GET("/health", h.Live)
}
