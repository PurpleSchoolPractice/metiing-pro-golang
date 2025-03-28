package secret

import (
	"errors"
	"time"
	"unicode"

	"gorm.io/gorm"
)

type Secret struct {
	gorm.Model
	CurrentPassword   string             `json:"currentPassword"`
	PreviousPasswords []PreviousPassword `gorm:"foreignKey:SecretID"`
}

type PreviousPassword struct {
	ID        uint `gorm:"primarykey"`
	SecretID  uint
	Password  string
	CreatedAt time.Time
}

func NewSecret(password string) (*Secret, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	return &Secret{
		CurrentPassword: password,
	}, nil
}

// ValidatePassword проверяет соответствие пароля политике безопасности
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("пароль должен содержать не менее 12 символов")
	}

	var hasDigit, hasUpper, hasSpecial bool
	for _, char := range password {
		switch {
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasDigit {
		return errors.New("пароль должен содержать хотя бы одну цифру")
	}
	if !hasUpper {
		return errors.New("пароль должен содержать хотя бы одну заглавную букву")
	}
	if !hasSpecial {
		return errors.New("пароль должен содержать хотя бы один специальный символ")
	}

	return nil
}

// IsDifferentFromPrevious проверяет, отличается ли пароль от предыдущих
func (s *Secret) IsDifferentFromPrevious(password string) bool {
	if s.CurrentPassword == password {
		return false
	}

	for _, prev := range s.PreviousPasswords {
		if prev.Password == password {
			return false
		}
	}

	return true
}
