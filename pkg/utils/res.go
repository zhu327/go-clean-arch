package utils

import "strings"

type Response struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Errors  interface{} `json:"error,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code"`
}

func SuccessResponse(status int, message string, data ...interface{}) Response {
	return Response{
		Status:  status,
		Message: message,
		Errors:  nil,
		Data:    data,
	}
}

func ErrorResponse(status int, message string, err string, data interface{}) Response {
	splittedError := strings.Split(err, "\n")
	return Response{
		Status:  status,
		Message: message,
		Errors:  splittedError,
		Data:    data,
	}
}

type ReponseUsers struct {
	ID         uint   `json:"id"`
	Username   string `json:"username"`
	Email      string `json:"email"`
	Created_at string `json:"created_at"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}