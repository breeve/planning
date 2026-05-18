package dto

import "time"

type CreateProductRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateSubsystemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	ProductID   uint   `json:"product_id" binding:"required"`
}

type UpdateSubsystemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateApplicationRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateApplicationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type CreateComponentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Type        string `json:"type"`
	RepoName    string `json:"repo_name"`
	RepoBranch  string `json:"repo_branch"`
	RepoUser    string `json:"repo_user"`
	RepoPasswd  string `json:"repo_passwd"`
}

type UpdateComponentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Type        string `json:"type"`
	RepoName    string `json:"repo_name"`
	RepoBranch  string `json:"repo_branch"`
	RepoUser    string `json:"repo_user"`
	RepoPasswd  string `json:"repo_passwd"`
}

type CreateArtifactRequest struct {
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	ComponentID uint      `json:"component_id" binding:"required"`
	Version     string    `json:"version" binding:"required"`
	BuiltAt     time.Time `json:"built_at"`
	Registry    string    `json:"registry"`
	RepoName    string    `json:"repo_name"`
	RepoBranch  string    `json:"repo_branch"`
	RepoCommit  string    `json:"repo_commit"`
}

type CreateComponentVersionRequest struct {
	ComponentID uint   `json:"component_id" binding:"required"`
	ArtifactID  uint   `json:"artifact_id" binding:"required"`
	Version     string `json:"version" binding:"required"`
	CreatedBy   string `json:"created_by"`
}

type CreateProductVersionRequest struct {
	ProductID           uint   `json:"product_id" binding:"required"`
	Version             string `json:"version" binding:"required"`
	ComponentVersionIDs []uint `json:"component_version_ids"`
	CreatedBy           string `json:"created_by"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateDeliveryPlanRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateDeliveryPlanRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type AssociationRequest struct {
	ID uint `json:"id" binding:"required"`
}
