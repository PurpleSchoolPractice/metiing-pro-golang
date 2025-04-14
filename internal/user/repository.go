package user

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
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
func (repo *UserRepository) Create(user *User) (*User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{})
	result := repo.DataBase.DB.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	return user, nil
}

// FindByEmail находит пользователя по указанному адресу электронной почты в базе данных.
func (repo *UserRepository) FindByEmail(email string) (*User, error) {
	repo.DataBase.DB = repo.DataBase.DB.Model(&User{})
	var user User
	result := repo.DataBase.DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
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
