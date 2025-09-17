package auth_test

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/auth"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/passwordReset"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func setupAuthService(t *testing.T) (*auth.AuthService, sqlmock.Sqlmock, func()) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	database := &db.Db{DB: gormDB}
	userRepo := user.NewUserRepository(database)

	// Создаем обычный логгер
	log := logger.NewLogger(&configs.Config{})
	secretRepo := secret.NewSecretRepository(database, log)

	jwtService := jwt.NewJWT("test-secret")

	passwordReset := passwordReset.NewPasswordResetRepository(database, log)

	authService := auth.NewAuthService(userRepo, secretRepo, passwordReset, jwtService)
	return authService, mockDB, cleanup
}

func TestRegisterSuccess(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{}))
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mockDB.ExpectCommit()
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "secrets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mockDB.ExpectCommit()

	user, err := authService.Register("test@example.com", "Password123!", "testuser")
	if err != nil {
		t.Fatalf("Register error = %v", err)
	}
	if user.Email != "test@example.com" {
		t.Fatalf("email = %s, want %s", user.Email, "test@example.com")
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRegisterUserExists(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Настройка ожиданий для БД - пользователь существует
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, "test@example.com", "hashedpassword", "testuser"))

	_, err := authService.Register("test@example.com", "Password123!", "testuser")

	// Проверка, что получили ожидаемую ошибку
	if err == nil {
		t.Fatal("Expected error for existing user, got nil")
	}
	if err.Error() != auth.ErrUserExists {
		t.Fatalf("Expected error message %s, got %s", auth.ErrUserExists, err.Error())
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Хэшируем пароль для имитации сохраненного пароля
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, "test@example.com", string(hashedPassword), "testuser"))

	email, err := authService.Login("test@example.com", "Password123!")
	if err != nil {
		t.Fatalf("Login error = %v", err)
	}
	if email.Email != "test@example.com" {
		t.Fatalf("email = %s, want %s", email.Email, "test@example.com")
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginWrongPassword(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Хэшируем пароль для имитации сохраненного пароля
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, "test@example.com", string(hashedPassword), "testuser"))

	_, err := authService.Login("test@example.com", "WrongPassword123!")

	// Проверка, что получили ожидаемую ошибку
	if err == nil {
		t.Fatal("Expected error for wrong password, got nil")
	}
	if err.Error() != auth.ErrWrongCredentials {
		t.Fatalf("Expected error message %s, got %s", auth.ErrWrongCredentials, err.Error())
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginUserNotFound(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Настройка ожиданий для БД - пользователь не найден
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	_, err := authService.Login("nonexistent@example.com", "Password123!")

	// Проверка, что получили ожидаемую ошибку
	if err == nil {
		t.Fatal("Expected error for nonexistent user, got nil")
	}
	if err.Error() != auth.ErrWrongCredentials {
		t.Fatalf("Expected error message %s, got %s", auth.ErrWrongCredentials, err.Error())
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRefreshTokensSuccess(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Создаем валидный refresh token
	validEmail := "test@example.com"
	jwtService := jwt.NewJWT("test-secret")
	tokenPair, _ := jwtService.GenerateTokenPair(jwt.JWTData{Email: validEmail})

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, validEmail, "hashedpassword", "testuser"))

	// Вызываем метод обновления токенов
	newTokenPair, err := authService.RefreshTokens(tokenPair.AccessToken, tokenPair.RefreshToken)
	if err != nil {
		t.Fatalf("RefreshTokens error = %v", err)
	}
	if newTokenPair.AccessToken == "" {
		t.Fatal("AccessToken is empty")
	}
	if newTokenPair.RefreshToken == "" {
		t.Fatal("RefreshToken is empty")
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRefreshTokensInvalidToken(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Вызываем метод обновления токенов с невалидным токеном
	_, err := authService.RefreshTokens("invalid.access.token", "invalid.refresh.token")

	// Проверка, что получили ожидаемую ошибку
	if err == nil {
		t.Fatal("Expected error for invalid refresh token, got nil")
	}
	if err.Error() != auth.ErrInvalidRefreshToken {
		t.Fatalf("Expected error message %s, got %s", auth.ErrInvalidRefreshToken, err.Error())
	}

	// Не должно быть обращений к БД
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRefreshTokensUserNotFound(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	// Создаем валидный refresh token
	validEmail := "deleted@example.com"
	jwtService := jwt.NewJWT("test-secret")
	tokenPair, _ := jwtService.GenerateTokenPair(jwt.JWTData{Email: validEmail})

	// Настройка ожиданий для БД - пользователь не найден
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{}))

	// Вызываем метод обновления токенов
	_, err := authService.RefreshTokens(tokenPair.AccessToken, tokenPair.RefreshToken)

	// Проверка, что получили ожидаемую ошибку
	if err == nil {
		t.Fatal("Expected error for user not found, got nil")
	}
	if err.Error() != auth.ErrUserNotFound {
		t.Fatalf("Expected error message %s, got %s", auth.ErrUserNotFound, err.Error())
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestForgotPasswordSuccess(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	email := "test@example.com"

	// Настройка ожиданий для БД - поиск пользователя
	mockDB.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow(email))

	// Настройка ожиданий для для БД создание записи токена для сброса пароля
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "password_resets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mockDB.ExpectCommit()

	config := &configs.Config{}

	err := authService.ForgotPassword(config, email)
	if err != nil {
		t.Fatalf("ForgotPassword error = %v", err)
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestForgotPasswordUserNotExists(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	email := "test@example.com"

	// Настраиваем SELECT: вернёт пустой результат
	mockDB.ExpectQuery(`SELECT .* FROM "users" WHERE`).
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

	config := &configs.Config{}

	err := authService.ForgotPassword(config, email)
	if err == nil {
		t.Fatalf("expected error for not existing user, got nil")
	}
	if err.Error() != auth.ErrUserNotFound {
		t.Errorf("ожидали ошибку %q, получили %q", auth.ErrUserNotFound, err)
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestForgotPasswordErrorCreatingPassword(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	email := "test@example.com"

	// Настройка ожиданий для БД - поиск пользователя
	mockDB.ExpectQuery(`SELECT .* FROM "users"`).
		WillReturnRows(sqlmock.NewRows([]string{"email"}).
			AddRow(email))

	// Настройка ожиданий для для БД создание записи токена для сброса пароля
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "password_resets"`).
		WillReturnError(errors.New("db insert error"))
	mockDB.ExpectRollback()

	config := &configs.Config{}

	err := authService.ForgotPassword(config, email)
	if err == nil {
		t.Fatal("Expected error for password reset creation, got nil")
	}
	if err.Error() != auth.ErrCreatePasswordReset {
		t.Errorf("ожидали ошибку %q, получили %q", auth.ErrCreatePasswordReset, err)
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestResetPasswordSuccess(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	token := "fdf77f4394224b3436b30b449bcbd99c8f930f5e7bee01bee5eecea6dccbb1e8"
	userID := uint(1)
	newPassword := "TestPassword1!"
	secretID := uint(10)
	passwordResetID := uint(1)

	// === password_resets (SELECT active token) ===
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}).AddRow(
			passwordResetID,
			userID,
			token,
			false,
			time.Now().Add(10*time.Minute),
			time.Now(),
			time.Now(),
		))

	// === users (SELECT by id) ===
	mockDB.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password", "created_at", "updated_at",
		}).AddRow(
			userID,
			"test@example.com",
			"old_hashed_password",
			time.Now(),
			time.Now(),
		))

	// === users (UPDATE password) ===
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "users"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			userID,
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// === secrets (SELECT by user_id) - первый запрос ===
	mockDB.ExpectQuery(`SELECT \* FROM "secrets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "secret", "created_at", "updated_at", "current_password",
		}).AddRow(
			secretID,
			userID,
			"old_secret_value",
			time.Now(),
			time.Now(),
			"old_current_password_hash",
		))

	// === secrets (SELECT by id) - второй запрос ===
	mockDB.ExpectQuery(`SELECT \* FROM "secrets" WHERE "secrets"."id" = \$1 AND "secrets"."deleted_at" IS NULL ORDER BY "secrets"."id" LIMIT \$2`).
		WithArgs(secretID, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "secret", "created_at", "updated_at", "current_password",
		}).AddRow(
			secretID,
			userID,
			"old_secret_value",
			time.Now(),
			time.Now(),
			"old_current_password_hash",
		))

	// === previous_passwords (SELECT by secret_id) - проверка истории ===
	mockDB.ExpectQuery(`SELECT \* FROM "previous_passwords" WHERE "previous_passwords"."secret_id" = \$1`).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "secret_id", "password", "created_at", "updated_at",
		})) // Пустой результат - нет предыдущих паролей

	// === ТРАНЗАКЦИЯ ДЛЯ ОБНОВЛЕНИЯ SECRET ===
	mockDB.ExpectBegin() // Начало транзакции для Update секрета

	// INSERT в previous_passwords (выполняется ПЕРВЫМ в транзакции)
	mockDB.ExpectQuery(`INSERT INTO "previous_passwords" \("secret_id","password","created_at"\) VALUES \(\$1,\$2,\$3\) RETURNING "id"`).
		WithArgs(
			secretID,
			"old_current_password_hash",
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// === previous_passwords (COUNT) - проверка количества ПОСЛЕ INSERT ===
	mockDB.ExpectQuery(`SELECT count\(\*\) FROM "previous_passwords" WHERE secret_id = \$1`).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(4)) // Теперь 4 (3 было + 1 новый)

	// UPDATE secrets
	mockDB.ExpectExec(`UPDATE "secrets" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"user_id"=\$4,"current_password"=\$5 WHERE "secrets"."deleted_at" IS NULL AND "id" = \$6`).
		WithArgs(
			sqlmock.AnyArg(), // created_at
			sqlmock.AnyArg(), // updated_at
			nil,              // deleted_at
			userID,           // user_id
			sqlmock.AnyArg(), // new current_password hash
			secretID,         // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectCommit() // Конец транзакции для Update секрета

	// === password_resets (UPDATE used = true) ===
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "password_resets" SET "used"=\$1,"updated_at"=\$2 WHERE \(id = \$3 AND expires_at > \$4\) AND "password_resets"."deleted_at" IS NULL`).
		WithArgs(
			true,             // used
			sqlmock.AnyArg(), // updated_at
			passwordResetID,  // id
			sqlmock.AnyArg(), // expires_at (текущее время)
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// Вызов метода
	updatedUser, err := authService.ResetPassword(token, newPassword)
	require.NoError(t, err)
	require.NotNil(t, updatedUser)
	require.Equal(t, userID, updatedUser.ID)

	// Проверяем, что все ожидания были выполнены
	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestResetPasswordInvalidToken(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	token := "invalid_token"

	// password_resets (SELECT - токен не найден)
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		})) // Пустой результат

	// Вызов метода
	updatedUser, err := authService.ResetPassword(token, "NewPassword1!")
	require.Error(t, err)
	require.Nil(t, updatedUser)
	require.Contains(t, err.Error(), "Token Expired")

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestResetPasswordInvalidNewPassword(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	token := "valid_token"
	userID := uint(1)

	// password_resets (SELECT активный токен)
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}).AddRow(
			1,
			userID,
			token,
			false,
			time.Now().Add(10*time.Minute),
			time.Now(),
			time.Now(),
		))

	// users (SELECT by id)
	mockDB.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password", "created_at", "updated_at",
		}).AddRow(
			userID,
			"test@example.com",
			"old_hashed_password",
			time.Now(),
			time.Now(),
		))

	// Вызов метода с невалидным паролем
	updatedUser, err := authService.ResetPassword(token, "short")
	require.Error(t, err)
	require.Nil(t, updatedUser)
	require.Contains(t, err.Error(), "The password must contain 12 characters including numbers and special characters.") // или конкретное сообщение об ошибке

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestResetPasswordUserUpdateError(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	token := "valid_token"
	userID := uint(1)

	// password_resets (SELECT активный токен)
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}).AddRow(
			1,
			userID,
			token,
			false,
			time.Now().Add(10*time.Minute),
			time.Now(),
			time.Now(),
		))

	// users (SELECT by id)
	mockDB.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password", "created_at", "updated_at",
		}).AddRow(
			userID,
			"test@example.com",
			"old_hashed_password",
			time.Now(),
			time.Now(),
		))

	// users (UPDATE password) - ошибка
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "users"`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			userID,
		).
		WillReturnError(errors.New("database error"))
	mockDB.ExpectRollback() // Rollback при ошибке

	// Вызов метода
	updatedUser, err := authService.ResetPassword(token, "TestPassword1!")
	require.Error(t, err)
	require.Nil(t, updatedUser)
	require.Contains(t, err.Error(), "database error")

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestResetPasswordSecretUpdateError(t *testing.T) {
	authService, mockDB, cleanup := setupAuthService(t)
	defer cleanup()

	token := "valid_token"
	userID := uint(1)
	secretID := uint(10)
	passwordResetID := uint(1)

	// === password_resets (SELECT активный токен) ===
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}).AddRow(
			passwordResetID,
			userID,
			token,
			false,
			time.Now().Add(10*time.Minute),
			time.Now(),
			time.Now(),
		))

	// === users (SELECT by id) ===
	mockDB.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password", "created_at", "updated_at",
		}).AddRow(
			userID,
			"test@example.com",
			"old_hashed_password",
			time.Now(),
			time.Now(),
		))

	// === users (UPDATE password) ===
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "users"`).
		WithArgs(
			sqlmock.AnyArg(), // email
			sqlmock.AnyArg(), // password
			sqlmock.AnyArg(), // username
			sqlmock.AnyArg(), // updated_at
			userID,           // id
		).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// === secrets (SELECT by user_id) - первый запрос ===
	mockDB.ExpectQuery(`SELECT \* FROM "secrets"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "secret", "created_at", "updated_at", "current_password",
		}).AddRow(
			secretID,
			userID,
			"old_secret_value",
			time.Now(),
			time.Now(),
			"old_current_password_hash",
		))

	// === secrets (SELECT by id) - второй запрос ===
	mockDB.ExpectQuery(`SELECT \* FROM "secrets" WHERE "secrets"."id" = \$1 AND "secrets"."deleted_at" IS NULL ORDER BY "secrets"."id" LIMIT \$2`).
		WithArgs(secretID, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "secret", "created_at", "updated_at", "current_password",
		}).AddRow(
			secretID,
			userID,
			"old_secret_value",
			time.Now(),
			time.Now(),
			"old_current_password_hash",
		))

	// === previous_passwords (SELECT by secret_id) - проверка истории ===
	mockDB.ExpectQuery(`SELECT \* FROM "previous_passwords" WHERE "previous_passwords"."secret_id" = \$1`).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "secret_id", "password", "created_at", "updated_at",
		})) // Пустой результат

	// === ТРАНЗАКЦИЯ ДЛЯ ОБНОВЛЕНИЯ SECRET ===
	mockDB.ExpectBegin()

	// INSERT в previous_passwords
	mockDB.ExpectQuery(`INSERT INTO "previous_passwords" \("secret_id","password","created_at"\) VALUES \(\$1,\$2,\$3\) RETURNING "id"`).
		WithArgs(
			secretID,
			"old_current_password_hash",
			sqlmock.AnyArg(),
		).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// previous_passwords (COUNT)
	mockDB.ExpectQuery(`SELECT count\(\*\) FROM "previous_passwords" WHERE secret_id = \$1`).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// UPDATE secrets - ошибка
	mockDB.ExpectExec(`UPDATE "secrets" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"user_id"=\$4,"current_password"=\$5 WHERE "secrets"."deleted_at" IS NULL AND "id" = \$6`).
		WithArgs(
			sqlmock.AnyArg(),
			sqlmock.AnyArg(),
			nil,
			userID,
			sqlmock.AnyArg(),
			secretID,
		).
		WillReturnError(errors.New("update error"))

	mockDB.ExpectRollback() // Rollback при ошибке

	// Вызов метода
	updatedUser, err := authService.ResetPassword(token, "TestPassword1!")
	require.Error(t, err)
	require.Nil(t, updatedUser)
	require.Contains(t, err.Error(), "update error")

	require.NoError(t, mockDB.ExpectationsWereMet())
}
