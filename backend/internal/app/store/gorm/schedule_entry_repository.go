package gorm

import (
	"context"
	"fmt"
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	orderStartTimeASC = "start_time ASC"
	orderEndTimeASC   = "end_time ASC"
)

type scheduleEntryRepository struct {
	db *gorm.DB
}

func (s *Store) NewScheduleEntryRepository() repository.ScheduleEntryRepository {
	return NewScheduleEntryRepository(s.db)
}

func NewScheduleEntryRepository(db *gorm.DB) repository.ScheduleEntryRepository {
	return &scheduleEntryRepository{db: db}
}

func (r *scheduleEntryRepository) Save(ctx context.Context, scheduleEntry *domain.ScheduleEntry) error {
	dbObj := fromDomainScheduleEntry(scheduleEntry)

	err := r.db.WithContext(ctx).Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(dbObj).Error
	if err == nil {
		scheduleEntry.ID = dbObj.ID
	}

	return err
}

func (r *scheduleEntryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&dbScheduleEntry{}, id).Error
}

func (r *scheduleEntryRepository) List(
	ctx context.Context,
	params repository.ScheduleEntryListParams,
	filters repository.ScheduleEntryListFilters,
) ([]domain.ScheduleEntry, int64, error) {
	var (
		dbEntries []dbScheduleEntry
		total     int64
	)

	query := r.db.WithContext(ctx).Model(&dbScheduleEntry{}).Preload("Category").Preload("Location")

	query = r.applyFilters(query, filters)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: params.Sort},
		Desc:   params.Order == "DESC",
	})

	if params.Sort != "start_time" {
		query = query.Order(orderStartTimeASC)
	}

	if params.Sort != "end_time" {
		query = query.Order(orderEndTimeASC)
	}

	err = query.Offset(params.Offset).
		Limit(params.Limit).
		Find(&dbEntries).
		Error
	if err != nil {
		return nil, 0, err
	}

	entries := make([]domain.ScheduleEntry, len(dbEntries))
	for i, dbEntry := range dbEntries {
		entries[i] = *toDomainScheduleEntry(&dbEntry)
	}

	return entries, total, nil
}

func (r *scheduleEntryRepository) GetByID(ctx context.Context, id int64) (*domain.ScheduleEntry, error) {
	var dbEntry dbScheduleEntry

	err := r.db.WithContext(ctx).Preload("Category").Preload("Location").First(&dbEntry, id).Error
	if err != nil {
		return nil, err
	}

	entry := toDomainScheduleEntry(&dbEntry)

	return entry, nil
}

func (r *scheduleEntryRepository) GetAllPreloaded(ctx context.Context) (domain.TimeTable, error) {
	var dbEntries []dbScheduleEntry

	err := r.db.WithContext(ctx).
		Preload("Category").
		Preload("Location").
		Order(orderStartTimeASC).
		Order(orderEndTimeASC).
		Find(&dbEntries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch timetable: %w", err)
	}

	var result domain.TimeTable
	for _, dbE := range dbEntries {
		result = append(result, toDomainScheduleEntry(&dbE))
	}

	return result, nil
}

func (r *scheduleEntryRepository) GetByCategoryID(
	ctx context.Context,
	categoryID int64,
) ([]domain.ScheduleEntry, error) {
	var dbEntries []dbScheduleEntry

	err := r.db.WithContext(ctx).
		Where("category_id = ?", categoryID).
		Preload("Category").
		Preload("Location").
		Order(orderStartTimeASC).
		Order(orderEndTimeASC).
		Find(&dbEntries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule entries by category ID: %w", err)
	}

	var result []domain.ScheduleEntry
	for _, dbE := range dbEntries {
		result = append(result, *toDomainScheduleEntry(&dbE))
	}

	return result, nil
}

func (r *scheduleEntryRepository) GetByLocationID(
	ctx context.Context,
	locationID int64,
) ([]domain.ScheduleEntry, error) {
	var dbEntries []dbScheduleEntry

	err := r.db.WithContext(ctx).
		Where("location_id = ?", locationID).
		Preload("Category").
		Preload("Location").
		Order(orderStartTimeASC).
		Order(orderEndTimeASC).
		Find(&dbEntries).Error
	if err != nil {
		return nil, fmt.Errorf("failed to fetch schedule entries by location ID: %w", err)
	}

	var result []domain.ScheduleEntry
	for _, dbE := range dbEntries {
		result = append(result, *toDomainScheduleEntry(&dbE))
	}

	return result, nil
}

func (r *scheduleEntryRepository) applyFilters(db *gorm.DB, f repository.ScheduleEntryListFilters) *gorm.DB {
	if f.Query != nil {
		likeQuery := fmt.Sprintf("%%%s%%", *f.Query)
		db = db.Where("title LIKE ? OR description LIKE ?", likeQuery, likeQuery)
	}

	if f.CategoryID != nil {
		db = db.Where("category_id = ?", *f.CategoryID)
	}

	if f.LocationID != nil {
		db = db.Where("location_id = ?", *f.LocationID)
	}

	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}

	if f.HidePast {
		now := time.Now().UTC()
		db = db.Where("end_time > ?", now)
	}

	return db
}
