package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/repository"
	"github.com/flynnzhang/planning/backend/internal/service"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	r := gin.Default()

	productRepo := repository.NewProductRepository(db)
	subsystemRepo := repository.NewSubsystemRepository(db)
	appRepo := repository.NewApplicationRepository(db)
	componentRepo := repository.NewComponentRepository(db)
	artifactRepo := repository.NewArtifactRepository(db)
	cvRepo := repository.NewComponentVersionRepository(db)
	pvRepo := repository.NewProductVersionRepository(db)
	dpRepo := repository.NewDeliveryPlanRepository(db)

	cvSvc := service.NewComponentVersionService(cvRepo, componentRepo, artifactRepo)
	pvSvc := service.NewProductVersionService(pvRepo, productRepo, subsystemRepo, appRepo, componentRepo)

	productH := NewProductHandler(productRepo)
	subsystemH := NewSubsystemHandler(subsystemRepo)
	appH := NewApplicationHandler(appRepo)
	componentH := NewComponentHandler(componentRepo)
	artifactH := NewArtifactHandler(artifactRepo)
	cvH := NewComponentVersionHandler(cvSvc)
	pvH := NewProductVersionHandler(pvSvc)
	dpH := NewDeliveryPlanHandler(dpRepo)

	v1 := r.Group("/api/v1")
	{
		products := v1.Group("/products")
		{
			products.GET("", productH.List)
			products.POST("", productH.Create)
			products.GET("/:id", productH.Get)
			products.PUT("/:id", productH.Update)
			products.DELETE("/:id", productH.Delete)
		}

		subsystems := v1.Group("/subsystems")
		{
			subsystems.GET("", subsystemH.List)
			subsystems.POST("", subsystemH.Create)
			subsystems.GET("/:id", subsystemH.Get)
			subsystems.PUT("/:id", subsystemH.Update)
			subsystems.DELETE("/:id", subsystemH.Delete)
		}

		applications := v1.Group("/applications")
		{
			applications.GET("", appH.List)
			applications.POST("", appH.Create)
			applications.GET("/:id", appH.Get)
			applications.PUT("/:id", appH.Update)
			applications.DELETE("/:id", appH.Delete)
			applications.POST("/:id/subsystems", appH.AddSubsystem)
			applications.DELETE("/:id/subsystems/:sid", appH.RemoveSubsystem)
		}

		components := v1.Group("/components")
		{
			components.GET("", componentH.List)
			components.POST("", componentH.Create)
			components.GET("/:id", componentH.Get)
			components.PUT("/:id", componentH.Update)
			components.DELETE("/:id", componentH.Delete)
			components.POST("/:id/applications", componentH.AddApplication)
			components.DELETE("/:id/applications/:aid", componentH.RemoveApplication)
		}

		artifacts := v1.Group("/artifacts")
		{
			artifacts.GET("", artifactH.List)
			artifacts.POST("", artifactH.Create)
			artifacts.GET("/:id", artifactH.Get)
			artifacts.DELETE("/:id", artifactH.Delete)
		}

		componentVersions := v1.Group("/component-versions")
		{
			componentVersions.GET("", cvH.List)
			componentVersions.POST("", cvH.Create)
			componentVersions.GET("/:id", cvH.Get)
		}

		productVersions := v1.Group("/product-versions")
		{
			productVersions.GET("", pvH.List)
			productVersions.POST("", pvH.Create)
			productVersions.GET("/:id", pvH.Get)
			productVersions.PUT("/:id/status", pvH.UpdateStatus)
		}

		deliveryPlans := v1.Group("/delivery-plans")
		{
			deliveryPlans.GET("", dpH.List)
			deliveryPlans.POST("", dpH.Create)
			deliveryPlans.GET("/:id", dpH.Get)
			deliveryPlans.PUT("/:id", dpH.Update)
			deliveryPlans.DELETE("/:id", dpH.Delete)
			deliveryPlans.POST("/:id/product-versions", dpH.AddProductVersion)
			deliveryPlans.DELETE("/:id/product-versions/:pvid", dpH.RemoveProductVersion)
		}
	}

	return r
}
