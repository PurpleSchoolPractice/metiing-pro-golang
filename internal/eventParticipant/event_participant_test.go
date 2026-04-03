package eventParticipant

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestAddParticipant(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "event_participants"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			uint(1), uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mockDB.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.EventParticipant{}) // Установка модели таблицы
	err := repo.AddParticipant(uint(1), uint(1))
	require.NoError(t, err, "Add participant failed")
	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestRemoveParticipant(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`DELETE FROM "event_participants" WHERE event_id = $1 AND user_id = $2`).
		WithArgs(uint(1), uint(1)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.EventParticipant{}) // Установка модели таблицы
	err := repo.RemoveParticipant(uint(1), uint(1))
	require.NoError(t, err, "Remove participant failed")
	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetEventParticipants(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// Мокаем запрос на получение участников события
	mockDB.ExpectQuery(`SELECT \* FROM "event_participants" WHERE event_id = \$1`).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "event_id", "user_id", "status",
		}).
			AddRow(1, fixedTime, fixedTime, nil, 1, 10, "Принято").
			AddRow(2, fixedTime, fixedTime, nil, 1, 11, "Принято"))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)

	participants, err := repo.GetUsersWithInvites(uint(1))
	require.NoError(t, err, "Get event participants failed")
	require.Len(t, participants, 2)

	require.Equal(t, uint(1), participants[0].EventID)
	require.Equal(t, uint(10), participants[0].UserID)
	require.Equal(t, models.EventStatus("Принято"), participants[0].Status)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestGetUserEvents(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	require.NoError(t, err)

	// Мокаем запрос на получение событий пользователя через таблицу event_participants
	mockDB.ExpectQuery(`SELECT events.* FROM "events" 
                      JOIN event_participants ON events.id = event_participants.event_id 
                      WHERE event_participants.user_id = \$1`).
		WithArgs(uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "created_at", "updated_at", "deleted_at", "title", "description", "event_date", "creator_id",
		}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", date, 1))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)

	events, err := repo.GetUserEvents(uint(1))
	require.NoError(t, err, "Get user events failed")
	require.Len(t, events, 1)
	require.Equal(t, uint(1), events[0].ID)
	require.Equal(t, "testevent", events[0].Title)

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestIsParticipant(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)

	// Тест для случая, когда пользователь является участником
	mockDB.ExpectQuery(`SELECT count(*) FROM "event_participants" WHERE event_id = $1 AND user_id = $2`).
		WithArgs(uint(1), uint(1)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.EventParticipant{}) // Установка модели таблицы
	isParticipant, err := repo.IsParticipant(uint(1), uint(1))
	require.NoError(t, err, "Check participant failed")
	require.True(t, isParticipant, "User should be a participant")
	require.NoError(t, mockDB.ExpectationsWereMet())

	// Тест для случая, когда пользователь не является участником
	mockDB.ExpectQuery(`SELECT count(*) FROM "event_participants" WHERE event_id = $1 AND user_id = $2`).
		WithArgs(uint(1), uint(2)).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	isParticipant, err = repo.IsParticipant(uint(1), uint(2))
	require.NoError(t, err, "Check participant failed")
	require.False(t, isParticipant, "User should not be a participant")
	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestDeleteEventParticipant(t *testing.T) {

	//Тестовые данные
	testUserID := uint(42)
	testParticipantID := uint(54)
	testEventID := uint(33)
	testEmail := "test@example.com"

	//Создаем моковую базу данных с ожиданием выборки из events и удаления из таблицы event_participants
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	t.Cleanup(cleanup)

	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "events" WHERE (id = $1 AND creator_id = $2) AND "events"."deleted_at" IS NULL`)).
		WithArgs(testEventID, testUserID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))

	mockDB.ExpectBegin()

	mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "event_participants" SET "deleted_at"=$1 WHERE (event_id = $2 AND user_id = $3) AND "event_participants"."deleted_at" IS NULL`)).
		WithArgs(sqlmock.AnyArg(), testEventID, testParticipantID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventParticipantRepository(dbWrapper)
	repo.DataBase.DB = repo.DataBase.DB.Model(&models.EventParticipant{})

	// Создаем jwt
	newJWT := jwt.NewJWT("secret")

	//Созданим обработчик на основании моковой базы и jwt
	handler := &EventParticipantHandler{
		EventParticipantRepository: repo,
		JWTService:                 newJWT,
	}

	// Выполнение запроса
	r := chi.NewRouter()
	r.Delete("/event-participant/{id}/event/{event_id}", handler.DeleteEventParticipant())
	reader := bytes.NewReader([]byte{})
	w := httptest.NewRecorder()

	deleteMetod := fmt.Sprintf("/event-participant/%d/event/%d", testParticipantID, testEventID)
	req := httptest.NewRequest(http.MethodDelete, deleteMetod, reader)
	ctx := context.WithValue(req.Context(), middleware.ContextEmailKey, testEmail)
	ctx = context.WithValue(ctx, middleware.ContextUserIDKey, testUserID)
	req = req.WithContext(ctx)

	r.ServeHTTP(w, req)

	require.Equal(t, w.Code, 200)
	require.NoError(t, mockDB.ExpectationsWereMet())
}
