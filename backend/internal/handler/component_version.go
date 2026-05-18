package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/service"
)

type ComponentVersionHandler struct {
	svc *service.ComponentVersionService
}

func NewComponentVersionHandler(svc *service.ComponentVersionService) *ComponentVersionHandler {
	return &ComponentVersionHandler{svc: svc}
}

func (h *ComponentVersionHandler) Create(c *gin.Context) {
	var req dto.CreateComponentVersionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	cv, err := h.svc.Create(req.ComponentID, req.ArtifactID, req.Version, req.CreatedBy)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, cv)
}

func (h *ComponentVersionHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.svc.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *ComponentVersionHandler) Get(c *gin.Context) {
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
