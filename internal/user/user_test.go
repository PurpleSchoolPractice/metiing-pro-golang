package user_test

import (
	"github.com/PurpleSchoolPractice/metiing-pro-golang/internal/user"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock, func()) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("Not created sqlmock: %v", err)
	}
	dialector := postgres.New(postgres.Config{
		Conn:                 sqlDB,
		PreferSimpleProtocol: true,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("Not opened DB: %v", err)
	}
	cleanup := func() {
		sqlDB.Close()
	}
	return gormDB, mock, cleanup
}

func TestCreateUser(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "testuser", "password", "email@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)

	newUser := user.NewUser("email@example.com", "password", "testuser")
	createdUser, err := repo.Create(newUser)
	if err != nil {
		t.Errorf("Error creating user: %v", err)
	}
	if createdUser.Email != "email@example.com" {
		t.Errorf("Exepected email email@example.com, got %s", createdUser.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying expectations: %v", err)
	}
}

func TestFindByEmail(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
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
	if err != nil {
		t.Errorf("Error fetching user: %v", err)
	}
	if foundUser.Email != "email@example.com" {
		t.Errorf("Expected email email@example.com, got %s", foundUser.Email)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying expectations: %v", err)
	}
}

func TestFindAllUsers(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
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
	if err != nil {
		t.Errorf("Error fetching users: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying expectations: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(`UPDATE "users" SET`)).
		WithArgs(sqlmock.AnyArg(), "newusername", "newpassword", "email@example.com", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()

	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)

	userToUpdate := &user.User{
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
	gormDB, mock, cleanup := setupMockDB(t)
	defer cleanup()
	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "deleted_at"=`)).
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	dbWrapper := &db.Db{DB: gormDB}
	repo := user.NewUserRepository(dbWrapper)
	userToDelete := &user.User{
		Model:    gorm.Model{ID: 1},
		Username: "testuser",
		Password: "password",
		Email:    "email@example.com",
	}
	err := repo.Delete(userToDelete)
	if err != nil {
		t.Errorf("Error deleting user: %v", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Error verifying expectations: %v", err)
	}
}
