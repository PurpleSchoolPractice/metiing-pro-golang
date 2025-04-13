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
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	participant := NewEventParticipant(eventID, userID)
	result := repo.DataBase.DB.Model(&EventParticipant{}).Create(participant)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// RemoveParticipant удаляет пользователя из события
func (repo *EventParticipantRepository) RemoveParticipant(eventID, userID uint) error {
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	result := repo.DataBase.DB.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&EventParticipant{})
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetEventParticipants возвращает список участников события
func (repo *EventParticipantRepository) GetEventParticipants(eventID uint) ([]user.User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Table("users") // Установка таблицы
	var users []user.User
	result := repo.DataBase.DB.
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
	repo.DataBase.DB = repo.DataBase.DB.Table("events") // Установка таблицы
	var events []types.Event
	result := repo.DataBase.DB.
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
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{})
	var count int64
	result := repo.DataBase.DB.
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// IsEventCreatorById проверяет, является ли пользователь создателем события
func (repo *EventParticipantRepository) IsEventCreatorById(eventID, userID uint) (bool, error) {
	repo.DataBase.DB = repo.DataBase.DB.Table("events")
	var count int64
	result := repo.DataBase.DB.
		Where("id = ? AND creator_id = ?", eventID, userID).
		Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}
