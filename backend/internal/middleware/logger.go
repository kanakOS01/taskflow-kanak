package middleware

import (
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// Process request
		c.Next()

		end := time.Now()
		latency := end.Sub(start)
		status := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		if len(c.Errors) > 0 {
			for _, e := range c.Errors {
				slog.Error("request error",
					"method", method,
					"path", path,
					"query", query,
					"status", status,
					"ip", clientIP,
					"latency", latency,
					"error", e.Error(),
				)
			}
		} else {
			slog.Info("request",
				"method", method,
				"path", path,
				"query", query,
				"status", status,
				"ip", clientIP,
				"latency", latency,
			)
		}
	}
}
