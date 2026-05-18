package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type ProductVersionRepository struct {
	db *gorm.DB
}

func NewProductVersionRepository(db *gorm.DB) *ProductVersionRepository {
	return &ProductVersionRepository{db: db}
}

func (r *ProductVersionRepository) Create(pv *model.ProductVersion) error {
	return r.db.Create(pv).Error
}

func (r *ProductVersionRepository) CreateWithComponentVersions(pv *model.ProductVersion, cvIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(pv).Error; err != nil {
			return err
		}
		if len(cvIDs) > 0 {
			var cvs []model.ComponentVersion
			if err := tx.Where("id IN ?", cvIDs).Find(&cvs).Error; err != nil {
				return err
			}
			return tx.Model(pv).Association("ComponentVersions").Replace(cvs)
		}
		return nil
	})
}

func (r *ProductVersionRepository) FindAll(offset, limit int) ([]model.ProductVersion, int64, error) {
	var items []model.ProductVersion
	var total int64
	r.db.Model(&model.ProductVersion{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ProductVersionRepository) FindByID(id uint) (*model.ProductVersion, error) {
	var item model.ProductVersion
	err := r.db.Preload("ComponentVersions").Preload("Product").First(&item, id).Error
	return &item, err
}

func (r *ProductVersionRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&model.ProductVersion{}).Where("id = ?", id).Update("status", status).Error
}
