package app

import (
	"google.golang.org/grpc"

	authpb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"
	"go.uber.org/zap"

	grpcHandlers "github.com/Yu-Leo/web-app-security/backend-api/internal/handlers/grpc"
	eventLogsRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/event_logs"
	mlModelsRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/ml_models"
	requestLogsRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/request_logs"
	resourcesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/resources"
	securityProfilesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/security_profiles"
	securityRulesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/security_rules"
	trafficProfilesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/traffic_profiles"
	trafficRulesRepo "github.com/Yu-Leo/web-app-security/backend-api/internal/repository/traffic_rules"
	envoyFilter "github.com/Yu-Leo/web-app-security/backend-api/internal/usecase/envoy_filter"
)

func (srv *Service) SetupGRPCAPIHandlers() error {
	srv.grpcServer = grpc.NewServer()

	resourcesRepository := resourcesRepo.NewRepository(srv.txProvider)
	securityProfilesRepository := securityProfilesRepo.NewRepository(srv.txProvider)
	trafficProfilesRepository := trafficProfilesRepo.NewRepository(srv.txProvider)
	securityRulesRepository := securityRulesRepo.NewRepository(srv.txProvider)
	trafficRulesRepository := trafficRulesRepo.NewRepository(srv.txProvider)
	mlModelsRepository := mlModelsRepo.NewRepository(srv.txProvider)
	requestLogsRepository := requestLogsRepo.NewRepository(srv.txProvider)
	eventLogsRepository := eventLogsRepo.NewRepository(srv.txProvider)
	var mlScorer envoyFilter.MLScorer = envoyFilter.NewUnavailableMLScorer()
	onnxScorer, err := envoyFilter.NewONNXMLScorer()
	if err != nil {
		srv.logger.Warn("failed to initialize onnx runtime scorer; fallback scorer will be used", zap.Error(err))
	} else {
		mlScorer = onnxScorer
	}

	envoyUsecase := envoyFilter.NewUsecase(
		resourcesRepository,
		securityProfilesRepository,
		trafficProfilesRepository,
		securityRulesRepository,
		trafficRulesRepository,
		mlModelsRepository,
		mlScorer,
		requestLogsRepository,
		eventLogsRepository,
	)
	authzHandler := grpcHandlers.NewAuthzHandler(envoyUsecase)
	authpb.RegisterAuthorizationServer(srv.grpcServer, authzHandler)

	return nil
}
