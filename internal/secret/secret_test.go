package secret_test

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func mockLogger() *logger.Logger {
	var cfg configs.Config
	loggerMock := logger.NewLogger(&cfg)
	return loggerMock
}

func TestCreateSecret(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Тест на успешное создание
	t.Run("Valid password", func(t *testing.T) {
		mockDB.ExpectBegin()
		mockDB.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "secrets"`)).
			WithArgs(
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				sqlmock.AnyArg(),
				"ValidPass123!",
			).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mockDB.ExpectCommit()

		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		createdSecret, err := repo.Create("ValidPass123!")
		require.NoError(t, err)
		require.Equal(t, uint(1), createdSecret.ID)
		require.Equal(t, "ValidPass123!", createdSecret.CurrentPassword)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	// Тест на невалидный пароль
	t.Run("Invalid password", func(t *testing.T) {
		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		_, err := repo.Create("short")
		require.Error(t, err)
		require.Contains(t, err.Error(), "не менее 12 символов")
	})
}

func TestGetSecretByID(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "current_password"}).
		AddRow(1, fixedTime, fixedTime, nil, "Password123!")

	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "secrets"`)).
		WithArgs(1).
		WillReturnRows(rows)

	// Тест на успешное получение
	t.Run("Existing secret", func(t *testing.T) {
		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		foundSecret, err := repo.GetByID(1)
		require.NoError(t, err)
		require.Equal(t, uint(1), foundSecret.ID)
		require.Equal(t, "Password123!", foundSecret.CurrentPassword)
	})

	// Тест на несуществующий секрет
	t.Run("Not found", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "secrets"`)).
			WithArgs(999).
			WillReturnError(gorm.ErrRecordNotFound)

		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		_, err := repo.GetByID(999)
		require.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}

func TestUpdateSecret(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)

	// Подготовка тестовых данных
	secretRows := sqlmock.NewRows([]string{"id", "current_password"}).
		AddRow(1, "OldPassword123!")

	prevPasswordsRows := sqlmock.NewRows([]string{"id", "password"}).
		AddRow(1, "OldPassword123!").
		AddRow(2, "OldPassword456!").
		AddRow(3, "OldPassword789!").
		AddRow(4, "OldPassword000!").
		AddRow(5, "OldPassword111!")

	t.Run("Successful update with history cleanup", func(t *testing.T) {
		// Mock GetByID
		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "secrets"`)).
			WithArgs(1).
			WillReturnRows(secretRows)
		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "previous_passwords"`)).
			WillReturnRows(prevPasswordsRows)

		// Transaction
		mockDB.ExpectBegin()
		// Insert previous password
		mockDB.ExpectExec(regexp.QuoteMeta(`INSERT INTO "previous_passwords"`)).
			WillReturnResult(sqlmock.NewResult(6, 1))
		// Delete oldest password
		mockDB.ExpectExec(regexp.QuoteMeta(`DELETE FROM "previous_passwords"`)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		// Update secret
		mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "secrets"`)).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.ExpectCommit()

		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		_, err := repo.Update(1, "NewPassword123!")
		require.NoError(t, err)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})

	t.Run("Duplicate password", func(t *testing.T) {
		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "secrets"`)).
			WithArgs(1).
			WillReturnRows(secretRows)
		mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "previous_passwords"`)).
			WillReturnRows(prevPasswordsRows)

		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		_, err := repo.Update(1, "OldPassword123!")
		require.Error(t, err)
		require.Contains(t, err.Error(), "должен отличаться")
	})
}

func TestDeleteSecret(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)

	t.Run("Successful delete", func(t *testing.T) {
		mockDB.ExpectBegin()
		mockDB.ExpectExec(regexp.QuoteMeta(`UPDATE "secrets" SET "deleted_at"=`)).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(1, 1))
		mockDB.ExpectCommit()

		loggerMock := mockLogger()
		repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

		err := repo.Delete(1)
		require.NoError(t, err)
		require.NoError(t, mockDB.ExpectationsWereMet())
	})
}

func TestListSecrets(t *testing.T) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)

	rows := sqlmock.NewRows([]string{"id", "current_password"}).
		AddRow(1, "Pass1").
		AddRow(2, "Pass2")

	mockDB.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "secrets"`)).
		WithArgs(10, 0).
		WillReturnRows(rows)

	loggerMock := mockLogger()
	repo := secret.NewSecretRepository(&db.Db{DB: gormDB}, loggerMock)

	secrets, err := repo.List(10, 0)
	require.NoError(t, err)
	require.Len(t, secrets, 2)
	require.Equal(t, "Pass1", secrets[0].CurrentPassword)
}
