package repository

import (
	"github.com/flynnzhang/planning/backend/internal/model"
	"gorm.io/gorm"
)

type ApplicationRepository struct {
	db *gorm.DB
}

func NewApplicationRepository(db *gorm.DB) *ApplicationRepository {
	return &ApplicationRepository{db: db}
}

func (r *ApplicationRepository) Create(a *model.Application) error {
	return r.db.Create(a).Error
}

func (r *ApplicationRepository) FindAll(offset, limit int) ([]model.Application, int64, error) {
	var items []model.Application
	var total int64
	r.db.Model(&model.Application{}).Count(&total)
	err := r.db.Offset(offset).Limit(limit).Find(&items).Error
	return items, total, err
}

func (r *ApplicationRepository) FindByID(id uint) (*model.Application, error) {
	var item model.Application
	err := r.db.Preload("Subsystems").First(&item, id).Error
	return &item, err
}

func (r *ApplicationRepository) Update(a *model.Application) error {
	return r.db.Save(a).Error
}

func (r *ApplicationRepository) Delete(id uint) error {
	return r.db.Delete(&model.Application{}, id).Error
}

func (r *ApplicationRepository) AddSubsystem(appID, subsystemID uint) error {
	return r.db.Exec(
		"INSERT OR IGNORE INTO application_subsystems (application_id, subsystem_id) VALUES (?, ?)",
		appID, subsystemID,
	).Error
}

func (r *ApplicationRepository) RemoveSubsystem(appID, subsystemID uint) error {
	return r.db.Exec(
		"DELETE FROM application_subsystems WHERE application_id = ? AND subsystem_id = ?",
		appID, subsystemID,
	).Error
}

func (r *ApplicationRepository) FindBySubsystemID(subsystemID uint) ([]model.Application, error) {
	var items []model.Application
	err := r.db.
		Joins("JOIN application_subsystems ON application_subsystems.application_id = applications.id").
		Where("application_subsystems.subsystem_id = ?", subsystemID).
		Find(&items).Error
	return items, err
}
