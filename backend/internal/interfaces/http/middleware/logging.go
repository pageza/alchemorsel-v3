package middleware

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

// Logging logs request method, path and duration.
func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		method := c.Request.Method
		path := c.Request.URL.Path

		c.Next()

		duration := time.Since(start)
		log.Printf("%s %s %v", method, path, duration)
	}
}
