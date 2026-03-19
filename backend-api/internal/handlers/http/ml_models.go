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

func (h *Handler) ListMLModels(c *gin.Context) {
	items, err := h.mlModels.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list ml models", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list ml models"})
		return
	}

	response := make([]service.MLModel, 0, len(items))
	for _, item := range items {
		dto, err := mlModelToDTO(item)
		if err != nil {
			h.logError(c, "failed to serialize ml model", err, zap.Int64("model_id", item.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize ml model"})
			return
		}
		response = append(response, dto)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateMLModel(c *gin.Context) {
	var req service.CreateMLModelJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	model, err := h.mlModels.Create(c.Request.Context(), models.MLModelToCreate{
		Name:      req.Name,
		ModelData: req.ModelData,
	})
	if err != nil {
		h.logError(c, "failed to create ml model", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create ml model"})
		return
	}

	response, err := mlModelToDTO(model)
	if err != nil {
		h.logError(c, "failed to serialize ml model", err, zap.Int64("model_id", model.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize ml model"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetMLModel(c *gin.Context, id int64) {
	model, err := h.mlModels.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ml model not found"})
			return
		}
		h.logError(c, "failed to get ml model", err, zap.Int64("model_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get ml model"})
		return
	}

	response, err := mlModelToDTO(model)
	if err != nil {
		h.logError(c, "failed to serialize ml model", err, zap.Int64("model_id", model.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize ml model"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateMLModel(c *gin.Context, id int64) {
	var req service.UpdateMLModelJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	model, err := h.mlModels.Update(c.Request.Context(), models.MLModelToUpdate{
		ID:        id,
		Name:      req.Name,
		ModelData: req.ModelData,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ml model not found"})
			return
		}
		h.logError(c, "failed to update ml model", err, zap.Int64("model_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update ml model"})
		return
	}

	response, err := mlModelToDTO(model)
	if err != nil {
		h.logError(c, "failed to serialize ml model", err, zap.Int64("model_id", model.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize ml model"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteMLModel(c *gin.Context, id int64) {
	if err := h.mlModels.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "ml model not found"})
			return
		}
		h.logError(c, "failed to delete ml model", err, zap.Int64("model_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete ml model"})
		return
	}

	c.Status(http.StatusNoContent)
}

func mlModelToDTO(model models.MLModel) (service.MLModel, error) {
	return service.MLModel{
		Id:        model.ID,
		Name:      model.Name,
		ModelData: model.ModelData,
	}, nil
}
