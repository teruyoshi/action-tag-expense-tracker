package handlers

import (
	"net/http/httptest"
	"testing"
)

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name     string
		data     any
		wantBody string
	}{
		{
			name:     "map",
			data:     map[string]int{"total": 1000},
			wantBody: "{\"total\":1000}\n",
		},
		{
			name:     "slice",
			data:     []string{"a", "b"},
			wantBody: "[\"a\",\"b\"]\n",
		},
		{
			name:     "struct",
			data:     struct{ Name string }{Name: "test"},
			wantBody: "{\"Name\":\"test\"}\n",
		},
		{
			name:     "empty slice",
			data:     []string{},
			wantBody: "[]\n",
		},
		{
			name:     "nil",
			data:     nil,
			wantBody: "null\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			writeJSON(w, tt.data)

			if ct := w.Header().Get("Content-Type"); ct != "application/json" {
				t.Errorf("Content-Type = %q, want %q", ct, "application/json")
			}
			if got := w.Body.String(); got != tt.wantBody {
				t.Errorf("body = %q, want %q", got, tt.wantBody)
			}
		})
	}
}
