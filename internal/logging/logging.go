package logging

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func DetailedLoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {

		path := c.Request.URL.Path

		// 🔍 Bypass logging untuk health check
		if path == "/health" {
			c.Next()
			return
		}

		// Eksekusi request
		c.Next()

	}
}
