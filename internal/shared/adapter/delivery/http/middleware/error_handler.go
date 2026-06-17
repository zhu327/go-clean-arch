package middleware

import (
	"errors"
	"net/http"

	"go-clean-arch/internal/user/domain"
	"go-clean-arch/internal/user/usecase"
	"go-clean-arch/pkg/log"
	"go-clean-arch/pkg/utils"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func jsonError(c *gin.Context, code int, message string) {
	c.JSON(code, errorResponse{Code: code, Message: message})
}

func mapDomainError(err error) (int, string, bool) {
	switch {
	case errors.Is(err, domain.ErrUserAlreadyExists):
		return http.StatusConflict, domain.ErrUserAlreadyExists.Error(), true
	case errors.Is(err, domain.ErrUserNotFound):
		return http.StatusNotFound, domain.ErrUserNotFound.Error(), true
	case errors.Is(err, usecase.ErrInvalidCredentials):
		return http.StatusUnauthorized, usecase.ErrInvalidCredentials.Error(), true
	case errors.Is(err, domain.ErrInvalidEmail),
		errors.Is(err, domain.ErrUsernameTooShort),
		errors.Is(err, domain.ErrUsernameTooLong),
		errors.Is(err, domain.ErrEmptyPassword):
		return http.StatusBadRequest, err.Error(), true
	default:
		return 0, "", false
	}
}

// ErrorHandler is a middleware that converts errors from c.Error() into JSON responses.
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) == 0 {
			return
		}

		err := c.Errors.Last().Err

		var appErr *utils.AppError
		if errors.As(err, &appErr) {
			if appErr.Internal != "" {
				log.Error("request error",
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
					"detail", appErr.Internal,
				)
			}
			jsonError(c, appErr.Code, appErr.Message)
			return
		}

		if code, msg, ok := mapDomainError(err); ok {
			jsonError(c, code, msg)
			return
		}

		log.Error("unhandled error",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"error", err.Error(),
		)
		jsonError(c, http.StatusInternalServerError, "internal server error")
	}
}
