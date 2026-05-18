package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type DeliveryPlanRepository struct {
	db *gorm.DB
}

func NewDeliveryPlanRepository(db *gorm.DB) *DeliveryPlanRepository {
	return &DeliveryPlanRepository{db: db}
}

func (r *DeliveryPlanRepository) Create(dp *model.DeliveryPlan) error {
	return r.db.Create(dp).Error
}

func (r *DeliveryPlanRepository) FindAll(offset, limit int) ([]model.DeliveryPlan, int64, error) {
	var items []model.DeliveryPlan
	var total int64
	r.db.Model(&model.DeliveryPlan{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *DeliveryPlanRepository) FindByID(id uint) (*model.DeliveryPlan, error) {
	var item model.DeliveryPlan
	err := r.db.Preload("ProductVersions").First(&item, id).Error
	return &item, err
}

func (r *DeliveryPlanRepository) Update(dp *model.DeliveryPlan) error {
	return r.db.Save(dp).Error
}

func (r *DeliveryPlanRepository) Delete(id uint) error {
	return r.db.Delete(&model.DeliveryPlan{}, id).Error
}

func (r *DeliveryPlanRepository) AddProductVersion(planID, pvID uint) error {
	plan := &model.DeliveryPlan{ID: planID}
	pv := &model.ProductVersion{ID: pvID}
	return r.db.Model(plan).Association("ProductVersions").Append(pv)
}

func (r *DeliveryPlanRepository) RemoveProductVersion(planID, pvID uint) error {
	plan := &model.DeliveryPlan{ID: planID}
	pv := &model.ProductVersion{ID: pvID}
	return r.db.Model(plan).Association("ProductVersions").Delete(pv)
}
