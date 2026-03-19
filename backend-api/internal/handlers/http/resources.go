package http

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/generated/service"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ListResources(c *gin.Context) {
	items, err := h.resources.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list resources", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list resources"})
		return
	}

	response := make([]service.Resource, 0, len(items))
	for _, item := range items {
		response = append(response, resourceToDTO(item))
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateResource(c *gin.Context) {
	var req service.CreateResourceJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resource, err := h.resources.Create(c.Request.Context(), models.ResourceToCreate{
		Name:              req.Name,
		URLPattern:        req.UrlPattern,
		SecurityProfileID: req.SecurityProfileId,
		TrafficProfileID:  req.TrafficProfileId,
	})
	if err != nil {
		h.logError(c, "failed to create resource", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create resource"})
		return
	}

	c.JSON(http.StatusCreated, resourceToDTO(resource))
}

func (h *Handler) GetResource(c *gin.Context, id int64) {
	resource, err := h.resources.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		h.logError(c, "failed to get resource", err, zap.Int64("resource_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get resource"})
		return
	}

	c.JSON(http.StatusOK, resourceToDTO(resource))
}

func (h *Handler) UpdateResource(c *gin.Context, id int64) {
	var req service.UpdateResourceJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	resource, err := h.resources.Update(c.Request.Context(), models.ResourceToUpdate{
		ID:                id,
		Name:              req.Name,
		URLPattern:        req.UrlPattern,
		SecurityProfileID: req.SecurityProfileId,
		TrafficProfileID:  req.TrafficProfileId,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		h.logError(c, "failed to update resource", err, zap.Int64("resource_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update resource"})
		return
	}

	c.JSON(http.StatusOK, resourceToDTO(resource))
}

func (h *Handler) DeleteResource(c *gin.Context, id int64) {
	if err := h.resources.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		h.logError(c, "failed to delete resource", err, zap.Int64("resource_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete resource"})
		return
	}

	c.Status(http.StatusNoContent)
}

func resourceToDTO(resource models.Resource) service.Resource {
	return service.Resource{
		Id:                resource.ID,
		Name:              resource.Name,
		UrlPattern:        resource.URLPattern,
		SecurityProfileId: resource.SecurityProfileID,
		TrafficProfileId:  resource.TrafficProfileID,
		CreatedAt:         resource.CreatedAt,
		UpdatedAt:         resource.UpdatedAt,
	}
}
