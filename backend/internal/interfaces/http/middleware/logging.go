package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"alchemorsel/backend/internal/pkg/logger"
)

// Logging logs request method, path and duration.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		duration := time.Since(start)

		l := logger.FromContext(c.Request.Context())
		l.Infow("request completed",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"duration", duration.String(),
		)
	}
}
