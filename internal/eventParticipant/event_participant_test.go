package eventParticipant

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
)

func TestAddParticipant(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "event_participants"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			uint(1), uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	err := repo.AddParticipant(uint(1), uint(1))
	require.NoError(t, err, "Add participant failed")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestRemoveParticipant(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "event_participants" WHERE event_id = $1 AND user_id = $2`).
		WithArgs(uint(1), uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	err := repo.RemoveParticipant(uint(1), uint(1))
	require.NoError(t, err, "Remove participant failed")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEventParticipants(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// Мокаем запрос на получение участников
	mock.ExpectQuery(`SELECT users.* FROM "users" JOIN event_participants ON users.id = event_participants.user_id WHERE event_participants.event_id = $1`).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "username", "email",
		}).AddRow(1, fixedTime, fixedTime, nil, "testuser", "test@example.com"))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	participants, err := repo.GetEventParticipants(uint(1))
	require.NoError(t, err, "Get event participants failed")
	require.Len(t, participants, 1)
	require.Equal(t, uint(1), participants[0].ID)
	require.Equal(t, "testuser", participants[0].Username)
	require.Equal(t, "test@example.com", participants[0].Email)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserEvents(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	if err != nil {
		t.Fatalf("Not possible to parse date: %v", err)
	}

	// Мокаем запрос на получение событий пользователя
	mock.ExpectQuery(`SELECT events.* FROM "events" JOIN event_participants ON events.id = event_participants.event_id WHERE event_participants.user_id = $1`).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "title", "description", "event_date", "creator_id",
		}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", date, 1))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	events, err := repo.GetUserEvents(uint(1))
	require.NoError(t, err, "Get user events failed")
	require.Len(t, events, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestIsParticipant(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)

	// Тест для случая, когда пользователь является участником
	mock.ExpectQuery(`SELECT count(*) FROM "event_participants" WHERE event_id = $1 AND user_id = $2`).
		WithArgs(uint(1), uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&EventParticipant{}) // Установка модели таблицы
	isParticipant, err := repo.IsParticipant(uint(1), uint(1))
	require.NoError(t, err, "Check participant failed")
	require.True(t, isParticipant, "User should be a participant")
	require.NoError(t, mock.ExpectationsWereMet())

	// Тест для случая, когда пользователь не является участником
	mock.ExpectQuery(`SELECT count(*) FROM "event_participants" WHERE event_id = $1 AND user_id = $2`).
		WithArgs(uint(1), uint(2)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isParticipant, err = repo.IsParticipant(uint(1), uint(2))
	require.NoError(t, err, "Check participant failed")
	require.False(t, isParticipant, "User should not be a participant")
	require.NoError(t, mock.ExpectationsWereMet())
}
