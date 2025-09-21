// migrations_integration_test.go
package migrations_test

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/migrations"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockLogger struct{}

func (l *MockLogger) Info(msg string, args ...any)  {}
func (l *MockLogger) Error(msg string, args ...any) {}
func TestUserModelInit(t *testing.T) {
	db, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	mockLog := &MockLogger{}
	//Проверяем созданы ли таблицы
	mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
		WithArgs("users", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		//создаем таблицы
	mock.ExpectExec(`CREATE TABLE "users"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
		//создаем индекс
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS "idx_users_deleted_at"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
		//проверяем таблицу
	mock.ExpectQuery(`SELECT count\(\*\) FROM "users" WHERE "users"\."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))
		//ожидаем
	mock.ExpectBegin()
	//создаем дефолтных
	mock.ExpectQuery(`INSERT INTO "users" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	mock.ExpectCommit()

	users, err := migrations.UserModelInit(db, mockLog)

	require.NoError(t, err)
	require.NotEmpty(t, users)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestSecretInit(t *testing.T) {
	db, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	mockLog := &MockLogger{}
	mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
		WithArgs("secrets", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(`CREATE TABLE "secrets"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS "idx_secrets_user_id"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS "idx_secrets_deleted_at"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectQuery(`SELECT count\(\*\) FROM "secrets" WHERE "secrets"\."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()

	mock.ExpectQuery(`INSERT INTO "secrets" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	mock.ExpectCommit()

	hash1, _ := bcrypt.GenerateFromPassword([]byte("Test1Test1!2021"), bcrypt.DefaultCost) //user Test1
	hash2, _ := bcrypt.GenerateFromPassword([]byte("Test2Test2!2022"), bcrypt.DefaultCost) //user Test2
	users := []*models.User{
		{

			Username: "Test1",
			Password: string(hash1),
			Email:    "test1@test1.ru",
		},
		{

			Username: "Test2",
			Password: string(hash2),
			Email:    "test2@test2.ru",
		},
	}
	err := migrations.SecretModelInit(db, mockLog, users)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestEventInit(t *testing.T) {
	db, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()
	mockLog := &MockLogger{}
	// Проверка существования таблицы
	mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
		WithArgs("events", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Создание таблицы
	mock.ExpectExec(`CREATE TABLE "events"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
		// Создание индекса
	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS "idx_events_deleted_at"`).
		WillReturnResult(sqlmock.NewResult(0, 0))
	// проверяем количества записей
	mock.ExpectQuery(`SELECT count\(\*\) FROM "events" WHERE "events"\."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()

	// Вставка новых событий
	mock.ExpectQuery(`INSERT INTO "events" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	mock.ExpectCommit()
	users := []*models.User{
		{

			Username: "Test1",
			Email:    "test1@test1.ru",
		},
		{

			Username: "Test2",
			Email:    "test2@test2.ru",
		},
	}
	_, err := migrations.EventModelInit(db, mockLog, users)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
func TestEventParticipantInit(t *testing.T) {
	db, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()
	mockLog := &MockLogger{}

	mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
		WithArgs("event_participants", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectExec(`CREATE TABLE "event_participants"`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectExec(`CREATE INDEX IF NOT EXISTS "idx_event_participants_deleted_at"`).
		WillReturnResult(sqlmock.NewResult(0, 0))

	mock.ExpectQuery(`SELECT count\(\*\) FROM "event_participants" WHERE "event_participants"\."deleted_at" IS NULL`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "event_participants" .* RETURNING "id"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1).AddRow(2))

	mock.ExpectCommit()
	events := []*models.Event{
		{
			Title:       "Test title",
			Description: "Test about description my testing",
			StartDate:   time.Now(),
			Duration:    20,
			CreatorID:   1,
		},
		{
			Title:       "Head of comunication",
			Description: "Meet with workers in my company for test",
			StartDate:   time.Now(),
			Duration:    30,
			CreatorID:   2,
		},
	}
	err := migrations.EventParticipantModelInit(db, mockLog, events)

	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPreviousPasswordInit(t *testing.T) {
	db, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	mockLog := &MockLogger{}

	mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
		WithArgs("previous_passwords", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// checkTable вызовет SELECT count(*) FROM "previous_passwords"
	mock.ExpectQuery(`SELECT count\(\*\) FROM "previous_passwords"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Вставка записи (gorm делает INSERT ... RETURNING "id")
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "previous_passwords" .* RETURNING "id"`).
		WithArgs(1, "testPassword", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Вызов миграции
	err := migrations.PreviousPasswordModelInit(db, mockLog)
	require.NoError(t, err)

	// Проверка, что все ожидания выполнены
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestPasswordResetModelInit(t *testing.T) {
	gdb, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	mockLog := &MockLogger{}

	// Проверка, есть ли таблица
	mock.ExpectQuery(`SELECT count\(\*\) FROM information_schema.tables WHERE table_schema = CURRENT_SCHEMA\(\) AND table_name = \$1 AND table_type = \$2`).
		WithArgs("password_resets", "BASE TABLE").
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Проверка количества записей в таблице
	mock.ExpectQuery(`SELECT count\(\*\) FROM "password_resets"`).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

	// Вставка записи (7 аргументов из-за gorm.Model)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "password_resets" .* RETURNING "id"`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			nil,              // deleted_at
			2,                // user_id
			"testToken",      // token
			false,            // used
			sqlmock.AnyArg(), // expires_at
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	// Вызов миграции
	err := migrations.PasswordResetModelInit(gdb, mockLog)
	require.NoError(t, err)

	// Проверка, что все ожидания выполнены
	require.NoError(t, mock.ExpectationsWereMet())
}
