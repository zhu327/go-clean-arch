package middleware

import (
	"errors"
	"net/http"

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

		log.Error("unhandled error",
			"path", c.Request.URL.Path,
			"method", c.Request.Method,
			"error", err.Error(),
		)
		jsonError(c, http.StatusInternalServerError, "internal server error")
	}
}
