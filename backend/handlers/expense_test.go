package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExpenseHandler_Create(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		repo     *mockExpenseRepo
		balRepo  *mockBalanceRepo
		wantCode int
	}{
		{
			name:     "creates expense and subtracts balance",
			body:     `{"event_id":1,"item":"電車賃","amount":500}`,
			repo:     &mockExpenseRepo{},
			balRepo:  &mockBalanceRepo{},
			wantCode: http.StatusCreated,
		},
		{
			name:     "creates expense without balance repo",
			body:     `{"event_id":1,"item":"電車賃","amount":500}`,
			repo:     &mockExpenseRepo{},
			balRepo:  nil,
			wantCode: http.StatusCreated,
		},
		{
			name:     "rejects zero amount",
			body:     `{"event_id":1,"item":"電車賃","amount":0}`,
			repo:     &mockExpenseRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects negative amount",
			body:     `{"event_id":1,"item":"電車賃","amount":-100}`,
			repo:     &mockExpenseRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects missing event_id",
			body:     `{"item":"電車賃","amount":500}`,
			repo:     &mockExpenseRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects invalid json",
			body:     `{invalid`,
			repo:     &mockExpenseRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 500 on expense db error",
			body:     `{"event_id":1,"item":"電車賃","amount":500}`,
			repo:     &mockExpenseRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "returns 500 on balance subtract error",
			body:     `{"event_id":1,"item":"電車賃","amount":500}`,
			repo:     &mockExpenseRepo{},
			balRepo:  &mockBalanceRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ExpenseHandler{Repo: tt.repo}
			if tt.balRepo != nil {
				h.BalanceRepo = tt.balRepo
			}
			r := httptest.NewRequest(http.MethodPost, "/expenses", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.Create(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
			if tt.balRepo != nil && tt.wantCode == http.StatusCreated && tt.balRepo.subtracted != 500 {
				t.Errorf("subtracted = %d, want 500", tt.balRepo.subtracted)
			}
		})
	}
}
