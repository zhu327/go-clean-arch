package auth

import "time"

type TokenService interface {
	GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error)
}

type GenerateTokenRequest struct {
	UserID   uint
	ExpireAt time.Time
}

type GenerateTokenResponse struct {
	TokenID     string
	TokenString string
}
