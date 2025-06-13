package handlers

import "github.com/gin-gonic/gin"

// Health responds with a simple health status.
func Health(c *gin.Context) {
	c.JSON(200, gin.H{"status": "healthy"})
}
