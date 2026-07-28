package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"go-clean-arch/pkg/auth"
	"go-clean-arch/pkg/config"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

type accessVerifierStub struct {
	userID uint
	err    error
}

func (s accessVerifierStub) ValidateAccessToken(string) (uint, error) { return s.userID, s.err }

func TestAuthMiddlewareRejectsRefreshToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	service := auth.NewTokenService(config.Config{SecretKey: "12345678901234567890123456789012"})
	refresh, err := service.IssueRefreshToken(7, time.Now().Add(time.Minute))
	if err != nil {
		t.Fatalf("IssueRefreshToken returned error: %v", err)
	}
	r := gin.New()
	r.Use(ErrorHandler())
	r.GET(
		"/protected",
		AuthMiddleware(NewAccessTokenVerifier(service)),
		func(c *gin.Context) { c.Status(http.StatusNoContent) },
	)
	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer "+refresh)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", w.Code)
	}
}

func TestAuthMiddlewareUsesAccessVerifier(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	r.GET(
		"/protected",
		AuthMiddleware(accessVerifierStub{userID: 7}),
		func(c *gin.Context) { c.Status(http.StatusNoContent) },
	)

	request := httptest.NewRequest(http.MethodGet, "/protected", nil)
	request.Header.Set("Authorization", "Bearer access-token")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, request)
	if w.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", w.Code)
	}
}

func TestErrorHandlerDoesNotWriteAfterResponseIsCommitted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	r.GET("/written", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"ok": true})
		_ = c.Error(http.ErrAbortHandler)
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/written", nil))
	if w.Code != http.StatusCreated || w.Body.String() != `{"ok":true}` {
		t.Fatalf("response was rewritten: %d %s", w.Code, w.Body.String())
	}
}

func TestRecoveryDoesNotWriteAfterResponseIsCommitted(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery())
	r.GET("/written-panic", func(c *gin.Context) {
		c.JSON(http.StatusCreated, gin.H{"ok": true})
		panic("boom")
	})
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/written-panic", nil))
	if w.Code != http.StatusCreated || w.Body.String() != `{"ok":true}` {
		t.Fatalf("response was rewritten: %d %s", w.Code, w.Body.String())
	}
}

func TestRecoveryUsesErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(Recovery())
	r.GET("/panic", func(*gin.Context) { panic("boom") })
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/panic", nil))
	if w.Code != http.StatusInternalServerError ||
		w.Body.String() != `{"code":"internal","message":"internal server error"}` {
		t.Fatalf("panic response = %d %s", w.Code, w.Body.String())
	}
}

func TestAuthEndpointProtectionDoesNotTrustForwardedForByDefault(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	if err := r.SetTrustedProxies(nil); err != nil {
		t.Fatalf("set trusted proxies: %v", err)
	}
	r.Use(ErrorHandler())
	r.POST(
		"/login",
		AuthEndpointProtection(AuthEndpointProtectionConfig{RequestsPerWindow: 1}),
		func(c *gin.Context) { c.Status(http.StatusNoContent) },
	)
	for _, forwardedFor := range []string{"198.51.100.1", "203.0.113.2"} {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.Header.Set("X-Forwarded-For", forwardedFor)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		if forwardedFor == "203.0.113.2" && w.Code != http.StatusTooManyRequests {
			t.Fatalf("untrusted forwarded header bypassed rate limit: %d", w.Code)
		}
	}
}

func TestAuthEndpointProtectionRejectsChunkedOversizeBodyWithErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.Use(ErrorHandler())
	r.POST("/login", AuthEndpointProtection(AuthEndpointProtectionConfig{MaxBodyBytes: 4}), func(c *gin.Context) {
		var payload any
		if err := c.ShouldBindJSON(&payload); err != nil {
			if _, tooLarge := err.(*http.MaxBytesError); tooLarge {
				_ = c.Error(utils.PayloadTooLargeError("request body too large"))
				return
			}
			_ = c.Error(utils.BadRequestError(err.Error()))
			return
		}
		c.Status(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(`{"value":"too large"}`))
	req.ContentLength = -1
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	if w.Code != http.StatusRequestEntityTooLarge ||
		w.Body.String() != `{"code":"payload_too_large","message":"request body too large"}` {
		t.Fatalf("chunked oversize response = %d %s", w.Code, w.Body.String())
	}
}

func TestAuthEndpointProtectionEvictsOnlyOneOldestEntryWhenFull(t *testing.T) {
	now := time.Now()
	limiter := newAuthRateLimiter(AuthEndpointProtectionConfig{
		RequestsPerWindow: 1,
		Window:            time.Minute,
		MaxEntries:        2,
		Now:               func() time.Time { return now },
	})

	if !limiter.allow("192.0.2.1", now) || !limiter.allow("192.0.2.2", now) {
		t.Fatal("initial IPs should be allowed")
	}
	now = now.Add(time.Minute)
	if !limiter.allow("192.0.2.3", now) {
		t.Fatal("new IP should replace the oldest expired entry")
	}
	if len(limiter.entries) != 2 {
		t.Fatalf("entry count = %d, want 2", len(limiter.entries))
	}
	if _, exists := limiter.entries["192.0.2.1"]; exists {
		t.Fatal("oldest entry was not evicted")
	}
	if _, exists := limiter.entries["192.0.2.2"]; !exists {
		t.Fatal("a normal request performed a full-table expired-entry cleanup")
	}
}

func TestAuthEndpointProtectionExpiresEntriesAndCapsNewIPs(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	r := gin.New()
	r.Use(ErrorHandler())
	r.POST("/login", AuthEndpointProtection(AuthEndpointProtectionConfig{
		RequestsPerWindow: 1, Window: time.Minute, MaxEntries: 1, Now: func() time.Time { return now },
	}), func(c *gin.Context) { c.Status(http.StatusNoContent) })
	request := func(ip string) *httptest.ResponseRecorder {
		req := httptest.NewRequest(http.MethodPost, "/login", nil)
		req.RemoteAddr = ip + ":1234"
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		return w
	}
	if w := request("192.0.2.1"); w.Code != http.StatusNoContent {
		t.Fatalf("first IP = %d", w.Code)
	}
	if w := request("192.0.2.2"); w.Code != http.StatusTooManyRequests {
		t.Fatalf("capacity overflow = %d", w.Code)
	}
	now = now.Add(time.Minute)
	if w := request("192.0.2.2"); w.Code != http.StatusNoContent {
		t.Fatalf("expired entry was not cleaned = %d", w.Code)
	}
}

func TestAuthEndpointProtectionLimitsBodyAndRequestsPerIP(t *testing.T) {
	gin.SetMode(gin.TestMode)
	now := time.Now()
	r := gin.New()
	r.Use(ErrorHandler())
	r.POST("/login", AuthEndpointProtection(AuthEndpointProtectionConfig{
		MaxBodyBytes: 4, RequestsPerWindow: 1, Window: time.Minute, Now: func() time.Time { return now },
	}), func(c *gin.Context) { c.Status(http.StatusNoContent) })

	for _, body := range []string{"ok", "ok"} {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString(body)))
		if body == "ok" && w.Code != http.StatusNoContent && w.Code != http.StatusTooManyRequests {
			t.Fatalf("status = %d", w.Code)
		}
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodPost, "/login", bytes.NewBufferString("12345")))
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d, want 413", w.Code)
	}
}
