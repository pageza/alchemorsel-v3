package middleware

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)
	var buf bytes.Buffer
	orig := log.Writer()
	log.SetOutput(&buf)
	log.SetFlags(0)
	defer log.SetOutput(orig)

	r := gin.New()
	r.Use(Logging())
	r.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	out := buf.String()
	if !strings.Contains(out, "GET /ping") {
		t.Fatalf("expected log to contain GET /ping, got %s", out)
	}
}
