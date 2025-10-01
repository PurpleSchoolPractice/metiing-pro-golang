package event

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateEvent(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "events"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"testevent", "description", "2025-05-09 11:02", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.Event{}) // Установка модели таблицы
	date, err := time.Parse("2006-01-02 15:04", "2025-05-09 11:02")
	require.NoError(t, err)

	newEvent := models.NewEvent("testevent", "description", 30, 1, date)
	createdEvent, err := repo.Create(newEvent)
	require.NoError(t, err, "Create event failed")
	require.NotEmpty(t, createdEvent, "Create event should not be empty")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindEventByID(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"title",
		"description",
		"event_date",
		"creator_id",
	}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", date, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE id = $1`)).
		WithArgs(uint(1)).WillReturnRows(rows)

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.Event{}) // Установка модели таблицы
	foundEvent, err := repo.FindById(uint(1))
	require.NoError(t, err, "Error finding event")
	require.NotNil(t, foundEvent, "Expected to find event")
	require.Equal(t, uint(1), foundEvent.ID)
	require.Equal(t, "testevent", foundEvent.Title)
	require.Equal(t, "description", foundEvent.Description)
	require.Equal(t, date, foundEvent.StartDate)
	require.Equal(t, uint(1), foundEvent.CreatorID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByCreatorID(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	require.NoError(t, err)
	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"title",
		"description",
		"event_date",
		"creator_id",
	}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", date, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE creator_id = $1`)).
		WithArgs(uint(1)).WillReturnRows(rows)
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.Event{}) // Установка модели таблицы
	findedEvents, err := repo.FindAllByCreatorId(uint(1))
	require.NoError(t, err, "Error finding events")
	require.Len(t, findedEvents, 1)
	require.Equal(t, uint(1), findedEvents[0].ID)
	require.Equal(t, "testevent", findedEvents[0].Title)
	require.Equal(t, "description", findedEvents[0].Description)
	require.Equal(t, date, findedEvents[0].StartDate)
	require.Equal(t, uint(1), findedEvents[0].CreatorID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEvent(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "events" SET`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"newtestevent", "newdescription", "2025-05-10 11:02", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	date, err := time.Parse("2006-01-02 15:01", "2025-05-10 11:02")
	require.NoError(t, err)

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.Event{}) // Установка модели таблицы
	eventUpdate := &models.Event{
		Title:       "newtestevent",
		Description: "newdescription",
		CreatorID:   1,
		StartDate:   date,
	}
	updatedEvent, err := repo.Update(eventUpdate)
	require.NoError(t, err, "Error updating event")
	require.Equal(t, uint(1), updatedEvent.ID)
	require.Equal(t, "newtestevent", updatedEvent.Title)
	require.Equal(t, "newdescription", updatedEvent.Description)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteEventByID(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "events" WHERE id = $1`).
		WithArgs(uint(1)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.Event{}) // Установка модели таблицы
	err := repo.DeleteById(uint(1))
	require.NoError(t, err, "Error deleting event")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEventWithCreator(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	require.NoError(t, err)

	// Мокаем запрос на получение события
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE id = $1`)).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "title", "description", "event_date", "creator_id",
		}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", date, 1))

	// Мокаем запрос на получение создателя
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."id" = $1`)).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "username", "email",
		}).AddRow(1, fixedTime, fixedTime, nil, "testuser", "test@example.com"))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.Event{}) // Установка модели таблицы
	event, err := repo.GetEventWithCreator(uint(1))
	require.NoError(t, err, "Error getting event with creator")
	require.NotNil(t, event, "Event should not be nil")
	require.Equal(t, uint(1), event.ID)
	require.Equal(t, "testevent", event.Title)
	require.Equal(t, uint(1), event.CreatorID)
	require.NoError(t, mock.ExpectationsWereMet())
}
