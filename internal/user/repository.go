package user

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
)

type UserRepository struct {
	DataBase *db.Db
}

func NewUserRepository(dataBase *db.Db) *UserRepository {
	return &UserRepository{
		DataBase: dataBase,
	}
}

func (repo *UserRepository) Create(user *User) (*User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{}) // Установка модели таблицы
	result := repo.DataBase.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{}) // Установка модели таблицы
	var user User
	result := repo.DataBase.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (repo *UserRepository) FindAllUsers() ([]User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{}) // Установка модели таблицы
	var users []User
	result := repo.DataBase.DB.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}
	return users, nil
}

func (repo *UserRepository) Update(user *User) (*User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{}) // Установка модели таблицы
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

func (repo *UserRepository) Delete(user *User) error {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{}) // Установка модели таблицы
	result := repo.DataBase.DB.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
