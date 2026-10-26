// Package http serves the HTML welcome page and interactive OpenAPI docs.
package http

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/savrin-sharif/go-clean-template/api"
	"github.com/savrin-sharif/go-clean-template/internal/config"
)

//go:embed landing.html swagger.html
var templatesFS embed.FS

// Handler serves landing and docs pages.
type Handler struct {
	cfg         config.Config
	landingTmpl *template.Template
	swaggerTmpl *template.Template
}

type pageData struct {
	AppName     string
	Environment string
}

// NewHandler constructs the docs HTTP handler.
func NewHandler(cfg config.Config) *Handler {
	landing := template.Must(template.ParseFS(templatesFS, "landing.html"))
	swagger := template.Must(template.ParseFS(templatesFS, "swagger.html"))
	return &Handler{
		cfg:         cfg,
		landingTmpl: landing,
		swaggerTmpl: swagger,
	}
}

// Landing handles GET /.
func (h *Handler) Landing(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	_ = h.landingTmpl.Execute(c.Writer, pageData{
		AppName:     h.cfg.AppName,
		Environment: string(h.cfg.Env),
	})
}

// Docs handles GET /docs and GET /docs/.
func (h *Handler) Docs(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Status(http.StatusOK)
	_ = h.swaggerTmpl.Execute(c.Writer, pageData{
		AppName:     h.cfg.AppName,
		Environment: string(h.cfg.Env),
	})
}

// OpenAPI handles GET /openapi.yaml.
func (h *Handler) OpenAPI(c *gin.Context) {
	c.Header("Content-Type", "application/yaml; charset=utf-8")
	c.Data(http.StatusOK, "application/yaml; charset=utf-8", api.OpenAPI)
}

// Register mounts docs routes on the engine root.
func (h *Handler) Register(engine *gin.Engine) {
	engine.GET("/", h.Landing)
	engine.GET("/docs", h.Docs)
	engine.GET("/docs/", h.Docs)
	engine.GET("/openapi.yaml", h.OpenAPI)
}
