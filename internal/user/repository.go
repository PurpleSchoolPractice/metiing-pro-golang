package user

import (
	"errors"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"gorm.io/gorm"
)

type UserRepository struct {
	DataBase *db.Db
}

// NewUserRepository создает новый экземпляр UserRepository с заданным объектом базы данных.
func NewUserRepository(dataBase *db.Db) *UserRepository {
	return &UserRepository{
		DataBase: dataBase,
	}
}

// Create создает новую запись в базе данных
func (r *UserRepository) Create(u *User) (*User, error) {
	if err := r.DataBase.
		Session(&gorm.Session{NewDB: true}).
		Create(u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail находит пользователя по указанному адресу электронной почты в базе данных.
func (r *UserRepository) FindByEmail(email string) (*User, error) {
	var u User
	err := r.DataBase.
		Session(&gorm.Session{NewDB: true}).
		Where("email = ?", email).
		First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// FindAllUsers находит всех пользователей в базе данных.
func (repo *UserRepository) FindAllUsers() ([]User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{})
	var users []User
	result := repo.DataBase.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

// Update обновляет информацию о пользователе в базе данных.
func (repo *UserRepository) Update(user *User) (*User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{})
	result := repo.DataBase.DB.Model(&User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
		"username": user.Username,
		"password": user.Password,
		"email":    user.Email,
	})
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// Delete удаляет пользователя из базы данных.
func (repo *UserRepository) Delete(user *User) error {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{})
	result := repo.DataBase.DB.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) FindByid(id uint) (*User, error) {
	var u User
	err := r.DataBase.
		Session(&gorm.Session{NewDB: true}).
		Where("id = ?", id).
		First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
