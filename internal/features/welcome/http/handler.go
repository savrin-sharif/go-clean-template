// Package http exposes welcome HTTP handlers.
package http

import (
	nethttp "net/http"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/internal/features/welcome/domain"
	"github.com/savrin-sharif/go-clean-template/internal/httpserver"
)

// Handler serves welcome endpoints.
type Handler struct {
	getWelcome *domain.GetWelcomeInfo
}

// NewHandler constructs the welcome HTTP handler.
func NewHandler(getWelcome *domain.GetWelcomeInfo) *Handler {
	return &Handler{getWelcome: getWelcome}
}

type welcomeResponse struct {
	AppName     string `json:"app_name"`
	Environment string `json:"environment"`
	Message     string `json:"message"`
}

// Get handles GET /api/v1/welcome.
func (h *Handler) Get(c *gin.Context) {
	info, err := h.getWelcome.Execute(c.Request.Context())
	if err != nil {
		httpserver.WriteError(c, err)
		return
	}
	httpserver.WriteJSON(c, nethttp.StatusOK, welcomeResponse{
		AppName:     info.AppName,
		Environment: info.Environment,
		Message:     info.Message,
	})
}

// Register mounts welcome routes on the given router group.
func (h *Handler) Register(rg *gin.RouterGroup) {
	rg.GET("/welcome", h.Get)
}
