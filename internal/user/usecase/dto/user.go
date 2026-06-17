package dto

// SignUpParams is the input for the sign-up use case.
type SignUpParams struct {
	Username string
	Email    string
	Password string
}

// LoginParams is the input for the login use case.
type LoginParams struct {
	Email    string
	Password string
}

// AuthTokens holds the tokens returned after successful login.
type AuthTokens struct {
	AccessToken  string
	RefreshToken string
}
