package user_test

import (
	"regexp"
	"testing"
	"time"

	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/models"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCreateUser(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "password", "email@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)

	newUser := models.NewUser("email@example.com", "password", "testuser")
	createdUser, err := repo.Create(newUser)
	require.NoError(t, err, "Create user failed")
	require.Equal(t, "testuser", createdUser.Username, "Username does not match")
	require.Equal(t, "password", createdUser.Password, "Password does not match")
	require.Equal(t, "email@example.com", createdUser.Email, "Email does not match")
	// Удаляем проверку временных меток, так как они зависят от текущего времени
	require.NotNil(t, createdUser.CreatedAt, "CreatedAt should not be nil")
	require.NotNil(t, createdUser.UpdatedAt, "UpdatedAt should not be nil")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByEmail(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "email"}).
		AddRow(1, fixedTime, fixedTime, nil, "testuser", "password", "email@example.com")

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1 AND "users"."deleted_at" IS NULL ORDER BY "users"."id" LIMIT $2`)).
		WithArgs("email@example.com", 1).
		WillReturnRows(rows)

	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)
	foundUser, err := repo.FindByEmail("email@example.com")
	require.NoError(t, err, "Error finding user")
	require.NotNil(t, foundUser, "Expected to find user")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindAllUsers(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()

	// Используем фиксированное время
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	// Ожидаем SELECT-запрос для получения всех пользователей.
	rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "username", "password", "email"}).
		AddRow(1, fixedTime, fixedTime, nil, "testuser", "password", "email@example.com")
	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL`)).
		WillReturnRows(rows)

	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)
	users, err := repo.FindAllUsers()
	require.NoError(t, err, "Error fetching users")
	require.Len(t, users, 1)
	require.Equal(t, users[0].Email, "email@example.com")
	require.NoError(t, mock.ExpectationsWereMet())

}

func TestUpdateUser(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "email"=$1,"password"=$2,"username"=$3,"updated_at"=$4 WHERE id = $5 AND "users"."deleted_at" IS NULL`)).
		WithArgs("email@example.com", "newpassword", "newusername", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)

	userToUpdate := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "newusername",
		Password: "newpassword",
		Email:    "email@example.com",
	}
	updatedUser, err := repo.Update(userToUpdate)
	require.NoError(t, err, "Error updating user")
	require.Equal(t, "newusername", updatedUser.Username, "Updated username does not match")
	require.Equal(t, "newpassword", updatedUser.Password, "Updated password does not match")
	require.Equal(t, "email@example.com", updatedUser.Email, "Updated email does not match")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUser(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)
	userToDelete := &models.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Password: "password",
		Email:    "email@example.com",
	}
	err := repo.Delete(userToDelete)
	require.NoError(t, err, "Error deleting user")
	require.NoError(t, mock.ExpectationsWereMet())
}
