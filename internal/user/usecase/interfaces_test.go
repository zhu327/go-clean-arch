package usecase_test

import (
	"errors"
	"testing"

	userUsecase "go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/crypto"
)

func TestBcryptHasherImplementsPasswordHasher(t *testing.T) {
	var _ userUsecase.PasswordHasher = crypto.NewBcryptHasher()
}

func TestJWTServiceImplementsTokenIssuer(t *testing.T) {
	var _ userUsecase.TokenIssuer = auth.NewTokenService(config.Config{})
}

func TestApplicationErrorExposesStableClassification(t *testing.T) {
	cause := errors.New("cause")
	err := userUsecase.NewApplicationError("invalid_credentials", 401, "invalid credentials", cause)

	if err.Code != "invalid_credentials" || err.HTTPStatus != 401 || err.Message != "invalid credentials" {
		t.Errorf("unexpected application error: %+v", err)
	}
	if !errors.Is(err, cause) {
		t.Fatal("application error must retain its cause")
	}
	if err.HTTPStatusCode() != 401 || err.ErrorCode() != "invalid_credentials" ||
		err.ErrorMessage() != "invalid credentials" {
		t.Errorf("application error does not implement the transport error contract: %+v", err)
	}
}
