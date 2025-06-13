package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperrors "alchemorsel/backend/internal/pkg/errors"
	"alchemorsel/backend/internal/pkg/logger"
)

// ErrorHandler converts returned errors into JSON HTTP responses.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) == 0 {
			return
		}
		err := c.Errors.Last().Err
		if appErr, ok := err.(*apperrors.AppError); ok {
			logger.FromContext(c.Request.Context()).Errorw(appErr.Message, "code", appErr.Code)
			resp := gin.H{"error": appErr.Code, "message": appErr.Message}
			if appErr.Details != nil {
				resp["details"] = appErr.Details
			}
			c.AbortWithStatusJSON(appErr.Status, resp)
			return
		}
		logger.FromContext(c.Request.Context()).Errorf("unhandled error: %v", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
