package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	return r
}

func TestErrorHandler_MapsDomainErrUserAlreadyExistsTo409(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(domain.ErrUserAlreadyExists)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("expected 409, got %d", w.Code)
	}
}

func TestErrorHandler_MapsDomainErrUserNotFoundTo404(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(domain.ErrUserNotFound)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}

func TestErrorHandler_MapsErrInvalidCredentialsTo401(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(usecase.ErrInvalidCredentials)
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", w.Code)
	}
}

func TestErrorHandler_MapsDomainValidationErrorsTo400(t *testing.T) {
	errs := map[string]error{
		"ErrInvalidEmail":     domain.ErrInvalidEmail,
		"ErrUsernameTooShort": domain.ErrUsernameTooShort,
		"ErrUsernameTooLong":  domain.ErrUsernameTooLong,
		"ErrEmptyPassword":    domain.ErrEmptyPassword,
	}
	for name, err := range errs {
		t.Run(name, func(t *testing.T) {
			r := setupRouter()
			r.GET("/test", func(c *gin.Context) {
				_ = c.Error(err)
			})

			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/test", nil)
			r.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400, got %d", w.Code)
			}
		})
	}
}

func TestErrorHandler_StillHandlesAppError(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(utils.BadRequestError("bad request"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestErrorHandler_MapsUnknownErrorTo500(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(errors.New("something unexpected"))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", w.Code)
	}
}

func TestErrorHandler_WrappedDomainErrUserNotFoundMapsTo404(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(fmt.Errorf("wrapped: %w", domain.ErrUserNotFound))
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Code)
	}
}
