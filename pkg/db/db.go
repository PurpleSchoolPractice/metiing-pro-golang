package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
)

// Storage определяет интерфейс для работы с хранилищем данных
type Storage interface {
	Create(value interface{}) error
	Find(dest interface{}, conditions ...interface{}) error
	First(dest interface{}, conditions ...interface{}) error
	Update(value interface{}) error
	Delete(value interface{}) error
}

// Db реализует интерфейс Storage
type Db struct {
	*gorm.DB
}

// Убедимся, что Db реализует интерфейс Storage
var _ Storage = (*Db)(nil)

// NewDB создает новое подключение к базе данных
func NewDB(conf *configs.Config) *Db {
	db, err := gorm.Open(postgres.Open(conf.Database.DSN), &gorm.Config{})
	if err != nil {
		return nil
		//TODO установить логирование
	}
	return &Db{db}
}

// Create создает новую запись в базе данных
func (d *Db) Create(value interface{}) error {
	return d.DB.Create(value).Error
}

// Find находит записи, соответствующие условиям
func (d *Db) Find(dest interface{}, conditions ...interface{}) error {
	return d.DB.Find(dest, conditions...).Error
}

// First находит первую запись, соответствующую условиям
func (d *Db) First(dest interface{}, conditions ...interface{}) error {
	return d.DB.First(dest, conditions...).Error
}

// Update обновляет запись в базе данных
func (d *Db) Update(value interface{}) error {
	return d.DB.Save(value).Error
}

// Delete удаляет запись из базы данных
func (d *Db) Delete(value interface{}) error {
	return d.DB.Delete(value).Error
}
