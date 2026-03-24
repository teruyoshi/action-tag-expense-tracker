package repositories

import (
	"testing"
	"time"

	"action-tag-expense-tracker/backend/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestIncomeRepository_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "returns incomes for month",
			setup: func(mock sqlmock.Sqlmock) {
				categoryRows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "給与")
				incomeRows := sqlmock.NewRows([]string{"id", "income_category_id", "date", "description", "amount"}).
					AddRow(1, 1, time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), "3月分給与", 300000).
					AddRow(2, 1, time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), "副業", 50000)
				mock.ExpectQuery("SELECT \\* FROM `incomes` WHERE YEAR\\(date\\) = \\? AND MONTH\\(date\\) = \\?").
					WithArgs(2026, 3).
					WillReturnRows(incomeRows)
				mock.ExpectQuery("SELECT \\* FROM `income_categories` WHERE `income_categories`.`id` = \\?").
					WithArgs(1).
					WillReturnRows(categoryRows)
			},
			want: 2,
		},
		{
			name: "returns empty when no incomes",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "income_category_id", "date", "description", "amount"})
				mock.ExpectQuery("SELECT \\* FROM `incomes` WHERE YEAR\\(date\\) = \\? AND MONTH\\(date\\) = \\?").
					WithArgs(2026, 1).
					WillReturnRows(rows)
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &IncomeRepository{DB: db}

			var year, month int
			if tt.want > 0 {
				year, month = 2026, 3
			} else {
				year, month = 2026, 1
			}
			incomes, err := repo.FindAll(year, month)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(incomes) != tt.want {
				t.Fatalf("FindAll() got %d incomes, want %d", len(incomes), tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIncomeRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "found",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				incomeRows := sqlmock.NewRows([]string{"id", "income_category_id", "date", "description", "amount"}).
					AddRow(1, 1, time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), "3月分給与", 300000)
				mock.ExpectQuery("SELECT \\* FROM `incomes` WHERE `incomes`.`id` = \\?").
					WithArgs(1, 1).
					WillReturnRows(incomeRows)
				categoryRows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "給与")
				mock.ExpectQuery("SELECT \\* FROM `income_categories` WHERE `income_categories`.`id` = \\?").
					WithArgs(1).
					WillReturnRows(categoryRows)
			},
			want: 300000,
		},
		{
			name: "not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "income_category_id", "date", "description", "amount"})
				mock.ExpectQuery("SELECT \\* FROM `incomes` WHERE `incomes`.`id` = \\?").
					WithArgs(999, 1).
					WillReturnRows(rows)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &IncomeRepository{DB: db}

			income, err := repo.FindByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err == gorm.ErrRecordNotFound && tt.wantErr {
				// expected
			} else if !tt.wantErr && income.Amount != tt.want {
				t.Fatalf("FindByID() amount = %d, want %d", income.Amount, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIncomeRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `incomes`").
		WithArgs(1, sqlmock.AnyArg(), "3月分給与", 300000).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &IncomeRepository{DB: db}
	income := &models.Income{
		IncomeCategoryID: 1,
		Date:             time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		Description:      "3月分給与",
		Amount:           300000,
	}

	if err := repo.Create(income); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIncomeRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `incomes`").
		WithArgs(1, sqlmock.AnyArg(), "更新後", 350000, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &IncomeRepository{DB: db}
	income := &models.Income{
		ID:               1,
		IncomeCategoryID: 1,
		Date:             time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
		Description:      "更新後",
		Amount:           350000,
	}

	if err := repo.Update(income); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIncomeRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `incomes` WHERE `incomes`.`id` = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := &IncomeRepository{DB: db}

	if err := repo.Delete(1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
