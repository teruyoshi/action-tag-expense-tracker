package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
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

func TestParseDate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    time.Time
		wantErr bool
	}{
		{
			name:  "valid date",
			input: "2026-03-14",
			want:  time.Date(2026, 3, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "first day of year",
			input: "2026-01-01",
			want:  time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:  "last day of year",
			input: "2025-12-31",
			want:  time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
		{
			name:    "wrong format slash",
			input:   "2026/03/14",
			wantErr: true,
		},
		{
			name:    "wrong format no zero padding",
			input:   "2026-3-14",
			wantErr: true,
		},
		{
			name:    "invalid date",
			input:   "not-a-date",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDate(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
				return
			}
			if !tt.wantErr && !got.Equal(tt.want) {
				t.Errorf("parseDate(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseYearMonth(t *testing.T) {
	tests := []struct {
		name      string
		query     url.Values
		wantYear  int
		wantMonth int
		wantErr   bool
	}{
		{
			name:      "valid year and month",
			query:     url.Values{"year": {"2026"}, "month": {"3"}},
			wantYear:  2026,
			wantMonth: 3,
		},
		{
			name:      "january",
			query:     url.Values{"year": {"2025"}, "month": {"1"}},
			wantYear:  2025,
			wantMonth: 1,
		},
		{
			name:      "december",
			query:     url.Values{"year": {"2025"}, "month": {"12"}},
			wantYear:  2025,
			wantMonth: 12,
		},
		{
			name:    "missing year",
			query:   url.Values{"month": {"3"}},
			wantErr: true,
		},
		{
			name:    "missing month",
			query:   url.Values{"year": {"2026"}},
			wantErr: true,
		},
		{
			name:    "empty query",
			query:   url.Values{},
			wantErr: true,
		},
		{
			name:    "non-numeric year",
			query:   url.Values{"year": {"abc"}, "month": {"3"}},
			wantErr: true,
		},
		{
			name:    "non-numeric month",
			query:   url.Values{"year": {"2026"}, "month": {"abc"}},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{URL: &url.URL{RawQuery: tt.query.Encode()}}
			year, month, err := parseYearMonth(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseYearMonth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if year != tt.wantYear {
					t.Errorf("parseYearMonth() year = %d, want %d", year, tt.wantYear)
				}
				if month != tt.wantMonth {
					t.Errorf("parseYearMonth() month = %d, want %d", month, tt.wantMonth)
				}
			}
		})
	}
}
