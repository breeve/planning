package repository_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flynnzhang/planning/backend/internal/database"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

func TestProductRepository_CRUD(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	repo := repository.NewProductRepository(db)

	product := &model.Product{Name: "TestProduct", Description: "desc"}
	require.NoError(t, repo.Create(product))
	assert.NotZero(t, product.ID)

	found, err := repo.FindByID(product.ID)
	require.NoError(t, err)
	assert.Equal(t, "TestProduct", found.Name)

	items, total, err := repo.FindAll(0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)

	found.Name = "Updated"
	require.NoError(t, repo.Update(found))
	found2, _ := repo.FindByID(product.ID)
	assert.Equal(t, "Updated", found2.Name)

	require.NoError(t, repo.Delete(product.ID))
	_, err = repo.FindByID(product.ID)
	assert.Error(t, err)
}

func TestSubsystemRepository_CRUD(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	repo := repository.NewSubsystemRepository(db)

	product := &model.Product{Name: "P1"}
	require.NoError(t, productRepo.Create(product))

	sub := &model.Subsystem{Name: "Sub1", ProductID: product.ID}
	require.NoError(t, repo.Create(sub))
	assert.NotZero(t, sub.ID)

	subs, err := repo.FindByProductID(product.ID)
	require.NoError(t, err)
	assert.Len(t, subs, 1)
}

func TestApplicationRepository_ManyToMany(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	subsystemRepo := repository.NewSubsystemRepository(db)
	appRepo := repository.NewApplicationRepository(db)

	product := &model.Product{Name: "P1"}
	require.NoError(t, productRepo.Create(product))

	sub := &model.Subsystem{Name: "Sub1", ProductID: product.ID}
	require.NoError(t, subsystemRepo.Create(sub))

	app := &model.Application{Name: "App1"}
	require.NoError(t, appRepo.Create(app))

	require.NoError(t, appRepo.AddSubsystem(app.ID, sub.ID))

	found, err := appRepo.FindByID(app.ID)
	require.NoError(t, err)
	assert.Len(t, found.Subsystems, 1)
	assert.Equal(t, sub.ID, found.Subsystems[0].ID)

	apps, err := appRepo.FindBySubsystemID(sub.ID)
	require.NoError(t, err)
	assert.Len(t, apps, 1)

	require.NoError(t, appRepo.RemoveSubsystem(app.ID, sub.ID))
	found, _ = appRepo.FindByID(app.ID)
	assert.Len(t, found.Subsystems, 0)
}

func TestComponentRepository_ManyToMany(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	appRepo := repository.NewApplicationRepository(db)
	compRepo := repository.NewComponentRepository(db)

	app := &model.Application{Name: "App1"}
	require.NoError(t, appRepo.Create(app))

	comp := &model.Component{Name: "Comp1", Type: "helm"}
	require.NoError(t, compRepo.Create(comp))

	require.NoError(t, compRepo.AddApplication(comp.ID, app.ID))

	found, err := compRepo.FindByID(comp.ID)
	require.NoError(t, err)
	assert.Len(t, found.Applications, 1)

	comps, err := compRepo.FindByApplicationID(app.ID)
	require.NoError(t, err)
	assert.Len(t, comps, 1)

	require.NoError(t, compRepo.RemoveApplication(comp.ID, app.ID))
	found, _ = compRepo.FindByID(comp.ID)
	assert.Len(t, found.Applications, 0)
}

func TestArtifactRepository_CRUD(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	compRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)

	comp := &model.Component{Name: "C1", Type: "helm"}
	require.NoError(t, compRepo.Create(comp))

	art := &model.Artifact{Name: "A1", ComponentID: comp.ID, Version: "1.0.0", Registry: "reg/c1:1.0.0"}
	require.NoError(t, artifactRepo.Create(art))
	assert.NotZero(t, art.ID)

	found, err := artifactRepo.FindByID(art.ID)
	require.NoError(t, err)
	assert.Equal(t, "A1", found.Name)

	items, total, err := artifactRepo.FindAll(0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)

	byComp, err := artifactRepo.FindByComponentID(comp.ID)
	require.NoError(t, err)
	assert.Len(t, byComp, 1)

	require.NoError(t, artifactRepo.Delete(art.ID))
	_, err = artifactRepo.FindByID(art.ID)
	assert.Error(t, err)
}

