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
	CategoryID  *int64
	Category    *dbCategory `gorm:"foreignKey:CategoryID"`
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
		CategoryID:  s.CategoryID,
	}
}

func toDomainScheduleEntry(db *dbScheduleEntry) *domain.ScheduleEntry {
	var category *domain.Category
	if db.Category != nil {
		category = toDomainCategory(db.Category)
	}

	return &domain.ScheduleEntry{
		ID:          db.ID,
		Title:       db.Title,
		Description: db.Description,
		StartTime:   db.StartTime,
		EndTime:     db.EndTime,
		Location:    db.Location,
		CategoryID:  db.CategoryID,
		Category:    category,
	}
}
