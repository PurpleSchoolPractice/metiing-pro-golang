package event

import "time"

// EventRequest используется для создания или обновления события.
type EventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	CreatorID   uint      `json:"creator_id" binding:"required"`
}

// EventResponse используется для возврата информации о событии клиенту.
type EventResponse struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatorID   uint      `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
