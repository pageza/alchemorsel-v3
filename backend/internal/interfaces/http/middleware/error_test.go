package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	apperrors "alchemorsel/backend/internal/pkg/errors"
	"alchemorsel/backend/internal/pkg/logger"
)

func TestErrorHandler_AppError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, recorded := observer.New(zap.ErrorLevel)
	logger.SetLogger(zap.New(core).Sugar())

	r := gin.New()
	r.Use(ErrorHandler())
	r.GET("/fail", func(c *gin.Context) {
		c.Error(apperrors.ErrInvalidInput)
	})

	req := httptest.NewRequest(http.MethodGet, "/fail", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), apperrors.ErrInvalidInput.Code) {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}

	if recorded.Len() == 0 {
		t.Fatalf("no error logged")
	}
}

func TestErrorHandler_GenericError(t *testing.T) {
	gin.SetMode(gin.TestMode)

	core, recorded := observer.New(zap.ErrorLevel)
	logger.SetLogger(zap.New(core).Sugar())

	r := gin.New()
	r.Use(ErrorHandler())
	r.GET("/boom", func(c *gin.Context) {
		c.Error(errors.New("boom"))
	})

	req := httptest.NewRequest(http.MethodGet, "/boom", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "internal server error") {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}

	if recorded.Len() == 0 {
		t.Fatalf("no error logged")
	}
}
