package models

import "gorm.io/gorm"

// EventParticipantRepository определяет интерфейс для работы с участниками событий
type EventParticipantRepository interface {
	AddParticipant(eventID, userID uint) error
	RemoveParticipant(eventID, userID uint) error
	GetEventParticipants(eventID uint) ([]User, error)
	GetUserEvents(userID uint) ([]Event, error)
	IsParticipant(eventID, userID uint) (bool, error)
}

// EventParticipant представляет связь "многие ко многим" между событиями и пользователями через внешний ключ
type EventParticipant struct {
	gorm.Model
	EventID uint        `json:"event_id" gorm:"not null"`
	UserID  uint        `json:"user_id" gorm:"not null"`
	Status  EventStatus `json:"status" gorm:"type:varchar(255);default:'Принято'"`
	// Связи
	Event *Event `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	User  *User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// NewEventParticipant создает новую связь между событием и участником
func NewEventParticipant(eventID, userID uint) *EventParticipant {
	return &EventParticipant{
		EventID: eventID,
		UserID:  userID,
	}
}
