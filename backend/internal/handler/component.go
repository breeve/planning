package handler

import (
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type ComponentHandler struct {
	repo *repository.ComponentRepository
}

func NewComponentHandler(repo *repository.ComponentRepository) *ComponentHandler {
	return &ComponentHandler{repo: repo}
}

func (h *ComponentHandler) Create(c *gin.Context) {
	var req dto.CreateComponentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	item := &model.Component{
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		RepoName:    req.RepoName,
		RepoBranch:  req.RepoBranch,
		RepoUser:    req.RepoUser,
		RepoPasswd:  req.RepoPasswd,
	}
	if item.Type == "" {
		item.Type = "helm"
	}
	if err := h.repo.Create(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, item)
}

func (h *ComponentHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.repo.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *ComponentHandler) Get(c *gin.Context) {
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

func (h *ComponentHandler) Update(c *gin.Context) {
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
	var req dto.UpdateComponentRequest
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
	if req.Type != "" {
		item.Type = req.Type
	}
	if req.RepoName != "" {
		item.RepoName = req.RepoName
	}
	if req.RepoBranch != "" {
		item.RepoBranch = req.RepoBranch
	}
	if req.RepoUser != "" {
		item.RepoUser = req.RepoUser
	}
	if req.RepoPasswd != "" {
		item.RepoPasswd = req.RepoPasswd
	}
	if err := h.repo.Update(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, item)
}

func (h *ComponentHandler) Delete(c *gin.Context) {
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

func (h *ComponentHandler) AddApplication(c *gin.Context) {
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
	if err := h.repo.AddApplication(id, req.ID); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}

func (h *ComponentHandler) RemoveApplication(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		dto.Err(c, 400, "invalid id")
		return
	}
	aid, err := strconv.ParseUint(c.Param("aid"), 10, 32)
	if err != nil {
		dto.Err(c, 400, "invalid application id")
		return
	}
	if err := h.repo.RemoveApplication(id, uint(aid)); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OK(c, nil)
}
