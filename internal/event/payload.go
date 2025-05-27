package event

import "time"

// EventRequest представляет данные для создания или обновления события
type EventRequest struct {
	Title       string    `json:"title" validate:"required"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" validate:"required"`
	CreatorID   uint      `json:"creator_id" validate:"required"`
}

// EventResponse представляет данные для ответа о событии
type EventResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatorID   uint      `json:"creator_id"`
}
type DeleteResponse struct {
	Delete bool `json:"delete"`
}
