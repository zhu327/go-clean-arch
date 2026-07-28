package middleware

import (
	"container/list"
	"net/http"
	"sync"
	"time"

	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

const (
	defaultMaxBodyBytes      int64 = 1 << 20
	defaultRequestsPerWindow       = 10
	defaultWindow                  = time.Minute
	defaultMaxEntries              = 10000
)

// AuthEndpointProtectionConfig configures in-memory, per-IP abuse controls.
type AuthEndpointProtectionConfig struct {
	MaxBodyBytes      int64
	RequestsPerWindow int
	Window            time.Duration
	MaxEntries        int
	Now               func() time.Time
}

type rateLimitEntry struct {
	started  time.Time
	requests int
	element  *list.Element
}

type authRateLimiter struct {
	mu      sync.Mutex
	entries map[string]*rateLimitEntry
	order   list.List
	config  AuthEndpointProtectionConfig
}

func newAuthRateLimiter(config AuthEndpointProtectionConfig) *authRateLimiter {
	if config.RequestsPerWindow <= 0 {
		config.RequestsPerWindow = defaultRequestsPerWindow
	}
	if config.Window <= 0 {
		config.Window = defaultWindow
	}
	if config.MaxEntries <= 0 {
		config.MaxEntries = defaultMaxEntries
	}
	if config.Now == nil {
		config.Now = time.Now
	}
	return &authRateLimiter{entries: make(map[string]*rateLimitEntry), config: config}
}

// allow is O(1): it only examines the requesting IP, or the oldest entry when full.
func (l *authRateLimiter) allow(ip string, now time.Time) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if entry, exists := l.entries[ip]; exists {
		if now.Sub(entry.started) >= l.config.Window {
			entry.started = now
			entry.requests = 0
			l.order.MoveToBack(entry.element)
		}
		entry.requests++
		return entry.requests <= l.config.RequestsPerWindow
	}
	if len(l.entries) >= l.config.MaxEntries {
		oldest := l.order.Front()
		oldestIP := oldest.Value.(string)
		entry := l.entries[oldestIP]
		if now.Sub(entry.started) < l.config.Window {
			return false
		}
		delete(l.entries, oldestIP)
		l.order.Remove(oldest)
	}
	entry := &rateLimitEntry{started: now, requests: 1}
	entry.element = l.order.PushBack(ip)
	l.entries[ip] = entry
	return true
}

// AuthEndpointProtection limits authentication request bodies and request rate.
func AuthEndpointProtection(config AuthEndpointProtectionConfig) gin.HandlerFunc {
	if config.MaxBodyBytes <= 0 {
		config.MaxBodyBytes = defaultMaxBodyBytes
	}
	limiter := newAuthRateLimiter(config)
	return func(c *gin.Context) {
		if c.Request.ContentLength > config.MaxBodyBytes {
			_ = c.Error(utils.PayloadTooLargeError("request body too large"))
			c.Abort()
			return
		}
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, config.MaxBodyBytes)
		if !limiter.allow(c.ClientIP(), limiter.config.Now()) {
			_ = c.Error(utils.TooManyRequestsError("too many requests"))
			c.Abort()
			return
		}
		c.Next()
	}
}
