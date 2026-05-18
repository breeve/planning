package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type ComponentVersionRepository struct {
	db *gorm.DB
}

func NewComponentVersionRepository(db *gorm.DB) *ComponentVersionRepository {
	return &ComponentVersionRepository{db: db}
}

func (r *ComponentVersionRepository) Create(cv *model.ComponentVersion) error {
	return r.db.Create(cv).Error
}

func (r *ComponentVersionRepository) FindAll(offset, limit int) ([]model.ComponentVersion, int64, error) {
	var items []model.ComponentVersion
	var total int64
	r.db.Model(&model.ComponentVersion{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Preload("Artifact").Find(&items).Error
	return items, total, err
}

func (r *ComponentVersionRepository) FindByID(id uint) (*model.ComponentVersion, error) {
	var item model.ComponentVersion
	err := r.db.Preload("Artifact").Preload("Component").First(&item, id).Error
	return &item, err
}

func (r *ComponentVersionRepository) FindByIDs(ids []uint) ([]model.ComponentVersion, error) {
	var items []model.ComponentVersion
	err := r.db.Where("id IN ?", ids).Find(&items).Error
	return items, err
}

func (r *ComponentVersionRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.ComponentVersion{}).Where("id = ?", id).Update("status", status).Error
}
