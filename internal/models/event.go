package models

import (
	"time"

	"gorm.io/gorm"
)

// Статусы приглашений
type EventStatus string

const (
	StatusAccepted EventStatus = "Принято"
	StatusBusy     EventStatus = "Занят"
	StatusDecline  EventStatus = "Отклонено"
	StatusSent     EventStatus = "Отправлено"
)

type UserStatus struct {
	UserId   uint        //для добавления участников в обработчике
	UserName string      `json:"user_name"`
	Status   EventStatus `json:"status"`
}
type Event struct {
	gorm.Model
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	StartDate   time.Time `json:"start_date" `
	Duration    int       `json:"duration_min"`
	CreatorID   uint      `json:"creator_id" gorm:"not null"`

	// Связи
	Creator *User `gorm:"foreignKey:CreatorID;constraint:OnDelete:CASCADE"` //для API чтобы в некоторых случаях было NULL, а не пустые поля.
}

// NewEvent создает новый объект события
func NewEvent(title, description string, duration int, creatorID uint, startDate time.Time) *Event {
	return &Event{
		Title:       title,
		Description: description,
		StartDate:   startDate,
		Duration:    duration,
		CreatorID:   creatorID,
	}
}

// EventRepository определяет интерфейс для работы с событиями
type EventRepository interface {
	Create(event *Event) (*Event, error)
	FindById(id uint) (*Event, error)
	FindAllByCreatorId(id uint) ([]Event, error)
	Update(event *Event) (*Event, error)
	DeleteById(id uint) error
	GetEventWithCreator(eventID, userID uint) (*Event, error)
	GetEventsWithCreators() ([]Event, error)
}
