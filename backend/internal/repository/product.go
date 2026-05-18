package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(p *model.Product) error {
	return r.db.Create(p).Error
}

func (r *ProductRepository) FindAll(offset, limit int) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64
	r.db.Model(&model.Product{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&products).Error
	return products, total, err
}

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	return &product, err
}

func (r *ProductRepository) Update(p *model.Product) error {
	return r.db.Save(p).Error
}

func (r *ProductRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}
