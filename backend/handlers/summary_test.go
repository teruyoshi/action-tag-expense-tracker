package handlers

import (
	"net/http"
	"net/url"
	"testing"
)

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
