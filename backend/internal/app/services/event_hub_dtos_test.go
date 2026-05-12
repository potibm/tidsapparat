package services

import (
	"testing"
	"time"

	"github.com/potibm/protokolapparat/pkg/schedule"
	"github.com/potibm/tidsapparat/internal/app/domain"

	"github.com/stretchr/testify/assert"
)

func TestMapToEntryDTO(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	address := "123 Main St"

	tests := []struct {
		name     string
		entry    *domain.ScheduleEntry
		expected schedule.Entry
	}{
		{
			name: "full entry with category and location",
			entry: &domain.ScheduleEntry{
				ID:          1,
				Title:       "Opening Ceremony",
				Description: "The grand opening",
				ExternalURL: "https://example.com/opening",
				StartTime:   now,
				EndTime:     now.Add(2 * time.Hour),
				Category: &domain.Category{
					ID:    10,
					Name:  "Ceremony",
					Color: "#FF0000",
				},
				Location: &domain.Location{
					ID:      20,
					Name:    "Main Hall",
					Address: &address,
				},
			},
			expected: schedule.Entry{
				ID:          1,
				Title:       "Opening Ceremony",
				Description: "The grand opening",
				ExternalURL: "https://example.com/opening",
				StartTime:   "2024-06-15T10:30:00Z",
				EndTime:     "2024-06-15T12:30:00Z",
				Hidden:      false,
				Category: &schedule.Category{
					Name:  "Ceremony",
					Color: "#FF0000",
				},
				Location: &schedule.Location{
					Name:    "Main Hall",
					Address: "123 Main St",
				},
			},
		},
		{
			name: "entry without category and location",
			entry: &domain.ScheduleEntry{
				ID:          2,
				Title:       "Lunch Break",
				Description: "Free time",
				ExternalURL: "",
				StartTime:   now,
				EndTime:     now.Add(1 * time.Hour),
			},
			expected: schedule.Entry{
				ID:          2,
				Title:       "Lunch Break",
				Description: "Free time",
				ExternalURL: "",
				StartTime:   "2024-06-15T10:30:00Z",
				EndTime:     "2024-06-15T11:30:00Z",
				Hidden:      false,
				Category:    nil,
				Location:    nil,
			},
		},
		{
			name: "entry with empty address",
			entry: &domain.ScheduleEntry{
				ID:        3,
				Title:     "Quick Meeting",
				StartTime: now,
				EndTime:   now.Add(30 * time.Minute),
				Location: &domain.Location{
					ID:      30,
					Name:    "Room A",
					Address: strPtr(""),
				},
			},
			expected: schedule.Entry{
				ID:        3,
				Title:     "Quick Meeting",
				StartTime: "2024-06-15T10:30:00Z",
				EndTime:   "2024-06-15T11:00:00Z",
				Hidden:    false,
				Location: &schedule.Location{
					Name: "Room A",
				},
			},
		},
		{
			name: "entry with nil address",
			entry: &domain.ScheduleEntry{
				ID:        4,
				Title:     "Outdoor Event",
				StartTime: now,
				EndTime:   now.Add(45 * time.Minute),
				Location: &domain.Location{
					ID:   40,
					Name: "Park",
				},
			},
			expected: schedule.Entry{
				ID:        4,
				Title:     "Outdoor Event",
				StartTime: "2024-06-15T10:30:00Z",
				EndTime:   "2024-06-15T11:15:00Z",
				Hidden:    false,
				Location: &schedule.Location{
					Name: "Park",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToEventPayload(tt.entry)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToTimeTablePayload(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)

	entries := domain.TimeTable{
		{
			ID:        1,
			Title:     "First",
			StartTime: now,
			EndTime:   now.Add(time.Hour),
		},
		{
			ID:        2,
			Title:     "Second",
			StartTime: now.Add(2 * time.Hour),
			EndTime:   now.Add(3 * time.Hour),
		},
	}

	result := mapToTimeTablePayload(entries)

	assert.Len(t, result, 2)
	assert.Equal(t, schedule.Entry{
		ID:        1,
		Title:     "First",
		StartTime: "2024-06-15T10:00:00Z",
		EndTime:   "2024-06-15T11:00:00Z",
		Hidden:    false,
	}, result[0])
	assert.Equal(t, schedule.Entry{
		ID:        2,
		Title:     "Second",
		StartTime: "2024-06-15T12:00:00Z",
		EndTime:   "2024-06-15T13:00:00Z",
		Hidden:    false,
	}, result[1])
}

func TestMapToTimeTablePayload_Empty(t *testing.T) {
	result := mapToTimeTablePayload(domain.TimeTable{})

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func strPtr(s string) *string {
	return &s
}
