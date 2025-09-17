package models

import (
	"time"

	"gorm.io/gorm"
)

type PasswordReset struct {
	gorm.Model
	UserID    uint      `json:"user_id" gorm:"index"`
	Token     string    `json:"token" gorm:"uniqueIndex"`
	Used      bool      `json:"used"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewPasswordReset(userID uint, token string) *PasswordReset {
	return &PasswordReset{
		UserID:    userID,
		Token:     token,
		Used:      false,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}
}
