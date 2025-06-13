package deepseek

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGenerateRecipe(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/generate" {
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
		if auth := r.Header.Get("Authorization"); auth != "Bearer test" {
			t.Fatalf("unexpected auth header: %s", auth)
		}
		var payload struct {
			Prompt string `json:"prompt"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode payload: %v", err)
		}
		if !strings.Contains(payload.Prompt, "Style: italian") {
			t.Fatalf("prompt missing style: %s", payload.Prompt)
		}
		resp := map[string]any{
			"recipe":    map[string]any{"title": "Caprese"},
			"embedding": []float64{1, 2, 3},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer ts.Close()

	c := Client{apiKey: "test", apiURL: ts.URL, httpClient: ts.Client()}
	rec, emb, err := c.GenerateRecipe(context.Background(), GenerateRecipeRequest{Style: "italian"})
	if err != nil {
		t.Fatalf("GenerateRecipe returned error: %v", err)
	}
	if rec.Title != "Caprese" {
		t.Fatalf("unexpected recipe title: %s", rec.Title)
	}
	if len(emb) != 3 || emb[0] != 1 {
		t.Fatalf("unexpected embeddings: %v", emb)
	}
}

func TestGenerateRecipe_HTTPError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "bad", http.StatusBadRequest)
	}))
	defer ts.Close()

	c := Client{apiKey: "k", apiURL: ts.URL, httpClient: ts.Client()}
	if _, _, err := c.GenerateRecipe(context.Background(), GenerateRecipeRequest{}); err == nil {
		t.Fatalf("expected error")
	}
}
