package hub

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/potibm/tidsapparat/internal/app/config"
	"github.com/stretchr/testify/assert"
)

func TestMapToPublicConfig(t *testing.T) {
	internalCfg := &config.Config{
		App: config.AppConfig{
			Version:            "1.2.3",
			Environment:        "production",
			EnvironmentMessage: "Hello World",
		},
		Sentry: config.SentryConfig{
			DSN:         "https://secret@sentry.io/123",
			Environment: "prod",
			Version:     "v1",
		},
	}

	public := mapToPublicConfig(internalCfg)

	// Verification
	assert.Equal(t, "1.2.3", public.Version)
	assert.Equal(t, "production", public.Environment)
	assert.Equal(t, "Hello World", public.EnvironmentMessage)
	assert.Equal(t, "https://secret@sentry.io/123", public.Sentry.DSN)

	payload, err := json.Marshal(public)
	assert.NoError(t, err)
	assert.NotContains(t, string(payload), "admin_api_key")
}

func TestGeneratePartyDays(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
		want      []PartyDaysPublic
		wantErr   bool
	}{
		{
			name:      "Short period (<= 7 days) uses only weekdays",
			startDate: "2026-05-01", // 1st May 2026 is a Friday
			endDate:   "2026-05-03",
			want: []PartyDaysPublic{
				{ID: "2026-05-01", Name: "Friday"},
				{ID: "2026-05-02", Name: "Saturday"},
				{ID: "2026-05-03", Name: "Sunday"},
			},
			wantErr: false,
		},
		{
			name:      "Long period (> 7 days) adds ordinal suffixes",
			startDate: "2026-05-01",
			endDate:   "2026-05-08", // 8 days total
			want: []PartyDaysPublic{
				{ID: "2026-05-01", Name: "Friday, 1st May"},
				{ID: "2026-05-02", Name: "Saturday, 2nd May"},
				{ID: "2026-05-03", Name: "Sunday, 3rd May"},
				{ID: "2026-05-04", Name: "Monday, 4th May"},
				{ID: "2026-05-05", Name: "Tuesday, 5th May"},
				{ID: "2026-05-06", Name: "Wednesday, 6th May"},
				{ID: "2026-05-07", Name: "Thursday, 7th May"},
				{ID: "2026-05-08", Name: "Friday, 8th May"},
			},
			wantErr: false,
		},
		{
			name:      "Start date after end date returns empty slice",
			startDate: "2026-05-10",
			endDate:   "2026-05-01",
			want:      []PartyDaysPublic{},
			wantErr:   false,
		},
		{
			name:      "Invalid start date",
			startDate: "invalid-date",
			endDate:   "2026-05-01",
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "Invalid end date",
			startDate: "2026-05-01",
			endDate:   "01-05-2026", // Wrong layout
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GeneratePartyDays(tt.startDate, tt.endDate)
			if (err != nil) != tt.wantErr {
				t.Errorf("GeneratePartyDays() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GeneratePartyDays() \ngot  = %v\nwant = %v", got, tt.want)
			}
		})
	}
}

func TestGeneratePartyDaysOrEmpty(t *testing.T) {
	tests := []struct {
		name      string
		startDate string
		endDate   string
		wantLen   int
	}{
		{
			name:      "Valid dates return slice with elements",
			startDate: "2026-05-01",
			endDate:   "2026-05-03",
			wantLen:   3,
		},
		{
			name:      "Invalid dates are swallowed and return empty slice",
			startDate: "broken",
			endDate:   "2026-05-03",
			wantLen:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generatePartyDaysOrEmpty(tt.startDate, tt.endDate)

			if got == nil {
				t.Error("generatePartyDaysOrEmpty() returned nil, want empty slice []PartyDaysPublic{}")
			}

			if len(got) != tt.wantLen {
				t.Errorf("generatePartyDaysOrEmpty() returned %d elements, want %d", len(got), tt.wantLen)
			}
		})
	}
}

func TestFormatOrdinal(t *testing.T) {
	tests := []struct {
		name string
		day  int
		want string
	}{
		{"1st", 1, "1st May"},
		{"2nd", 2, "2nd May"},
		{"3rd", 3, "3rd May"},
		{"4th", 4, "4th May"},
		{"11th", 11, "11th May"},
		{"21st", 21, "21st May"},
		{"22nd", 22, "22nd May"},
		{"23rd", 23, "23rd May"},
		{"31st", 31, "31st May"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDate := time.Date(2026, time.May, tt.day, 12, 0, 0, 0, time.UTC)
			got := formatOrdinal(testDate)

			if got != tt.want {
				t.Errorf("formatOrdinal() for day %d = %v, want %v", tt.day, got, tt.want)
			}
		})
	}
}
