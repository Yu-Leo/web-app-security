package http

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/generated/service"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (h *Handler) ListSecurityRules(c *gin.Context) {
	items, err := h.securityRules.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list security rules", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list security rules"})
		return
	}

	response := make([]service.SecurityRule, 0, len(items))
	for _, item := range items {
		dto, err := securityRuleToDTO(item)
		if err != nil {
			h.logError(c, "failed to serialize security rule", err, zap.Int64("rule_id", item.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize security rule"})
			return
		}
		response = append(response, dto)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) CreateSecurityRule(c *gin.Context) {
	var req service.CreateSecurityRuleJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if _, err := h.securityProfiles.Get(c.Request.Context(), req.ProfileId); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "security profile not found"})
			return
		}
		h.logError(c, "failed to validate security profile", err, zap.Int64("profile_id", req.ProfileId))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate security profile"})
		return
	}
	if err := validateSecurityRuleRequest(string(req.RuleType), string(req.Action), req.MlModelId, req.MlThreshold, req.Conditions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validateSecurityRuleMLModel(c.Request.Context(), string(req.RuleType), req.MlModelId); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ml model not found"})
			return
		}
		h.logError(c, "failed to validate ml model", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate ml model"})
		return
	}

	conditions, err := marshalSecurityConditions(req.Conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conditions"})
		return
	}

	rule, err := h.securityRules.Create(c.Request.Context(), models.SecurityRuleToCreate{
		ProfileID:   req.ProfileId,
		Name:        req.Name,
		Description: req.Description,
		Priority:    req.Priority,
		RuleType:    models.SecurityRuleType(req.RuleType),
		Action:      models.SecurityRuleAction(req.Action),
		Conditions:  conditions,
		MLModelID:   req.MlModelId,
		MLThreshold: int32ToInt16Ptr(req.MlThreshold),
		DryRun:      req.DryRun,
		IsEnabled:   req.IsEnabled,
	})
	if err != nil {
		h.logError(c, "failed to create security rule", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create security rule"})
		return
	}

	response, err := securityRuleToDTO(rule)
	if err != nil {
		h.logError(c, "failed to serialize security rule", err, zap.Int64("rule_id", rule.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize security rule"})
		return
	}

	c.JSON(http.StatusCreated, response)
}

func (h *Handler) GetSecurityRule(c *gin.Context, id int64) {
	rule, err := h.securityRules.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "security rule not found"})
			return
		}
		h.logError(c, "failed to get security rule", err, zap.Int64("rule_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get security rule"})
		return
	}

	response, err := securityRuleToDTO(rule)
	if err != nil {
		h.logError(c, "failed to serialize security rule", err, zap.Int64("rule_id", rule.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize security rule"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) UpdateSecurityRule(c *gin.Context, id int64) {
	var req service.UpdateSecurityRuleJSONRequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if _, err := h.securityProfiles.Get(c.Request.Context(), req.ProfileId); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "security profile not found"})
			return
		}
		h.logError(c, "failed to validate security profile", err, zap.Int64("profile_id", req.ProfileId))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate security profile"})
		return
	}
	if err := validateSecurityRuleRequest(string(req.RuleType), string(req.Action), req.MlModelId, req.MlThreshold, req.Conditions); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.validateSecurityRuleMLModel(c.Request.Context(), string(req.RuleType), req.MlModelId); err != nil {
		if errors.Is(err, models.ErrNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ml model not found"})
			return
		}
		h.logError(c, "failed to validate ml model", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate ml model"})
		return
	}

	conditions, err := marshalSecurityConditions(req.Conditions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid conditions"})
		return
	}

	rule, err := h.securityRules.Update(c.Request.Context(), models.SecurityRuleToUpdate{
		ID:          id,
		ProfileID:   req.ProfileId,
		Name:        req.Name,
		Description: req.Description,
		Priority:    req.Priority,
		RuleType:    models.SecurityRuleType(req.RuleType),
		Action:      models.SecurityRuleAction(req.Action),
		Conditions:  conditions,
		MLModelID:   req.MlModelId,
		MLThreshold: int32ToInt16Ptr(req.MlThreshold),
		DryRun:      req.DryRun,
		IsEnabled:   req.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "security rule not found"})
			return
		}
		h.logError(c, "failed to update security rule", err, zap.Int64("rule_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update security rule"})
		return
	}

	response, err := securityRuleToDTO(rule)
	if err != nil {
		h.logError(c, "failed to serialize security rule", err, zap.Int64("rule_id", rule.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize security rule"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) DeleteSecurityRule(c *gin.Context, id int64) {
	if err := h.securityRules.Delete(c.Request.Context(), id); err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "security rule not found"})
			return
		}
		h.logError(c, "failed to delete security rule", err, zap.Int64("rule_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete security rule"})
		return
	}

	c.Status(http.StatusNoContent)
}

