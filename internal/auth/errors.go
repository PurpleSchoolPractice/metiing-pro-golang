package auth

const (
	ErrUserExists          = "User already exists"
	ErrWrongCredentials    = "Wrong credentials"
	ErrCreateSecret        = "Failed to create user secret"
	ErrInvalidRefreshToken = "Invalid refresh token"
	ErrAccessToken         = "Access token is valid"
	ErrUserNotFound        = "User not found"
	ErrGenerateToken       = "Failed to generate token"
	InvalidPassword        = "The password must contain 12 characters including numbers and special characters."
	InValidPasOrEmail      = "Invalid passwod or email"
)
