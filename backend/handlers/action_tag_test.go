package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"action-tag-expense-tracker/backend/models"

	"github.com/go-chi/chi/v5"
)

func TestActionTagHandler_List(t *testing.T) {
	tests := []struct {
		name       string
		repo       *mockActionTagRepo
		wantCode   int
		wantBody   string
	}{
		{
			name:     "returns tags",
			repo:     &mockActionTagRepo{tags: []models.ActionTag{{ID: 1, Name: "通勤"}, {ID: 2, Name: "外食"}}},
			wantCode: http.StatusOK,
			wantBody: `[{"id":1,"name":"通勤"},{"id":2,"name":"外食"}]`,
		},
		{
			name:     "returns empty array when no tags",
			repo:     &mockActionTagRepo{tags: []models.ActionTag{}},
			wantCode: http.StatusOK,
			wantBody: `[]`,
		},
		{
			name:     "returns 500 on db error",
			repo:     &mockActionTagRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ActionTagHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodGet, "/tags", nil)
			w := httptest.NewRecorder()

			h.List(w, r)

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

func TestActionTagHandler_Create(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		repo     *mockActionTagRepo
		wantCode int
	}{
		{
			name:     "creates tag",
			body:     `{"name":"通勤"}`,
			repo:     &mockActionTagRepo{},
			wantCode: http.StatusCreated,
		},
		{
			name:     "rejects empty name",
			body:     `{"name":""}`,
			repo:     &mockActionTagRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "rejects invalid json",
			body:     `{invalid`,
			repo:     &mockActionTagRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 500 on db error",
			body:     `{"name":"通勤"}`,
			repo:     &mockActionTagRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ActionTagHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodPost, "/tags", strings.NewReader(tt.body))
			w := httptest.NewRecorder()

			h.Create(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestActionTagHandler_Update(t *testing.T) {
	tagsInDB := []models.ActionTag{{ID: 1, Name: "通勤"}}

	tests := []struct {
		name     string
		id       string
		body     string
		repo     *mockActionTagRepo
		wantCode int
	}{
		{
			name:     "updates tag",
			id:       "1",
			body:     `{"name":"外食"}`,
			repo:     &mockActionTagRepo{tags: tagsInDB},
			wantCode: http.StatusOK,
		},
		{
			name:     "rejects invalid id",
			id:       "abc",
			body:     `{"name":"外食"}`,
			repo:     &mockActionTagRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 404 for missing tag",
			id:       "99",
			body:     `{"name":"外食"}`,
			repo:     &mockActionTagRepo{tags: tagsInDB},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "rejects empty name",
			id:       "1",
			body:     `{"name":""}`,
			repo:     &mockActionTagRepo{tags: tagsInDB},
			wantCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ActionTagHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodPut, "/tags/"+tt.id, strings.NewReader(tt.body))
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			h.Update(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}

func TestActionTagHandler_Delete(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		repo     *mockActionTagRepo
		wantCode int
	}{
		{
			name:     "deletes tag",
			id:       "1",
			repo:     &mockActionTagRepo{},
			wantCode: http.StatusNoContent,
		},
		{
			name:     "rejects invalid id",
			id:       "abc",
			repo:     &mockActionTagRepo{},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "returns 500 on db error",
			id:       "1",
			repo:     &mockActionTagRepo{err: errDB},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &ActionTagHandler{Repo: tt.repo}
			r := httptest.NewRequest(http.MethodDelete, "/tags/"+tt.id, nil)
			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", tt.id)
			r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
			w := httptest.NewRecorder()

			h.Delete(w, r)

			if w.Code != tt.wantCode {
				t.Errorf("status = %d, want %d", w.Code, tt.wantCode)
			}
		})
	}
}
