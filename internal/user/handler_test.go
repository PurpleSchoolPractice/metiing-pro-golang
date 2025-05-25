package user_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

func setupDBUsersHandler(t *testing.T) (*user.UserHandler, sqlmock.Sqlmock, func()) {
	gormDB, mockDB, cleanup := mock.SetupMockDB(t)
	database := &db.Db{DB: gormDB}
	userRepo := user.NewUserRepository(database)

	jwtService := jwt.NewJWT("test-secret")

	handler := &user.UserHandler{
		UserRepository: userRepo,
		JWTService:     jwtService,
	}

	return handler, mockDB, cleanup
}

func TestGetAllUser(t *testing.T) {
	handler, mockDB, clean := setupDBUsersHandler(t)
	defer clean()
	//Создаем данные для теста в моковой БД
	rows := sqlmock.NewRows([]string{"id", "email", "password", "username"}).
		AddRow(1, "test@test.ru", "TestTest1254!", "TestName").
		AddRow(2, "test2@test2.ru", "Test2Test!1254!", "TestName2")
	mockDB.ExpectQuery("SELECT").WillReturnRows(rows)
	r := chi.NewRouter()
	r.Get("/users", handler.GetAllUsers())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	r.ServeHTTP(w, req)
	require.Equal(t, w.Code, 200)
	require.NoError(t, mockDB.ExpectationsWereMet())
}
func TestUserByID(t *testing.T) {
	handler, mockDB, clean := setupDBUsersHandler(t)
	defer clean()
	// Хэшируем пароль для имитации сохраненного пароля
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("Test1Test1!2021"), bcrypt.DefaultCost)
	//Создаем данные для теста в моковой БД
	rows := sqlmock.NewRows([]string{"id", "email", "password", "username"}).
		AddRow(1, "test2@test2.ru", string(hashedPassword), "Test2")
	mockDB.ExpectQuery("SELECT").WillReturnRows(rows)
	r := chi.NewRouter()
	r.Get("/users/{id}", handler.GetUserByID())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/users/1", nil)

	r.ServeHTTP(w, req)
	t.Logf("Response code: %d, body: %s", w.Code, w.Body.String())
	require.Equal(t, w.Code, 200)
	require.NoError(t, mockDB.ExpectationsWereMet())

}
func TestUserUpdate(t *testing.T) {
	handler, mockDB, clean := setupDBUsersHandler(t)
	defer clean()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Test1Test1!2021"), bcrypt.DefaultCost)
	require.NoError(t, err)
	// Обновляем пользователя
	mockDB.ExpectBegin()
	mockDB.ExpectExec(`UPDATE "users" SET "email"=$1,"password"=$2,"username"=$3 WHERE "id" = $4`).
		WithArgs("test2@test2.ru", string(hashedPassword), "Test", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mockDB.ExpectCommit()
	data, err := json.Marshal(&user.UserUpdateRequest{
		Username: "Test",
		Password: "Test1Test1!2021",
		Email:    "test2@test2.ru",
	})
	require.NoError(t, err)
	r := chi.NewRouter()
	r.Put("/user/{id}", handler.UpdateDataUser())
	reader := bytes.NewReader(data)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPut, "/user/1", reader)
	r.ServeHTTP(w, req)
	t.Logf("Response code: %d, body: %s", w.Code, w.Body.String())
	require.Equal(t, w.Code, 200)
	require.NoError(t, mockDB.ExpectationsWereMet())
}

func TestUserDelete(t *testing.T) {
	handler, mock, clean := setupDBUsersHandler(t)
	defer clean()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("Test1Test1!2021"), bcrypt.DefaultCost)
	require.NoError(t, err)

	sqlmock.NewRows([]string{"id", "email", "password", "username"}).
		AddRow(1, "test2@test2.ru", string(hashedPassword), "Test2")
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r := chi.NewRouter()
	r.Delete("/user/{id}", handler.DeleteUser())
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodDelete, "/user/1", nil)
	r.ServeHTTP(w, req)
	t.Logf("Response code: %d, body: %s", w.Code, w.Body.String())
	require.Equal(t, w.Code, 200)
	require.NoError(t, mock.ExpectationsWereMet())
}
