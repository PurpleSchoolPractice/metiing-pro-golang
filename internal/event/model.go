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
	OwnerID     uint      `json:"owner_id"`

	// Связи
	Creator *user.User `gorm:"foreignKey:CreatorID"` //для API чтобы в некоторых случаях было NULL, а не пустые поля.
}

// NewEvent создает новый объект события
func NewEvent(title, description string, creatorID uint, eventDate time.Time) *Event {
	return &Event{
		Title:       title,
		Description: description,
		EventDate:   eventDate,
		CreatorID:   creatorID,
	}
}
