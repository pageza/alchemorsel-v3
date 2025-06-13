package handlers

import "github.com/gin-gonic/gin"

// SearchRecipes searches for recipes.
func SearchRecipes(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// CreateRecipe creates a new recipe.
func CreateRecipe(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// GetRecipe gets a recipe by ID.
func GetRecipe(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// AddFavorite marks a recipe as favorite.
func AddFavorite(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// RemoveFavorite unmarks a recipe as favorite.
func RemoveFavorite(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}

// GenerateRecipe generates a recipe via LLM.
func GenerateRecipe(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
