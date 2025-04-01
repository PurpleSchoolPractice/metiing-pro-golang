package auth

const (
	ErrUserExists          = "User already exists"
	ErrWrongCredentials    = "Wrong credentials"
	ErrCreateSecret        = "Failed to create user secret"
	ErrInvalidRefreshToken = "Invalid refresh token"
	ErrUserNotFound        = "User not found"
	ErrGenerateToken       = "Failed to generate token"
)
