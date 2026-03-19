package app

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/config"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/logging"
	txprovider "github.com/Yu-Leo/web-app-security/backend-api/internal/storages/transaction/provider"
)

type Service struct {
	name string

	config     *config.Config
	db         *sql.DB
	txProvider *txprovider.Provider
	httpRouter *gin.Engine
	httpServer *http.Server
	grpcServer *grpc.Server
	logger     *zap.Logger
}

func New() *Service {
	return &Service{name: "WAS"}
}

func (srv *Service) Name() string {
	return srv.name
}

func (srv *Service) ConfigureService() error {
	var err error

	if srv.config, err = srv.initConfig(); err != nil {
		return err
	}

	if err = srv.ConfigureComponents(); err != nil {
		return err
	}

	if err := srv.SetupHTTPAPIHandlers(); err != nil {
		return err
	}

	if err := srv.SetupGRPCAPIHandlers(); err != nil {
		return err
	}

	return nil
}

func (srv *Service) initConfig() (*config.Config, error) {
	prefix := srv.Name()

	var cfg config.Config
	if err := envconfig.Process(prefix, &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config from env: %w", err)
	}

	return &cfg, nil
}

func (srv *Service) ConfigureComponents() error {
	logger, err := logging.New()
	if err != nil {
		return fmt.Errorf("failed to init logger: %w", err)
	}
	srv.logger = logger

	dbConn, err := sql.Open("postgres", srv.config.PostgresMaster)
	if err != nil {
		return fmt.Errorf("failed to open postgres: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := dbConn.PingContext(ctx); err != nil {
		return fmt.Errorf("failed to ping postgres: %w", err)
	}

	srv.db = dbConn
	srv.txProvider = txprovider.New(dbConn)

	return nil
}

func (srv *Service) getGinMode() string {
	switch srv.config.Env {
	case config.EnvLocal, config.EnvTests:
		return gin.DebugMode
	}

	return gin.ReleaseMode
}
