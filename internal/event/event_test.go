package event

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db"
	"github.com/PurpleSchoolPractice/metiing-pro-golang/pkg/db/mock"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"time"
)

func TestCreateEvent(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectQuery(`INSERT INTO "events"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"testevent", "description", "2025-05-09", 1, 2).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	if err != nil {
		t.Fatalf("Not possible to parse date: %v", err)
	}
	newEvent := NewEvent("testevent", "description", 2, 1, date)
	createdEvent, err := repo.Create(newEvent)
	require.NoError(t, err, "Create event failed")
	require.NotEmpty(t, createdEvent, "Create event should not be empty")
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindEventByID(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	if err != nil {
		t.Fatalf("Not possible to parse date: %v", err)
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"title",
		"description",
		"ownerID",
		"event_date",
		"creatorID",
	}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", 2, date, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE id = $1`)).
		WithArgs(uint(1)).WillReturnRows(rows)

	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	foundEvent, err := repo.FindById(uint(1))
	require.NoError(t, err, "Error finding event")
	require.NotNil(t, foundEvent, "Expected to find event")
	require.Equal(t, uint(1), foundEvent.ID)
	require.Equal(t, "testevent", foundEvent.Title)
	require.Equal(t, "description", foundEvent.Description)
	require.Equal(t, "2025-05-09", foundEvent.EventDate)
	require.Equal(t, 1, foundEvent.CreatorID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestFindByCreatorID(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	fixedTime := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	date, err := time.Parse("2006-01-02", "2025-05-09")
	if err != nil {
		t.Fatalf("Not possible to parse date: %v", err)
	}
	rows := sqlmock.NewRows([]string{
		"id",
		"created_at",
		"updated_at",
		"deleted_at",
		"title",
		"description",
		"ownerID",
		"event_date",
		"creatorID",
	}).AddRow(1, fixedTime, fixedTime, nil, "testevent", "description", 2, date, 1)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "events" WHERE creatorID = $1`)).
		WithArgs(uint(1)).WillReturnRows(rows)
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	findedEvents, err := repo.FindAllByCreatorId(uint(1))
	require.NoError(t, err, "Error finding events")
	require.Len(t, findedEvents, 1)
	require.Equal(t, uint(1), findedEvents[0].ID)
	require.Equal(t, "testevent", findedEvents[0].Title)
	require.Equal(t, "description", findedEvents[0].Description)
	require.Equal(t, "2025-05-09", findedEvents[0].EventDate)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEvent(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectExec(`UPDATE "events" SET`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(),
			"newtestevent", "newdescription", "2025-05-10", 1, 2).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	date, err := time.Parse("2006-01-02", "2025-05-09")
	if err != nil {
		t.Fatalf("Not possible to parse date: %v", err)
	}
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	eventUpdate := &Event{
		Title:       "newtestevent",
		Description: "newdescription",
		CreatorID:   1,
		EventDate:   date,
		OwnerID:     2,
	}
	updatedEvent, err := repo.Update(eventUpdate)
	require.NoError(t, err, "Error updating event")
	require.Equal(t, uint(1), updatedEvent.ID)
	require.Equal(t, "newtestevent", updatedEvent.Title)
	require.Equal(t, "newdescription", updatedEvent.Description)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteEventByID(t *testing.T) {
	gormDB, mock, cleanup := mock.SetupMockDB(t)
	defer t.Cleanup(cleanup)
	mock.ExpectBegin()
	mock.ExpectExec(`DELETE FROM "events" WHERE id = $1`).
		WithArgs(uint(1)).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()
	dbWrapper := &db.Db{DB: gormDB}
	repo := NewEventRepository(dbWrapper)
	err := repo.DeleteById(uint(1))
	require.NoError(t, err, "Error deleting event")
	require.NoError(t, mock.ExpectationsWereMet())
}
