package service

import (
	"encoding/json"
	"fmt"

	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type ComponentVersionService struct {
	repo          *repository.ComponentVersionRepository
	componentRepo *repository.ComponentRepository
	artifactRepo  *repository.ArtifactRepository
}

func NewComponentVersionService(
	repo *repository.ComponentVersionRepository,
	componentRepo *repository.ComponentRepository,
	artifactRepo *repository.ArtifactRepository,
) *ComponentVersionService {
	return &ComponentVersionService{
		repo:          repo,
		componentRepo: componentRepo,
		artifactRepo:  artifactRepo,
	}
}

func (s *ComponentVersionService) Create(componentID, artifactID uint, version, createdBy string) (*model.ComponentVersion, error) {
	component, err := s.componentRepo.FindByID(componentID)
	if err != nil {
		return nil, fmt.Errorf("component not found: %w", err)
	}

	if _, err := s.artifactRepo.FindByID(artifactID); err != nil {
		return nil, fmt.Errorf("artifact not found: %w", err)
	}

	snapshot := map[string]interface{}{
		"name":        component.Name,
		"description": component.Description,
		"type":        component.Type,
		"repo_name":   component.RepoName,
		"repo_branch": component.RepoBranch,
	}
	snapshotJSON, err := json.Marshal(snapshot)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal snapshot: %w", err)
	}

	cv := &model.ComponentVersion{
		ComponentID: componentID,
		ArtifactID:  artifactID,
		Version:     version,
		Snapshot:    snapshotJSON,
		Status:      "built",
		CreatedBy:   createdBy,
	}
	if err := s.repo.Create(cv); err != nil {
		return nil, err
	}
	return cv, nil
}

func (s *ComponentVersionService) FindAll(offset, limit int) ([]model.ComponentVersion, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *ComponentVersionService) FindByID(id uint) (*model.ComponentVersion, error) {
	return s.repo.FindByID(id)
}
