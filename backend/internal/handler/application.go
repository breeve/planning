package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type ApplicationHandler struct {
	repo *repository.ApplicationRepository
}

func NewApplicationHandler(repo *repository.ApplicationRepository) *ApplicationHandler {
	return &ApplicationHandler{repo: repo}
}

func (h *ApplicationHandler) Create(c *gin.Context) {
	var req dto.CreateApplicationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	item := &model.Application{Name: req.Name, Description: req.Description}
	if err := h.repo.Create(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, item)
}

func (h *ApplicationHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.repo.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *ApplicationHandler) Get(c *gin.Context) {
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

func (h *ApplicationHandler) Update(c *gin.Context) {
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
	var req dto.UpdateApplicationRequest
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

func (h *ApplicationHandler) Delete(c *gin.Context) {
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

func (h *ApplicationHandler) AddSubsystem(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	var req dto.AssociationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	if err := h.repo.AddSubsystem(id, req.ID); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}

func (h *ApplicationHandler) RemoveSubsystem(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	sid, err := strconv.ParseUint(c.Param("sid"), 10, 32)
	if err != nil {
		dto.Err(c, 400, "invalid subsystem id")
		return
	}
	if err := h.repo.RemoveSubsystem(id, uint(sid)); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}
