package handlers

import "github.com/gin-gonic/gin"

// GetProfile returns the current user profile.
func GetProfile(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// UpdateProfile updates the current user profile.
func UpdateProfile(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// UploadProfilePicture uploads a profile picture.
func UploadProfilePicture(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
