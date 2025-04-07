package event

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	OwnerID     uint      `json:"ownerID"`
	CreatorID   uint      `json:"creatorID" gorm:"not null"`
}

func NewEvent(title, description string, ownerID, creatorID uint, eventDate time.Time) *Event {
	return &Event{
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		OwnerID:     ownerID,
		CreatorID:   creatorID,
	}
}
