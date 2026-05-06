package gorm

import (
	"time"

	"github.com/potibm/billedapparat/internal/app/domain"
	"gorm.io/gorm"
)

type dbLocation struct {
	GormModel

	Name    string
	Address *string
}

func (dbLocation) TableName() string {
	return "locations"
}

func fromDomainLocation(l *domain.Location) *dbLocation {
	dbObj := &dbLocation{
		GormModel: GormModel{
			ID:        l.ID,
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		},
		Name:    l.Name,
		Address: l.Address,
	}

	if l.DeletedAt != nil {
		dbObj.DeletedAt = gorm.DeletedAt{
			Time:  *l.DeletedAt,
			Valid: true,
		}
	}

	return dbObj
}

func toDomainLocation(db *dbLocation) *domain.Location {
	var deletedAt *time.Time
	if db.DeletedAt.Valid {
		deletedAt = &db.DeletedAt.Time
	}

	return &domain.Location{
		ID:        db.ID,
		CreatedAt: db.CreatedAt,
		UpdatedAt: db.UpdatedAt,
		DeletedAt: deletedAt,
		Name:      db.Name,
		Address:   db.Address,
	}
}
