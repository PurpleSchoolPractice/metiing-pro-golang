package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/passwordReset"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/sendmail"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserRepository          *user.UserRepository
	SecretRepository        *secret.SecretRepository
	passwordResetRepository *passwordReset.PasswordResetRepository
	JWT                     *jwt.JWT
}

// NewAuthService - конструктор сервиса авторизации
func NewAuthService(
	userRepository *user.UserRepository,
	secretRepository *secret.SecretRepository,
	passwordResetRepository *passwordReset.PasswordResetRepository,
	jwtService *jwt.JWT,
) *AuthService {
	return &AuthService{
		UserRepository:          userRepository,
		SecretRepository:        secretRepository,
		passwordResetRepository: passwordResetRepository,
		JWT:                     jwtService,
	}
}

// Register - регистрация пользователя
func (service *AuthService) Register(email, password, username string) (*models.User, error) {
	existUser, _ := service.UserRepository.FindByEmail(email)
	if existUser != nil {
		return nil, errors.New(ErrUserExists)
	}
	if err := secret.ValidatePassword(password); err != nil {
		return nil, errors.New(InvalidPassword)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	newUser := models.User{
		Email:    email,
		Password: string(hashedPassword),
		Username: username,
	}
	createdUser, err := service.UserRepository.Create(&newUser)
	if err != nil {
		return nil, errors.New(ErrUserExists)
	}
	_, err = service.SecretRepository.Create(string(hashedPassword), createdUser.ID)
	if err != nil {
		return nil, errors.New(ErrCreateSecret)
	}
	return &newUser, nil
}

// Login - авторизация пользователя
func (service *AuthService) Login(email, password string) (jwt.JWTData, error) {
	user, err := service.UserRepository.FindByEmail(email)
	if err != nil || user == nil {
		return jwt.JWTData{}, errors.New(InValidPasOrEmail)
	}
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password), []byte(password),
	); err != nil {
		return jwt.JWTData{}, errors.New(InValidPasOrEmail)
	}

	// Передаём и ID, и email
	return jwt.JWTData{
		UserID: user.ID,
		Email:  user.Email,
	}, nil
}

// RefreshTokens - обновление токенов
func (service *AuthService) RefreshTokens(refreshToken, accessToken string) (*jwt.TokenPair, error) {
	accessTokenValid, _ := service.JWT.ParseToken(accessToken)
	if accessTokenValid {
		return nil, errors.New(ErrAccessToken)
	}
	refreshTokenValid, data := service.JWT.ParseRefreshToken(refreshToken)
	if !refreshTokenValid || data == nil {
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

func (service *AuthService) ForgotPassword(config *configs.Config, email string) error {
	existUser, err := service.UserRepository.FindByEmail(email)
	if err != nil || existUser == nil {
		return errors.New(ErrUserNotFound)
	}
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return err
	}
	token := hex.EncodeToString(bytes)
	passwordReset := models.NewPasswordReset(existUser.ID, token)
	err = service.passwordResetRepository.Create(passwordReset)
	if err != nil {
		return errors.New(ErrCreatePasswordReset)
	}
	sendmail.SendEmailPasswordReset(config, email, token)
	return nil
}

func (service *AuthService) ResetPassword(token, newPassword string) (*models.User, error) {
	// Находим активный токен
	activeToken, err := service.passwordResetRepository.GetActiveToken(token)
	if err != nil {
		return nil, err
	}
	if activeToken == nil {
		return nil, errors.New(ErrTokenExpired)
	}

	// Находим пользователя по ID из токена
	user, err := service.UserRepository.FindByid(activeToken.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New(ErrUserNotFound)
	}

	// Валидируем пароль
	if err := secret.ValidatePassword(newPassword); err != nil {
		return nil, errors.New(InvalidPassword)
	}

	// Хэшируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Обновляем пароль пользователя
	user.Password = string(hashedPassword)
	updatedUser, err := service.UserRepository.Update(user)
	if err != nil {
		return nil, err
	}

	// Обновляем секрет
	oldSecret, err := service.SecretRepository.GetByUserID(updatedUser.ID)
	if err != nil {
		return nil, err
	}
	_, err = service.SecretRepository.Update(oldSecret.ID, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	// Помечаем токен как использованный
	if err := service.passwordResetRepository.TokenUsed(activeToken.ID); err != nil {
		return nil, err
	}

	return updatedUser, nil
}
