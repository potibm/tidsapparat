package gorm

import (
	"context"
	"fmt"

	"github.com/potibm/billedapparat/internal/app/domain"
	"github.com/potibm/billedapparat/internal/app/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type locationRepository struct {
	db *gorm.DB
}

func (s *Store) NewLocationRepository() repository.LocationRepository {
	return NewLocationRepository(s.db)
}

func NewLocationRepository(db *gorm.DB) repository.LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(ctx context.Context, location *domain.Location) error {
	dbObj := fromDomainLocation(location)

	err := r.db.WithContext(ctx).Create(dbObj).Error
	if err == nil {
		location.ID = dbObj.ID
	}

	return err
}

func (r *locationRepository) Update(ctx context.Context, location *domain.Location) error {
	dbObj := fromDomainLocation(location)

	return r.db.WithContext(ctx).Save(dbObj).Error
}

func (r *locationRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&dbLocation{}, id).Error
}

func (r *locationRepository) List(
	ctx context.Context,
	params repository.LocationListParams,
	filters repository.LocationListFilters,
) ([]domain.Location, int64, error) {
	var (
		dbLocations []dbLocation
		total       int64
	)

	query := r.db.WithContext(ctx).Model(&dbLocation{})

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
		Find(&dbLocations).
		Error
	if err != nil {
		return nil, 0, err
	}

	locations := make([]domain.Location, len(dbLocations))
	for i, dbLocation := range dbLocations {
		locations[i] = *toDomainLocation(&dbLocation)
	}

	return locations, total, nil
}

func (r *locationRepository) GetByID(ctx context.Context, id int64) (*domain.Location, error) {
	var dbLocation dbLocation

	err := r.db.WithContext(ctx).First(&dbLocation, id).Error
	if err != nil {
		return nil, err
	}

	location := toDomainLocation(&dbLocation)

	return location, nil
}

func (r *locationRepository) applyFilters(db *gorm.DB, f repository.LocationListFilters) *gorm.DB {
	if f.Query != nil {
		likeQuery := fmt.Sprintf("%%%s%%", *f.Query)
		db = db.Where("name LIKE ?", likeQuery)
	}

	return db
}
