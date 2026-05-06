package repository

import (
	"context"

	"github.com/potibm/billedapparat/internal/app/domain"
)

type CategoryListFilters struct {
	Query *string
}

type CategoryListParams struct {
	Offset int
	Limit  int
	Sort   string
	Order  string
}

type CategoryRepository interface {
	Create(ctx context.Context, category *domain.Category) error
	Update(ctx context.Context, category *domain.Category) error
	Delete(ctx context.Context, id int64) error
	List(
		ctx context.Context,
		params CategoryListParams,
		filters CategoryListFilters,
	) ([]domain.Category, int64, error)
	GetByID(ctx context.Context, id int64) (*domain.Category, error)
}
