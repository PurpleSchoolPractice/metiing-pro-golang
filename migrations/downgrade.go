package migrations

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/eventParticipant"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
)

// удаляем дефолтные значения
func Migrate() {
	database, logging, err := InitDBWithLogger()
	if err != nil {
		logging.Error((err.Error()))
		return
	}
	if err := database.Unscoped().Where("email IN", []string{"test1@test1.ru", "test2@test2.ru"}).Delete(&user.User{}).Error; err != nil {
		logging.Error("Failed to delete users with email test1@test1.ru and test2@test2.ru")
		return
	}
	if err := database.Unscoped().Where("id IN", []int{1, 2}).Delete(&secret.Secret{}).Error; err != nil {
		logging.Error("Failed to delete default secrets")
		return
	}
	if err := database.Unscoped().Where("id IN", []int{1, 2}).Delete(&event.Event{}).Error; err != nil {
		logging.Error("Failed to delete default events")
		return
	}
	if err := database.Unscoped().Where("id IN", []int{1, 2}).Delete(&eventParticipant.EventParticipant{}).Error; err != nil {
		logging.Error("Failed to delete default eventParticipant")
		return
	}
}
