package migrations

import (
	"os"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDBWithLogger() (*gorm.DB, *logger.Logger, error) {
	logging := logger.NewLogger(configs.LoadConfig())
	err := godotenv.Load(".env")
	if err != nil {
		logging.Error(err.Error())
		return nil, nil, err
	}
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		logging.Error(err.Error())
		return nil, nil, err
	}
	return database, logging, nil
}

// миграция таблиц с создаными дефолтными значениями
func InitModelMigration() {
	database, logging, err := InitDBWithLogger()
	if err != nil {
		logging.Error((err.Error()))
		return
	}
	//удаляет все записи с БД
	//database.Migrator().DropTable(&user.User{}, &secret.Secret{}, &event.Event{}, &eventParticipant.EventParticipant{})

	userSecret1, userSecret2, err := UserModelInit(database, logging)
	if err != nil {
		logging.Error(err.Error())
	}
	err = SecretModelInit(database, logging, userSecret1, userSecret2)
	if err != nil {
		logging.Error(err.Error())
	}
	err = EventModelInit(database, logging)
	if err != nil {
		logging.Error(err.Error())
	}
	err = EventParticipantModelInit(database, logging)
	if err != nil {
		logging.Error(err.Error())
	}
	//Автомиграция таблиц
	//err = database.AutoMigrate(
	//	&user.User{},
	//	&secret.Secret{},
	//	&event.Event{},
	//	&eventParticipant.EventParticipant{},
	//)
	//if err != nil {
	//	logging.Info(err.Error())
	//}
}
