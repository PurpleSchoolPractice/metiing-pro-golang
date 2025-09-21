package passwordReset

import (
	"errors"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"gorm.io/gorm"
)

type PasswordResetRepository struct {
	DataBase *db.Db
	logger   *logger.Logger
}

func NewPasswordResetRepository(db *db.Db, logger *logger.Logger) *PasswordResetRepository {
	return &PasswordResetRepository{
		DataBase: db,
		logger:   logger,
	}
}

// Create создает новый временный токен для отправки на почту
func (r *PasswordResetRepository) Create(passwordReset *models.PasswordReset) error {
	if err := r.DataBase.Create(passwordReset); err != nil {
		r.logger.Error("Error saving temporary token", "error", err.Error())
		return err
	}
	r.logger.Info("Added new temporary token", "passwordReset_id", passwordReset.ID, "user_id", passwordReset.UserID)
	return nil
}

// Ищет активный временный токен
func (r *PasswordResetRepository) GetActiveToken(token string) (*models.PasswordReset, error) {
	var passwordReset models.PasswordReset
	result := r.DataBase.DB.Take(
		&passwordReset,
		"token = ? AND used = ? AND expires_at > ?", token, false, time.Now(),
	)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &passwordReset, nil
}

// Помечает временный токен как использованный
func (r *PasswordResetRepository) TokenUsed(tokenID uint) error {
	result := r.DataBase.DB.Model(models.PasswordReset{}).
		Where("id = ? AND expires_at > ?", tokenID, time.Now()).
		Update("used", true)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no active tokens found")
	}
	return nil
}
