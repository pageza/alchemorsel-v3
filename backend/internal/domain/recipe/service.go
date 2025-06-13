package recipe

import (
	"context"

	"github.com/google/uuid"
)

// Service defines business logic for recipes.
type Service interface {
	Create(ctx context.Context, req CreateRequest) (*Recipe, error)
	Get(ctx context.Context, id uuid.UUID) (*Recipe, error)
	Search(ctx context.Context, params SearchParams) (*SearchResult, error)
	Generate(ctx context.Context, req GenerateRequest) (*Recipe, error)
	AddFavorite(ctx context.Context, userID, recipeID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, recipeID uuid.UUID) error
	GetFavorites(ctx context.Context, userID uuid.UUID) ([]*Recipe, error)
}

type CreateRequest struct {
	Title             string
	Description       string
	Ingredients       []Ingredient
	Instructions      []string
	PrepTime          int
	CookTime          int
	Servings          int
	Category          string
	DietaryCategories []string
	Allergens         []string
}

type GenerateRequest struct {
	Style        string
	Ingredients  []string
	CookingTime  int
	Servings     int
	CustomPrompt string
}
