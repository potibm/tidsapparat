package domain

import (
	"time"
)

type ScheduleEntry struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	CategoryID  *int64    `json:"category_id,omitempty"`
	Category    *Category `json:"category,omitempty"`
	LocationID  *int64    `json:"location_id,omitempty"`
	Location    *Location `json:"location,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type TimeTable []*ScheduleEntry
