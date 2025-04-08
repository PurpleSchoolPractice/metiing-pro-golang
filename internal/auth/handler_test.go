package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/configs"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/auth"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/logger"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/secret"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
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

	authService := auth.NewAuthService(userRepo, secretRepo, jwtService)
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
