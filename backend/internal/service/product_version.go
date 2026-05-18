package service

import (
	"encoding/json"
	"fmt"

	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

var validTransitions = map[string][]string{
	"draft":   {"testing"},
	"testing": {"tested"},
	"tested":  {"released"},
}

type ProductVersionService struct {
	repo          *repository.ProductVersionRepository
	productRepo   *repository.ProductRepository
	subsystemRepo *repository.SubsystemRepository
	appRepo       *repository.ApplicationRepository
	componentRepo *repository.ComponentRepository
}

func NewProductVersionService(
	repo *repository.ProductVersionRepository,
	productRepo *repository.ProductRepository,
	subsystemRepo *repository.SubsystemRepository,
	appRepo *repository.ApplicationRepository,
	componentRepo *repository.ComponentRepository,
) *ProductVersionService {
	return &ProductVersionService{
		repo:          repo,
		productRepo:   productRepo,
		subsystemRepo: subsystemRepo,
		appRepo:       appRepo,
		componentRepo: componentRepo,
	}
}

type treeSnapshotProduct struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Subsystems  []treeSnapshotSubsystem `json:"subsystems"`
}

type treeSnapshotSubsystem struct {
	ID           uint                      `json:"id"`
	Name         string                    `json:"name"`
	Description  string                    `json:"description"`
	Applications []treeSnapshotApplication `json:"applications"`
}

type treeSnapshotApplication struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Components  []treeSnapshotComponent `json:"components"`
}

type treeSnapshotComponent struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Type        string `json:"type"`
	RepoName    string `json:"repo_name"`
	RepoBranch  string `json:"repo_branch"`
}

func (s *ProductVersionService) Create(productID uint, version string, componentVersionIDs []uint, createdBy string) (*model.ProductVersion, error) {
	product, err := s.productRepo.FindByID(productID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %w", err)
	}

	tree, err := s.buildTreeSnapshot(product)
	if err != nil {
		return nil, fmt.Errorf("failed to build tree snapshot: %w", err)
	}

	treeJSON, err := json.Marshal(tree)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tree snapshot: %w", err)
	}

	pv := &model.ProductVersion{
		ProductID:    productID,
		Version:      version,
		Status:       "draft",
		TreeSnapshot: treeJSON,
		CreatedBy:    createdBy,
	}

	if err := s.repo.CreateWithComponentVersions(pv, componentVersionIDs); err != nil {
		return nil, err
	}
	return pv, nil
}

func (s *ProductVersionService) buildTreeSnapshot(product *model.Product) (*treeSnapshotProduct, error) {
	subsystems, err := s.subsystemRepo.FindByProductID(product.ID)
	if err != nil {
		return nil, err
	}

	tree := &treeSnapshotProduct{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
	}

	for _, sub := range subsystems {
		subSnap := treeSnapshotSubsystem{
			ID:          sub.ID,
			Name:        sub.Name,
			Description: sub.Description,
		}

		apps, err := s.appRepo.FindBySubsystemID(sub.ID)
		if err != nil {
			return nil, err
		}

		for _, app := range apps {
			appSnap := treeSnapshotApplication{
				ID:          app.ID,
				Name:        app.Name,
				Description: app.Description,
			}

			components, err := s.componentRepo.FindByApplicationID(app.ID)
			if err != nil {
				return nil, err
			}

			for _, comp := range components {
				appSnap.Components = append(appSnap.Components, treeSnapshotComponent{
					ID:         comp.ID,
					Name:       comp.Name,
					Type:       comp.Type,
					RepoName:   comp.RepoName,
					RepoBranch: comp.RepoBranch,
				})
			}
			subSnap.Applications = append(subSnap.Applications, appSnap)
		}
		tree.Subsystems = append(tree.Subsystems, subSnap)
	}
	return tree, nil
}

func (s *ProductVersionService) FindAll(offset, limit int) ([]model.ProductVersion, int64, error) {
	return s.repo.FindAll(offset, limit)
}

func (s *ProductVersionService) FindByID(id uint) (*model.ProductVersion, error) {
	return s.repo.FindByID(id)
}

func (s *ProductVersionService) UpdateStatus(id uint, newStatus string) error {
	pv, err := s.repo.FindByID(id)
	if err != nil {
		return fmt.Errorf("product version not found: %w", err)
	}

	allowed, exists := validTransitions[pv.Status]
	if !exists {
		return fmt.Errorf("no transitions available from status: %s", pv.Status)
	}

	valid := false
	for _, s := range allowed {
		if s == newStatus {
			valid = true
			break
		}
	}
	if !valid {
		return fmt.Errorf("invalid transition: %s → %s", pv.Status, newStatus)
	}

	return s.repo.UpdateStatus(id, newStatus)
}
