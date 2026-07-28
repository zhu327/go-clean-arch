package auth

import (
	"testing"
	"time"

	"go-clean-arch/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "12345678901234567890123456789012"

func newTestTokenService() TypedTokenService {
	return NewTokenService(config.Config{SecretKey: testSecret})
}

func TestTokenService_ValidatesOnlyRequestedTokenType(t *testing.T) {
	service := newTestTokenService()
	refresh, err := service.GenerateToken(
		GenerateTokenRequest{UserID: 42, ExpireAt: time.Now().Add(time.Hour), Type: RefreshToken},
	)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	if _, err := service.ValidateToken(refresh.TokenString); err == nil {
		t.Fatal("ValidateToken accepted a refresh token as an access token")
	}
	claims, err := service.ValidateTokenOfType(refresh.TokenString, RefreshToken)
	if err != nil {
		t.Fatalf("ValidateTokenOfType returned error: %v", err)
	}
	if claims.UserID != 42 || claims.Type != RefreshToken {
		t.Errorf("unexpected claims: %+v", claims)
	}
}

func TestTokenService_IssuesAccessAndRefreshTokensExplicitly(t *testing.T) {
	service := newTestTokenService()

	access, err := service.IssueAccessToken(42, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("IssueAccessToken returned error: %v", err)
	}
	refresh, err := service.IssueRefreshToken(42, time.Now().Add(time.Hour))
	if err != nil {
		t.Fatalf("IssueRefreshToken returned error: %v", err)
	}
	if _, err := service.ValidateTokenOfType(access, AccessToken); err != nil {
		t.Fatalf("access token validation failed: %v", err)
	}
	if _, err := service.ValidateTokenOfType(refresh, RefreshToken); err != nil {
		t.Fatalf("refresh token validation failed: %v", err)
	}
}

func TestTokenService_RejectsUnexpectedSigningAlgorithm(t *testing.T) {
	service := newTestTokenService()
	claims := jwt.MapClaims{
		"exp":        time.Now().Add(time.Hour).Unix(),
		"user_id":    42,
		"jti":        "id",
		"token_type": string(AccessToken),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS384, claims)
	tokenString, err := token.SignedString([]byte(testSecret))
	if err != nil {
		t.Fatal(err)
	}

	if _, err := service.ValidateToken(tokenString); err == nil {
		t.Fatal("ValidateToken accepted a token signed with HS384")
	}
}

func TestTokenService_RejectsMissingOrInvalidTokenType(t *testing.T) {
	service := newTestTokenService()
	for _, tokenType := range []TokenType{"", "unknown"} {
		t.Run(string(tokenType), func(t *testing.T) {
			if _, err := service.GenerateToken(GenerateTokenRequest{UserID: 42, ExpireAt: time.Now().Add(time.Hour), Type: tokenType}); err == nil {
				t.Fatal("GenerateToken succeeded, want error")
			}
		})
	}
}
