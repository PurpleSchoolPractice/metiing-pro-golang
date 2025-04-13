package event

import (
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"gorm.io/gorm"
)

type Event struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date" binding:"required"`
	CreatorID   uint      `json:"creator_id" gorm:"not null"`

	// Связи
	Creator user.User `gorm:"foreignKey:CreatorID"`
}

func NewEvent(title, description string, creatorID uint, eventDate time.Time) *Event {
	return &Event{
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		CreatorID:   creatorID,
	}
}
