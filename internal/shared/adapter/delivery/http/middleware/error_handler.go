package middleware

import (
	"errors"
	"net/http"

	"go-clean-arch/pkg/log"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

// ErrorResponse is the JSON error response shared by HTTP endpoints and Swagger.
type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// TransportError is a transport-neutral, stable error classification.
type TransportError interface {
	error
	HTTPStatusCode() int
	ErrorCode() string
	ErrorMessage() string
}

func jsonError(c *gin.Context, status int, code, message string) {
	if c.Writer.Written() {
		return
	}
	c.JSON(status, ErrorResponse{Code: code, Message: message})
}

func statusCode(status int) string {
	switch status {
	case http.StatusBadRequest:
		return "bad_request"
	case http.StatusUnauthorized:
		return "unauthorized"
	case http.StatusForbidden:
		return "forbidden"
	case http.StatusNotFound:
		return "not_found"
	case http.StatusConflict:
		return "conflict"
	case http.StatusRequestEntityTooLarge:
		return "payload_too_large"
	case http.StatusTooManyRequests:
		return "rate_limited"
	default:
		return "internal"
	}
}

// Recovery converts panics to the JSON error contract.
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Error("panic recovered", "path", c.Request.URL.Path, "method", c.Request.Method, "panic", recovered)
				jsonError(c, http.StatusInternalServerError, "internal", "internal server error")
				c.Abort()
			}
		}()
		c.Next()
	}
}

// ErrorHandler is a middleware that converts errors from c.Error() into JSON responses.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 || c.Writer.Written() {
			return
		}

		err := c.Errors.Last().Err
		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			if appErr.Internal != "" {
				log.Error(
					"request error",
					"path",
					c.Request.URL.Path,
					"method",
					c.Request.Method,
					"detail",
					appErr.Internal,
				)
			}
			jsonError(c, appErr.Code, statusCode(appErr.Code), appErr.Message)
			return
		}

		var transportErr TransportError
		if errors.As(err, &transportErr) {
			jsonError(c, transportErr.HTTPStatusCode(), transportErr.ErrorCode(), transportErr.ErrorMessage())
			return
		}

		log.Error("unhandled error", "path", c.Request.URL.Path, "method", c.Request.Method, "error", err.Error())
		jsonError(c, http.StatusInternalServerError, "internal", "internal server error")
	}
}
