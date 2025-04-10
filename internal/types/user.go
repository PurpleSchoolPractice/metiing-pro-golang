package types

import "time"

type User struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRepository interface {
	Create(user *User) (*User, error)
	FindById(id uint) (*User, error)
	FindByEmail(email string) (*User, error)
	Update(user *User) (*User, error)
	DeleteById(id uint) error
} 