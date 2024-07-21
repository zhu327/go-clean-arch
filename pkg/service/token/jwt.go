package token

import (
	"fmt"
	"go-wire/pkg/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtAuth struct {
	userSecretKey string
}


func NewTokenService (cfg config.Config) TokenService {
	return &jwtAuth{
		userSecretKey: cfg.SECRET_KEY,
	}
}

type jwtClaims struct {
	TokenID string
	UserID uint
	ExpireAt time.Time
}

// Generate a new JWT token string
func (c *jwtAuth) GenerateToken(req GenerateTokenRequest)(GenerateTokenResponse, error){
	tokenID := uuid.NewString()
	claims := jwtClaims{
		TokenID: tokenID,
		UserID: req.UserID,
		ExpireAt: req.ExpireAt,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	var (
		tokenString string
		err			error
	)
	tokenString, err = token.SignedString([]byte(c.userSecretKey))
	if err != nil {
		return GenerateTokenResponse{}, fmt.Errorf("failed to sign the token \nerror:%w", err)
	}

	response := GenerateTokenResponse{
		TokenID: tokenID,
		TokenString: tokenString,
	}

	return response, nil
}