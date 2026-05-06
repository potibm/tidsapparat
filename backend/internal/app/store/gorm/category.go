package gorm

import (
	"github.com/potibm/billedapparat/internal/app/domain"
)

type dbCategory struct {
	GormModel

	Name  string `gorm:"not null"`
	Color string `gorm:"default:'#888888'"`
}

func (dbCategory) TableName() string {
	return "categories"
}

func fromDomainCategory(c *domain.Category) *dbCategory {
	dbObj := &dbCategory{
		GormModel: GormModel{
			ID: c.ID,
		},
		Name:  c.Name,
		Color: c.Color,
	}

	return dbObj
}

func toDomainCategory(db *dbCategory) *domain.Category {
	return &domain.Category{
		ID:    db.ID,
		Name:  db.Name,
		Color: db.Color,
	}
}
