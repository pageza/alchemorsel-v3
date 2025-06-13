package deepseek

import (
	"context"
	"net/http"

	"alchemorsel/backend/internal/domain/recipe"
)

// Client provides access to the DeepSeek API.
type Client struct {
	apiKey     string
	apiURL     string
	httpClient *http.Client
}

// GenerateRecipeRequest defines parameters for recipe generation.
type GenerateRecipeRequest struct {
	UserPreferences   any
	RecipeConstraints any
	Style             string
}

// GenerateRecipe generates a recipe using the DeepSeek API.
func (c *Client) GenerateRecipe(ctx context.Context, req GenerateRecipeRequest) (*recipe.Recipe, error) {
	// TODO: create prompt based on request
	// TODO: call DeepSeek API
	// TODO: parse response into recipe entity
	// TODO: generate embeddings
	return nil, nil
}
