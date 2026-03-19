package http

import (
	eventLogs "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/event_logs"
	mlModels "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/ml_models"
	requestLogs "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/request_logs"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/resources"
	securityProfiles "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/security_profiles"
	securityRules "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/security_rules"
	trafficProfiles "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/traffic_profiles"
	trafficRules "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/traffic_rules"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Handler struct {
	resources        *resources.Usecase
	securityProfiles *securityProfiles.Usecase
	securityRules    *securityRules.Usecase
	trafficProfiles  *trafficProfiles.Usecase
	trafficRules     *trafficRules.Usecase
	mlModels         *mlModels.Usecase
	requestLogs      *requestLogs.Usecase
	eventLogs        *eventLogs.Usecase
	logger           *zap.Logger
}

func New(
	resourcesUsecase *resources.Usecase,
	securityProfilesUsecase *securityProfiles.Usecase,
	securityRulesUsecase *securityRules.Usecase,
	trafficProfilesUsecase *trafficProfiles.Usecase,
	trafficRulesUsecase *trafficRules.Usecase,
	mlModelsUsecase *mlModels.Usecase,
	requestLogsUsecase *requestLogs.Usecase,
	eventLogsUsecase *eventLogs.Usecase,
	logger *zap.Logger,
) *Handler {
	return &Handler{
		resources:        resourcesUsecase,
		securityProfiles: securityProfilesUsecase,
		securityRules:    securityRulesUsecase,
		trafficProfiles:  trafficProfilesUsecase,
		trafficRules:     trafficRulesUsecase,
		mlModels:         mlModelsUsecase,
		requestLogs:      requestLogsUsecase,
		eventLogs:        eventLogsUsecase,
		logger:           logger,
	}
}

func (h *Handler) logError(c *gin.Context, message string, err error, fields ...zap.Field) {
	requestFields := []zap.Field{
		zap.String("method", c.Request.Method),
		zap.String("path", c.Request.URL.Path),
	}
	fields = append(fields, requestFields...)
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	h.logger.Error(message, fields...)
}
