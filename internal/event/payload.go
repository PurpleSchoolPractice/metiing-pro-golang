package event

import "time"

// EventRequest используется для создания или обновления события.
type EventRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	OwnerID     uint      `json:"owner_id" binding:"required"`
	CreatorID   uint      `json:"creator_id" binding:"required"`
}

// EventResponse используется для возврата информации о событии клиенту.
type EventResponse struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	OwnerID     uint      `json:"owner_id"`
	CreatorID   uint      `json:"creator_id"`
}
