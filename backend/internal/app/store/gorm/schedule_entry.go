package gorm

import (
	"time"

	"github.com/potibm/tidsapparat/internal/app/domain"
)

type dbScheduleEntry struct {
	GormModel
	AuditModel

	Title       string
	Description string
	ExternalURL string
	StartTime   time.Time
	EndTime     time.Time
	Hidden      bool `gorm:"not null;default:false"`
	CategoryID  *int64
	Category    *dbCategory `gorm:"foreignKey:CategoryID"`
	LocationID  *int64
	Location    *dbLocation `gorm:"foreignKey:LocationID"`
}

func (dbScheduleEntry) TableName() string {
	return "schedule_entries"
}

func fromDomainScheduleEntry(s *domain.ScheduleEntry) *dbScheduleEntry {
	return &dbScheduleEntry{
		GormModel: GormModel{ID: s.ID},

		Title:       s.Title,
		Description: s.Description,
		ExternalURL: s.ExternalURL,
		StartTime:   s.StartTime,
		EndTime:     s.EndTime,
		Hidden:      s.Hidden,
		CategoryID:  s.CategoryID,
		LocationID:  s.LocationID,
	}
}

func toDomainScheduleEntry(db *dbScheduleEntry) *domain.ScheduleEntry {
	var category *domain.Category
	if db.Category != nil {
		category = toDomainCategory(db.Category)
	}

	var location *domain.Location
	if db.Location != nil {
		location = toDomainLocation(db.Location)
	}

	return &domain.ScheduleEntry{
		ID:          db.ID,
		Title:       db.Title,
		Description: db.Description,
		ExternalURL: db.ExternalURL,
		StartTime:   db.StartTime,
		EndTime:     db.EndTime,
		Hidden:      db.Hidden,
		CategoryID:  db.CategoryID,
		Category:    category,
		LocationID:  db.LocationID,
		Location:    location,
		CreatedAt:   db.CreatedAt,
		UpdatedAt:   db.UpdatedAt,
	}
}

func toDomainScheduleEntries(db *[]dbScheduleEntry) domain.TimeTable {
	entries := make(domain.TimeTable, len(*db))
	for i, dbEntry := range *db {
		entries[i] = toDomainScheduleEntry(&dbEntry)
	}

	return entries
}
