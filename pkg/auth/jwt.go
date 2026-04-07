package auth

import (
	"fmt"

	"go-clean-arch/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	secretKey string
}

// NewTokenService creates a new JWT token service.
func NewTokenService(cfg config.Config) TokenService {
	return &jwtService{
		secretKey: cfg.SecretKey,
	}
}

// GenerateToken creates a new signed JWT token.
func (s *jwtService) GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error) {
	tokenID := uuid.NewString()
	claims := jwt.MapClaims{
		"exp":     req.ExpireAt.Unix(),
		"user_id": req.UserID,
		"jti":     tokenID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return GenerateTokenResponse{}, fmt.Errorf("failed to sign token: %w", err)
	}

	return GenerateTokenResponse{
		TokenID:     tokenID,
		TokenString: tokenString,
	}, nil
}

// ValidateToken parses and validates a JWT token string, returning the claims.
func (s *jwtService) ValidateToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id claim")
	}

	jti, _ := claims["jti"].(string)

	return &TokenClaims{
		UserID:  uint(userIDFloat),
		TokenID: jti,
	}, nil
}
