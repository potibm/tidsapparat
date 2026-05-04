package gorm

import (
	"context"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

	query := r.db.WithContext(ctx).Model(&dbScheduleEntry{})

	query = r.applyFilters(query, filters)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order(clause.OrderByColumn{Column: clause.Column{Name: params.Sort}, Desc: params.Order == "desc"}).
		Offset(params.Offset).
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

	err := r.db.WithContext(ctx).First(&dbEntry, id).Error
	if err != nil {
		return nil, err
	}

	entry := toDomainScheduleEntry(&dbEntry)

	return entry, nil
}

func (r *scheduleEntryRepository) applyFilters(db *gorm.DB, f repository.ScheduleEntryListFilters) *gorm.DB {
	if f.Query != nil {
		db = db.Where("title ILIKE ?", "%"+*f.Query+"%")
	}

	if f.Category != nil {
		db = db.Where("category = ?", *f.Category)
	}

	if f.ID != nil {
		db = db.Where("id = ?", *f.ID)
	}

	return db
}
