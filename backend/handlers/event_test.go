package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestEventHandler_Create(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		repo     *mockEventRepo
		wantCode int
	}{
		{
			name:     "creates event",
			body:     `{"date":"2026-03-15","action_tag_id":1}`,
			repo:     &mockEventRepo{},
			wantCode: http.StatusCreated,
		},
		{
			name:     "rejects invalid json",
			body:     `{invalid`,
			repo:     &mockEventRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects invalid date format",
			body:     `{"date":"2026/03/15","action_tag_id":1}`,
			repo:     &mockEventRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects empty date",
			body:     `{"date":"","action_tag_id":1}`,
			repo:     &mockEventRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 500 on db error",
			body:     `{"date":"2026-03-15","action_tag_id":1}`,
			repo:     &mockEventRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &EventHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodPost, "/events", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.Create(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}
