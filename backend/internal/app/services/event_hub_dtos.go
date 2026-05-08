package services

import (
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type ActionType string

const (
	ActionCreate ActionType = "create"
	ActionUpdate ActionType = "update"
	ActionDelete ActionType = "delete"
	ActionSync   ActionType = "sync"
)

type ScheduleSyncEventDTO struct {
	Action    ActionType         `json:"action"`
	Timestamp int64              `json:"timestamp"`
	Payload   []ScheduleEntryDTO `json:"payload"`
}

type ScheduleEventDTO struct {
	Action    ActionType       `json:"action"`
	Timestamp int64            `json:"timestamp"`
	Payload   ScheduleEntryDTO `json:"payload"`
}

type ScheduleEntryDTO struct {
	ID          int64        `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	StartTime   string       `json:"start_time"` // RFC3339
	EndTime     string       `json:"end_time"`   // RFC3339
	Category    *CategoryDTO `json:"category,omitempty"`
	Location    *LocationDTO `json:"location,omitempty"`
}

type CategoryDTO struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type LocationDTO struct {
	Name    string `json:"name"`
	Address string `json:"address,omitempty"`
}

func mapToEntryDTO(entry *domain.ScheduleEntry) ScheduleEntryDTO {
	dto := ScheduleEntryDTO{
		ID:          entry.ID,
		Title:       entry.Title,
		Description: entry.Description,
		StartTime:   entry.StartTime.Format(time.RFC3339),
		EndTime:     entry.EndTime.Format(time.RFC3339),
	}

	if entry.Category != nil {
		dto.Category = &CategoryDTO{
			Name:  entry.Category.Name,
			Color: entry.Category.Color,
		}
	}

	if entry.Location != nil {
		dto.Location = &LocationDTO{
			Name: entry.Location.Name,
		}

		if entry.Location.Address != nil && *entry.Location.Address != "" {
			dto.Location.Address = *entry.Location.Address
		}
	}

	return dto
}

func mapToEventDTO(entry *domain.ScheduleEntry, action ActionType) ScheduleEventDTO {
	return ScheduleEventDTO{
		Action:    action,
		Timestamp: time.Now().Unix(),
		Payload:   mapToEntryDTO(entry),
	}
}

func mapToTimeTableDTO(entries domain.TimeTable) []ScheduleEntryDTO {
	dtos := make([]ScheduleEntryDTO, 0, len(entries))
	for _, entry := range entries {
		dtos = append(dtos, mapToEntryDTO(entry))
	}

	return dtos
}
