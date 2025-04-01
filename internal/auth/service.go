package auth

import (
	"errors"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository   *user.UserRepository
	SecretRepository *secret.SecretRepository
	JWT              *jwt.JWT
}

func NewAuthService(
	userRepository *user.UserRepository,
	secretRepository *secret.SecretRepository,
	jwtService *jwt.JWT,
) *AuthService {
	return &AuthService{
		UserRepository:   userRepository,
		SecretRepository: secretRepository,
		JWT:              jwtService,
	}
}

func (service *AuthService) Register(email, password, username string) (string, error) {
	existUser, _ := service.UserRepository.FindByEmail(email)
	if existUser != nil {
		return "", errors.New(ErrUserExists)
	}
	if err := secret.ValidatePassword(password); err != nil {
		return "", errors.New(ErrUserExists)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	newUser := user.User{
		Email:    email,
		Password: string(hashedPassword),
		Username: username,
	}
	createdUser, err := service.UserRepository.Create(&newUser)
	if err != nil {
		return "", errors.New(ErrUserExists)
	}
	_, err = service.SecretRepository.Create(string(hashedPassword), createdUser.ID)
	if err != nil {
		return "", errors.New(ErrCreateSecret)
	}
	return newUser.Email, nil
}

func (service *AuthService) Login(email, password string) (string, error) {
	existUser, _ := service.UserRepository.FindByEmail(email)
	if existUser == nil {
		return "", errors.New(ErrWrongCredentials)
	}
	err := bcrypt.CompareHashAndPassword([]byte(existUser.Password), []byte(password))
	if err != nil {
		return "", errors.New(ErrWrongCredentials)
	}
	return existUser.Email, nil
}

func (service *AuthService) RefreshTokens(refreshToken string) (*jwt.TokenPair, error) {
	valid, data := service.JWT.ParseRefreshToken(refreshToken)
	if !valid || data == nil {
		return nil, errors.New(ErrInvalidRefreshToken)
	}

	existUser, err := service.UserRepository.FindByEmail(data.Email)
	if err != nil || existUser == nil {
		return nil, errors.New(ErrUserNotFound)
	}

	tokenPair, err := service.JWT.GenerateTokenPair(jwt.JWTData{Email: existUser.Email})
	if err != nil {
		return nil, errors.New(ErrGenerateToken)
	}

	return tokenPair, nil
}
