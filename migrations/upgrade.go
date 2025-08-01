package migrations

import (
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func checkTable[T any](db *gorm.DB) (int64, error) {
	var payload T
	if !db.Migrator().HasTable(&payload) {
		db.AutoMigrate(&payload)
	}
	//проверяем записи в таблице
	var count int64
	if err := db.Model(&payload).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// проверка и заполнение тестовыми данными таблицы User
func UserModelInit(db *gorm.DB, logger logger.LoggerInterface) ([]*models.User, error) {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[models.User](db)
	if err != nil {
		return nil, err
	}

	//если нет записей заполняем таблицу
	hashedPassword1, _ := bcrypt.GenerateFromPassword([]byte("Test1Test1!2021"), bcrypt.DefaultCost) //user Test1
	hashedPassword2, _ := bcrypt.GenerateFromPassword([]byte("Test2Test2!2022"), bcrypt.DefaultCost) //user Test2
	var users []*models.User
	if count == 0 {
		users = []*models.User{
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
			return nil, result.Error
		}
		logger.Info("Table User created")
	} else {
		logger.Info("Table User has records already")
		db.Find(&users)
	}
	return users, nil
}

// проверка и заполнение тестовыми данными таблицы Secret
func SecretModelInit(db *gorm.DB, logger logger.LoggerInterface, users []*models.User) error {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[secret.Secret](db)
	if err != nil {
		return err
	}
	var secrets []secret.Secret
	//если нет записей заполняем таблицу
	for _, user := range users {

		secrets = []secret.Secret{
			{
				UserID:          user.ID,
				CurrentPassword: user.Password,
			},
		}
	}
	if count == 0 {
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
func EventModelInit(db *gorm.DB, logger logger.LoggerInterface, users []*models.User) ([]*models.Event, error) {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[models.Event](db)
	if err != nil {
		return nil, err
	}
	var events []*models.Event
	//если нет записей заполняем таблицу
	for i := range users {

		events = []*models.Event{
			{
				Title:       "Test title",
				Description: "Test about description my testing",
				StartDate:   time.Now(),
				Duration:    20,
				CreatorID:   users[i].ID,
			},
			{
				Title:       "Head of comunication",
				Description: "Meet with workers in my company for test",
				StartDate:   time.Now(),
				Duration:    30,
				CreatorID:   users[i+1].ID,
			},
		}
		break
	}
	if count == 0 {
		if result := db.Create(&events); result.Error != nil {
			return nil, result.Error
		}
		logger.Info("Table Event created")
	} else {
		logger.Info("Table Events has records already")
		db.Find(&events)
	}

	return events, nil
}

// проверка и заполнение тестовыми данными таблицы  EventParticipant
func EventParticipantModelInit(db *gorm.DB, logger logger.LoggerInterface, events []*models.Event) error {
	//проверка существования таблицы, если нет то создаем
	count, err := checkTable[models.EventParticipant](db)
	if err != nil {
		return err
	}
	//если нет записей заполняем таблицу
	var eventsPart []*models.EventParticipant
	for i := range events {

		eventsPart = []*models.EventParticipant{
			{
				EventID: events[i].ID,
				UserID:  events[i].CreatorID,
				Status:  models.StatusAccepted,
			},
			{
				EventID: events[i+1].ID,
				UserID:  events[i+1].CreatorID,
				Status:  models.StatusAccepted,
			},
		}
		break
	}
	if count == 0 {
		if result := db.Create(&eventsPart); result.Error != nil {
			return result.Error
		}
		logger.Info("Table EventPart created")
	} else {
		logger.Info("Table EventPart has records already")
	}

	return nil
}
