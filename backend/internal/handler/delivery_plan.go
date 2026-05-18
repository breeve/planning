package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type DeliveryPlanHandler struct {
	repo *repository.DeliveryPlanRepository
}

func NewDeliveryPlanHandler(repo *repository.DeliveryPlanRepository) *DeliveryPlanHandler {
	return &DeliveryPlanHandler{repo: repo}
}

func (h *DeliveryPlanHandler) Create(c *gin.Context) {
	var req dto.CreateDeliveryPlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	item := &model.DeliveryPlan{Name: req.Name, Description: req.Description}
	if err := h.repo.Create(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, item)
}

func (h *DeliveryPlanHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.repo.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *DeliveryPlanHandler) Get(c *gin.Context) {
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

func (h *DeliveryPlanHandler) Update(c *gin.Context) {
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
	var req dto.UpdateDeliveryPlanRequest
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

func (h *DeliveryPlanHandler) Delete(c *gin.Context) {
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

func (h *DeliveryPlanHandler) AddProductVersion(c *gin.Context) {
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
	if err := h.repo.AddProductVersion(id, req.ID); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}

func (h *DeliveryPlanHandler) RemoveProductVersion(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	pvid, err := strconv.ParseUint(c.Param("pvid"), 10, 32)
	if err != nil {
		dto.Err(c, 400, "invalid product version id")
		return
	}
	if err := h.repo.RemoveProductVersion(id, uint(pvid)); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}
