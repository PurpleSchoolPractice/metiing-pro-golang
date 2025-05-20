package event

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (repo *EventRepository) Create(event *Event) (*Event, error) {
	db := repo.DataBase.DB.
		Session(&gorm.Session{NewDB: true}).
		Model(&Event{})
	result := db.Model(&Event{}).Create(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

// FindById находит событие по его ID
func (repo *EventRepository) FindById(id uint) (*Event, error) {
	if repo.DataBase.DB.Statement.Model == nil {
		repo.DataBase.DB = repo.DataBase.DB.Model(&Event{})
	}
	var event Event
	result := repo.DataBase.DB.First(&event, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// FindAllByCreatorId находит все события, созданные пользователем с указанным ID
func (repo *EventRepository) FindAllByCreatorId(id uint) ([]Event, error) {
	if repo.DataBase.DB.Statement.Model == nil {
		repo.DataBase.DB = repo.DataBase.DB.Model(&Event{})
	}
	var events []Event
	result := repo.DataBase.DB.Where("creator_id = ?", id).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

// Update обновляет информацию о событии в базе данных
func (repo *EventRepository) Update(event *Event) (*Event, error) {
	if repo.DataBase.DB.Statement.Model == nil {
		repo.DataBase.DB = repo.DataBase.DB.Model(&Event{})
	}
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

// DeleteById удаляет событие по его ID из базы данных
func (repo *EventRepository) DeleteById(id uint) error {
	if repo.DataBase.DB.Statement.Model == nil {
		repo.DataBase.DB = repo.DataBase.DB.Model(&Event{})
	}
	result := repo.DataBase.DB.Delete(&Event{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetEventWithCreator получает событие вместе с информацией о создателе
func (repo *EventRepository) GetEventWithCreator(id uint) (*Event, error) {
	if repo.DataBase.DB.Statement.Model == nil {
		repo.DataBase.DB = repo.DataBase.DB.Model(&Event{})
	}
	var event Event
	result := repo.DataBase.DB.Preload("Creator").First(&event, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

// GetEventsWithCreators получает список событий с информацией о создателях
func (repo *EventRepository) GetEventsWithCreators() ([]Event, error) {
	if repo.DataBase.DB.Statement.Model == nil {
		repo.DataBase.DB = repo.DataBase.DB.Model(&Event{})
	}
	var events []Event
	result := repo.DataBase.DB.Preload("Creator").Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}
