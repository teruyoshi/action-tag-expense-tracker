package repositories

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBalanceRepository_Get(t *testing.T) {
	tests := []struct {
		name  string
		setup func(sqlmock.Sqlmock)
		want  int
	}{
		{
			name: "existing balance",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "amount", "updated_at"}).
					AddRow(1, 50000, time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC))
				mock.ExpectQuery("SELECT \\* FROM `balances`").WillReturnRows(rows)
			},
			want: 50000,
		},
		{
			name: "no balance creates default",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "amount", "updated_at"})
				mock.ExpectQuery("SELECT \\* FROM `balances`").WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `balances`").
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &BalanceRepository{DB: db}

			balance, err := repo.Get()
			if err != nil {
				t.Fatalf("Get() error = %v", err)
			}
			if balance.Amount != tt.want {
				t.Fatalf("Get() amount = %d, want %d", balance.Amount, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestBalanceRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "amount", "updated_at"}).
		AddRow(1, 50000, time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery("SELECT \\* FROM `balances`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `balances`").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := &BalanceRepository{DB: db}

	balance, err := repo.Update(40000)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if balance.Amount != 40000 {
		t.Fatalf("Update() amount = %d, want 40000", balance.Amount)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestBalanceRepository_Subtract(t *testing.T) {
	db, mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "amount", "updated_at"}).
		AddRow(1, 50000, time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC))
	mock.ExpectQuery("SELECT \\* FROM `balances`").WillReturnRows(rows)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `balances`").
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := &BalanceRepository{DB: db}

	if err := repo.Subtract(1000); err != nil {
		t.Fatalf("Subtract() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
