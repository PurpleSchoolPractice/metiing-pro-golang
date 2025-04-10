package types

import "time"

// Event представляет базовую структуру события
type Event struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventDate   time.Time `json:"event_date"`
	CreatorID   uint      `json:"creator_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// EventParticipant представляет связь между событием и участником
type EventParticipant struct {
	ID        uint      `json:"id"`
	EventID   uint      `json:"event_id"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// EventRepository определяет интерфейс для работы с событиями
type EventRepository interface {
	Create(event *Event) (*Event, error)
	FindById(id uint) (*Event, error)
	FindAllByCreatorId(id uint) ([]Event, error)
	Update(event *Event) (*Event, error)
	DeleteById(id uint) error
	GetEventWithCreator(id uint) (*Event, error)
	GetEventsWithCreators() ([]Event, error)
}

// EventParticipantRepository определяет интерфейс для работы с участниками событий
type EventParticipantRepository interface {
	AddParticipant(eventID, userID uint) error
	RemoveParticipant(eventID, userID uint) error
	GetEventParticipants(eventID uint) ([]User, error)
	GetUserEvents(userID uint) ([]Event, error)
	IsParticipant(eventID, userID uint) (bool, error)
}
