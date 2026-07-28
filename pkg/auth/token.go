package auth

import "time"

// TokenType identifies the purpose for which a JWT was issued.
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// TokenService defines the legacy JWT operations consumed by current adapters.
type TokenService interface {
	GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// TypedTokenService additionally validates a token's intended purpose. It is
// suitable for adapting to the narrower application ports introduced by callers.
type TypedTokenService interface {
	TokenService
	ValidateTokenOfType(tokenString string, tokenType TokenType) (*TokenClaims, error)
	IssueAccessToken(userID uint, expireAt time.Time) (string, error)
	IssueRefreshToken(userID uint, expireAt time.Time) (string, error)
}

// GenerateTokenRequest is the input for token generation.
type GenerateTokenRequest struct {
	UserID   uint
	ExpireAt time.Time
	Type     TokenType
}

// GenerateTokenResponse is the output of token generation.
type GenerateTokenResponse struct {
	TokenID     string
	TokenString string
}

// TokenClaims holds the claims extracted from a validated JWT.
type TokenClaims struct {
	UserID  uint
	TokenID string
	Type    TokenType
}
