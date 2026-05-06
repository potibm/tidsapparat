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
	Location    string    `json:"location"`
	CategoryID  *int64    `json:"category_id,omitempty"`
	Category    *Category `json:"category,omitempty"`
}
