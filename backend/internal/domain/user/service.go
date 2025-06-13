package user

import (
	"context"
	"io"

	"github.com/google/uuid"
)

// Service defines business logic for user management.
type Service interface {
	Register(ctx context.Context, req RegisterRequest) (*User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req UpdateRequest) error
	UploadProfilePicture(ctx context.Context, userID uuid.UUID, file io.Reader) error
}

// RegisterRequest represents input for user registration.
type RegisterRequest struct {
	Email              string
	Username           string
	Password           string
	Name               string
	DietaryPreferences []string
	Allergies          []string
}

// UpdateRequest represents profile update data.
type UpdateRequest struct {
	Name               string
	DietaryPreferences []string
	Allergies          []string
}
