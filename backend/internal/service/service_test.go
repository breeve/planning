package service_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flynnzhang/planning/backend/internal/database"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
	"github.com/flynnzhang/planning/backend/internal/service"
)

func TestComponentVersionService_Create(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	compRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	cvRepo := repository.NewComponentVersionRepository(db)

	comp := &model.Component{Name: "auth-service", Type: "helm", RepoName: "auth", RepoBranch: "main"}
	require.NoError(t, compRepo.Create(comp))

	artifact := &model.Artifact{
		Name: "auth-v1", ComponentID: comp.ID, Version: "1.0.0",
		BuiltAt: time.Now(), Registry: "registry.example.com/auth:1.0.0",
		RepoName: "auth", RepoBranch: "main", RepoCommit: "abc123",
	}
	require.NoError(t, artifactRepo.Create(artifact))

	svc := service.NewComponentVersionService(cvRepo, compRepo, artifactRepo)

	cv, err := svc.Create(comp.ID, artifact.ID, "1.0.0", "admin")
	require.NoError(t, err)
	assert.NotZero(t, cv.ID)
	assert.Equal(t, "built", cv.Status)
	assert.NotEmpty(t, cv.Snapshot)
}

func TestComponentVersionService_Create_ComponentNotFound(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	compRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	cvRepo := repository.NewComponentVersionRepository(db)

	svc := service.NewComponentVersionService(cvRepo, compRepo, artifactRepo)

	_, err = svc.Create(999, 1, "1.0.0", "admin")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "component not found")
}

func TestProductVersionService_Create(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	subsystemRepo := repository.NewSubsystemRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	compRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	cvRepo := repository.NewComponentVersionRepository(db)
	pvRepo := repository.NewProductVersionRepository(db)

	product := &model.Product{Name: "WorkOrder"}
	require.NoError(t, productRepo.Create(product))

	sub := &model.Subsystem{Name: "Dispatch", ProductID: product.ID}
	require.NoError(t, subsystemRepo.Create(sub))

	app := &model.Application{Name: "dispatch-api"}
	require.NoError(t, appRepo.Create(app))
	require.NoError(t, appRepo.AddSubsystem(app.ID, sub.ID))

	comp := &model.Component{Name: "dispatch-helm", Type: "helm"}
	require.NoError(t, compRepo.Create(comp))
	require.NoError(t, compRepo.AddApplication(comp.ID, app.ID))

	artifact := &model.Artifact{
		Name: "dispatch-v1", ComponentID: comp.ID, Version: "1.0.0",
		BuiltAt: time.Now(), Registry: "reg/dispatch:1.0.0",
	}
	require.NoError(t, artifactRepo.Create(artifact))

	cvSvc := service.NewComponentVersionService(cvRepo, compRepo, artifactRepo)
	cv, err := cvSvc.Create(comp.ID, artifact.ID, "1.0.0", "dev")
	require.NoError(t, err)

	pvSvc := service.NewProductVersionService(pvRepo, productRepo, subsystemRepo, appRepo, compRepo)
	pv, err := pvSvc.Create(product.ID, "2.0.0", []uint{cv.ID}, "pm")
	require.NoError(t, err)
	assert.NotZero(t, pv.ID)
	assert.Equal(t, "draft", pv.Status)
	assert.NotEmpty(t, pv.TreeSnapshot)
}

func TestProductVersionService_StatusTransition(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	subsystemRepo := repository.NewSubsystemRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	compRepo := repository.NewComponentRepository(db)
	pvRepo := repository.NewProductVersionRepository(db)

	product := &model.Product{Name: "TestProd"}
	require.NoError(t, productRepo.Create(product))

	pvSvc := service.NewProductVersionService(pvRepo, productRepo, subsystemRepo, appRepo, compRepo)
	pv, err := pvSvc.Create(product.ID, "1.0.0", nil, "pm")
	require.NoError(t, err)

	require.NoError(t, pvSvc.UpdateStatus(pv.ID, "testing"))
	require.NoError(t, pvSvc.UpdateStatus(pv.ID, "tested"))
	require.NoError(t, pvSvc.UpdateStatus(pv.ID, "released"))

	err = pvSvc.UpdateStatus(pv.ID, "draft")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no transitions available")
}

func TestProductVersionService_InvalidTransition(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	subsystemRepo := repository.NewSubsystemRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	compRepo := repository.NewComponentRepository(db)
	pvRepo := repository.NewProductVersionRepository(db)

	product := &model.Product{Name: "TestProd"}
	require.NoError(t, productRepo.Create(product))

	pvSvc := service.NewProductVersionService(pvRepo, productRepo, subsystemRepo, appRepo, compRepo)
	pv, err := pvSvc.Create(product.ID, "1.0.0", nil, "pm")
	require.NoError(t, err)

	err = pvSvc.UpdateStatus(pv.ID, "released")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid transition: draft → released")
}
