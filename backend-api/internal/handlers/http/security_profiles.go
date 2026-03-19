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

func (h *Handler) ListSecurityProfiles(c *gin.Context) {
	items, err := h.securityProfiles.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list security profiles", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list security profiles"})
		return
	}

	response := make([]service.SecurityProfile, 0, len(items))
	for _, item := range items {
		response = append(response, securityProfileToDTO(item))
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateSecurityProfile(c *gin.Context) {
	var req service.CreateSecurityProfileJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if !isSupportedSecurityProfileAction(string(req.BaseAction)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported base_action"})
		return
	}

	profile, err := h.securityProfiles.Create(c.Request.Context(), models.SecurityProfileToCreate{
		Name:        req.Name,
		Description: req.Description,
		BaseAction:  models.SecurityProfileBaseAction(req.BaseAction),
		LogEnabled:  req.LogEnabled,
		IsEnabled:   req.IsEnabled,
	})
	if err != nil {
		h.logError(c, "failed to create security profile", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create security profile"})
		return
	}

	c.JSON(http.StatusCreated, securityProfileToDTO(profile))
}

func (h *Handler) GetSecurityProfile(c *gin.Context, id int64) {
	profile, err := h.securityProfiles.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "security profile not found"})
			return
		}
		h.logError(c, "failed to get security profile", err, zap.Int64("profile_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get security profile"})
		return
	}

	c.JSON(http.StatusOK, securityProfileToDTO(profile))
}

func (h *Handler) UpdateSecurityProfile(c *gin.Context, id int64) {
	var req service.UpdateSecurityProfileJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if !isSupportedSecurityProfileAction(string(req.BaseAction)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported base_action"})
		return
	}

	profile, err := h.securityProfiles.Update(c.Request.Context(), models.SecurityProfileToUpdate{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		BaseAction:  models.SecurityProfileBaseAction(req.BaseAction),
		LogEnabled:  req.LogEnabled,
		IsEnabled:   req.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "security profile not found"})
			return
		}
		h.logError(c, "failed to update security profile", err, zap.Int64("profile_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update security profile"})
		return
	}

	c.JSON(http.StatusOK, securityProfileToDTO(profile))
}

func (h *Handler) DeleteSecurityProfile(c *gin.Context, id int64) {
	if err := h.securityProfiles.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "security profile not found"})
			return
		}
		h.logError(c, "failed to delete security profile", err, zap.Int64("profile_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete security profile"})
		return
	}

	c.Status(http.StatusNoContent)
}

func securityProfileToDTO(profile models.SecurityProfile) service.SecurityProfile {
	response := service.SecurityProfile{
		Id:         profile.ID,
		Name:       profile.Name,
		BaseAction: service.SecurityProfileBaseAction(profile.BaseAction),
		LogEnabled: profile.LogEnabled,
		IsEnabled:  profile.IsEnabled,
		CreatedAt:  profile.CreatedAt,
		UpdatedAt:  profile.UpdatedAt,
	}

	response.Description = profile.Description

	return response
}

func isSupportedSecurityProfileAction(value string) bool {
	return value == string(models.SecurityProfileBaseActionAllow) || value == string(models.SecurityProfileBaseActionBlock)
}
