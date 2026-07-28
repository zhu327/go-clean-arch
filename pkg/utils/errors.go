package utils

import (
	"fmt"
	"net/http"
)

// Stable, machine-readable error codes shared across transport layers.
const (
	CodeBadRequest      = "bad_request"
	CodeUnauthorized    = "unauthorized"
	CodeForbidden       = "forbidden"
	CodeNotFound        = "not_found"
	CodeConflict        = "conflict"
	CodePayloadTooLarge = "payload_too_large"
	CodeRateLimited     = "rate_limited"
	CodeInternal        = "internal"
)

// AppError is the unified business error type. Message is client-facing; Internal
// is for server-side logging only. It satisfies the middleware.TransportError
// interface so the error handler has a single, consistent code path.
type AppError struct {
	Code     int    // HTTP status code
	ErrCode  string // stable machine-readable error code
	Message  string // client-facing message
	Internal string // server-side detail, never serialized to clients
}

func (e *AppError) Error() string {
	if e.Internal != "" {
		return e.Internal
	}
	return e.Message
}

// HTTPStatusCode returns the HTTP status for this error.
func (e *AppError) HTTPStatusCode() int { return e.Code }

// ErrorCode returns the stable machine-readable error code.
func (e *AppError) ErrorCode() string { return e.ErrCode }

// ErrorMessage returns the client-safe error message.
func (e *AppError) ErrorMessage() string { return e.Message }

func appError(code int, errCode, message string) *AppError {
	return &AppError{Code: code, ErrCode: errCode, Message: message}
}

// NewAppError creates an AppError with an explicit HTTP status, error code, and message.
func NewAppError(code int, errCode, message string) *AppError {
	return appError(code, errCode, message)
}

// NotFoundError creates a 404 error.
func NotFoundError(msg string) *AppError {
	return appError(http.StatusNotFound, CodeNotFound, msg)
}

// BadRequestError creates a 400 error.
func BadRequestError(msg string) *AppError {
	return appError(http.StatusBadRequest, CodeBadRequest, msg)
}

// ForbiddenError creates a 403 error.
func ForbiddenError(msg string) *AppError {
	return appError(http.StatusForbidden, CodeForbidden, msg)
}

// InternalError creates a 500 error.
func InternalError(msg string) *AppError {
	return appError(http.StatusInternalServerError, CodeInternal, msg)
}

// UnauthorizedError creates a 401 error.
func UnauthorizedError(msg string) *AppError {
	return appError(http.StatusUnauthorized, CodeUnauthorized, msg)
}

// ConflictError creates a 409 error.
func ConflictError(msg string) *AppError {
	return appError(http.StatusConflict, CodeConflict, msg)
}

// PayloadTooLargeError creates a 413 error.
func PayloadTooLargeError(msg string) *AppError {
	return appError(http.StatusRequestEntityTooLarge, CodePayloadTooLarge, msg)
}

// TooManyRequestsError creates a 429 error.
func TooManyRequestsError(msg string) *AppError {
	return appError(http.StatusTooManyRequests, CodeRateLimited, msg)
}

// WrapError wraps an internal error. The client only sees msg; the full error chain is logged server-side.
func WrapError(err error, msg string) *AppError {
	return &AppError{
		Code:     http.StatusInternalServerError,
		ErrCode:  CodeInternal,
		Message:  msg,
		Internal: fmt.Sprintf("%s: %v", msg, err),
	}
}
