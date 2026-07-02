package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	apperrors "github.com/refda/backend/internal/pkg/errors"
	jwtpkg "github.com/refda/backend/internal/pkg/jwt"
	"github.com/refda/backend/internal/pkg/response"
)

const UserIDKey = "user_id"
const PhoneKey = "phone"

func Auth(jwtMgr *jwtpkg.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			response.Fail(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := jwtMgr.Parse(token)
		if err != nil {
			response.Fail(c, apperrors.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Set(PhoneKey, claims.Phone)
		c.Next()
	}
}

func GetUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(UserIDKey)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
