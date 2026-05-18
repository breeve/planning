package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type SubsystemRepository struct {
	db *gorm.DB
}

func NewSubsystemRepository(db *gorm.DB) *SubsystemRepository {
	return &SubsystemRepository{db: db}
}

func (r *SubsystemRepository) Create(s *model.Subsystem) error {
	return r.db.Create(s).Error
}

func (r *SubsystemRepository) FindAll(offset, limit int) ([]model.Subsystem, int64, error) {
	var items []model.Subsystem
	var total int64
	r.db.Model(&model.Subsystem{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *SubsystemRepository) FindByID(id uint) (*model.Subsystem, error) {
	var item model.Subsystem
	err := r.db.First(&item, id).Error
	return &item, err
}

func (r *SubsystemRepository) FindByProductID(productID uint) ([]model.Subsystem, error) {
	var items []model.Subsystem
	err := r.db.Where("product_id = ?", productID).Find(&items).Error
	return items, err
}

func (r *SubsystemRepository) Update(s *model.Subsystem) error {
	return r.db.Save(s).Error
}

func (r *SubsystemRepository) Delete(id uint) error {
	return r.db.Delete(&model.Subsystem{}, id).Error
}
