package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"action-tag-expense-tracker/backend/repositories"
)

func TestSummaryHandler_MonthTotal(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		repo     *mockSummaryRepo
		wantCode int
		wantBody string
	}{
		{
			name:     "returns total",
			query:    "?year=2026&month=3",
			repo:     &mockSummaryRepo{total: 15000},
			wantCode: http.StatusOK,
			wantBody: `{"total":15000}`,
		},
		{
			name:     "returns zero total",
			query:    "?year=2026&month=3",
			repo:     &mockSummaryRepo{total: 0},
			wantCode: http.StatusOK,
			wantBody: `{"total":0}`,
		},
		{
			name:     "rejects missing params",
			query:    "",
			repo:     &mockSummaryRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 500 on db error",
			query:    "?year=2026&month=3",
			repo:     &mockSummaryRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SummaryHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodGet, "/summary/month"+tt.query, nil)
			w := httptest.NewRecorder()

			h.MonthTotal(w, r)

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

func TestSummaryHandler_TagMonthTotals(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		repo     *mockSummaryRepo
		wantCode int
		wantBody string
	}{
		{
			name:  "returns tag totals",
			query: "?year=2026&month=3",
			repo: &mockSummaryRepo{tagTotals: []repositories.TagSummary{
				{TagID: 1, Tag: "通勤", Total: 10000},
				{TagID: 2, Tag: "外食", Total: 5000},
			}},
			wantCode: http.StatusOK,
			wantBody: `[{"tag_id":1,"tag":"通勤","total":10000},{"tag_id":2,"tag":"外食","total":5000}]`,
		},
		{
			name:     "rejects missing params",
			query:    "?year=2026",
			repo:     &mockSummaryRepo{},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SummaryHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodGet, "/summary/tag"+tt.query, nil)
			w := httptest.NewRecorder()

			h.TagMonthTotals(w, r)

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

func TestSummaryHandler_TagExpenseDetails(t *testing.T) {
	tests := []struct {
		name     string
		query    string
		repo     *mockSummaryRepo
		wantCode int
		wantBody string
	}{
		{
			name:  "returns details",
			query: "?year=2026&month=3&tag_id=1",
			repo: &mockSummaryRepo{tagDetails: []repositories.TagExpenseDetail{
				{ID: 1, Date: "2026-03-15", Item: "電車賃", Amount: 500},
			}},
			wantCode: http.StatusOK,
			wantBody: `[{"id":1,"date":"2026-03-15","item":"電車賃","amount":500}]`,
		},
		{
			name:     "rejects missing tag_id",
			query:    "?year=2026&month=3",
			repo:     &mockSummaryRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects missing year/month",
			query:    "?tag_id=1",
			repo:     &mockSummaryRepo{},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &SummaryHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodGet, "/summary/tag/details"+tt.query, nil)
			w := httptest.NewRecorder()

			h.TagExpenseDetails(w, r)

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
