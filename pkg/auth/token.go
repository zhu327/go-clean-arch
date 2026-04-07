package auth

import "time"

// TokenService defines the interface for JWT token operations.
type TokenService interface {
	GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}

// GenerateTokenRequest is the input for token generation.
type GenerateTokenRequest struct {
	UserID   uint
	ExpireAt time.Time
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
}
