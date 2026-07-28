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
	"go-clean-arch/internal/user/usecase"
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

func TestSignUp_ShortUsername_Returns400WithApplicationError(t *testing.T) {
	uc := &spyUserUseCase{
		signUpFn: func(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
			return nil, usecase.NewApplicationError(
				usecase.ErrorCodeInvalidArgument,
				usecase.StatusBadRequest,
				domain.ErrUsernameTooShort.Error(),
				domain.ErrUsernameTooShort,
			)
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

func TestSignUp_EmptyPassword_Returns400WithApplicationError(t *testing.T) {
	uc := &spyUserUseCase{
		signUpFn: func(ctx context.Context, params dto.SignUpParams) (*domain.User, error) {
			return nil, usecase.NewApplicationError(
				usecase.ErrorCodeInvalidArgument,
				usecase.StatusBadRequest,
				domain.ErrEmptyPassword.Error(),
				domain.ErrEmptyPassword,
			)
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

func TestLogin_ValidRequest_Returns200WithTokens(t *testing.T) {
	uc := &spyUserUseCase{
		loginFn: func(ctx context.Context, params dto.LoginParams) (*dto.AuthTokens, error) {
			return &dto.AuthTokens{AccessToken: "access-token", RefreshToken: "refresh-token"}, nil
		},
	}

	router := setupRouter(uc)

	body := `{"email":"test@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", w.Body.String())
	}
	if resp["access_token"] != "access-token" || resp["refresh_token"] != "refresh-token" {
		t.Errorf("unexpected tokens: %v", resp)
	}
}

func TestLogin_InvalidCredentials_Returns401(t *testing.T) {
	uc := &spyUserUseCase{
		loginFn: func(ctx context.Context, params dto.LoginParams) (*dto.AuthTokens, error) {
			return nil, usecase.NewApplicationError(
				usecase.ErrorCodeInvalidCredentials,
				usecase.StatusUnauthorized,
				"invalid credentials",
				usecase.ErrInvalidCredentials,
			)
		},
	}

	router := setupRouter(uc)

	body := `{"email":"test@test.com","password":"wrong"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d; body=%s", w.Code, w.Body.String())
	}
}

func TestLogin_MalformedJSON_Returns400(t *testing.T) {
	router := setupRouter(&spyUserUseCase{})

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(`{not json`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d; body=%s", w.Code, w.Body.String())
	}
}

func TestSignUp_OversizeBody_Returns413(t *testing.T) {
	router := setupRouter(&spyUserUseCase{})

	body := `{"username":"validuser","email":"test@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/signup", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	req.Body = http.MaxBytesReader(w, req.Body, 4)

	router.ServeHTTP(w, req)

	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("expected status 413, got %d; body=%s", w.Code, w.Body.String())
	}
}

func setupMeRouter(uc UserUseCase, userID any) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(middleware.ErrorHandler())
	h := NewUserHandler(uc)
	if userID != nil {
		r.GET("/api/user/me", func(c *gin.Context) {
			c.Set(middleware.UserIDKey, userID)
			h.Me(c)
		})
	} else {
		r.GET("/api/user/me", h.Me)
	}
	return r
}

func TestMe_AuthenticatedUser_Returns200(t *testing.T) {
	uc := &spyUserUseCase{
		findByIDFn: func(ctx context.Context, id uint) (*domain.User, error) {
			return &domain.User{ID: id, Username: "testuser", Email: "test@test.com"}, nil
		},
	}

	router := setupMeRouter(uc, uint(7))

	req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d; body=%s", w.Code, w.Body.String())
	}
	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("failed to parse response: %v", w.Body.String())
	}
	if resp["username"] != "testuser" {
		t.Errorf("expected username 'testuser', got %v", resp["username"])
	}
}

func TestMe_MissingUserID_Returns401(t *testing.T) {
	router := setupMeRouter(&spyUserUseCase{}, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d; body=%s", w.Code, w.Body.String())
	}
}

func TestMe_UserNotFound_Returns404(t *testing.T) {
	uc := &spyUserUseCase{
		findByIDFn: func(ctx context.Context, id uint) (*domain.User, error) {
			return nil, usecase.NewApplicationError(
				usecase.ErrorCodeUserNotFound,
				usecase.StatusNotFound,
				domain.ErrUserNotFound.Error(),
				domain.ErrUserNotFound,
			)
		},
	}

	router := setupMeRouter(uc, uint(1))

	req := httptest.NewRequest(http.MethodGet, "/api/user/me", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d; body=%s", w.Code, w.Body.String())
	}
}
