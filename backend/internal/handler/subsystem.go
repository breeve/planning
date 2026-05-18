package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type SubsystemHandler struct {
	repo *repository.SubsystemRepository
}

func NewSubsystemHandler(repo *repository.SubsystemRepository) *SubsystemHandler {
	return &SubsystemHandler{repo: repo}
}

func (h *SubsystemHandler) Create(c *gin.Context) {
	var req dto.CreateSubsystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	item := &model.Subsystem{Name: req.Name, Description: req.Description, ProductID: req.ProductID}
	if err := h.repo.Create(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, item)
}

func (h *SubsystemHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.repo.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *SubsystemHandler) Get(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	item, err := h.repo.FindByID(id)
	if err != nil {
		dto.Err(c, 404, "not found")
		return
	}
	dto.OK(c, item)
}

func (h *SubsystemHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	item, err := h.repo.FindByID(id)
	if err != nil {
		dto.Err(c, 404, "not found")
		return
	}
	var req dto.UpdateSubsystemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	if req.Name != "" {
		item.Name = req.Name
	}
	if req.Description != "" {
		item.Description = req.Description
	}
	if err := h.repo.Update(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, item)
}

func (h *SubsystemHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	if err := h.repo.Delete(id); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}
