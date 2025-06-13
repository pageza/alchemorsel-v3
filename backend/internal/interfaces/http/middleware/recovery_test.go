package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestRecovery ensures that the middleware converts panics to 500 errors.
func TestRecovery(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	r.Use(Recovery())
	r.GET("/panic", func(c *gin.Context) {
		panic("boom")
	})

	req := httptest.NewRequest(http.MethodGet, "/panic", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
	if expected := "internal server error"; !strings.Contains(w.Body.String(), expected) {
		t.Fatalf("expected body to contain %q, got %s", expected, w.Body.String())
	}
}
