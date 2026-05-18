package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/flynnzhang/planning/backend/internal/dto"
	"github.com/flynnzhang/planning/backend/internal/model"
	"github.com/flynnzhang/planning/backend/internal/repository"
)

type ArtifactHandler struct {
	repo *repository.ArtifactRepository
}

func NewArtifactHandler(repo *repository.ArtifactRepository) *ArtifactHandler {
	return &ArtifactHandler{repo: repo}
}

func (h *ArtifactHandler) Create(c *gin.Context) {
	var req dto.CreateArtifactRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		dto.Err(c, 400, err.Error())
		return
	}
	item := &model.Artifact{
		Name:        req.Name,
		Description: req.Description,
		ComponentID: req.ComponentID,
		Version:     req.Version,
		BuiltAt:     req.BuiltAt,
		Registry:    req.Registry,
		RepoName:    req.RepoName,
		RepoBranch:  req.RepoBranch,
		RepoCommit:  req.RepoCommit,
	}
	if err := h.repo.Create(item); err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.Created(c, item)
}

func (h *ArtifactHandler) List(c *gin.Context) {
	page, limit := parsePagination(c)
	items, total, err := h.repo.FindAll((page-1)*limit, limit)
	if err != nil {
		dto.Err(c, 500, err.Error())
		return
	}
	dto.OKList(c, items, total, page, limit)
}

func (h *ArtifactHandler) Get(c *gin.Context) {
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

func (h *ArtifactHandler) Delete(c *gin.Context) {
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
