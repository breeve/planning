package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type ComponentRepository struct {
	db *gorm.DB
}

func NewComponentRepository(db *gorm.DB) *ComponentRepository {
	return &ComponentRepository{db: db}
}

func (r *ComponentRepository) Create(c *model.Component) error {
	return r.db.Create(c).Error
}

func (r *ComponentRepository) FindAll(offset, limit int) ([]model.Component, int64, error) {
	var items []model.Component
	var total int64
	r.db.Model(&model.Component{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ComponentRepository) FindByID(id uint) (*model.Component, error) {
	var item model.Component
	err := r.db.Preload("Applications").First(&item, id).Error
	return &item, err
}

func (r *ComponentRepository) Update(c *model.Component) error {
	return r.db.Save(c).Error
}

func (r *ComponentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Component{}, id).Error
}

func (r *ComponentRepository) AddApplication(componentID, appID uint) error {
	return r.db.Exec(
		"INSERT OR IGNORE INTO component_applications (component_id, application_id) VALUES (?, ?)",
		componentID, appID,
	).Error
}

func (r *ComponentRepository) RemoveApplication(componentID, appID uint) error {
	return r.db.Exec(
		"DELETE FROM component_applications WHERE component_id = ? AND application_id = ?",
		componentID, appID,
	).Error
}

func (r *ComponentRepository) FindByApplicationID(appID uint) ([]model.Component, error) {
	var items []model.Component
	err := r.db.
		Joins("JOIN component_applications ON component_applications.component_id = components.id").
		Where("component_applications.application_id = ?", appID).
		Find(&items).Error
	return items, err
}
