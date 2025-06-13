package recipe

import (
	"context"

	"github.com/google/uuid"
)

// SearchParams represents parameters for recipe search.
type SearchParams struct{}

// SearchResult is a placeholder for search results.
type SearchResult struct{}

// Repository defines persistence operations for recipes.
type Repository interface {
	Create(ctx context.Context, r *Recipe) error
	GetByID(ctx context.Context, id uuid.UUID) (*Recipe, error)
	Search(ctx context.Context, params SearchParams) (*SearchResult, error)
	GetUserFavorites(ctx context.Context, userID uuid.UUID) ([]*Recipe, error)
	AddFavorite(ctx context.Context, userID, recipeID uuid.UUID) error
	RemoveFavorite(ctx context.Context, userID, recipeID uuid.UUID) error
}
