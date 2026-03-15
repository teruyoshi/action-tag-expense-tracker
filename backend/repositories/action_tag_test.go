package repositories

import (
	"testing"

	"action-tag-expense-tracker/backend/models"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
)

func TestActionTagRepository_FindAll(t *testing.T) {
	tests := []struct {
		name    string
		setup   func(sqlmock.Sqlmock)
		want    int
		wantErr bool
	}{
		{
			name: "returns all tags",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).
					AddRow(1, "通勤").
					AddRow(2, "外食")
				mock.ExpectQuery("SELECT \\* FROM `action_tags`").WillReturnRows(rows)
			},
			want: 2,
		},
		{
			name: "returns empty when no tags",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT \\* FROM `action_tags`").WillReturnRows(rows)
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &ActionTagRepository{DB: db}

			tags, err := repo.FindAll()
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if len(tags) != tt.want {
				t.Fatalf("FindAll() got %d tags, want %d", len(tags), tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestActionTagRepository_Create(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("INSERT INTO `action_tags`").
		WithArgs("買い物").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &ActionTagRepository{DB: db}
	tag := &models.ActionTag{Name: "買い物"}

	if err := repo.Create(tag); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestActionTagRepository_Update(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("UPDATE `action_tags`").
		WithArgs("食費", 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repo := &ActionTagRepository{DB: db}
	tag := &models.ActionTag{ID: 1, Name: "食費"}

	if err := repo.Update(tag); err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestActionTagRepository_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM `action_tags` WHERE `action_tags`.`id` = ?").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := &ActionTagRepository{DB: db}

	if err := repo.Delete(1); err != nil {
		t.Fatalf("Delete() error = %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestActionTagRepository_FindByID(t *testing.T) {
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
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "通勤")
				mock.ExpectQuery("SELECT \\* FROM `action_tags` WHERE `action_tags`.`id` = \\?").
					WithArgs(1, 1).
					WillReturnRows(rows)
			},
			want: "通勤",
		},
		{
			name: "not found",
			id:   999,
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT \\* FROM `action_tags` WHERE `action_tags`.`id` = \\?").
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
			repo := &ActionTagRepository{DB: db}

			tag, err := repo.FindByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && err == gorm.ErrRecordNotFound && tt.wantErr {
				// expected not-found error
			} else if !tt.wantErr && tag.Name != tt.want {
				t.Fatalf("FindByID() name = %q, want %q", tag.Name, tt.want)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func TestActionTagRepository_FindOrCreateByName(t *testing.T) {
	tests := []struct {
		name    string
		tagName string
		setup   func(sqlmock.Sqlmock)
		wantID  uint
	}{
		{
			name:    "existing tag found",
			tagName: "通勤",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "通勤")
				mock.ExpectQuery("SELECT \\* FROM `action_tags` WHERE name = \\?").
					WithArgs("通勤", 1).
					WillReturnRows(rows)
			},
			wantID: 1,
		},
		{
			name:    "new tag created",
			tagName: "新規",
			setup: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name"})
				mock.ExpectQuery("SELECT \\* FROM `action_tags` WHERE name = \\?").
					WithArgs("新規", 1).
					WillReturnRows(rows)
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO `action_tags`").
					WithArgs("新規").
					WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectCommit()
			},
			wantID: 2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock := setupMockDB(t)
			tt.setup(mock)
			repo := &ActionTagRepository{DB: db}

			tag, err := repo.FindOrCreateByName(tt.tagName)
			if err != nil {
				t.Fatalf("FindOrCreateByName() error = %v", err)
			}
			if tag.ID != tt.wantID {
				t.Fatalf("FindOrCreateByName() ID = %d, want %d", tag.ID, tt.wantID)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
