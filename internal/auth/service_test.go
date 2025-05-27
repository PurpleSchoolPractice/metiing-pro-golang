package auth_test

import (
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

func setupAuthService(t *testing.T) (*auth.AuthService, sqlmock.Sqlmock, func()) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	database := &db.Db{DB: gormDB}
	userRepo := user.NewUserRepository(database)

	// Создаем обычный логгер
	log := logger.NewLogger(&configs.Config{})
	secretRepo := secret.NewSecretRepository(database, log)

	jwtService := jwt.NewJWT("test-secret")

	authService := auth.NewAuthService(userRepo, secretRepo, jwtService)
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
