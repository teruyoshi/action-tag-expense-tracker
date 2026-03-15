package repositories

import (
	"testing"
	"time"

	"action-tag-expense-tracker/backend/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestEventRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `events`").
		WithArgs(sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &EventRepository{DB: db}
	event := &models.Event{
		Date:        time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC),
		ActionTagID: 1,
	}

	if err := repo.Create(event); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
