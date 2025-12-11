package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"unique index" json:"email"`
}

type UserResponse struct {
	ID        uint       `json:"id"`
    CreatedAt time.Time  `json:"created_at"`
    UpdatedAt time.Time  `json:"updated_at"`
    Username  string     `json:"username"`
    Email     string     `json:"email"`
}

func NewUser(email string, password string, name string) *User {
	return &User{
		Email:    email,
		Password: password,
		Username: name,
	}
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindById(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	FindAllUsers(limit, offset int, search string) ([]UserResponse, int64, error)
	Update(user *User) (*User, error)
	DeleteById(id uint) error
}
