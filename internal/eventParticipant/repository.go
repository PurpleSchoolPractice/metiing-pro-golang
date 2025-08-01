package eventParticipant

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
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
		Model(&models.EventParticipant{})
	participant := models.NewEventParticipant(eventID, userID)
	if err := db.Create(participant).Error; err != nil {
		return err
	}
	return nil
}

// AddParticipant обновляет статус пользователя
func (repo *EventParticipantRepository) UpdateParticipant(eventPart *models.EventParticipant) (*models.EventParticipant, error) {
	result := repo.DataBase.DB.Save(eventPart)
	if result.Error != nil {
		return nil, result.Error
	}
	return eventPart, nil
}

// RemoveParticipant удаляет пользователя из события
func (repo *EventParticipantRepository) RemoveParticipant(eventID, userID uint) error {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&models.EventParticipant{})
	if err := db.Where("event_id = ? AND user_id = ?", eventID, userID).
		Delete(&models.EventParticipant{}).Error; err != nil {
		return err
	}
	return nil
}

// GetEventParticipants возвращает список участников события
func (repo *EventParticipantRepository) GetEventParticipants(eventID uint) ([]models.User, error) {
	var users []models.User
	err := repo.DataBase.DB.
		Table("event_participants").
		Joins("JOIN users ON users.id = event_participants.user_id").
		Where("event_participants.event_id = ? AND event_participants.deleted_at IS NULL AND users.deleted_at IS NULL", eventID).
		Select("users.id, users.username, users.email"). // Выбираем только нужные поля
		Scan(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserEvents возвращает список событий, в которых участвует пользователь
func (repo *EventParticipantRepository) GetUserEvents(userID uint) ([]models.Event, error) {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&models.EventParticipant{})
	var events []models.Event
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
		Model(&models.EventParticipant{})
	var count int64
	err := db.Model(&models.EventParticipant{}).
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
		Model(&models.EventParticipant{})
	var count int64
	err := db.Table("events").
		Where("id = ? AND creator_id = ?", eventID, userID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// получаем пользователей с приглашениями по событию
func (repo EventParticipantRepository) GetUsersWithInvites(eventID uint) (*models.EventParticipant, error) {
	var inviteUsers *models.EventParticipant
	result := repo.DataBase.DB.Where("event_id = ?", eventID).Find(&inviteUsers)
	if result.Error != nil {
		return nil, result.Error
	}
	return inviteUsers, nil
}
