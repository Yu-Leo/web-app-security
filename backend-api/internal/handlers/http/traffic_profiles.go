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

func (h *Handler) ListTrafficProfiles(c *gin.Context) {
	items, err := h.trafficProfiles.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list traffic profiles", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list traffic profiles"})
		return
	}

	response := make([]service.TrafficProfile, 0, len(items))
	for _, item := range items {
		response = append(response, trafficProfileToDTO(item))
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateTrafficProfile(c *gin.Context) {
	var req service.CreateTrafficProfileJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.trafficProfiles.Create(c.Request.Context(), models.TrafficProfileToCreate{
		Name:        req.Name,
		Description: req.Description,
		IsEnabled:   req.IsEnabled,
	})
	if err != nil {
		h.logError(c, "failed to create traffic profile", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create traffic profile"})
		return
	}

	c.JSON(http.StatusCreated, trafficProfileToDTO(profile))
}

func (h *Handler) GetTrafficProfile(c *gin.Context, id int64) {
	profile, err := h.trafficProfiles.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traffic profile not found"})
			return
		}
		h.logError(c, "failed to get traffic profile", err, zap.Int64("profile_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get traffic profile"})
		return
	}

	c.JSON(http.StatusOK, trafficProfileToDTO(profile))
}

func (h *Handler) UpdateTrafficProfile(c *gin.Context, id int64) {
	var req service.UpdateTrafficProfileJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	profile, err := h.trafficProfiles.Update(c.Request.Context(), models.TrafficProfileToUpdate{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		IsEnabled:   req.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traffic profile not found"})
			return
		}
		h.logError(c, "failed to update traffic profile", err, zap.Int64("profile_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update traffic profile"})
		return
	}

	c.JSON(http.StatusOK, trafficProfileToDTO(profile))
}

func (h *Handler) DeleteTrafficProfile(c *gin.Context, id int64) {
	if err := h.trafficProfiles.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traffic profile not found"})
			return
		}
		h.logError(c, "failed to delete traffic profile", err, zap.Int64("profile_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete traffic profile"})
		return
	}

	c.Status(http.StatusNoContent)
}

func trafficProfileToDTO(profile models.TrafficProfile) service.TrafficProfile {
	return service.TrafficProfile{
		Id:          profile.ID,
		Name:        profile.Name,
		Description: profile.Description,
		IsEnabled:   profile.IsEnabled,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}