func TestDeliveryPlanRepository_CRUD(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	dpRepo := repository.NewDeliveryPlanRepository(db)
	productRepo := repository.NewProductRepository(db)
	pvRepo := repository.NewProductVersionRepository(db)

	dp := &model.DeliveryPlan{Name: "Enterprise", Description: "desc"}
	require.NoError(t, dpRepo.Create(dp))
	assert.NotZero(t, dp.ID)

	found, err := dpRepo.FindByID(dp.ID)
	require.NoError(t, err)
	assert.Equal(t, "Enterprise", found.Name)

	items, total, err := dpRepo.FindAll(0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)

	found.Name = "Standard"
	require.NoError(t, dpRepo.Update(found))

	product := &model.Product{Name: "P1"}
	require.NoError(t, productRepo.Create(product))

	pv := &model.ProductVersion{ProductID: product.ID, Version: "1.0.0", Status: "released"}
	require.NoError(t, pvRepo.Create(pv))

	require.NoError(t, dpRepo.AddProductVersion(dp.ID, pv.ID))
	found, _ = dpRepo.FindByID(dp.ID)
	assert.Len(t, found.ProductVersions, 1)

	require.NoError(t, dpRepo.RemoveProductVersion(dp.ID, pv.ID))
	found, _ = dpRepo.FindByID(dp.ID)
	assert.Len(t, found.ProductVersions, 0)

	require.NoError(t, dpRepo.Delete(dp.ID))
	_, err = dpRepo.FindByID(dp.ID)
	assert.Error(t, err)
}

func TestComponentVersionRepository_CRUD(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	compRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	cvRepo := repository.NewComponentVersionRepository(db)

	comp := &model.Component{Name: "C1", Type: "helm"}
	require.NoError(t, compRepo.Create(comp))

	art := &model.Artifact{Name: "A1", ComponentID: comp.ID, Version: "1.0.0"}
	require.NoError(t, artifactRepo.Create(art))

	cv := &model.ComponentVersion{ComponentID: comp.ID, ArtifactID: art.ID, Version: "1.0.0", Status: "built"}
	require.NoError(t, cvRepo.Create(cv))
	assert.NotZero(t, cv.ID)

	found, err := cvRepo.FindByID(cv.ID)
	require.NoError(t, err)
	assert.Equal(t, "built", found.Status)

	items, total, err := cvRepo.FindAll(0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)

	byIDs, err := cvRepo.FindByIDs([]uint{cv.ID})
	require.NoError(t, err)
	assert.Len(t, byIDs, 1)

	require.NoError(t, cvRepo.UpdateStatus(cv.ID, "integrated"))
	found, _ = cvRepo.FindByID(cv.ID)
	assert.Equal(t, "integrated", found.Status)
}

func TestProductVersionRepository_CRUD(t *testing.T) {
	db, err := database.NewTestDB()
	require.NoError(t, err)

	productRepo := repository.NewProductRepository(db)
	pvRepo := repository.NewProductVersionRepository(db)
	compRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	cvRepo := repository.NewComponentVersionRepository(db)

	product := &model.Product{Name: "P1"}
	require.NoError(t, productRepo.Create(product))

	comp := &model.Component{Name: "C1", Type: "helm"}
	require.NoError(t, compRepo.Create(comp))
	art := &model.Artifact{Name: "A1", ComponentID: comp.ID, Version: "1.0.0"}
	require.NoError(t, artifactRepo.Create(art))
	cv := &model.ComponentVersion{ComponentID: comp.ID, ArtifactID: art.ID, Version: "1.0.0", Status: "built"}
	require.NoError(t, cvRepo.Create(cv))

	pv := &model.ProductVersion{ProductID: product.ID, Version: "1.0.0", Status: "draft"}
	require.NoError(t, pvRepo.CreateWithComponentVersions(pv, []uint{cv.ID}))
	assert.NotZero(t, pv.ID)

	found, err := pvRepo.FindByID(pv.ID)
	require.NoError(t, err)
	assert.Equal(t, "draft", found.Status)
	assert.Len(t, found.ComponentVersions, 1)

	items, total, err := pvRepo.FindAll(0, 10)
	require.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, items, 1)

	require.NoError(t, pvRepo.UpdateStatus(pv.ID, "testing"))
	found, _ = pvRepo.FindByID(pv.ID)
	assert.Equal(t, "testing", found.Status)
}
