package repositories

import (
	"testing"

	"action-tag-expense-tracker/backend/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestIncomeCategoryRepository_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "returns all categories",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "給与").
					AddRow(2, "副業")
				mock.ExpectQuery("SELECT \\* FROM `income_categories`").WillReturnRows(rows)
			},
			want: 2,
		},
		{
			name: "returns empty when no categories",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT \\* FROM `income_categories`").WillReturnRows(rows)
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &IncomeCategoryRepository{DB: db}

			categories, err := repo.FindAll()
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(categories) != tt.want {
				t.Fatalf("FindAll() got %d categories, want %d", len(categories), tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestIncomeCategoryRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `income_categories`").
		WithArgs("給与").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &IncomeCategoryRepository{DB: db}
	category := &models.IncomeCategory{Name: "給与"}

	if err := repo.Create(category); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIncomeCategoryRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `income_categories`").
		WithArgs("副業", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &IncomeCategoryRepository{DB: db}
	category := &models.IncomeCategory{ID: 1, Name: "副業"}

	if err := repo.Update(category); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIncomeCategoryRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `income_categories` WHERE `income_categories`.`id` = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := &IncomeCategoryRepository{DB: db}

	if err := repo.Delete(1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestIncomeCategoryRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		id      uint
		setup   func(sqlmock.Sqlmock)
		want    string
		wantErr bool
	}{
		{
			name: "found",
			id:   1,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "給与")
				mock.ExpectQuery("SELECT \\* FROM `income_categories` WHERE `income_categories`.`id` = \\?").
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			want: "給与",
		},
		{
			name: "not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT \\* FROM `income_categories` WHERE `income_categories`.`id` = \\?").
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
			repo := &IncomeCategoryRepository{DB: db}

			category, err := repo.FindByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err == gorm.ErrRecordNotFound && tt.wantErr {
				// expected not-found error
			} else if !tt.wantErr && category.Name != tt.want {
				t.Fatalf("FindByID() name = %q, want %q", category.Name, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
