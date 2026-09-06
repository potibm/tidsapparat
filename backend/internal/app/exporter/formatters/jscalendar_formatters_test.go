package formatters

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/potibm/tidsapparat/internal/app/domain"

	"github.com/airtrafik/jscal"
)

func TestJsCalendarFormatter_RFC8984Compliance(t *testing.T) {
	formatter := NewJsCalendarFormatter(
		"-//Tidsapparat//Test//EN",
		"Europe/Berlin",
		"AbenteuerHallenKALK, Köln",
	)

	startTime := time.Date(2026, 8, 21, 18, 0, 0, 0, time.UTC)
	endTime := startTime.Add(2 * time.Hour)

	entries := domain.TimeTable{
		{
			ID:          1,
			Title:       "Evoke 2026 Opening",
			Description: "Welcome to the demoparty!",
			StartTime:   startTime,
			EndTime:     endTime,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Location: &domain.Location{
				Name: "Main Hall",
			},
		},
	}

	jsonData, err := formatter.Format(entries)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	var parsedGroup jscal.Group

	err = json.Unmarshal(jsonData, &parsedGroup)
	if err != nil {
		t.Fatalf(
			"airtrafik/jscal could not parse the JSON (RFC violation?): %v\nGenerated JSON:\n%s",
			err,
			string(jsonData),
		)
	}

	if parsedGroup.Type != "Group" {
		t.Errorf("Expected Group-Type 'Group', got '%s'", parsedGroup.Type)
	}

	if *parsedGroup.ProdId != "-//Tidsapparat//Test//EN" {
		t.Errorf("ProdID was not correctly transferred")
	}

	uid := "event-1@demoparty.org"

	var parsedEvent *jscal.Event

	for _, e := range parsedGroup.Entries {
		if e.GetUID() == uid {
			event, ok := e.(*jscal.Event)
			if !ok {
				t.Fatalf(
					"Object with UID '%s' was found, but is of type %T, not *jscal.Event",
					uid,
					e,
				)
			}

			parsedEvent = event

			break
		}
	}

	if parsedEvent == nil {
		t.Fatalf("Event with UID '%s' was not found in the JSON by the external library", uid)
	}

	if *parsedEvent.Title != "Evoke 2026 Opening" {
		t.Errorf("Expected title 'Evoke 2026 Opening', got '%s'", *parsedEvent.Title)
	}

	if *parsedEvent.TimeZone != "Europe/Berlin" {
		t.Errorf("Expected TimeZone 'Europe/Berlin', got '%s'", *parsedEvent.TimeZone)
	}

	if *parsedEvent.Duration != "PT2H" {
		t.Errorf("Expected duration 'PT2H', got '%s'", *parsedEvent.Duration)
	}

	expectedStart := "2026-08-21T18:00:00"
	if parsedEvent.Start.String() != expectedStart {
		t.Errorf("Expected start '%s', got '%s'", expectedStart, *parsedEvent.Start)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"zero duration", 0, "PT0S"},
		{"negative duration", -5 * time.Minute, "PT0S"},
		{"only seconds", 45 * time.Second, "PT45S"},
		{"only minutes", 30 * time.Minute, "PT30M"},
		{"only hours", 2 * time.Hour, "PT2H"},
		{"hours and minutes", 1*time.Hour + 30*time.Minute, "PT1H30M"},
		{"hours minutes seconds", 1*time.Hour + 30*time.Minute + 45*time.Second, "PT1H30M45S"},
		{"minutes and seconds", 5*time.Minute + 30*time.Second, "PT5M30S"},
		{"hours and seconds", 2*time.Hour + 15*time.Second, "PT2H15S"},
		{"large duration", 24 * time.Hour, "PT24H"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatDuration(tt.duration)
			if result != tt.expected {
				t.Errorf("formatDuration(%v) = %q, want %q", tt.duration, result, tt.expected)
			}
		})
	}
}
