package gorm

import (
	"context"
	"fmt"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type categoryRepository struct {
	db *gorm.DB
}

func (s *Store) NewCategoryRepository() repository.CategoryRepository {
	return NewCategoryRepository(s.db)
}

func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) Create(ctx context.Context, category *domain.Category) error {
	dbObj := fromDomainCategory(category)

	err := r.db.WithContext(ctx).Create(dbObj).Error
	if err == nil {
		category.ID = dbObj.ID
	}

	return err
}

func (r *categoryRepository) Update(ctx context.Context, category *domain.Category) error {
	dbObj := fromDomainCategory(category)

	return r.db.WithContext(ctx).Save(dbObj).Error
}

func (r *categoryRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&dbCategory{}, id).Error
}

func (r *categoryRepository) List(
	ctx context.Context,
	params repository.CategoryListParams,
	filters repository.CategoryListFilters,
) ([]domain.Category, int64, error) {
	var (
		dbCategories []dbCategory
		total        int64
	)

	query := r.db.WithContext(ctx).Model(&dbCategory{})

	query = r.applyFilters(query, filters)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	query = query.Order(clause.OrderByColumn{
		Column: clause.Column{Name: params.Sort},
		Desc:   params.Order == "DESC",
	})

	err = query.Offset(params.Offset).
		Limit(params.Limit).
		Find(&dbCategories).
		Error
	if err != nil {
		return nil, 0, err
	}

	categories := make([]domain.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		categories[i] = *toDomainCategory(&dbCategory)
	}

	return categories, total, nil
}

func (r *categoryRepository) GetByID(ctx context.Context, id int64) (*domain.Category, error) {
	var dbCategory dbCategory

	err := r.db.WithContext(ctx).First(&dbCategory, id).Error
	if err != nil {
		return nil, err
	}

	category := toDomainCategory(&dbCategory)

	return category, nil
}

func (r *categoryRepository) applyFilters(db *gorm.DB, f repository.CategoryListFilters) *gorm.DB {
	if f.Query != nil {
		likeQuery := fmt.Sprintf("%%%s%%", *f.Query)
		db = db.Where("name LIKE ?", likeQuery, likeQuery)
	}

	return db
}