func securityRuleToDTO(rule models.SecurityRule) (service.SecurityRule, error) {
	conditions, err := unmarshalSecurityConditions(rule.Conditions)
	if err != nil {
		return service.SecurityRule{}, err
	}

	return service.SecurityRule{
		Id:          rule.ID,
		ProfileId:   rule.ProfileID,
		Name:        rule.Name,
		Description: rule.Description,
		Priority:    rule.Priority,
		RuleType:    service.SecurityRuleRuleType(rule.RuleType),
		Action:      service.SecurityRuleAction(rule.Action),
		Conditions:  conditions,
		MlModelId:   rule.MLModelID,
		MlThreshold: int16ToInt32Ptr(rule.MLThreshold),
		DryRun:      rule.DryRun,
		IsEnabled:   rule.IsEnabled,
		CreatedAt:   rule.CreatedAt,
		UpdatedAt:   rule.UpdatedAt,
	}, nil
}

func marshalSecurityConditions(value *service.SecurityRuleConditions) (json.RawMessage, error) {
	if value == nil {
		return nil, nil
	}
	return json.Marshal(value)
}

func unmarshalSecurityConditions(value json.RawMessage) (*service.SecurityRuleConditions, error) {
	return normalizeSecurityConditions(value)
}

func int32ToInt16Ptr(value *int32) *int16 {
	if value == nil {
		return nil
	}
	result := int16(*value)
	return &result
}

func int16ToInt32Ptr(value *int16) *int32 {
	if value == nil {
		return nil
	}
	result := int32(*value)
	return &result
}

func validateSecurityRuleRequest(
	ruleType string,
	action string,
	modelID *int64,
	threshold *int32,
	conditions *service.SecurityRuleConditions,
) error {
	if !isSupportedSecurityRuleType(ruleType) {
		return fmt.Errorf("unsupported rule_type: %s", ruleType)
	}
	if !isSupportedSecurityAction(action) {
		return fmt.Errorf("unsupported action: %s", action)
	}

	switch ruleType {
	case string(models.SecurityRuleTypeML):
		if modelID == nil || threshold == nil {
			return errors.New("ml rules require ml_model_id and ml_threshold")
		}
	case string(models.SecurityRuleTypeDeterministic):
		if modelID != nil || threshold != nil {
			return errors.New("deterministic rules must not contain ml_model_id or ml_threshold")
		}
	}

	if threshold != nil && (*threshold < 0 || *threshold > 100) {
		return fmt.Errorf("ml_threshold must be between 0 and 100: got %d", *threshold)
	}

	return validateSecurityConditions(conditions)
}

func isSupportedSecurityRuleType(value string) bool {
	return value == string(models.SecurityRuleTypeDeterministic) || value == string(models.SecurityRuleTypeML)
}

func isSupportedSecurityAction(value string) bool {
	return value == string(models.SecurityRuleActionAllow) || value == string(models.SecurityRuleActionBlock)
}

func (h *Handler) validateSecurityRuleMLModel(ctx context.Context, ruleType string, modelID *int64) error {
	if ruleType != string(models.SecurityRuleTypeML) || modelID == nil {
		return nil
	}

	_, err := h.mlModels.Get(ctx, *modelID)
	return err
}
