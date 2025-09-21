package migrations

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
)

// удаляем дефолтные значения
func Migrate() error {
	database, logging, err := InitDBWithLogger()
	if err != nil {
		logging.Error(err.Error())
		return err
	}

	id := []int{1, 2}
	if err := database.Unscoped().Where("id IN (?)", id).Delete(&models.EventParticipant{}).Error; err != nil {
		logging.Error("Failed to delete default eventParticipant")
		return err
	}
	if err := database.Unscoped().Where("id IN (?)", []int{1, 2}).Delete(&models.Event{}).Error; err != nil {
		logging.Error("Failed to delete default events")
		return err
	}
	users := []string{"test1@test1.ru", "test2@test2.ru"}
	if err := database.Unscoped().Where("email IN (?)", users).Delete(&models.User{}).Error; err != nil {
		logging.Error("Failed to delete users with email test1@test1.ru and test2@test2.ru")
		return err
	}
	if err := database.Unscoped().Where("id IN (?)", []int{1, 2}).Delete(&secret.Secret{}).Error; err != nil {
		logging.Error("Failed to delete default secrets")
		return err
	}
	if err := database.Unscoped().Where("user_id = ?", 2).Delete(&models.PasswordReset{}).Error; err != nil {
		logging.Error("Failed to delete default passwordReset")
		return err
	}
	if err := database.Unscoped().Where("user_id = ?", 2).Delete(&secret.PreviousPassword{}).Error; err != nil {
		logging.Error("Failed to delete default passwordReset")
		return err
	}
	return nil
}
func DeleteAllTableWithDate() error {
	database, logging, err := InitDBWithLogger()
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	database.Migrator().DropTable(&models.User{}, &secret.Secret{}, &models.Event{}, &models.EventParticipant{}, &models.PasswordReset{})
	logging.Info("All tables has deleted")
	return nil
}
