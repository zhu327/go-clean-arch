package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-clean-arch/internal/shared/adapter/delivery/http/middleware"
	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase/dto"

	"github.com/gin-gonic/gin"
)

type spyUserUseCase struct {
	signUpFn   func(ctx context.Context, params dto.SignUpParams) (*domain.User, error)
	loginFn    func(ctx context.Context, params dto.LoginParams) (*dto.AuthTokens, error)
	findByIDFn func(ctx context.Context, id uint) (*domain.User, error)
}

func (s *spyUserUseCase) SignUp(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
	if s.signUpFn != nil {
		return s.signUpFn(ctx, params)
	}
	return &domain.User{ID: 1, Username: params.Username, Email: params.Email}, nil
}

func (s *spyUserUseCase) Login(ctx context.Context, params dto.LoginParams) (*dto.AuthTokens, error) {
	if s.loginFn != nil {
		return s.loginFn(ctx, params)
	}
	return &dto.AuthTokens{AccessToken: "access", RefreshToken: "refresh"}, nil
}

func (s *spyUserUseCase) FindByID(ctx context.Context, id uint) (*domain.User, error) {
	if s.findByIDFn != nil {
		return s.findByIDFn(ctx, id)
	}
	return &domain.User{ID: id}, nil
}

func setupRouter(uc UserUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	h := NewUserHandler(uc)
	r.POST("/api/auth/signup", h.SignUp)
	r.POST("/api/auth/login", h.Login)
	return r
}

func TestSignUp_ShortUsername_Returns400WithDomainError(t *testing.T) {
	uc := &spyUserUseCase{
		signUpFn: func(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
			return nil, domain.ErrUsernameTooShort
		},
	}

	router := setupRouter(uc)

	body := `{"username":"ab","email":"test@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d; body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", w.Body.String())
	}

	msg, _ := resp["message"].(string)
	if msg != domain.ErrUsernameTooShort.Error() {
		t.Errorf("expected message %q, got %q", domain.ErrUsernameTooShort.Error(), msg)
	}
}

func TestSignUp_EmptyPassword_Returns400WithDomainError(t *testing.T) {
	uc := &spyUserUseCase{
		signUpFn: func(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
			return nil, domain.ErrEmptyPassword
		},
	}

	router := setupRouter(uc)

	body := `{"username":"validuser","email":"test@test.com","password":"secret"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d; body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", w.Body.String())
	}

	msg, _ := resp["message"].(string)
	if msg != domain.ErrEmptyPassword.Error() {
		t.Errorf("expected message %q, got %q", domain.ErrEmptyPassword.Error(), msg)
	}
}

func TestSignUp_ValidRequest_Returns201(t *testing.T) {
	uc := &spyUserUseCase{
		signUpFn: func(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
			return &domain.User{
				ID:        1,
				Username:  params.Username,
				Email:     params.Email,
				Password:  "hashed",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	}

	router := setupRouter(uc)

	body := `{"username":"validuser","email":"test@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d; body=%s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", w.Body.String())
	}

	if resp["username"] != "validuser" {
		t.Errorf("expected username 'validuser', got %v", resp["username"])
	}
}
