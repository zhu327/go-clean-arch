package middleware

import (
	"strings"

	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

// AuthMiddleware validates the Bearer token and sets user_id in the Gin context.
func AuthMiddleware(tokenService auth.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.Error(utils.UnauthorizedError("missing authorization header"))
			c.Abort()
			return
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			c.Error(utils.UnauthorizedError("invalid authorization header format"))
			c.Abort()
			return
		}

		claims, err := tokenService.ValidateToken(parts[1])
		if err != nil {
			c.Error(utils.UnauthorizedError("invalid or expired token"))
			c.Abort()
			return
		}

		c.Set(UserIDKey, claims.UserID)
		c.Next()
	}
}
