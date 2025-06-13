package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"alchemorsel/backend/internal/pkg/logger"
)

func TestLogging(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, recorded := observer.New(zap.InfoLevel)
	logger.SetLogger(zap.New(core).Sugar())

	r := gin.New()
	r.Use(Logging())
	r.GET("/ping", func(c *gin.Context) { c.Status(http.StatusOK) })

	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if recorded.Len() == 0 {
		t.Fatalf("no log messages recorded")
	}
	entry := recorded.All()[0]
	if entry.ContextMap()["method"] != "GET" || entry.ContextMap()["path"] != "/ping" {
		t.Fatalf("unexpected log fields: %v", entry.Context)
	}
}
