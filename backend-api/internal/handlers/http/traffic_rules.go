package http

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/generated/service"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ListTrafficRules(c *gin.Context) {
	items, err := h.trafficRules.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list traffic rules", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list traffic rules"})
		return
	}

	response := make([]service.TrafficRule, 0, len(items))
	for _, item := range items {
		dto, err := trafficRuleToDTO(item)
		if err != nil {
			h.logError(c, "failed to serialize traffic rule", err, zap.Int64("rule_id", item.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize traffic rule"})
			return
		}
		response = append(response, dto)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateTrafficRule(c *gin.Context) {
	var req service.CreateTrafficRuleJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateTrafficRuleRequest(req.MatchAll, req.Conditions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conditions, err := marshalTrafficConditions(req.Conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conditions"})
		return
	}

	rule, err := h.trafficRules.Create(c.Request.Context(), models.TrafficRuleToCreate{
		ProfileID:     req.ProfileId,
		Name:          req.Name,
		Description:   req.Description,
		Priority:      req.Priority,
		DryRun:        req.DryRun,
		MatchAll:      req.MatchAll,
		RequestsLimit: req.RequestsLimit,
		PeriodSeconds: req.PeriodSeconds,
		Conditions:    conditions,
		IsEnabled:     req.IsEnabled,
	})
	if err != nil {
		h.logError(c, "failed to create traffic rule", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create traffic rule"})
		return
	}

	response, err := trafficRuleToDTO(rule)
	if err != nil {
		h.logError(c, "failed to serialize traffic rule", err, zap.Int64("rule_id", rule.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize traffic rule"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetTrafficRule(c *gin.Context, id int64) {
	rule, err := h.trafficRules.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traffic rule not found"})
			return
		}
		h.logError(c, "failed to get traffic rule", err, zap.Int64("rule_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get traffic rule"})
		return
	}

	response, err := trafficRuleToDTO(rule)
	if err != nil {
		h.logError(c, "failed to serialize traffic rule", err, zap.Int64("rule_id", rule.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize traffic rule"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateTrafficRule(c *gin.Context, id int64) {
	var req service.UpdateTrafficRuleJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if err := validateTrafficRuleRequest(req.MatchAll, req.Conditions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conditions, err := marshalTrafficConditions(req.Conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conditions"})
		return
	}

	rule, err := h.trafficRules.Update(c.Request.Context(), models.TrafficRuleToUpdate{
		ID:            id,
		ProfileID:     req.ProfileId,
		Name:          req.Name,
		Description:   req.Description,
		Priority:      req.Priority,
		DryRun:        req.DryRun,
		MatchAll:      req.MatchAll,
		RequestsLimit: req.RequestsLimit,
		PeriodSeconds: req.PeriodSeconds,
		Conditions:    conditions,
		IsEnabled:     req.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traffic rule not found"})
			return
		}
		h.logError(c, "failed to update traffic rule", err, zap.Int64("rule_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update traffic rule"})
		return
	}

	response, err := trafficRuleToDTO(rule)
	if err != nil {
		h.logError(c, "failed to serialize traffic rule", err, zap.Int64("rule_id", rule.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize traffic rule"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteTrafficRule(c *gin.Context, id int64) {
	if err := h.trafficRules.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "traffic rule not found"})
			return
		}
		h.logError(c, "failed to delete traffic rule", err, zap.Int64("rule_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete traffic rule"})
		return
	}

	c.Status(http.StatusNoContent)
}

func trafficRuleToDTO(rule models.TrafficRule) (service.TrafficRule, error) {
	conditions, err := unmarshalTrafficConditions(rule.Conditions)
	if err != nil {
		return service.TrafficRule{}, err
	}

	return service.TrafficRule{
		Id:            rule.ID,
		ProfileId:     rule.ProfileID,
		Name:          rule.Name,
		Description:   rule.Description,
		Priority:      rule.Priority,
		DryRun:        rule.DryRun,
		MatchAll:      rule.MatchAll,
		RequestsLimit: rule.RequestsLimit,
		PeriodSeconds: rule.PeriodSeconds,
		Conditions:    conditions,
		IsEnabled:     rule.IsEnabled,
		CreatedAt:     rule.CreatedAt,
		UpdatedAt:     rule.UpdatedAt,
	}, nil
}

func marshalTrafficConditions(value *service.SecurityRuleConditions) (json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}
	return json.Marshal(value)
}

func unmarshalTrafficConditions(value json.RawMessage) (*service.SecurityRuleConditions, error) {
	return normalizeSecurityConditions(value)
}

func validateTrafficRuleRequest(matchAll bool, conditions *service.SecurityRuleConditions) error {
	if matchAll {
		return nil
	}
	return validateSecurityConditions(conditions)
}
