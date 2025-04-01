package auth

import (
	"errors"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
)

type AuthService struct {
	UserRepository *user.UserRepository
}

func NewAuthService(userRepository *user.UserRepository) *AuthService {
	return &AuthService{
		UserRepository: userRepository,
	}
}

func (service *AuthService) Register(email, password, username string) (string, error) {
	existUser, _ := service.UserRepository.FindByEmail(email)
	if existUser != nil {
		return "", errors.New(ErrUserExists)
	}
	newUser := user.User{
		Email:    email,
		Password: password,
		Username: username,
	}
	_, err := service.UserRepository.Create(&newUser)
	if err != nil {
		return "", err
	}
	return newUser.Email, nil
}
