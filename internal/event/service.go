package event

import (
	"errors"
	types2 "github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/types"
	"time"
)

type EventService struct {
	EventRepo            types2.EventRepository
	EventParticipantRepo types2.EventParticipantRepository
}

func NewEventService(
	eventRepo types2.EventRepository,
	eventParticipantRepo types2.EventParticipantRepository,
) *EventService {
	return &EventService{
		EventRepo:            eventRepo,
		EventParticipantRepo: eventParticipantRepo,
	}
}

// CreateEvent создает новое событие
func (s *EventService) CreateEvent(title, description string, creatorID uint, eventDate string) (*types2.Event, error) {
	date, err := time.Parse(time.RFC3339, eventDate)
	if err != nil {
		return nil, errors.New("invalid date format")
	}

	event := &types2.Event{
		Title:       title,
		Description: description,
		EventDate:   date,
		CreatorID:   creatorID,
	}
	return s.EventRepo.Create(event)
}

// GetEvent получает информацию о событии
func (s *EventService) GetEvent(id uint) (*types2.Event, error) {
	return s.EventRepo.GetEventWithCreator(id)
}

// GetEventsByCreator получает список событий создателя
func (s *EventService) GetEventsByCreator(creatorID uint) ([]types2.Event, error) {
	return s.EventRepo.FindAllByCreatorId(creatorID)
}

// UpdateEvent обновляет информацию о событии
func (s *EventService) UpdateEvent(id uint, title, description string, eventDate string) (*types2.Event, error) {
	event, err := s.EventRepo.FindById(id)
	if err != nil {
		return nil, err
	}

	if title != "" {
		event.Title = title
	}
	if description != "" {
		event.Description = description
	}
	if eventDate != "" {
		date, err := time.Parse(time.RFC3339, eventDate)
		if err != nil {
			return nil, errors.New("invalid date format")
		}
		event.EventDate = date
	}

	return s.EventRepo.Update(event)
}

// DeleteEvent удаляет событие
func (s *EventService) DeleteEvent(id uint) error {
	return s.EventRepo.DeleteById(id)
}

// AddParticipant добавляет участника к событию
func (s *EventService) AddParticipant(eventID, userID uint) error {
	return s.EventParticipantRepo.AddParticipant(eventID, userID)
}

// RemoveParticipant удаляет участника из события
func (s *EventService) RemoveParticipant(eventID, userID uint) error {
	return s.EventParticipantRepo.RemoveParticipant(eventID, userID)
}

// GetEventParticipants получает список участников события
func (s *EventService) GetEventParticipants(eventID uint) ([]types2.User, error) {
	return s.EventParticipantRepo.GetEventParticipants(eventID)
}

// GetUserEvents получает список событий пользователя
func (s *EventService) GetUserEvents(userID uint) ([]types2.Event, error) {
	return s.EventParticipantRepo.GetUserEvents(userID)
}
