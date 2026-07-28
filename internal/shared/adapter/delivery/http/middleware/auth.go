package middleware

import (
	"strings"

	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "user_id"

// AccessTokenVerifier is the minimal port required by protected HTTP routes.
type AccessTokenVerifier interface {
	ValidateAccessToken(string) (uint, error)
}

type accessTokenVerifier struct{ service auth.TypedTokenService }

func (v accessTokenVerifier) ValidateAccessToken(token string) (uint, error) {
	claims, err := v.service.ValidateTokenOfType(token, auth.AccessToken)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// NewAccessTokenVerifier adapts JWT validation to the protected-route port.
func NewAccessTokenVerifier(service auth.TypedTokenService) AccessTokenVerifier {
	return accessTokenVerifier{service: service}
}

// AuthMiddleware validates only access Bearer tokens and sets user_id in Gin context.
func AuthMiddleware(verifier AccessTokenVerifier) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		parts := strings.SplitN(header, " ", 2)
		if header == "" {
			_ = c.Error(utils.UnauthorizedError("missing authorization header"))
			c.Abort()
			return
		}
		if len(parts) != 2 || !strings.EqualFold(parts[0], "bearer") {
			_ = c.Error(utils.UnauthorizedError("invalid authorization header format"))
			c.Abort()
			return
		}
		userID, err := verifier.ValidateAccessToken(parts[1])
		if err != nil {
			_ = c.Error(utils.UnauthorizedError("invalid or expired access token"))
			c.Abort()
			return
		}
		if userID == 0 {
			_ = c.Error(utils.UnauthorizedError("invalid user id"))
			c.Abort()
			return
		}
		c.Set(UserIDKey, userID)
		c.Next()
	}
}
