package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

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
func (c *Client) GenerateRecipe(ctx context.Context, req GenerateRecipeRequest) (*recipe.Recipe, []float64, error) {
	prompt, err := buildPrompt(req)
	if err != nil {
		return nil, nil, err
	}

	payloadBytes, err := json.Marshal(map[string]string{"prompt": prompt})
	if err != nil {
		return nil, nil, err
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, strings.TrimRight(c.apiURL, "/")+"/generate", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, nil, fmt.Errorf("deepseek: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var dsResp struct {
		Recipe    recipe.Recipe `json:"recipe"`
		Embedding []float64     `json:"embedding"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&dsResp); err != nil {
		return nil, nil, err
	}

	return &dsResp.Recipe, dsResp.Embedding, nil
}

func buildPrompt(req GenerateRecipeRequest) (string, error) {
	b := &strings.Builder{}
	if req.Style != "" {
		b.WriteString("Style: ")
		b.WriteString(req.Style)
		b.WriteString(". ")
	}
	if req.UserPreferences != nil {
		prefs, err := json.Marshal(req.UserPreferences)
		if err != nil {
			return "", err
		}
		b.WriteString("Preferences: ")
		b.Write(prefs)
		b.WriteString(". ")
	}
	if req.RecipeConstraints != nil {
		cons, err := json.Marshal(req.RecipeConstraints)
		if err != nil {
			return "", err
		}
		b.WriteString("Constraints: ")
		b.Write(cons)
		b.WriteString(". ")
	}
	b.WriteString("Generate a recipe.")
	return b.String(), nil
}
