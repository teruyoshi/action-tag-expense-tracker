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
		name     string
		body     string
		repo     *mockBalanceRepo
		wantCode int
		wantBody string
	}{
		{
			name:     "updates balance",
			body:     `{"amount":100000}`,
			repo:     &mockBalanceRepo{},
			wantCode: http.StatusOK,
			wantBody: `{"id":1,"amount":100000,"updated_at":"0001-01-01T00:00:00Z"}`,
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
			h := &BalanceHandler{Repo: tt.repo}
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
