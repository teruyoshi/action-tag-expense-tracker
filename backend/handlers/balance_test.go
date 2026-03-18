package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/services"
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
		name     string
		body     string
		repo     *mockBalanceRepo
		service  *services.BalanceService
		wantCode int
		wantBody string
	}{
		{
			name: "updates balance via service",
			body: `{"amount":150000}`,
			repo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			service: &services.BalanceService{
				BalanceRepo:   &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
				ActionTagRepo: &mockActionTagRepo{},
				EventRepo:     &mockEventRepo{},
				ExpenseRepo:   &mockExpenseRepo{},
			},
			wantCode: http.StatusOK,
			wantBody: `{"id":1,"amount":150000,"updated_at":"0001-01-01T00:00:00Z"}`,
		},
		{
			name:     "rejects invalid json",
			body:     `{invalid`,
			repo:     &mockBalanceRepo{},
			service:  &services.BalanceService{},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "returns 500 on service error",
			body: `{"amount":100000}`,
			repo: &mockBalanceRepo{},
			service: &services.BalanceService{
				BalanceRepo:   &mockBalanceRepo{err: errDB},
				ActionTagRepo: &mockActionTagRepo{},
				EventRepo:     &mockEventRepo{},
				ExpenseRepo:   &mockExpenseRepo{},
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &BalanceHandler{
				Repo:    tt.repo,
				Service: tt.service,
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
		})
	}
}
