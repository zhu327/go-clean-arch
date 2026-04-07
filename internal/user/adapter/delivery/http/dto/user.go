package dto

// SignUpRequest is the registration request body.
type SignUpRequest struct {
	Username string `json:"username" binding:"required,min=3,max=15"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5,max=25"`
}

// LoginRequest is the login request body.
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=5,max=25"`
}

// UserResponse is the user response (excludes sensitive data).
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// TokenResponse is the authentication token response.
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
