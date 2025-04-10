package event

import (
	"time"

	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	CreatorID   uint      `json:"creatorID" gorm:"not null"`

	// Связи
	Creator interface{} `gorm:"foreignKey:CreatorID"`
}

func NewEvent(title, description string, creatorID uint, eventDate time.Time) *Event {
	return &Event{
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		CreatorID:   creatorID,
	}
}
