package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase/dto"
	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/utils"
)

type stubUserRepo struct {
	findByEmailFn func(ctx context.Context, email string) (domain.User, error)
	findByIDFn    func(ctx context.Context, id uint) (domain.User, error)
}

func (s *stubUserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	return user, nil
}

func (s *stubUserRepo) FindByID(ctx context.Context, id uint) (domain.User, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return domain.User{}, nil
}

func (s *stubUserRepo) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	if s.findByEmailFn != nil {
		return s.findByEmailFn(ctx, email)
	}
	return domain.User{}, domain.ErrUserNotFound
}

type stubTokenService struct {
	generateTokenFn func(req auth.GenerateTokenRequest) (auth.GenerateTokenResponse, error)
	validateTokenFn func(tokenString string) (*auth.TokenClaims, error)
}

func (s *stubTokenService) GenerateToken(req auth.GenerateTokenRequest) (auth.GenerateTokenResponse, error) {
	if s.generateTokenFn != nil {
		return s.generateTokenFn(req)
	}
	return auth.GenerateTokenResponse{TokenString: "token", TokenID: "id"}, nil
}

func (s *stubTokenService) ValidateToken(tokenString string) (*auth.TokenClaims, error) {
	if s.validateTokenFn != nil {
		return s.validateTokenFn(tokenString)
	}
	return &auth.TokenClaims{}, nil
}

type stubPasswordHasher struct {
	hashFn   func(password string) (string, error)
	verifyFn func(password, hash string) error
}

func (s *stubPasswordHasher) Hash(password string) (string, error) {
	if s.hashFn != nil {
		return s.hashFn(password)
	}
	return "$2a$14$stubhashstubhash", nil
}

func (s *stubPasswordHasher) Verify(password, hash string) error {
	if s.verifyFn != nil {
		return s.verifyFn(password, hash)
	}
	return nil
}

func newUserManager(repo UserRepository, ts auth.TokenService, hasher PasswordHasher) *UserManager {
	return NewUserManager(repo, ts, hasher, TokenTTLs{
		Access:  20 * time.Minute,
		Refresh: 7 * 24 * time.Hour,
	})
}

func TestUserManager_SignUpUsesPasswordHasher(t *testing.T) {
	var hashedPassword string
	hasher := &stubPasswordHasher{
		hashFn: func(password string) (string, error) {
			hashedPassword = "hashed:" + password
			return hashedPassword, nil
		},
	}

	repo := &stubUserRepo{}

	m := newUserManager(repo, &stubTokenService{}, hasher)

	user, err := m.SignUp(context.Background(), dto.SignUpParams{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "mypassword",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hashedPassword != "hashed:mypassword" {
		t.Errorf("expected hasher to be called with mypassword, got hash=%q", hashedPassword)
	}
	if user.Password != hashedPassword {
		t.Errorf("expected user.Password=%q, got %q", hashedPassword, user.Password)
	}
}

func TestUserManager_SignUpReturnsDomainErrorWhenUserExists(t *testing.T) {
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{ID: 1, Email: email}, nil
		},
	}

	m := newUserManager(repo, &stubTokenService{}, &stubPasswordHasher{})

	_, err := m.SignUp(context.Background(), dto.SignUpParams{
		Username: "testuser",
		Email:    "test@test.com",
		Password: "password",
	})
	if err == nil {
		t.Fatal("expected error when user already exists")
	}
	if !errors.Is(err, domain.ErrUserAlreadyExists) {
		t.Errorf("expected domain.ErrUserAlreadyExists, got %v", err)
	}
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		t.Error("error should not be an AppError, should be a plain domain sentinel")
	}
}

func TestUserManager_SignUpReturnsDomainValidationErrorForShortUsername(t *testing.T) {
	repo := &stubUserRepo{}
	m := newUserManager(repo, &stubTokenService{}, &stubPasswordHasher{})

	_, err := m.SignUp(context.Background(), dto.SignUpParams{
		Username: "ab",
		Email:    "test@test.com",
		Password: "password",
	})
	if err == nil {
		t.Fatal("expected error for short username")
	}
	if !errors.Is(err, domain.ErrUsernameTooShort) {
		t.Errorf("expected domain.ErrUsernameTooShort, got %v", err)
	}
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		t.Error("error should not be an AppError")
	}
}

func TestUserManager_LoginReturnsErrInvalidCredentialsForMissingUser(t *testing.T) {
	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{}, domain.ErrUserNotFound
		},
	}

	m := newUserManager(repo, &stubTokenService{}, &stubPasswordHasher{})

	_, err := m.Login(context.Background(), dto.LoginParams{
		Email:    "missing@test.com",
		Password: "password",
	})
	if err == nil {
		t.Fatal("expected error for missing user")
	}
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("expected ErrInvalidCredentials, got %v", err)
	}
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		t.Error("error should not be an AppError")
	}
}

func TestUserManager_FindByIDReturnsDomainErrUserNotFound(t *testing.T) {
	repo := &stubUserRepo{
		findByIDFn: func(ctx context.Context, id uint) (domain.User, error) {
			return domain.User{}, domain.ErrUserNotFound
		},
	}

	m := newUserManager(repo, &stubTokenService{}, &stubPasswordHasher{})

	_, err := m.FindByID(context.Background(), 999)
	if err == nil {
		t.Fatal("expected error for missing user")
	}
	if !errors.Is(err, domain.ErrUserNotFound) {
		t.Errorf("expected domain.ErrUserNotFound, got %v", err)
	}
	var appErr *utils.AppError
	if errors.As(err, &appErr) {
		t.Error("error should not be an AppError")
	}
}

func TestUserManager_LoginUsesPasswordHasher(t *testing.T) {
	var verified bool
	hasher := &stubPasswordHasher{
		verifyFn: func(password, hash string) error {
			verified = true
			if password == "correct" {
				return nil
			}
			return ErrInvalidCredentials
		},
	}

	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{ID: 1, Email: email, Password: "hashed"}, nil
		},
	}

	m := newUserManager(repo, &stubTokenService{}, hasher)

	_, err := m.Login(context.Background(), dto.LoginParams{Email: "test@test.com", Password: "correct"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !verified {
		t.Fatal("expected passwordHasher.Verify to be called")
	}
}

func TestUserManager_LoginRejectsWrongPassword(t *testing.T) {
	hasher := &stubPasswordHasher{
		verifyFn: func(password, hash string) error {
			return ErrInvalidCredentials
		},
	}

	repo := &stubUserRepo{
		findByEmailFn: func(ctx context.Context, email string) (domain.User, error) {
			return domain.User{ID: 1, Email: email, Password: "hashed"}, nil
		},
	}

	m := newUserManager(repo, &stubTokenService{}, hasher)

	_, err := m.Login(context.Background(), dto.LoginParams{Email: "test@test.com", Password: "wrong"})
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}
