package eventParticipant

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/types"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
)

type EventParticipantRepository struct {
	DataBase *db.Db
}

func NewEventParticipantRepository(dataBase *db.Db) *EventParticipantRepository {
	return &EventParticipantRepository{
		DataBase: dataBase,
	}
}

// AddParticipant добавляет пользователя к событию
func (repo *EventParticipantRepository) AddParticipant(eventID, userID uint) error {
	participant := NewEventParticipant(eventID, userID)
	result := repo.DataBase.DB.Create(participant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// RemoveParticipant удаляет пользователя из события
func (repo *EventParticipantRepository) RemoveParticipant(eventID, userID uint) error {
	result := repo.DataBase.DB.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&EventParticipant{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetEventParticipants возвращает список участников события
func (repo *EventParticipantRepository) GetEventParticipants(eventID uint) ([]user.User, error) {
	var users []user.User
	result := repo.DataBase.DB.Table("users").
		Joins("JOIN event_participants ON users.id = event_participants.user_id").
		Where("event_participants.event_id = ?", eventID).
		Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// GetUserEvents возвращает список событий, в которых участвует пользователь
func (repo *EventParticipantRepository) GetUserEvents(userID uint) ([]types.Event, error) {
	var events []types.Event
	result := repo.DataBase.DB.Table("events").
		Joins("JOIN event_participants ON events.id = event_participants.event_id").
		Where("event_participants.user_id = ?", userID).
		Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// IsParticipant проверяет, является ли пользователь участником события
func (repo *EventParticipantRepository) IsParticipant(eventID, userID uint) (bool, error) {
	var count int64
	result := repo.DataBase.DB.Model(&EventParticipant{}).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
