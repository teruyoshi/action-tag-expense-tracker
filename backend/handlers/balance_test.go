package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"action-tag-expense-tracker/backend/models"
)

func TestBalanceHandler_Get(t *testing.T) {
	tests := []struct {
		name     string
		repo     *mockBalanceRepo
		wantCode int
		wantBody string
	}{
		{
			name:     "returns balance",
			repo:     &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 50000}},
			wantCode: http.StatusOK,
			wantBody: `{"id":1,"amount":50000,"updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:     "returns 500 on db error",
			repo:     &mockBalanceRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BalanceHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodGet, "/balance", nil)
			w := httptest.NewRecorder()

			h.Get(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
			if tt.wantBody != "" {
				got := strings.TrimSpace(w.Body.String())
				if got != tt.wantBody {
					t.Errorf("body = %s, want %s", got, tt.wantBody)
				}
			}
		})
	}
}

func TestBalanceHandler_Update(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		repo           *mockBalanceRepo
		tagRepo        *mockActionTagRepo
		eventRepo      *mockEventRepo
		expenseRepo    *mockExpenseRepo
		wantCode       int
		wantBody       string
		wantEvent      bool
		wantExpense    bool
		wantExpenseAmt int
	}{
		{
			name:        "updates balance without expense when amount increases",
			body:        `{"amount":150000}`,
			repo:        &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			tagRepo:     &mockActionTagRepo{},
			eventRepo:   &mockEventRepo{},
			expenseRepo: &mockExpenseRepo{},
			wantCode:    http.StatusOK,
			wantBody:    `{"id":1,"amount":150000,"updated_at":"0001-01-01T00:00:00Z"}`,
			wantEvent:   false,
			wantExpense: false,
		},
		{
			name:           "creates expense when amount decreases",
			body:           `{"amount":80000}`,
			repo:           &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			tagRepo:        &mockActionTagRepo{tags: []models.ActionTag{{ID: 5, Name: "その他"}}},
			eventRepo:      &mockEventRepo{},
			expenseRepo:    &mockExpenseRepo{},
			wantCode:       http.StatusOK,
			wantBody:       `{"id":1,"amount":80000,"updated_at":"0001-01-01T00:00:00Z"}`,
			wantEvent:      true,
			wantExpense:    true,
			wantExpenseAmt: 20000,
		},
		{
			name:        "does not create expense when amount is same",
			body:        `{"amount":100000}`,
			repo:        &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			tagRepo:     &mockActionTagRepo{},
			eventRepo:   &mockEventRepo{},
			expenseRepo: &mockExpenseRepo{},
			wantCode:    http.StatusOK,
			wantBody:    `{"id":1,"amount":100000,"updated_at":"0001-01-01T00:00:00Z"}`,
			wantEvent:   false,
			wantExpense: false,
		},
		{
			name:     "rejects invalid json",
			body:     `{invalid`,
			repo:     &mockBalanceRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 500 on db error",
			body:     `{"amount":100000}`,
			repo:     &mockBalanceRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BalanceHandler{
				Repo:          tt.repo,
				ActionTagRepo: tt.tagRepo,
				EventRepo:     tt.eventRepo,
				ExpenseRepo:   tt.expenseRepo,
			}
			r := httptest.NewRequest(http.MethodPut, "/balance", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.Update(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
			if tt.wantBody != "" {
				got := strings.TrimSpace(w.Body.String())
				if got != tt.wantBody {
					t.Errorf("body = %s, want %s", got, tt.wantBody)
				}
			}
			if tt.wantEvent && tt.eventRepo != nil && tt.eventRepo.created == nil {
				t.Error("expected event to be created, but it was not")
			}
			if !tt.wantEvent && tt.eventRepo != nil && tt.eventRepo.created != nil {
				t.Error("expected no event, but one was created")
			}
			if tt.wantExpense && tt.expenseRepo != nil {
				if tt.expenseRepo.created == nil {
					t.Error("expected expense to be created, but it was not")
				} else {
					if tt.expenseRepo.created.Amount != tt.wantExpenseAmt {
						t.Errorf("expense amount = %d, want %d", tt.expenseRepo.created.Amount, tt.wantExpenseAmt)
					}
					if tt.expenseRepo.created.Item != "" {
						t.Errorf("expense item = %q, want empty string", tt.expenseRepo.created.Item)
					}
				}
			}
			if !tt.wantExpense && tt.expenseRepo != nil && tt.expenseRepo.created != nil {
				t.Error("expected no expense, but one was created")
			}
		})
	}
}
