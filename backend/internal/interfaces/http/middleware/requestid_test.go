package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"alchemorsel/backend/internal/pkg/logger"
)

func TestRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger.SetLogger(zap.NewNop().Sugar())

	r := gin.New()
	r.Use(RequestID())
	r.GET("/", func(c *gin.Context) {
		id := c.GetString(RequestIDKey)
		if id == "" {
			t.Errorf("request id missing from context")
		}
		l := logger.FromContext(c.Request.Context())
		if l == nil {
			t.Errorf("logger missing from context")
		}
		c.String(http.StatusOK, id)
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	id := w.Header().Get("X-Request-ID")
	if id == "" {
		t.Fatalf("X-Request-ID header missing")
	}
	if w.Body.String() != id {
		t.Fatalf("handler saw id %s, header %s", w.Body.String(), id)
	}

	// second request should receive a different id
	req2 := httptest.NewRequest(http.MethodGet, "/", nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)
	id2 := w2.Header().Get("X-Request-ID")
	if id2 == id {
		t.Fatalf("expected unique request id, got %s twice", id)
	}
}
