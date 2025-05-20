package eventParticipant

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/types"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"gorm.io/gorm"
)

type EventParticipantRepository struct {
	DataBase *db.Db
}

func NewEventParticipantRepository(dataBase *db.Db) *EventParticipantRepository {
	return &EventParticipantRepository{DataBase: dataBase}
}

// AddParticipant добавляет пользователя к событию
func (repo *EventParticipantRepository) AddParticipant(eventID, userID uint) error {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&EventParticipant{})
	participant := NewEventParticipant(eventID, userID)
	if err := db.Create(participant).Error; err != nil {
		return err
	}
	return nil
}

// RemoveParticipant удаляет пользователя из события
func (repo *EventParticipantRepository) RemoveParticipant(eventID, userID uint) error {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&EventParticipant{})
	if err := db.Where("event_id = ? AND user_id = ?", eventID, userID).
		Delete(&EventParticipant{}).Error; err != nil {
		return err
	}
	return nil
}

// GetEventParticipants возвращает список участников события
func (repo *EventParticipantRepository) GetEventParticipants(eventID uint) ([]user.User, error) {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&EventParticipant{})
	var users []user.User
	err := db.Table("users").
		Joins("JOIN event_participants ON users.id = event_participants.user_id").
		Where("event_participants.event_id = ?", eventID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserEvents возвращает список событий, в которых участвует пользователь
func (repo *EventParticipantRepository) GetUserEvents(userID uint) ([]types.Event, error) {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&EventParticipant{})
	var events []types.Event
	err := db.Table("events").
		Joins("JOIN event_participants ON events.id = event_participants.event_id").
		Where("event_participants.user_id = ?", userID).
		Find(&events).Error
	if err != nil {
		return nil, err
	}
	return events, nil
}

// IsParticipant проверяет, является ли пользователь участником события
func (repo *EventParticipantRepository) IsParticipant(eventID, userID uint) (bool, error) {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&EventParticipant{})
	var count int64
	err := db.Model(&EventParticipant{}).
		Where("event_id = ? AND user_id = ?", eventID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// IsEventCreatorById проверяет, является ли пользователь создателем события
func (repo *EventParticipantRepository) IsEventCreatorById(eventID, userID uint) (bool, error) {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&EventParticipant{})
	var count int64
	err := db.Table("events").
		Where("id = ? AND creator_id = ?", eventID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
