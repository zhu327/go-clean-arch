package auth

import (
	"fmt"
	"time"

	"go-clean-arch/pkg/config"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtService struct {
	secretKey []byte
}

// NewTokenService creates a new JWT token service.
func NewTokenService(cfg config.Config) TypedTokenService {
	return &jwtService{secretKey: []byte(cfg.SecretKey)}
}

// GenerateToken creates a new signed JWT token with an explicit type.
func (s *jwtService) GenerateToken(req GenerateTokenRequest) (GenerateTokenResponse, error) {
	if !validTokenType(req.Type) {
		return GenerateTokenResponse{}, fmt.Errorf("invalid token type: %q", req.Type)
	}
	tokenID := uuid.NewString()
	claims := jwt.MapClaims{
		"exp":        req.ExpireAt.Unix(),
		"user_id":    req.UserID,
		"jti":        tokenID,
		"token_type": string(req.Type),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return GenerateTokenResponse{}, fmt.Errorf("sign token: %w", err)
	}
	return GenerateTokenResponse{TokenID: tokenID, TokenString: tokenString}, nil
}

// IssueAccessToken issues an access token.
func (s *jwtService) IssueAccessToken(userID uint, expireAt time.Time) (string, error) {
	response, err := s.GenerateToken(GenerateTokenRequest{UserID: userID, ExpireAt: expireAt, Type: AccessToken})
	if err != nil {
		return "", err
	}
	return response.TokenString, nil
}

// IssueRefreshToken issues a refresh token.
func (s *jwtService) IssueRefreshToken(userID uint, expireAt time.Time) (string, error) {
	response, err := s.GenerateToken(GenerateTokenRequest{UserID: userID, ExpireAt: expireAt, Type: RefreshToken})
	if err != nil {
		return "", err
	}
	return response.TokenString, nil
}

// ValidateToken validates an access token.
func (s *jwtService) ValidateToken(tokenString string) (*TokenClaims, error) {
	return s.ValidateTokenOfType(tokenString, AccessToken)
}

// ValidateTokenOfType validates a JWT and verifies its intended purpose.
func (s *jwtService) ValidateTokenOfType(tokenString string, expectedType TokenType) (*TokenClaims, error) {
	if !validTokenType(expectedType) {
		return nil, fmt.Errorf("invalid expected token type: %q", expectedType)
	}
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.secretKey, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", err)
	}
	userID, ok := claims["user_id"].(float64)
	if !ok || userID < 0 || userID != float64(uint(userID)) {
		return nil, fmt.Errorf("invalid user_id claim")
	}
	tokenType, ok := claims["token_type"].(string)
	if !ok || TokenType(tokenType) != expectedType {
		return nil, fmt.Errorf("unexpected token type")
	}
	jti, ok := claims["jti"].(string)
	if !ok || jti == "" {
		return nil, fmt.Errorf("invalid jti claim")
	}
	return &TokenClaims{UserID: uint(userID), TokenID: jti, Type: TokenType(tokenType)}, nil
}

func validTokenType(tokenType TokenType) bool {
	return tokenType == AccessToken || tokenType == RefreshToken
}
