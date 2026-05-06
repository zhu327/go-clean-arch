package usecase_test

import (
	"testing"

	userUsecase "go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/crypto"
)

func TestBcryptHasherImplementsPasswordHasher(t *testing.T) {
	var _ userUsecase.PasswordHasher = crypto.NewBcryptHasher()
}
