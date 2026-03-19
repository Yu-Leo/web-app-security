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

func (h *Handler) ListEventLogs(c *gin.Context) {
	items, err := h.eventLogs.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list event logs", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list event logs"})
		return
	}

	response := make([]service.EventLog, 0, len(items))
	for _, item := range items {
		dto, err := eventLogToDTO(item)
		if err != nil {
			h.logError(c, "failed to serialize event log", err, zap.Int64("log_id", item.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize event log"})
			return
		}
		response = append(response, dto)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetEventLog(c *gin.Context, id int64) {
	logRecord, err := h.eventLogs.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "event log not found"})
			return
		}
		h.logError(c, "failed to get event log", err, zap.Int64("log_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get event log"})
		return
	}

	response, err := eventLogToDTO(logRecord)
	if err != nil {
		h.logError(c, "failed to serialize event log", err, zap.Int64("log_id", logRecord.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize event log"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func eventLogToDTO(logRecord models.EventLog) (service.EventLog, error) {
	metadata, err := unmarshalEventMetadata(logRecord.Metadata)
	if err != nil {
		return service.EventLog{}, err
	}

	return service.EventLog{
		Id:         logRecord.ID,
		ResourceId: logRecord.ResourceID,
		OccurredAt: logRecord.OccurredAt,
		EventType:  logRecord.EventType,
		Severity:   logRecord.Severity,
		Message:    logRecord.Message,
		RuleId:     logRecord.RuleID,
		ProfileId:  logRecord.ProfileID,
		Metadata:   metadata,
	}, nil
}

func unmarshalEventMetadata(value json.RawMessage) (*map[string]interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(value, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
