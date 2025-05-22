package user

type UserUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"unique index" json:"email" validate:"required,email"`
}
