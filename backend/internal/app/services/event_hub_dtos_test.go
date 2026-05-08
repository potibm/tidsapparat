package services

import (
	"testing"
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/stretchr/testify/assert"
)

func TestMapToEntryDTO(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 30, 0, 0, time.UTC)
	address := "123 Main St"

	tests := []struct {
		name     string
		entry    *domain.ScheduleEntry
		expected ScheduleEntryDTO
	}{
		{
			name: "full entry with category and location",
			entry: &domain.ScheduleEntry{
				ID:          1,
				Title:       "Opening Ceremony",
				Description: "The grand opening",
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
			expected: ScheduleEntryDTO{
				ID:          1,
				Title:       "Opening Ceremony",
				Description: "The grand opening",
				StartTime:   "2024-06-15T10:30:00Z",
				EndTime:     "2024-06-15T12:30:00Z",
				Category: &CategoryDTO{
					Name:  "Ceremony",
					Color: "#FF0000",
				},
				Location: &LocationDTO{
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
				StartTime:   now,
				EndTime:     now.Add(1 * time.Hour),
			},
			expected: ScheduleEntryDTO{
				ID:          2,
				Title:       "Lunch Break",
				Description: "Free time",
				StartTime:   "2024-06-15T10:30:00Z",
				EndTime:     "2024-06-15T11:30:00Z",
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
			expected: ScheduleEntryDTO{
				ID:        3,
				Title:     "Quick Meeting",
				StartTime: "2024-06-15T10:30:00Z",
				EndTime:   "2024-06-15T11:00:00Z",
				Location: &LocationDTO{
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
			expected: ScheduleEntryDTO{
				ID:        4,
				Title:     "Outdoor Event",
				StartTime: "2024-06-15T10:30:00Z",
				EndTime:   "2024-06-15T11:15:00Z",
				Location: &LocationDTO{
					Name: "Park",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapToEntryDTO(tt.entry)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestMapToEventDTO(t *testing.T) {
	now := time.Date(2024, 6, 15, 10, 0, 0, 0, time.UTC)
	entry := &domain.ScheduleEntry{
		ID:        1,
		Title:     "Test Event",
		StartTime: now,
		EndTime:   now.Add(time.Hour),
	}

	result := mapToEventDTO(entry, ActionCreate)

	assert.Equal(t, ActionCreate, result.Action)
	assert.Equal(t, ScheduleEntryDTO{
		ID:        1,
		Title:     "Test Event",
		StartTime: "2024-06-15T10:00:00Z",
		EndTime:   "2024-06-15T11:00:00Z",
	}, result.Payload)

	// Timestamp should be within the last second
	assert.InDelta(t, time.Now().Unix(), result.Timestamp, 1)
}

func TestMapToTimeTableDTO(t *testing.T) {
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

	result := mapToTimeTableDTO(entries)

	assert.Len(t, result, 2)
	assert.Equal(t, ScheduleEntryDTO{
		ID:        1,
		Title:     "First",
		StartTime: "2024-06-15T10:00:00Z",
		EndTime:   "2024-06-15T11:00:00Z",
	}, result[0])
	assert.Equal(t, ScheduleEntryDTO{
		ID:        2,
		Title:     "Second",
		StartTime: "2024-06-15T12:00:00Z",
		EndTime:   "2024-06-15T13:00:00Z",
	}, result[1])
}

func TestMapToTimeTableDTO_Empty(t *testing.T) {
	result := mapToTimeTableDTO(domain.TimeTable{})

	assert.NotNil(t, result)
	assert.Empty(t, result)
}

func strPtr(s string) *string {
	return &s
}
