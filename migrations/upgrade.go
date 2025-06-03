package migrations

import (
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/event"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/eventParticipant"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func checkTable[T any](db *gorm.DB) (int64, error) {
	var payload T
	if !db.Migrator().HasTable(&payload) {
		db.Migrator().CreateTable(&payload)
	}
	//проверяем записи в таблице
	var count int64
	if err := db.Model(&payload).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// проверка и заполнение тестовыми данными таблицы User
func UserModelInit(db *gorm.DB, logger *logger.Logger) (string, string, error) {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[user.User](db)
	if err != nil {
		return "", "", err
	}

	//если нет записей заполняем таблицу
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("Test1Test1!2021"), bcrypt.DefaultCost) //user Test1
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("Test2Test2!2022"), bcrypt.DefaultCost) //user Test2
	if count == 0 {
		users := []*user.User{
			{

				Username: "Test1",
				Password: string(hashedPassword1),
				Email:    "test1@test1.ru",
			},
			{

				Username: "Test2",
				Password: string(hashedPassword2),
				Email:    "test2@test2.ru",
			},
		}
		if result := db.Create(&users); result.Error != nil {
			return "", "", result.Error
		}
		logger.Info("Table User created")
	} else {
		logger.Info("Table User has records already")
	}
	return string(hashedPassword1), string(hashedPassword2), nil
}

// проверка и заполнение тестовыми данными таблицы Secret
func SecretModelInit(db *gorm.DB, logger *logger.Logger, userSecret1, userSecret2 string) error {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[secret.Secret](db)
	if err != nil {
		return err
	}
	//если нет записей заполняем таблицу
	if count == 0 {
		secrets := []*secret.Secret{
			{
				UserID:          1,
				CurrentPassword: userSecret1,
			},
			{
				UserID:          2,
				CurrentPassword: userSecret2,
			},
		}

		if result := db.Create(&secrets); result.Error != nil {
			return result.Error
		}
		logger.Info("Table Secret created")
	} else {
		logger.Info("Table Secret has records already")
	}
	return nil
}

// проверка и заполнение тестовыми данными таблицы Event
func EventModelInit(db *gorm.DB, logger *logger.Logger) error {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[event.Event](db)
	if err != nil {
		return err
	}
	//если нет записей заполняем таблицу
	if count == 0 {
		events := []*event.Event{
			{
				Title:       "Test title",
				Description: "Test about description my testing",
				EventDate:   time.Now(),
				CreatorID:   1,
				OwnerID:     1,
			},
			{
				Title:       "Head of comunication",
				Description: "Meet with workers in my company for test",
				EventDate:   time.Now(),
				CreatorID:   2,
				OwnerID:     2,
			},
		}
		if result := db.Create(&events); result.Error != nil {
			return result.Error
		}
		logger.Info("Table Event created")
	} else {
		logger.Info("Table Events has records already")
	}
	return nil
}

// проверка и заполнение тестовыми данными таблицы  EventParticipant
func EventParticipantModelInit(db *gorm.DB, logger *logger.Logger) error {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[eventParticipant.EventParticipant](db)
	if err != nil {
		return err
	}
	//если нет записей заполняем таблицу
	if count == 0 {
		eventsPart := []*eventParticipant.EventParticipant{
			{
				EventID: 1,
				UserID:  1,
			},
			{
				EventID: 2,
				UserID:  2,
			},
		}
		if result := db.Create(&eventsPart); result.Error != nil {
			return result.Error
		}
		logger.Info("Table EventPart created")
	} else {
		logger.Info("Table EventPart has records already")
	}
	return nil
}
