package app

import (
	"time"

	dto "github.com/Yu-Leo/web-app-security/backend-api/internal/generated/service"
	httpHandlers "github.com/Yu-Leo/web-app-security/backend-api/internal/handlers/http"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/monitoring"
	eventLogsRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/event_logs"
	mlModelsRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/ml_models"
	requestLogsRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/request_logs"
	resourcesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/resources"
	securityProfilesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/security_profiles"
	securityRulesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/security_rules"
	trafficProfilesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/traffic_profiles"
	trafficRulesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/traffic_rules"
	eventLogsUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/event_logs"
	mlModelsUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/ml_models"
	requestLogsUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/request_logs"
	resourcesUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/resources"
	securityProfilesUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/security_profiles"
	securityRulesUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/security_rules"
	trafficProfilesUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/traffic_profiles"
	trafficRulesUsecase "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/traffic_rules"
	"github.com/gin-contrib/cors"
	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (srv *Service) SetupHTTPAPIHandlers() error {
	gin.SetMode(srv.getGinMode())

	srv.httpRouter = gin.New()
	srv.httpRouter.Use(ginzap.Ginzap(srv.logger, time.RFC3339, true))
	srv.httpRouter.Use(ginzap.RecoveryWithZap(srv.logger, true))
	srv.httpRouter.Use(monitoring.GinMiddleware())
	srv.httpRouter.Use(cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},
		AllowHeaders: []string{
			"Origin",
			"Content-Length",
			"Content-Type",
			"Authorization",
			"Accept",
			"Cache-Control",
			"Pragma",
			"Expires",
			"X-Requested-With",
			"If-Modified-Since",
			"If-None-Match",
		},
		ExposeHeaders: []string{
			"Content-Length",
			"Content-Type",
			"Date",
			"ETag",
		},
		AllowCredentials: false,
		MaxAge:           24 * time.Hour,
	}))

	resourcesRepository := resourcesRepo.NewRepository(srv.txProvider)
	resourcesUC := resourcesUsecase.NewUsecase(resourcesRepository)
	securityProfilesRepository := securityProfilesRepo.NewRepository(srv.txProvider)
	securityProfilesUC := securityProfilesUsecase.NewUsecase(securityProfilesRepository)
	securityRulesRepository := securityRulesRepo.NewRepository(srv.txProvider)
	securityRulesUC := securityRulesUsecase.NewUsecase(securityRulesRepository)
	trafficProfilesRepository := trafficProfilesRepo.NewRepository(srv.txProvider)
	trafficProfilesUC := trafficProfilesUsecase.NewUsecase(trafficProfilesRepository)
	trafficRulesRepository := trafficRulesRepo.NewRepository(srv.txProvider)
	trafficRulesUC := trafficRulesUsecase.NewUsecase(trafficRulesRepository)
	mlModelsRepository := mlModelsRepo.NewRepository(srv.txProvider)
	mlModelsUC := mlModelsUsecase.NewUsecase(mlModelsRepository)
	requestLogsRepository := requestLogsRepo.NewRepository(srv.txProvider)
	requestLogsUC := requestLogsUsecase.NewUsecase(requestLogsRepository)
	eventLogsRepository := eventLogsRepo.NewRepository(srv.txProvider)
	eventLogsUC := eventLogsUsecase.NewUsecase(eventLogsRepository)
	handler := httpHandlers.New(
		resourcesUC,
		securityProfilesUC,
		securityRulesUC,
		trafficProfilesUC,
		trafficRulesUC,
		mlModelsUC,
		requestLogsUC,
		eventLogsUC,
		srv.logger,
	)

	srv.httpRouter.StaticFile("/docs/openapi.yaml", "./api/service/openapi.yaml")
	srv.httpRouter.GET("/metrics", gin.WrapH(monitoring.Handler()))
	srv.httpRouter.GET("/docs/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/openapi.yaml")))
	dto.RegisterHandlers(srv.httpRouter, handler)

	return nil
}
