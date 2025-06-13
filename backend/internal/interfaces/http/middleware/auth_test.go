package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuth(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Auth())
	r.GET("/", func(c *gin.Context) {
		if val, exists := c.Get("user"); exists {
			c.String(http.StatusOK, val.(string))
		} else {
			c.Status(http.StatusUnauthorized)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set("Authorization", "token123")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if body := w.Body.String(); body != "token123" {
		t.Fatalf("expected body token123, got %s", body)
	}
}
