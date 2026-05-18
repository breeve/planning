package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/service"
)

type ProductVersionHandler struct {
	svc *service.ProductVersionService
}

func NewProductVersionHandler(svc *service.ProductVersionService) *ProductVersionHandler {
	return &ProductVersionHandler{svc: svc}
}

func (h *ProductVersionHandler) Create(c *gin.Context) {
	var req dto.CreateProductVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	pv, err := h.svc.Create(req.ProductID, req.Version, req.ComponentVersionIDs, req.CreatedBy)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, pv)
}

func (h *ProductVersionHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.svc.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *ProductVersionHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	item, err := h.svc.FindByID(id)
	if err != nil {
		dto.Err(c, 404, "not found")
		return
	}
	dto.OK(c, item)
}

func (h *ProductVersionHandler) UpdateStatus(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	var req dto.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	if err := h.svc.UpdateStatus(id, req.Status); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	dto.OK(c, nil)
}
