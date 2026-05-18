package model

import "time"

type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;not null" json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Subsystem struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	ProductID   uint      `gorm:"not null;index" json:"product_id"`
	Product     Product   `gorm:"constraint:OnDelete:CASCADE" json:"product,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Application struct {
	ID          uint        `gorm:"primaryKey" json:"id"`
	Name        string      `gorm:"not null" json:"name"`
	Description string      `json:"description"`
	Subsystems  []Subsystem `gorm:"many2many:application_subsystems" json:"subsystems,omitempty"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type Component struct {
	ID           uint          `gorm:"primaryKey" json:"id"`
	Name         string        `gorm:"not null" json:"name"`
	Description  string        `json:"description"`
	Type         string        `gorm:"not null;default:'helm'" json:"type"`
	RepoName     string        `json:"repo_name"`
	RepoBranch   string        `json:"repo_branch"`
	RepoUser     string        `json:"repo_user"`
	RepoPasswd   string        `json:"-"`
	Applications []Application `gorm:"many2many:component_applications" json:"applications,omitempty"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type Artifact struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`
	Description string    `json:"description"`
	ComponentID uint      `gorm:"not null;index" json:"component_id"`
	Component   Component `gorm:"constraint:OnDelete:CASCADE" json:"component,omitempty"`
	Version     string    `gorm:"not null" json:"version"`
	BuiltAt     time.Time `json:"built_at"`
	Registry    string    `json:"registry"`
	RepoName    string    `json:"repo_name"`
	RepoBranch  string    `json:"repo_branch"`
	RepoCommit  string    `json:"repo_commit"`
	CreatedAt   time.Time `json:"created_at"`
}
