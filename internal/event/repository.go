package event

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"gorm.io/gorm/clause"
)

type EventRepository struct {
	DataBase *db.Db
}

func NewEventRepository(dataBase *db.Db) *EventRepository {
	return &EventRepository{
		DataBase: dataBase,
	}
}

func (repo *EventRepository) Create(event *Event) (*Event, error) {
	result := repo.DataBase.DB.Create(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

func (repo *EventRepository) FindById(id uint) (*Event, error) {
	var event Event
	result := repo.DataBase.DB.First(&event, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &event, nil
}

func (repo *EventRepository) FindAllByCreatorId(id uint) ([]Event, error) {
	var events []Event
	result := repo.DataBase.DB.Where("creatorID = ?", id).Find(&events)
	if result.Error != nil {
		return nil, result.Error
	}
	return events, nil
}

func (repo *EventRepository) Update(event *Event) (*Event, error) {
	result := repo.DataBase.DB.Clauses(clause.Returning{}).Updates(event)
	if result.Error != nil {
		return nil, result.Error
	}
	return event, nil
}

func (repo *EventRepository) DeleteById(id uint) error {
	result := repo.DataBase.DB.Delete(&Event{}, id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
