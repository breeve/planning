package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type ProductHandler struct {
	repo *repository.ProductRepository
}

func NewProductHandler(repo *repository.ProductRepository) *ProductHandler {
	return &ProductHandler{repo: repo}
}

func (h *ProductHandler) Create(c *gin.Context) {
	var req dto.CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	product := &model.Product{Name: req.Name, Description: req.Description}
	if err := h.repo.Create(product); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, product)
}

func (h *ProductHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.repo.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *ProductHandler) Get(c *gin.Context) {
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

func (h *ProductHandler) Update(c *gin.Context) {
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
	var req dto.UpdateProductRequest
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

func (h *ProductHandler) Delete(c *gin.Context) {
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

func parseID(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	return uint(id), err
}

func parsePagination(c *gin.Context) (int, int) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return page, limit
}
