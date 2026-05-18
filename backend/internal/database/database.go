package database

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	if err := db.AutoMigrate(
		&model.Product{},
		&model.Subsystem{},
		&model.Application{},
		&model.Component{},
		&model.Artifact{},
		&model.ComponentVersion{},
		&model.ProductVersion{},
		&model.DeliveryPlan{},
	); err != nil {
		return nil, err
	}
	return db, nil
}

func NewTestDB() (*gorm.DB, error) {
	return New(":memory:")
}
