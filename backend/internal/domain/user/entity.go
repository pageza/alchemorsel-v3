package user

import (
	"time"

	"github.com/google/uuid"
)

// User represents an application user.
type User struct {
	ID                 uuid.UUID  `json:"id"`
	Email              string     `json:"email"`
	Username           string     `json:"username"`
	PasswordHash       string     `json:"-"`
	Name               string     `json:"name"`
	ProfilePictureURL  *string    `json:"profile_picture_url,omitempty"`
	DietaryPreferences []string   `json:"dietary_preferences"`
	Allergies          []string   `json:"allergies"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
	DeletedAt          *time.Time `json:"deleted_at,omitempty"`
}
