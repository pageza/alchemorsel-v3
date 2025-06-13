package middleware

import "github.com/gin-gonic/gin"

// Auth extracts the Authorization header and stores user info in the context.
// In this placeholder implementation the header value itself is stored.
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		if auth := c.GetHeader("Authorization"); auth != "" {
			c.Set("user", auth)
		}
		c.Next()
	}
}
