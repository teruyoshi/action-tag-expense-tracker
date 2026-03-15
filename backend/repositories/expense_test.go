package repositories

import (
	"testing"

	"action-tag-expense-tracker/backend/models"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExpenseRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)

	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `expenses`").
		WithArgs(1, "ランチ", 800).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &ExpenseRepository{DB: db}
	expense := &models.Expense{EventID: 1, Item: "ランチ", Amount: 800}

	if err := repo.Create(expense); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestExpenseRepository_Update(t *testing.T) {
	tests := []struct {
		name   string
		expense *models.Expense
		setup  func(sqlmock.Sqlmock)
	}{
		{
			name:    "updates only item and amount",
			expense: &models.Expense{ID: 1, Item: "夕食", Amount: 1200},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `expenses` SET `item`=\\?,`amount`=\\? WHERE `id` = \\?").
					WithArgs("夕食", 1200, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
		{
			name:    "updates with empty item",
			expense: &models.Expense{ID: 2, Item: "", Amount: 500},
			setup: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec("UPDATE `expenses` SET `item`=\\?,`amount`=\\? WHERE `id` = \\?").
					WithArgs("", 500, 2).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &ExpenseRepository{DB: db}

			if err := repo.Update(tt.expense); err != nil {
				t.Fatalf("Update() error = %v", err)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
