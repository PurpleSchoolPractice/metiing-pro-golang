package event

import (
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
)

type EventRepository struct {
	DataBase *db.Db
}

// NewEventRepository создает новый репозиторий событий
func NewEventRepository(dataBase *db.Db) *EventRepository {
	return &EventRepository{
		DataBase: dataBase,
	}
}

// Create создает новое событие в базе данных
func (repo *EventRepository) Create(event *models.Event) (*models.Event, error) {

	result := repo.DataBase.DB.Create(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

// FindById находит событие по его ID
func (repo *EventRepository) FindById(id uint) (*models.Event, error) {
	var event models.Event
	result := repo.DataBase.DB.First(&event, id)
	if result.Error != nil {
		return nil, result.Error
	}

	return &event, nil
}

// FindAllByCreatorId находит все события, созданные пользователем с указанным ID
func (repo *EventRepository) FindAllByCreatorId(id uint) ([]models.Event, error) {
	var events []models.Event
	result := repo.DataBase.DB.Preload("Creator").Where("creator_id = ?", id).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// Update обновляет информацию о событии в базе данных
func (repo *EventRepository) Update(event *models.Event) (*models.Event, error) {
	result := repo.DataBase.DB.Updates(event)
	if result.Error != nil {
		return nil, result.Error
	}

	return event, nil
}

// DeleteById удаляет событие по его ID из базы данных
func (repo *EventRepository) DeleteById(id uint) error {
	result := repo.DataBase.DB.Delete(&models.Event{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetEventWithCreator получает событие вместе с информацией о создателе
func (repo *EventRepository) GetEventWithCreator(eventID, userID uint) (*models.Event, error) {
	var event models.Event

	result := repo.DataBase.DB.Preload("Creator").Where("id = ? AND creator_id = ?", eventID, userID).First(&event)

	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// GetEventsWithCreators получает список событий с информацией о создателях
func (repo *EventRepository) GetEventsWithCreators() ([]models.Event, error) {
	var events []models.Event
	result := repo.DataBase.DB.Preload("Creator").Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// ищем пересекающиеся события
func (r *EventRepository) IsUserBusy(userID uint, start time.Time, duration int) bool {
	end := start.Add(time.Duration(duration) * time.Minute)
	var count int64

	r.DataBase.DB.
		Model(&models.Event{}).
		Joins("JOIN event_participants ep ON ep.event_id = events.id").
		Where("ep.user_id = ?", userID).
		Where(`
			(events.start_date, events.start_date + (events.duration || ' minutes')::interval)
			OVERLAPS (?, ?)`,
			start, end).
		Count(&count)

	return count > 0
}
