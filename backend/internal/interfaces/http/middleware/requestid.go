package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"alchemorsel/backend/internal/pkg/logger"
)

// RequestIDKey is the gin context key storing the request id.
const RequestIDKey = "request_id"

// RequestID injects a unique request ID into the context and response headers.
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.New().String()
		c.Set(RequestIDKey, id)
		c.Writer.Header().Set("X-Request-ID", id)

		l := logger.Logger().With("request_id", id)
		c.Request = c.Request.WithContext(logger.ToContext(c.Request.Context(), l))

		c.Next()
	}
}
