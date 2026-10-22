package httpserver

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const requestIDHeader = "X-Request-ID"

// RequestID ensures every request has a correlation id.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.GetHeader(requestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		c.Writer.Header().Set(requestIDHeader, id)
		c.Set("request_id", id)
		c.Next()
	}
}

// Recover logs panics and returns a safe 500.
func Recover(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error("panic recovered",
					"panic", rec,
					"request_id", c.GetString("request_id"),
					"path", c.Request.URL.Path,
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, ErrorBody{
					Error: ErrorPayload{
						Kind:    "internal",
						Message: "Something went wrong. Please try again.",
					},
				})
			}
		}()
		c.Next()
	}
}

// AccessLog logs method, path, status, and latency.
func AccessLog(log *slog.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		log.Info("http request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", c.Writer.Status(),
			"latency_ms", time.Since(start).Milliseconds(),
			"request_id", c.GetString("request_id"),
			"client_ip", c.ClientIP(),
		)
	}
}

// AuthHook is a no-op placeholder for authentication middleware.
// Replace this with JWT/session validation when the API requires auth.
func AuthHook() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}
