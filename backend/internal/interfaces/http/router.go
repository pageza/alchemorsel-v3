package http

import (
	"github.com/gin-gonic/gin"

	"alchemorsel/backend/internal/interfaces/http/handlers"
	"alchemorsel/backend/internal/interfaces/http/middleware"
)

// SetupRouter configures all HTTP routes following the design docs.
func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(
		middleware.Recovery(),
		middleware.RequestID(),
		middleware.Logging(),
		middleware.CORS(),
	)

	api := r.Group("/api/v1")
	{
		api.GET("/health", handlers.Health)

		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
			auth.POST("/refresh", handlers.RefreshToken)
		}

		// Protected routes require authentication middleware.
		protected := api.Group("")
		protected.Use(middleware.Auth())
		{
			users := protected.Group("/users")
			{
				users.GET("/profile", handlers.GetProfile)
				users.PUT("/profile", handlers.UpdateProfile)
				users.POST("/profile/picture", handlers.UploadProfilePicture)
			}

			recipes := protected.Group("/recipes")
			{
				recipes.GET("/", handlers.SearchRecipes)
				recipes.POST("/", handlers.CreateRecipe)
				recipes.GET("/:id", handlers.GetRecipe)
				recipes.POST("/:id/favorite", handlers.AddFavorite)
				recipes.DELETE("/:id/favorite", handlers.RemoveFavorite)
			}

			protected.POST("/llm/generate", handlers.GenerateRecipe)
		}
	}

	return r
}
