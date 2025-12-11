package user

import "github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"

type UserUpdateRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"unique index" json:"email" validate:"required,email"`
}

type UserPaginatedResponse struct {
	Items  []models.UserResponse `json:"items"`
    Total  int64                 `json:"total"`
    Limit  int                   `json:"limit"`
    Offset int                   `json:"offset"`
}