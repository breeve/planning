package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/flynnzhang/planning/backend/internal/database"
	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/handler"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := database.NewTestDB()
	require.NoError(t, err)
	return handler.SetupRouter(db)
}

func TestProductAPI_CRUD(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateProductRequest{Name: "TestProd", Description: "desc"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	var resp dto.Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.True(t, resp.Success)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/products", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/products/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(dto.UpdateProductRequest{Name: "Updated"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/products/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/products/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestProductAPI_NotFound(t *testing.T) {
	r := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products/999", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 404, w.Code)
}

func TestProductAPI_InvalidID(t *testing.T) {
	r := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products/abc", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestProductAPI_InvalidBody(t *testing.T) {
	r := setupTestRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBufferString("{}"))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestSubsystemAPI_CRUD(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateProductRequest{Name: "P1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	body, _ = json.Marshal(dto.CreateSubsystemRequest{Name: "Sub1", ProductID: 1})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/subsystems", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/subsystems", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/subsystems/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(dto.UpdateSubsystemRequest{Name: "SubUpdated"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/subsystems/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/subsystems/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestApplicationAPI_CRUD(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateApplicationRequest{Name: "App1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/applications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/applications", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/applications/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(dto.UpdateApplicationRequest{Name: "AppUpdated"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/applications/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/applications/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestApplicationAPI_SubsystemAssociation(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateProductRequest{Name: "P1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateSubsystemRequest{Name: "Sub1", ProductID: 1})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/subsystems", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateApplicationRequest{Name: "App1"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/applications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.AssociationRequest{ID: 1})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/applications/1/subsystems", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/applications/1/subsystems/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestComponentAPI_CRUD(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateComponentRequest{Name: "Comp1", Type: "helm", RepoName: "repo", RepoBranch: "main"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/components", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/components", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/components/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(dto.UpdateComponentRequest{Name: "CompUpdated", RepoBranch: "dev"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/components/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/components/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestComponentAPI_ApplicationAssociation(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateApplicationRequest{Name: "App1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/applications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateComponentRequest{Name: "Comp1", Type: "helm"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/components", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.AssociationRequest{ID: 1})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/components/1/applications", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/components/1/applications/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestArtifactAPI_CRUD(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateComponentRequest{Name: "Comp1", Type: "helm"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/components", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateArtifactRequest{
		Name: "art1", ComponentID: 1, Version: "1.0.0", Registry: "reg/comp:1.0.0",
		RepoName: "repo", RepoBranch: "main", RepoCommit: "abc",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/artifacts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/artifacts", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/artifacts/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/artifacts/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestProductVersionAPI_FullFlow(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateProductRequest{Name: "Prod1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	body, _ = json.Marshal(dto.CreateProductVersionRequest{
		ProductID: 1, Version: "1.0.0", CreatedBy: "pm",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/product-versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/product-versions", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/product-versions/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	for _, status := range []string{"testing", "tested", "released"} {
		body, _ = json.Marshal(dto.UpdateStatusRequest{Status: status})
		w = httptest.NewRecorder()
		req, _ = http.NewRequest("PUT", "/api/v1/product-versions/1/status", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code, "failed transition to %s", status)
	}
}

func TestProductVersionAPI_InvalidTransition(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateProductRequest{Name: "Prod1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateProductVersionRequest{ProductID: 1, Version: "1.0.0"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/product-versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.UpdateStatusRequest{Status: "released"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/product-versions/1/status", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 400, w.Code)
}

func TestDeliveryPlanAPI_FullCRUD(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateDeliveryPlanRequest{Name: "Enterprise", Description: "Enterprise Edition"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/delivery-plans", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/delivery-plans", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/delivery-plans/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	body, _ = json.Marshal(dto.UpdateDeliveryPlanRequest{Name: "Standard"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("PUT", "/api/v1/delivery-plans/1", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/delivery-plans/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestDeliveryPlanAPI_ProductVersionAssociation(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateProductRequest{Name: "Prod1"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateProductVersionRequest{ProductID: 1, Version: "1.0.0"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/product-versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.CreateDeliveryPlanRequest{Name: "Plan1"})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/delivery-plans", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	body, _ = json.Marshal(dto.AssociationRequest{ID: 1})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/delivery-plans/1/product-versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("DELETE", "/api/v1/delivery-plans/1/product-versions/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestComponentVersionAPI_Full(t *testing.T) {
	r := setupTestRouter(t)

	body, _ := json.Marshal(dto.CreateComponentRequest{Name: "auth-helm", Type: "helm"})
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/v1/components", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	body, _ = json.Marshal(dto.CreateArtifactRequest{
		Name: "auth-v1", ComponentID: 1, Version: "1.0.0", Registry: "reg/auth:1.0.0",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/artifacts", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	require.Equal(t, 201, w.Code)

	body, _ = json.Marshal(dto.CreateComponentVersionRequest{
		ComponentID: 1, ArtifactID: 1, Version: "1.0.0", CreatedBy: "dev",
	})
	w = httptest.NewRecorder()
	req, _ = http.NewRequest("POST", "/api/v1/component-versions", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)
	assert.Equal(t, 201, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/component-versions", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	w = httptest.NewRecorder()
	req, _ = http.NewRequest("GET", "/api/v1/component-versions/1", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
}

func TestPagination(t *testing.T) {
	r := setupTestRouter(t)

	for i := 0; i < 5; i++ {
		body, _ := json.Marshal(dto.CreateProductRequest{Name: fmt.Sprintf("Prod%d", i)})
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/products", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/v1/products?page=1&limit=2", nil)
	r.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)

	var resp dto.Response
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	assert.Equal(t, int64(5), int64(resp.Meta.Total))
	assert.Equal(t, 1, resp.Meta.Page)
	assert.Equal(t, 2, resp.Meta.Limit)
}

func TestHandlerErrors_InvalidIDs(t *testing.T) {
	r := setupTestRouter(t)

	endpoints := []struct {
		method string
		path   string
	}{
		{"GET", "/api/v1/subsystems/abc"},
		{"PUT", "/api/v1/subsystems/abc"},
		{"DELETE", "/api/v1/subsystems/abc"},
		{"GET", "/api/v1/applications/abc"},
		{"PUT", "/api/v1/applications/abc"},
		{"DELETE", "/api/v1/applications/abc"},
		{"GET", "/api/v1/components/abc"},
		{"PUT", "/api/v1/components/abc"},
		{"DELETE", "/api/v1/components/abc"},
		{"GET", "/api/v1/artifacts/abc"},
		{"DELETE", "/api/v1/artifacts/abc"},
		{"GET", "/api/v1/component-versions/abc"},
		{"GET", "/api/v1/product-versions/abc"},
		{"PUT", "/api/v1/product-versions/abc/status"},
		{"GET", "/api/v1/delivery-plans/abc"},
		{"PUT", "/api/v1/delivery-plans/abc"},
		{"DELETE", "/api/v1/delivery-plans/abc"},
	}

	for _, ep := range endpoints {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest(ep.method, ep.path, nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, 400, w.Code, "%s %s should return 400", ep.method, ep.path)
	}
}
