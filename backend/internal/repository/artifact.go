package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type ArtifactRepository struct {
	db *gorm.DB
}

func NewArtifactRepository(db *gorm.DB) *ArtifactRepository {
	return &ArtifactRepository{db: db}
}

func (r *ArtifactRepository) Create(a *model.Artifact) error {
	return r.db.Create(a).Error
}

func (r *ArtifactRepository) FindAll(offset, limit int) ([]model.Artifact, int64, error) {
	var items []model.Artifact
	var total int64
	r.db.Model(&model.Artifact{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ArtifactRepository) FindByID(id uint) (*model.Artifact, error) {
	var item model.Artifact
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *ArtifactRepository) FindByComponentID(componentID uint) ([]model.Artifact, error) {
	var items []model.Artifact
	err := r.db.Where("component_id = ?", componentID).Find(&items).Error
	return items, err
}

func (r *ArtifactRepository) Delete(id uint) error {
	return r.db.Delete(&model.Artifact{}, id).Error
}
