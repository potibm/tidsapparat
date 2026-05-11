package calendar

import (
	"testing"
	"time"
)

func TestGetWeekdayRelative(t *testing.T) {
	tests := []struct {
		name          string
		now           time.Time    // Simuliertes "Heute"
		targetWeekday time.Weekday // Welchen Tag suchen wir?
		want          string       // Erwartetes Datum im Format YYYY-MM-DD
	}{
		{
			name:          "Today is monday, search friday (future)",
			now:           time.Date(2024, 5, 13, 10, 0, 0, 0, time.UTC), // Montag
			targetWeekday: time.Friday,
			want:          "2024-05-17",
		},
		{
			name:          "Today is monday, search sunday (future)",
			now:           time.Date(2024, 5, 13, 10, 0, 0, 0, time.UTC), // Montag
			targetWeekday: time.Sunday,
			want:          "2024-05-19",
		},
		{
			name:          "Today is wednesday, search friday (future)",
			now:           time.Date(2024, 5, 15, 10, 0, 0, 0, time.UTC), // Mittwoch
			targetWeekday: time.Friday,
			want:          "2024-05-17",
		},
		{
			name:          "Today is sunday, search friday (past)",
			now:           time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC), // Sonntag
			targetWeekday: time.Friday,
			want:          "2024-05-17",
		},
		{
			name:          "Today is sunday, search sunday (today)",
			now:           time.Date(2024, 5, 19, 10, 0, 0, 0, time.UTC), // Sonntag
			targetWeekday: time.Sunday,
			want:          "2024-05-19",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetWeekdayRelative(tt.now, tt.targetWeekday)

			gotFormat := got.Format("2006-01-02")
			if gotFormat != tt.want {
				t.Errorf("GetWeekdayRelative() = %v, want %v", gotFormat, tt.want)
			}
		})
	}
}
