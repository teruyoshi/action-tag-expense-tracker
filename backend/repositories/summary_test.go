package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestSummaryRepository_MonthTotal(t *testing.T) {
	db, mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"coalesce"}).AddRow(15000)
	mock.ExpectQuery("SELECT COALESCE\\(SUM\\(e.amount\\), 0\\)").
		WithArgs(2026, 3).
		WillReturnRows(rows)

	repo := &SummaryRepository{DB: db}

	total, err := repo.MonthTotal(2026, 3)
	if err != nil {
		t.Fatalf("MonthTotal() error = %v", err)
	}
	if total != 15000 {
		t.Fatalf("MonthTotal() = %d, want 15000", total)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSummaryRepository_TagMonthTotals(t *testing.T) {
	db, mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"tag_id", "tag", "total"}).
		AddRow(1, "通勤", 8000).
		AddRow(2, "外食", 7000)
	mock.ExpectQuery("SELECT at.id AS tag_id").
		WithArgs(2026, 3).
		WillReturnRows(rows)

	repo := &SummaryRepository{DB: db}

	results, err := repo.TagMonthTotals(2026, 3)
	if err != nil {
		t.Fatalf("TagMonthTotals() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("TagMonthTotals() got %d results, want 2", len(results))
	}
	if results[0].Tag != "通勤" || results[0].Total != 8000 {
		t.Fatalf("TagMonthTotals()[0] = %+v, want tag=通勤, total=8000", results[0])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}

func TestSummaryRepository_TagExpenseDetails(t *testing.T) {
	db, mock := setupMockDB(t)

	rows := sqlmock.NewRows([]string{"id", "date", "item", "amount"}).
		AddRow(1, "2026-03-10", "電車", 500).
		AddRow(2, "2026-03-10", "バス", 300)
	mock.ExpectQuery("SELECT e.id").
		WithArgs(2026, 3, 1).
		WillReturnRows(rows)

	repo := &SummaryRepository{DB: db}

	results, err := repo.TagExpenseDetails(2026, 3, 1)
	if err != nil {
		t.Fatalf("TagExpenseDetails() error = %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("TagExpenseDetails() got %d results, want 2", len(results))
	}
	if results[0].Item != "電車" || results[0].Amount != 500 {
		t.Fatalf("TagExpenseDetails()[0] = %+v, want item=電車, amount=500", results[0])
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatal(err)
	}
}
