package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

type transportError struct {
	status  int
	code    string
	message string
}

func (e transportError) Error() string        { return e.message }
func (e transportError) HTTPStatusCode() int  { return e.status }
func (e transportError) ErrorCode() string    { return e.code }
func (e transportError) ErrorMessage() string { return e.message }

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	return r
}

func TestErrorHandlerMapsTransportErrorAndUsesExportedResponse(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) {
		_ = c.Error(transportError{http.StatusConflict, "email_already_exists", "email already exists"})
	})

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want 409", w.Code)
	}
	if w.Body.String() != `{"code":"email_already_exists","message":"email already exists"}` {
		t.Fatalf("body = %s", w.Body.String())
	}
}

func TestErrorHandlerStillHandlesAppError(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) { _ = c.Error(utils.BadRequestError("bad request")) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want 400", w.Code)
	}
}

func TestErrorHandlerMapsUnknownErrorTo500(t *testing.T) {
	r := setupRouter()
	r.GET("/test", func(c *gin.Context) { _ = c.Error(errors.New("unexpected")) })

	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/test", nil))

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", w.Code)
	}
}
