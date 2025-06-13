package recipe

import (
	"time"

	"github.com/google/uuid"
)

type Recipe struct {
	ID                uuid.UUID       `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	Title             string          `json:"title"`
	Description       string          `json:"description"`
	Ingredients       []Ingredient    `json:"ingredients"`
	Instructions      []string        `json:"instructions"`
	PrepTime          int             `json:"prep_time"`
	CookTime          int             `json:"cook_time"`
	Servings          int             `json:"servings"`
	Category          string          `json:"category"`
	DietaryCategories []string        `json:"dietary_categories"`
	Allergens         []string        `json:"allergens"`
	NutritionalInfo   NutritionalInfo `json:"nutritional_info"`
	ImageURL          *string         `json:"image_url,omitempty"`
	IsPublic          bool            `json:"is_public"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

type Ingredient struct {
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Unit     string  `json:"unit"`
	Optional bool    `json:"optional"`
}

type NutritionalInfo struct {
	Calories      int `json:"calories"`
	Protein       int `json:"protein"`
	Carbohydrates int `json:"carbohydrates"`
	Fat           int `json:"fat"`
	Fiber         int `json:"fiber"`
	Sugar         int `json:"sugar"`
	Sodium        int `json:"sodium"`
}
