package secret

import (
	"errors"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
)

type SecretRepository struct {
	DataBase *db.Db
	logger   *logger.Logger
}

func NewSecretRepository(db *db.Db, logger *logger.Logger) *SecretRepository {
	return &SecretRepository{
		DataBase: db,
		logger:   logger,
	}
}

// Create создает новый секрет
func (r *SecretRepository) Create(password string, userID uint) (*Secret, error) {
	// Проверка пароля на соответствие политике
	if err := ValidatePassword(password); err != nil {
		r.logger.Error("Error validating password", "error", err.Error())
		return nil, err
	}

	secret, err := NewSecret(password, userID)
	if err != nil {
		r.logger.Error("Error creating secret", "error", err.Error())
		return nil, err
	}

	if err := r.DataBase.Create(secret); err != nil {
		r.logger.Error("Error saving secret", "error", err.Error())
		return nil, err
	}

	r.logger.Info("Added new secret", "secret_id", secret.ID, "user_id", userID)
	return secret, nil
}

// GetByID получает секрет по ID
func (r *SecretRepository) GetByID(id uint) (*Secret, error) {
	r.logger.Debug("Get secret by ID", "id", id)

	var secret Secret
	if err := r.DataBase.Preload("PreviousPasswords").First(&secret, id).Error; err != nil {
		r.logger.Error("Error getting secret", "id", id, "error", err.Error())
		return nil, err
	}

	r.logger.Info("Secret retrieved successfully", "id", id)
	return &secret, nil
}

// Update обновляет пароль и сохраняет предыдущий
func (r *SecretRepository) Update(id uint, newPassword string) (*Secret, error) {
	r.logger.Debug("Request to update password", "id", id)

	// Получаем текущий секрет
	secret, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Проверка пароля на соответствие политике
	if err := ValidatePassword(newPassword); err != nil {
		r.logger.Error("Error validating password", "error", err.Error())
		return nil, err
	}

	// Проверка, что пароль отличается от предыдущих
	if !secret.IsDifferentFromPrevious(newPassword) {
		err := errors.New("новый пароль должен отличаться от 5 предыдущих")
		r.logger.Error("Password not different from previous passwords", "id", id)
		return nil, err
	}

	// Сохраняем текущий пароль в историю
	prevPassword := PreviousPassword{
		SecretID:  secret.ID,
		Password:  secret.CurrentPassword,
		CreatedAt: time.Now(),
	}

	// Транзакция для атомарности операций
	r.logger.Debug("Beginning transaction", "id", id)
	tx := r.DataBase.Begin()

	// Добавляем текущий пароль в историю
	if err := tx.Create(&prevPassword).Error; err != nil {
		tx.Rollback()
		r.logger.Error("Error adding current password to history", "id", id, "error", err.Error())
		return nil, err
	}

	// Удаляем самые старые записи, если их больше 5
	var count int64
	tx.Model(&PreviousPassword{}).Where("secret_id = ?", secret.ID).Count(&count)

	if count >= 5 {
		r.logger.Debug("Deleting oldest password from history", "id", id, "count", count)
		var oldestPassword PreviousPassword
		if err := tx.Where("secret_id = ?", secret.ID).Order("created_at").First(&oldestPassword).Error; err == nil {
			tx.Delete(&oldestPassword)
		}
	}

	// Обновляем текущий пароль
	secret.CurrentPassword = newPassword
	if err := tx.Save(secret).Error; err != nil {
		tx.Rollback()
		r.logger.Error("Error updating current password", "id", id, "error", err.Error())
		return nil, err
	}

	// Подтверждаем транзакцию
	if err := tx.Commit().Error; err != nil {
		r.logger.Error("Error committing transaction", "id", id, "error", err.Error())
		return nil, err
	}

	r.logger.Info("Password updated successfully", "id", id)
	return secret, nil
}

// Delete удаляет секрет по ID
func (r *SecretRepository) Delete(id uint) error {
	result := r.DataBase.DB.Delete(id)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// List возвращает список всех секретов
func (r *SecretRepository) List(limit, offset int) ([]Secret, error) {
	r.logger.Debug("Request to list secrets", "limit", limit, "offset", offset)

	var secrets []Secret
	if err := r.DataBase.Limit(limit).Offset(offset).Find(&secrets).Error; err != nil {
		r.logger.Error("Error getting secrets list", "error", err.Error())
		return nil, err
	}

	r.logger.Info("Retrieved %d secrets", len(secrets), "limit", limit, "offset", offset)
	return secrets, nil
}
