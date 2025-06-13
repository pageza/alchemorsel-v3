package handlers

import "github.com/gin-gonic/gin"

// Register handles user registration.
func Register(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// Login handles user login.
func Login(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// RefreshToken refreshes JWT tokens.
func RefreshToken(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
