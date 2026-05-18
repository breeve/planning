package model

import (
	"time"

	"gorm.io/datatypes"
)

type ComponentVersion struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ComponentID uint           `gorm:"not null;index" json:"component_id"`
	Component   Component      `gorm:"constraint:OnDelete:CASCADE" json:"component,omitempty"`
	Version     string         `gorm:"not null" json:"version"`
	ArtifactID  uint           `gorm:"not null" json:"artifact_id"`
	Artifact    Artifact       `gorm:"constraint:OnDelete:CASCADE" json:"artifact,omitempty"`
	Snapshot    datatypes.JSON `gorm:"type:json" json:"snapshot"`
	Status      string         `gorm:"not null;default:'built'" json:"status"`
	CreatedBy   string         `json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
}

type ProductVersion struct {
	ID                uint               `gorm:"primaryKey" json:"id"`
	ProductID         uint               `gorm:"not null;index" json:"product_id"`
	Product           Product            `gorm:"constraint:OnDelete:CASCADE" json:"product,omitempty"`
	Version           string             `gorm:"not null" json:"version"`
	Status            string             `gorm:"not null;default:'draft'" json:"status"`
	TreeSnapshot      datatypes.JSON     `gorm:"type:json" json:"tree_snapshot"`
	ComponentVersions []ComponentVersion `gorm:"many2many:product_version_components" json:"component_versions,omitempty"`
	CreatedBy         string             `json:"created_by"`
	CreatedAt         time.Time          `json:"created_at"`
}

type DeliveryPlan struct {
	ID              uint             `gorm:"primaryKey" json:"id"`
	Name            string           `gorm:"not null" json:"name"`
	Description     string           `json:"description"`
	ProductVersions []ProductVersion `gorm:"many2many:delivery_plan_product_versions" json:"product_versions,omitempty"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
}
