package user

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"uniqueIndex" json:"email"`

	EventsCreated []event.Event `gorm:"foreignKey:CreatorID"`
	EventsOwned   []event.Event `gorm:"foreignKey:OwnerID"`
}

func NewUser(email string, password string, name string) *User {
	return &User{
		Email:    email,
		Password: password,
		Username: name,
	}
}
