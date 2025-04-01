package secret

import (
	"errors"
	"time"
	"unicode"

	"gorm.io/gorm"
)

type Secret struct {
	gorm.Model
	UserID            uint               `json:"user_id" gorm:"index"`
	CurrentPassword   string             `json:"currentPassword"`
	PreviousPasswords []PreviousPassword `gorm:"foreignKey:SecretID"`
}

type PreviousPassword struct {
	ID        uint `gorm:"primarykey"`
	SecretID  uint
	Password  string
	CreatedAt time.Time
}

func NewSecret(password string, userID uint) (*Secret, error) {
	if err := ValidatePassword(password); err != nil {
		return nil, err
	}

	return &Secret{
		UserID:          userID,
		CurrentPassword: password,
	}, nil
}

// ValidatePassword проверяет соответствие пароля политике безопасности
func ValidatePassword(password string) error {
	if len(password) < 12 {
		return errors.New("Password must be at least 12 characters long")
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
		return errors.New("Password must contain at least one digit")
	}
	if !hasUpper {
		return errors.New("Password must contain at least one uppercase")
	}
	if !hasSpecial {
		return errors.New("Password must contain at least one special character")
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
