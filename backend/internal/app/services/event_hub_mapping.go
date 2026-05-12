package services

import (
	"time"

	"github.com/potibm/protokolapparat/pkg/schedule"
	"github.com/potibm/tidsapparat/internal/app/domain"
)

func mapToEventPayload(entry *domain.ScheduleEntry) schedule.Entry {
	result := schedule.Entry{
		ID:          entry.ID,
		Title:       entry.Title,
		Description: entry.Description,
		ExternalURL: entry.ExternalURL,
		StartTime:   entry.StartTime.Format(time.RFC3339),
		EndTime:     entry.EndTime.Format(time.RFC3339),
		Hidden:      entry.Hidden,
	}

	if entry.Category != nil {
		result.Category = &schedule.Category{
			Name:  entry.Category.Name,
			Color: entry.Category.Color,
		}
	}

	if entry.Location != nil {
		result.Location = &schedule.Location{
			Name: entry.Location.Name,
		}

		if entry.Location.Address != nil && *entry.Location.Address != "" {
			result.Location.Address = *entry.Location.Address
		}
	}

	return result
}

func mapToTimeTablePayload(entries domain.TimeTable) []schedule.Entry {
	dtos := make([]schedule.Entry, 0, len(entries))
	for _, entry := range entries {
		dtos = append(dtos, mapToEventPayload(entry))
	}

	return dtos
}
