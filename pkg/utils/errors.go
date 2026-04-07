package utils

import (
	"fmt"
	"net/http"
)

// AppError is the unified business error type. Message is client-facing; Internal is for server-side logging only.
type AppError struct {
	Code     int    `json:"code"`
	Message  string `json:"message"`
	Internal string `json:"-"`
}

func (e *AppError) Error() string {
	if e.Internal != "" {
		return e.Internal
	}
	return e.Message
}

func appError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

// NewAppError creates an AppError with a custom status code and message.
func NewAppError(code int, message string) *AppError {
	return appError(code, message)
}

// NotFoundError creates a 404 error.
func NotFoundError(msg string) *AppError {
	return appError(http.StatusNotFound, msg)
}

// BadRequestError creates a 400 error.
func BadRequestError(msg string) *AppError {
	return appError(http.StatusBadRequest, msg)
}

// ForbiddenError creates a 403 error.
func ForbiddenError(msg string) *AppError {
	return appError(http.StatusForbidden, msg)
}

// InternalError creates a 500 error.
func InternalError(msg string) *AppError {
	return appError(http.StatusInternalServerError, msg)
}

// UnauthorizedError creates a 401 error.
func UnauthorizedError(msg string) *AppError {
	return appError(http.StatusUnauthorized, msg)
}

// ConflictError creates a 409 error.
func ConflictError(msg string) *AppError {
	return appError(http.StatusConflict, msg)
}

// WrapError wraps an internal error. The client only sees msg; the full error chain is logged server-side.
func WrapError(err error, msg string) *AppError {
	return &AppError{
		Code:     http.StatusInternalServerError,
		Message:  msg,
		Internal: fmt.Sprintf("%s: %v", msg, err),
	}
}
