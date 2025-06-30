package event

import "github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"

// приглашенные пользователи
type InviteUsers struct {
	UserName string `json:"user_name"`
	UserId   uint   `json:"user_id"`
}

// EventRequest представляет данные для создания или обновления события
type EventRequest struct {
	Title        string        `json:"title" validate:"required"`
	Description  string        `json:"description"`
	StartDate    string        `json:"start_date" `
	Duration     int           `json:"duration"`
	CreatorID    uint          `json:"creator_id" validate:"required"`
	InvatedUsers []InviteUsers `json:"invated_users"`
}

// EventResponse представляет данные для ответа о событии
type EventResponse struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	StartDate   string `json:"start_date" `
	Duration    int    `json:"duration"`
	Status      []models.UserStatus
}
type DeleteResponse struct {
	Delete bool `json:"delete"`
}
