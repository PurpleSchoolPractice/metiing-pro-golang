package passwordReset_test

import (
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/passwordReset"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestGetActiveToken(t *testing.T) {
	// Настроить моковую БД
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup() // Очистить ресурсы после теста

	// Определить ожидаемый запрос и ответ
	expectedQuery := regexp.QuoteMeta(`SELECT * FROM "password_resets" WHERE (token = $1 AND used = $2 AND expires_at > $3) AND "password_resets"."deleted_at" IS NULL LIMIT $4`)

	// Замокать данные строки (подстроить поля под модель PasswordReset)
	fixedTime := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	futureTime := fixedTime.Add(time.Hour)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "user_id", "token", "used", "expires_at"}).
		AddRow(1, fixedTime, fixedTime, nil, 123, "testToken", false, futureTime)

	// Настроить ожидание: запрос с конкретными аргументами и вернуть замоканные строки
	mock.ExpectQuery(expectedQuery).
		WithArgs("testToken", false, sqlmock.AnyArg(), 1).
		WillReturnRows(rows)

	// Обернуть моковый GORM DB в структуру db.Db
	dbWrapper := &db.Db{DB: gormDB}

	// Создать репозиторий на основе моковой БД (реальный конфиг/логгер не нужен)
	pr := passwordReset.NewPasswordResetRepository(dbWrapper, nil)

	// Вызов метода
	res, err := pr.GetActiveToken("testToken")

	// Проверка
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, "testToken", res.Token)

	// Проверка всех ожиданий
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetActiveToken_RecordNotFound(t *testing.T) {
	// Настроить моковую БД
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Определить ожидаемый запрос
	expectedQuery := regexp.QuoteMeta(`SELECT * FROM "password_resets" WHERE (token = $1 AND used = $2 AND expires_at > $3) AND "password_resets"."deleted_at" IS NULL LIMIT $4`)

	// Настроить ожидание возврата no rows (запись не найдена)
	mock.ExpectQuery(expectedQuery).
		WithArgs("nonexistentToken", false, sqlmock.AnyArg(), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// Обернуть моковый GORM DB в структуру db.Db
	dbWrapper := &db.Db{DB: gormDB}
	pr := passwordReset.NewPasswordResetRepository(dbWrapper, nil)

	// Выполнить вызов
	res, err := pr.GetActiveToken("nonexistentToken")

	// Проверка: GetActiveToken возвращает nil без ошибки при ErrRecordNotFound
	require.NoError(t, err)
	require.Nil(t, res)

	// Проверка всех ожиданий
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetActiveToken_DatabaseError(t *testing.T) {
	// Настроить моковую БД
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Определить ожидаемый запрос
	expectedQuery := regexp.QuoteMeta(`SELECT * FROM "password_resets" WHERE (token = $1 AND used = $2 AND expires_at > $3) AND "password_resets"."deleted_at" IS NULL LIMIT $4`)

	// Настроить ожидание возвращения ошибки БД
	dbError := errors.New("connection refused")
	mock.ExpectQuery(expectedQuery).
		WithArgs("testToken", false, sqlmock.AnyArg(), 1).
		WillReturnError(dbError)

	// Обернуть моковый GORM DB в структуру db.Db
	dbWrapper := &db.Db{DB: gormDB}
	pr := passwordReset.NewPasswordResetRepository(dbWrapper, nil)

	// Выполнить вызов
	res, err := pr.GetActiveToken("testToken")

	// Проверка: GetActiveToken возвращает ошибку и nil результат
	require.Error(t, err)
	require.Nil(t, res)
	require.Equal(t, "connection refused", err.Error())

	// Проверка всех ожиданий
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenUsed_Success(t *testing.T) {
	// Настроить моковую БД
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Определить ожидаемый UPDATE-запрос
	expectedQuery := regexp.QuoteMeta(`UPDATE "password_resets" SET "used"=$1,"updated_at"=$2 WHERE (id = $3 AND expires_at > $4) AND "password_resets"."deleted_at" IS NULL`)

	// Настроить ожидание успешного обновления (1 строка затронута)
	mock.ExpectBegin()
	mock.ExpectExec(expectedQuery).
		WithArgs(true, sqlmock.AnyArg(), uint(1), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Обернуть моковый GORM DB в структуру db.Db
	dbWrapper := &db.Db{DB: gormDB}
	pr := passwordReset.NewPasswordResetRepository(dbWrapper, nil)

	// Выполнить вызов
	err := pr.TokenUsed(1)

	// Проверка
	require.NoError(t, err)

	// Проверка всех ожиданий
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenUsed_NoRowsAffected(t *testing.T) {
	// Настроить моковую БД
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Определить ожидаемый UPDATE-запрос
	expectedQuery := regexp.QuoteMeta(`UPDATE "password_resets" SET "used"=$1,"updated_at"=$2 WHERE (id = $3 AND expires_at > $4) AND "password_resets"."deleted_at" IS NULL`)

	// Настроить ожидание без обновленных строк
	mock.ExpectBegin()
	mock.ExpectExec(expectedQuery).
		WithArgs(true, sqlmock.AnyArg(), uint(999), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 0))
	mock.ExpectCommit()

	// Обернуть моковый GORM DB в структуру db.Db
	dbWrapper := &db.Db{DB: gormDB}
	pr := passwordReset.NewPasswordResetRepository(dbWrapper, nil)

	// Выполнить вызов
	err := pr.TokenUsed(999)

	// Проверка: TokenUsed возвращает ошибку "no active tokens found"
	require.Error(t, err)
	require.Equal(t, "no active tokens found", err.Error())

	// Проверка всех ожиданий
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestTokenUsed_DatabaseError(t *testing.T) {
	// Настроить моковую БД
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Определить ожидаемый UPDATE-запрос
	expectedQuery := regexp.QuoteMeta(`UPDATE "password_resets" SET "used"=$1,"updated_at"=$2 WHERE (id = $3 AND expires_at > $4) AND "password_resets"."deleted_at" IS NULL`)

	// Настроить ожидание возврата ошибки БД
	dbError := errors.New("database connection failed")
	mock.ExpectBegin()
	mock.ExpectExec(expectedQuery).
		WithArgs(true, sqlmock.AnyArg(), uint(1), sqlmock.AnyArg()).
		WillReturnError(dbError)
	mock.ExpectRollback()

	// Обернуть моковый GORM DB в структуру db.Db
	dbWrapper := &db.Db{DB: gormDB}
	pr := passwordReset.NewPasswordResetRepository(dbWrapper, nil)

	// Выполнить вызов
	err := pr.TokenUsed(1)

	// Проверка
	require.Error(t, err)
	require.Equal(t, "database connection failed", err.Error())

	// Проверка всех ожиданий
	require.NoError(t, mock.ExpectationsWereMet())
}
