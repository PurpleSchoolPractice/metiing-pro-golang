package main

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

func main() {
	logging := logger.NewLogger(configs.LoadConfig())
	err := godotenv.Load(".env")
	if err != nil {
		logging.Info(err.Error())
	}
	database, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_DSN")), &gorm.Config{})
	if err != nil {
		logging.Info(err.Error())
	}
<<<<<<< HEAD
	//удаляет все записи с БД
	//database.Migrator().DropTable(&user.User{}, &secret.Secret{}, &event.Event{}, &eventParticipant.EventParticipant{})
=======

>>>>>>> master
	err = database.AutoMigrate(
		&user.User{},
		&secret.Secret{},
		&event.Event{},
		&eventParticipant.EventParticipant{},
	)
	if err != nil {
		logging.Info(err.Error())
	}
}
