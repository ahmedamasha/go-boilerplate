package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/refda/backend/internal/pkg/config"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	"github.com/refda/backend/internal/pkg/response"
)

func AdminAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-Admin-Key")
		if cfg.Admin.APIKey == "" || key != cfg.Admin.APIKey {
			response.Fail(c, apperrors.ErrInvalidAdminKey)
			c.Abort()
			return
		}
		c.Next()
	}
}
