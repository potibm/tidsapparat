package repository

import (
	"context"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type LocationListFilters struct {
	Query *string
}

type LocationListParams struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

type LocationRepository interface {
	Create(ctx context.Context, location *domain.Location) error
	Update(ctx context.Context, location *domain.Location) error
	Delete(ctx context.Context, id int64) error
	List(
		ctx context.Context,
		params LocationListParams,
		filters LocationListFilters,
	) ([]domain.Location, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Location, error)
}
