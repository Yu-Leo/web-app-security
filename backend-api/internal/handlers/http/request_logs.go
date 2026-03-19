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

func (h *Handler) ListRequestLogs(c *gin.Context) {
	items, err := h.requestLogs.List(c.Request.Context())
	if err != nil {
		h.logError(c, "failed to list request logs", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list request logs"})
		return
	}

	response := make([]service.RequestLog, 0, len(items))
	for _, item := range items {
		dto, err := requestLogToDTO(item)
		if err != nil {
			h.logError(c, "failed to serialize request log", err, zap.Int64("log_id", item.ID))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize request log"})
			return
		}
		response = append(response, dto)
	}

	c.JSON(http.StatusOK, response)
}

func (h *Handler) GetRequestLog(c *gin.Context, id int64) {
	logRecord, err := h.requestLogs.Get(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) || errors.Is(err, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "request log not found"})
			return
		}
		h.logError(c, "failed to get request log", err, zap.Int64("log_id", id))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get request log"})
		return
	}

	response, err := requestLogToDTO(logRecord)
	if err != nil {
		h.logError(c, "failed to serialize request log", err, zap.Int64("log_id", logRecord.ID))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize request log"})
		return
	}

	c.JSON(http.StatusOK, response)
}

func requestLogToDTO(logRecord models.RequestLog) (service.RequestLog, error) {
	metadata, err := unmarshalLogMetadata(logRecord.Metadata)
	if err != nil {
		return service.RequestLog{}, err
	}
	requestHeaders, err := unmarshalLogMetadata(logRecord.RequestHeaders)
	if err != nil {
		return service.RequestLog{}, err
	}
	sourceLabels, err := unmarshalLogMetadata(logRecord.SourceLabels)
	if err != nil {
		return service.RequestLog{}, err
	}
	destinationLabels, err := unmarshalLogMetadata(logRecord.DestinationLabels)
	if err != nil {
		return service.RequestLog{}, err
	}
	contextExtensions, err := unmarshalLogMetadata(logRecord.ContextExtensions)
	if err != nil {
		return service.RequestLog{}, err
	}
	metadataContext, err := unmarshalLogMetadata(logRecord.MetadataContext)
	if err != nil {
		return service.RequestLog{}, err
	}
	routeMetadataContext, err := unmarshalLogMetadata(logRecord.RouteMetadataCtx)
	if err != nil {
		return service.RequestLog{}, err
	}

	return service.RequestLog{
		Id:                   logRecord.ID,
		ResourceId:           logRecord.ResourceID,
		OccurredAt:           logRecord.OccurredAt,
		ClientIp:             logRecord.ClientIP,
		Method:               logRecord.Method,
		Path:                 logRecord.Path,
		StatusCode:           logRecord.StatusCode,
		Action:               logRecord.Action,
		RuleId:               logRecord.RuleID,
		ProfileId:            logRecord.ProfileID,
		UserAgent:            logRecord.UserAgent,
		Country:              logRecord.Country,
		LatencyMs:            logRecord.LatencyMs,
		RequestId:            logRecord.RequestID,
		Metadata:             metadata,
		Host:                 logRecord.Host,
		Scheme:               logRecord.Scheme,
		Protocol:             logRecord.Protocol,
		Authority:            logRecord.Authority,
		Query:                logRecord.Query,
		SourcePort:           logRecord.SourcePort,
		DestinationIp:        logRecord.DestinationIP,
		DestinationPort:      logRecord.DestinationPort,
		SourcePrincipal:      logRecord.SourcePrincipal,
		SourceService:        logRecord.SourceService,
		SourceLabels:         sourceLabels,
		DestinationService:   logRecord.DestinationService,
		DestinationLabels:    destinationLabels,
		RequestHttpId:        logRecord.RequestHTTPID,
		Fragment:             logRecord.Fragment,
		RequestHeaders:       requestHeaders,
		RequestBodySize:      logRecord.RequestBodySize,
		RequestBody:          logRecord.RequestBody,
		ContextExtensions:    contextExtensions,
		MetadataContext:      metadataContext,
		RouteMetadataContext: routeMetadataContext,
	}, nil
}

func unmarshalLogMetadata(value json.RawMessage) (*map[string]interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}
	var result map[string]interface{}
	if err := json.Unmarshal(value, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
