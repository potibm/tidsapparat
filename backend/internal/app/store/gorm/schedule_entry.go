package gorm

import (
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type dbScheduleEntry struct {
	GormModel

	Title       string
	Description string
	StartTime   time.Time
	EndTime     time.Time
	Location    string
	Category    string
}

func (dbScheduleEntry) TableName() string {
	return "schedule_entries"
}

func fromDomainScheduleEntry(s *domain.ScheduleEntry) *dbScheduleEntry {
	return &dbScheduleEntry{
		GormModel: GormModel{ID: s.ID},

		Title:       s.Title,
		Description: s.Description,
		StartTime:   s.StartTime,
		EndTime:     s.EndTime,
		Location:    s.Location,
		Category:    s.Category,
	}
}

func toDomainScheduleEntry(db *dbScheduleEntry) *domain.ScheduleEntry {
	return &domain.ScheduleEntry{
		ID:          db.ID,
		Title:       db.Title,
		Description: db.Description,
		StartTime:   db.StartTime,
		EndTime:     db.EndTime,
		Location:    db.Location,
		Category:    db.Category,
	}
}
