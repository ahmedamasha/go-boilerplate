package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/refda/backend/internal/pkg/jwt"
)

// OptionalAuth sets user_id when a valid Bearer token is present; does not abort otherwise.
func OptionalAuth(jwtMgr *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.Next()
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			c.Next()
			return
		}
		c.Set(UserIDKey, claims.UserID)
		c.Set(PhoneKey, claims.Phone)
		c.Next()
	}
}
