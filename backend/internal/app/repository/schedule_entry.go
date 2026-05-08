package repository

import (
	"context"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type ScheduleEntryListFilters struct {
	Query      *string
	CategoryID *int64
	LocationID *int64
	ID         *int64
	HidePast   bool
}

type ScheduleEntryListParams struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

type ScheduleEntryRepository interface {
	Save(ctx context.Context, scheduleEntry *domain.ScheduleEntry) error
	Delete(ctx context.Context, id int64) error
	List(
		ctx context.Context,
		params ScheduleEntryListParams,
		filters ScheduleEntryListFilters,
	) ([]domain.ScheduleEntry, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.ScheduleEntry, error)
	GetAllPreloaded(ctx context.Context) (domain.TimeTable, error)
	GetByCategoryID(ctx context.Context, categoryID int64) ([]domain.ScheduleEntry, error)
	GetByLocationID(ctx context.Context, locationID int64) ([]domain.ScheduleEntry, error)
}
