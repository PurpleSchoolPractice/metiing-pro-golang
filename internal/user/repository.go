package user

import (
	"errors"
	"strings"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
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
func (r *UserRepository) Create(u *models.User) (*models.User, error) {
	if err := r.DataBase.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

// FindByEmail находит пользователя по указанному адресу электронной почты в базе данных.
func (r *UserRepository) FindByEmail(email string) (*models.User, error) {
	var u models.User
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
func (repo *UserRepository) FindAllUsers(limit, offset int, search string) ([]models.UserResponse, int64, error) {
	result := repo.DataBase.DB.Model(&models.User{}).Where("deleted_at is null")

	if search != "" {
		search = strings.TrimSpace(search)
		search = "%" + search + "%"
		result = result.Where("username ILIKE ? OR email ILIKE ?", search, search)
	}

	var total int64
	result.Count(&total)

	result = result.Limit(limit).Offset(offset)

	var users []models.UserResponse
	selectedUsers := result.Omit("password", "deleted_at").Scan(&users)

	if selectedUsers.Error != nil {
		return nil, 0, selectedUsers.Error
	}
	return users, total, nil
}

// Update обновляет информацию о пользователе в базе данных.
func (repo *UserRepository) Update(user *models.User) (*models.User, error) {
	result := repo.DataBase.DB.Model(&models.User{}).Where("id = ?", user.ID).Updates(map[string]interface{}{
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
func (repo *UserRepository) Delete(user *models.User) error {
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.User{})
	result := repo.DataBase.DB.Delete(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
func (r *UserRepository) FindByid(id uint) (*models.User, error) {
	var u models.User
	err := r.DataBase.
		Session(&gorm.Session{NewDB: true}).
		Where("id = ?", id).
		First(&u).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return &u, nil
}
