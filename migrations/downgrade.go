package migrations

import (
	"os"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/eventParticipant"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// удаляем дефолтные значения
func Migrate() error {
	logging := logger.NewLogger(configs.LoadConfig())
	err := godotenv.Load(".env")
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		logging.Error(err.Error())
		return err
	}
	logging.Info("База данных запущена")
	id := []int{1, 2}
	logging.Info("Deleting event participants...")
	if err := database.Unscoped().Where("id IN (?)", id).Delete(&eventParticipant.EventParticipant{}).Error; err != nil {
		logging.Error("Failed to delete default eventParticipant")
		return err
	}
	logging.Info("Deleting events...")
	if err := database.Unscoped().Where("id IN (?)", []int{1, 2}).Delete(&event.Event{}).Error; err != nil {
		logging.Error("Failed to delete default events")
		return err
	}
	logging.Info("Emails to delete users")
	users := []string{"test1@test1.ru", "test2@test2.ru"}
	if err := database.Unscoped().Where("email IN (?)", users).Delete(&user.User{}).Error; err != nil {
		logging.Error("Failed to delete users with email test1@test1.ru and test2@test2.ru")
		return err
	}
	if err := database.Unscoped().Where("id IN (?)", []int{1, 2}).Delete(&secret.Secret{}).Error; err != nil {
		logging.Error("Failed to delete default secrets")
		return err
	}
	//либо сделать полностью удаление таблиц
	//сбрасываем последовтельность счетчика со всех таблиц
	if err := database.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
		logging.Error("Failed to reset users_id_seq sequence")
		return err
	}
	if err := database.Exec("ALTER SEQUENCE users_id_seq RESTART WITH 1").Error; err != nil {
		logging.Error("Failed to reset users_id_seq sequence")
		return err
	}
	if err := database.Exec("ALTER SEQUENCE secrets_id_seq RESTART WITH 1").Error; err != nil {
		logging.Error("Failed to reset secrets_id_seq sequence")
		return err
	}
	if err := database.Exec("ALTER SEQUENCE events_id_seq RESTART WITH 1").Error; err != nil {
		logging.Error("Failed to reset events_id_seq sequence")
		return err
	}
	if err := database.Exec("ALTER SEQUENCE event_participants_id_seq RESTART WITH 1").Error; err != nil {
		logging.Error("Failed to reset event_participants_id_seq sequence")
		return err
	}
	return nil
}
