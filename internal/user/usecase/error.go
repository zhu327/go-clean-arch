package usecase

const (
	// HTTP status values are defined locally to keep the use case independent of net/http.
	StatusBadRequest   = 400
	StatusUnauthorized = 401
	StatusNotFound     = 404
	StatusConflict     = 409

	ErrorCodeInvalidArgument       = "invalid_argument"
	ErrorCodeInvalidCredentials    = "invalid_credentials"
	ErrorCodeUserNotFound          = "user_not_found"
	ErrorCodeUserAlreadyExists     = "user_already_exists"
	ErrorCodeEmailAlreadyExists    = "email_already_exists"
	ErrorCodeUsernameAlreadyExists = "username_already_exists"
)

// ApplicationError is a stable, transport-neutral error classification.
type ApplicationError struct {
	Code       string
	HTTPStatus int
	Message    string
	cause      error
}

// NewApplicationError creates an application error with an optional cause.
func NewApplicationError(code string, httpStatus int, message string, cause error) *ApplicationError {
	return &ApplicationError{Code: code, HTTPStatus: httpStatus, Message: message, cause: cause}
}

func (e *ApplicationError) Error() string {
	return e.Message
}

// HTTPStatusCode returns the status for HTTP-like transports.
func (e *ApplicationError) HTTPStatusCode() int {
	return e.HTTPStatus
}

// ErrorCode returns a stable machine-readable error code.
func (e *ApplicationError) ErrorCode() string {
	return e.Code
}

// ErrorMessage returns the client-safe error message.
func (e *ApplicationError) ErrorMessage() string {
	return e.Message
}

// Unwrap returns the underlying cause.
func (e *ApplicationError) Unwrap() error {
	return e.cause
}
