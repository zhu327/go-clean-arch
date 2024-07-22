package token

import (
	"fmt"
	"go-wire/pkg/config"

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

// Generate a new JWT token string
func (c *jwtAuth) GenerateToken(req GenerateTokenRequest)(GenerateTokenResponse, error){
	tokenID := uuid.NewString()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := make(jwt.MapClaims)
	claims["exp"] = req.ExpireAt.Unix()
	claims["UserID"] = req.UserID
	token.Claims = claims

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