package handlers

import (
	"testing"
	"time"
)

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
