package eventParticipant

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"gorm.io/gorm"
)

// EventParticipant представляет связь "многие ко многим" между событиями и пользователями.
type EventParticipant struct {
	gorm.Model
	EventID uint `json:"event_id" gorm:"not null"`
	UserID  uint `json:"user_id" gorm:"not null"`

	// Связи
	Event *event.Event `gorm:"foreignKey:EventID"`
	User  *user.User   `gorm:"foreignKey:UserID"`
}

// NewEventParticipant создает новую связь между событием и участником
func NewEventParticipant(eventID, userID uint) *EventParticipant {
	return &EventParticipant{
		EventID: eventID,
		UserID:  userID,
	}
}
