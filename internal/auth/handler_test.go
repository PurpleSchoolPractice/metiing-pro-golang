package auth_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
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

func setupAuthHandler(t *testing.T) (*auth.AuthHandler, sqlmock.Sqlmock, func()) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	database := &db.Db{DB: gormDB}
	userRepo := user.NewUserRepository(database)

	// Создаем простой логгер для тестов
	log := logger.NewLogger(&configs.Config{})
	secretRepo := secret.NewSecretRepository(database, log)

	jwtService := jwt.NewJWT("test-secret")

	passwordReset := passwordReset.NewPasswordResetRepository(database, log)

	authService := auth.NewAuthService(userRepo, secretRepo, passwordReset, jwtService)
	handler := &auth.AuthHandler{
		Config: &configs.Config{
			Auth: configs.AuthConfig{
				Secret: "test-secret",
			},
		},
		AuthService: authService,
	}
	return handler, mockDB, cleanup
}

func TestRegisterHandlerSuccess(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
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

	// Подготовка запроса
	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "Password123!",
		Name:     "testuser",
	})
	reader := bytes.NewReader(data)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/register", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.Register()(w, req)

	// Проверка результата
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
		return
	}

	// Проверка структуры ответа
	var response auth.RegisterResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("AccessToken is empty")
	}

	if response.RefreshToken == "" {
		t.Error("RefreshToken is empty")
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRegisterHandlerUserExists(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Настройка ожиданий для БД - пользователь существует
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, "test@example.com", "hashedpassword", "testuser"))

	// Подготовка запроса
	data, _ := json.Marshal(&auth.RegisterRequest{
		Email:    "test@example.com",
		Password: "Password123!",
		Name:     "testuser",
	})
	reader := bytes.NewReader(data)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/register", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.Register()(w, req)

	// Проверка результата - ожидаем ошибку
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		return
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginHandlerSuccess(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Хэшируем пароль для имитации сохраненного пароля
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, "test@example.com", string(hashedPassword), "testuser"))

	// Подготовка запроса
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "Password123!",
	})
	reader := bytes.NewReader(data)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/login", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.Login()(w, req)

	// Проверка результата
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
		return
	}

	// Проверка структуры ответа
	var response auth.LoginResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("AccessToken is empty")
	}

	if response.RefreshToken == "" {
		t.Error("RefreshToken is empty")
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestLoginHandlerWrongCredentials(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Хэшируем пароль для имитации сохраненного пароля
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Password123!"), bcrypt.DefaultCost)

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, "test@example.com", string(hashedPassword), "testuser"))

	// Подготовка запроса
	data, _ := json.Marshal(&auth.LoginRequest{
		Email:    "test@example.com",
		Password: "WrongPassword123!",
	})
	reader := bytes.NewReader(data)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/login", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.Login()(w, req)

	// Проверка результата - ожидаем ошибку
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		return
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRefreshTokenHandlerSuccess(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Создаем валидный refresh token
	validEmail := "test@example.com"
	jwtService := jwt.NewJWT("test-secret")
	tokenPair, _ := jwtService.GenerateTokenPair(jwt.JWTData{Email: validEmail})

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT (.+) FROM "users" WHERE`).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "password", "username"}).
			AddRow(1, validEmail, "hashedpassword", "testuser"))

	// Подготовка запроса
	data, _ := json.Marshal(&auth.RefreshTokenRequest{
		RefreshToken: tokenPair.RefreshToken,
	})
	reader := bytes.NewReader(data)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/refresh", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.RefreshToken()(w, req)

	// Проверка результата
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
		return
	}

	// Проверка структуры ответа
	var response auth.RefreshTokenResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.AccessToken == "" {
		t.Error("AccessToken is empty")
	}

	if response.RefreshToken == "" {
		t.Error("RefreshToken is empty")
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestRefreshTokenHandlerInvalidToken(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Подготовка запроса с невалидным токеном
	data, _ := json.Marshal(&auth.RefreshTokenRequest{
		RefreshToken: "invalid.refresh.token",
	})
	reader := bytes.NewReader(data)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/refresh", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.RefreshToken()(w, req)

	// Проверка результата - ожидаем ошибку авторизации
	if w.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, w.Code)
		return
	}

	// Проверяем, что все ожидания были выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestInvalidJsonBody(t *testing.T) {
	handler, _, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Подготовка запроса с невалидным JSON
	invalidJson := []byte(`{"email": "test@example.com", "password": `)
	reader := bytes.NewReader(invalidJson)

	// Выполнение запроса
	w := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodPost, "/auth/login", reader)
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		t.Fatal(err)
		return
	}

	handler.Login()(w, req)

	// Ожидаем ошибку обработки JSON
	if w.Code == http.StatusOK {
		t.Errorf("Expected non-OK status code, got %d", w.Code)
		return
	}
}

func TestForgotPasswordHandlerSuccess(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	email := "test@example.com"
	userID := 1

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT .* FROM "users" WHERE .*email.*`).
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(userID, email))

	// INSERT password_resets
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "password_resets"`).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mockDB.ExpectCommit()

	// Подготовка запроса
	data, _ := json.Marshal(&auth.ForgotPasswordRequest{Email: email})
	req, err := http.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader(data))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	handler.ForgotPassword()(w, req)

	// Проверка результата
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
	}

	if !strings.Contains(w.Body.String(), "reset your password") {
		t.Errorf("Unexpected response: %s", w.Body.String())
	}

	// Проверяем, что все ожидания выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestForgotPasswordHandlerUserNotFound(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	email := "test@example.com"

	// Пользователь не найден
	mockDB.ExpectQuery(`SELECT .* FROM "users" WHERE .*email.*`).
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}))

	// Подготовка запроса
	data, _ := json.Marshal(&auth.ForgotPasswordRequest{Email: email})
	req, _ := http.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ForgotPassword()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "User not found") {
		t.Errorf("Unexpected response: %s", w.Body.String())
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestForgotPasswordHandlerDBInsertError(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	email := "test@example.com"
	userID := 1

	// Настройка ожиданий для БД
	mockDB.ExpectQuery(`SELECT .* FROM "users" WHERE .*email.*`).
		WithArgs(email, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email"}).
			AddRow(userID, email))

	// Ошибка при вставке в password_resets
	mockDB.ExpectBegin()
	mockDB.ExpectQuery(`INSERT INTO "password_resets"`).
		WillReturnError(errors.New("db insert error"))
	mockDB.ExpectRollback()

	data, _ := json.Marshal(&auth.ForgotPasswordRequest{Email: email})
	req, _ := http.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ForgotPassword()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if !strings.Contains(w.Body.String(), "Failed to create password reset") {
		t.Errorf("Unexpected response: %s", w.Body.String())
	}

	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestForgotPasswordHandlerInvalidJSON(t *testing.T) {
	handler, _, cleanup := setupAuthHandler(t)
	defer cleanup()

	req, _ := http.NewRequest(http.MethodPost, "/auth/forgot-password", bytes.NewReader([]byte(`{invalid json}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	handler.ForgotPassword()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestCheckTokenExpirationDateSuccess(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	token := "validToken"
	userID := uint(1)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Правильный мок для запроса - ожидаем 4 аргумента
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets" WHERE \(token = \$1 AND used = \$2 AND expires_at > \$3\) AND "password_resets"."deleted_at" IS NULL LIMIT \$4`).
		WithArgs(
			token,
			false,
			sqlmock.AnyArg(),
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}).AddRow(
			1,
			userID,
			token,
			false,
			expiresAt,
			time.Now(),
			time.Now(),
		))

	// Создаем HTTP запрос
	reqBody := auth.CheckTokenRequest{Token: token}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/check-token", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Записываем ответ
	rr := httptest.NewRecorder()

	// Вызываем handler
	handler.CheckTokenExpirationDate()(rr, req)

	// Проверяем статус код
	require.Equal(t, http.StatusOK, rr.Code)

	// Проверяем тело ответа
	var response auth.CheckTokenResponse
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)

	require.Equal(t, userID, response.UserID)
	require.Equal(t, token, response.Token)
	require.WithinDuration(t, expiresAt, response.ExpiresAt, time.Second)

	// Проверяем, что все ожидания выполнены
	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCheckTokenExpirationDateMissingToken(t *testing.T) {
	handler, _, cleanup := setupAuthHandler(t)
	defer cleanup()

	// Тело без токена
	reqBody := map[string]interface{}{}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/check-token", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CheckTokenExpirationDate()(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestCheckTokenExpirationDateTokenNotFound(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	token := "nonExistentToken"

	// Мок возвращает пустой результат
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets" WHERE \(token = \$1 AND used = \$2 AND expires_at > \$3\) AND "password_resets"."deleted_at" IS NULL LIMIT \$4`).
		WithArgs(
			token,
			false,
			sqlmock.AnyArg(),
			1,
		).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}))

	reqBody := auth.CheckTokenRequest{Token: token}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/check-token", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CheckTokenExpirationDate()(rr, req)

	require.Equal(t, http.StatusNotFound, rr.Code)
	require.Contains(t, rr.Body.String(), "Token not found")

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestCheckTokenExpirationDateDatabaseError(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	token := "testToken"

	// Мок возвращает ошибку базы данных
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets" WHERE \(token = \$1 AND used = \$2 AND expires_at > \$3\) AND "password_resets"."deleted_at" IS NULL LIMIT \$4`).
		WithArgs(
			token,
			false,
			sqlmock.AnyArg(),
			1,
		).
		WillReturnError(errors.New("database connection failed"))

	reqBody := auth.CheckTokenRequest{Token: token}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/check-token", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler.CheckTokenExpirationDate()(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	require.Contains(t, rr.Body.String(), "database connection failed")

	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestResetPasswordHandlerSuccess(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	token := "valid_token"
	newPassword := "NewPassword123!"
	userID := uint(1)
	secretID := uint(10)

	// 1. SELECT из password_resets
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}).AddRow(
			1, userID, token, false, time.Now().Add(10*time.Minute), time.Now(), time.Now(),
		))

	// 2. SELECT из users
	mockDB.ExpectQuery(`SELECT \* FROM "users"`).
		WithArgs(userID, sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "email", "password", "created_at", "updated_at",
		}).AddRow(
			userID, "test@example.com", "old_hash", time.Now(), time.Now(),
		))

	// 3. UPDATE users
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "users"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// 4. SELECT из secrets по user_id
	mockDB.ExpectQuery(`SELECT \* FROM "secrets" WHERE user_id = \$1 AND "secrets"."deleted_at" IS NULL ORDER BY "secrets"."id" LIMIT \$2`).
		WithArgs(userID, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "secret", "created_at", "updated_at", "current_password",
		}).AddRow(
			secretID, userID, "old_secret", time.Now(), time.Now(), "old_current_password_hash",
		))

	// 5. SELECT из secrets по id
	mockDB.ExpectQuery(`SELECT \* FROM "secrets" WHERE "secrets"."id" = \$1 AND "secrets"."deleted_at" IS NULL ORDER BY "secrets"."id" LIMIT \$2`).
		WithArgs(secretID, 1).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "secret", "created_at", "updated_at", "current_password",
		}).AddRow(
			secretID, userID, "old_secret", time.Now(), time.Now(), "old_current_password_hash",
		))

	// 6. SELECT из previous_passwords
	mockDB.ExpectQuery(`SELECT \* FROM "previous_passwords" WHERE "previous_passwords"."secret_id" = \$1`).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "secret_id", "password", "created_at", "updated_at",
		})) // Пустой результат

	// 7. Транзакция для обновления secret
	mockDB.ExpectBegin()

	// INSERT в previous_passwords
	mockDB.ExpectQuery(`INSERT INTO "previous_passwords" \("secret_id","password","created_at"\) VALUES \(\$1,\$2,\$3\) RETURNING "id"`).
		WithArgs(secretID, "old_current_password_hash", sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	// COUNT previous_passwords
	mockDB.ExpectQuery(`SELECT count\(\*\) FROM "previous_passwords" WHERE secret_id = \$1`).
		WithArgs(secretID).
		WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(3))

	// UPDATE secrets
	mockDB.ExpectExec(`UPDATE "secrets" SET "created_at"=\$1,"updated_at"=\$2,"deleted_at"=\$3,"user_id"=\$4,"current_password"=\$5 WHERE "secrets"."deleted_at" IS NULL AND "id" = \$6`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), nil, userID, sqlmock.AnyArg(), secretID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	mockDB.ExpectCommit()

	// 8. UPDATE password_resets
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "password_resets" SET "used"=\$1,"updated_at"=\$2 WHERE \(id = \$3 AND expires_at > \$4\) AND "password_resets"."deleted_at" IS NULL`).
		WithArgs(true, sqlmock.AnyArg(), 1, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()

	// Подготовка запроса
	reqBody := auth.ResetPasswordRequest{
		Token:       token,
		NewPassword: newPassword,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	handler.ResetPassword()(w, req)

	// Проверка результата
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		t.Errorf("Response body: %s", w.Body.String())
	}

	// Проверяем JSON ответ
	var response auth.ResetPasswordResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	if response.UserID != userID {
		t.Errorf("Expected user ID %d, got %d", userID, response.UserID)
	}

	// Проверяем, что все ожидания выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}

func TestResetPasswordHandlerInvalidToken(t *testing.T) {
	handler, mockDB, cleanup := setupAuthHandler(t)
	defer cleanup()

	token := "invalid_token"
	newPassword := "NewPassword123!"

	// SELECT из password_resets возвращает пустой результат
	mockDB.ExpectQuery(`SELECT \* FROM "password_resets"`).
		WithArgs(token, false, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{
			"id", "user_id", "token", "used", "expires_at", "created_at", "updated_at",
		}))

	// Подготовка запроса
	reqBody := auth.ResetPasswordRequest{
		Token:       token,
		NewPassword: newPassword,
	}
	reqBodyBytes, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/reset-password", bytes.NewReader(reqBodyBytes))
	req.Header.Set("Content-Type", "application/json")

	// Выполнение запроса
	w := httptest.NewRecorder()
	handler.ResetPassword()(w, req)

	// Проверка результата - должна быть ошибка
	if w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status %d, got %d", http.StatusInternalServerError, w.Code)
	}

	// Проверяем, что все ожидания выполнены
	if err := mockDB.ExpectationsWereMet(); err != nil {
		t.Errorf("Unfulfilled expectations: %s", err)
	}
}
