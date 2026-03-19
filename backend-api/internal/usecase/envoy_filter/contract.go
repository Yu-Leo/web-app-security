package envoy_filter

import (
	"context"

	authpb "github.com/envoyproxy/go-control-plane/envoy/service/auth/v3"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type ResourceRepository interface {
	List(ctx context.Context) ([]models.Resource, error)
}

type SecurityProfileRepository interface {
	Get(ctx context.Context, id int64) (models.SecurityProfile, error)
}

type TrafficProfileRepository interface {
	Get(ctx context.Context, id int64) (models.TrafficProfile, error)
}

type SecurityRuleRepository interface {
	List(ctx context.Context) ([]models.SecurityRule, error)
}

type TrafficRuleRepository interface {
	List(ctx context.Context) ([]models.TrafficRule, error)
}

type MLModelRepository interface {
	Get(ctx context.Context, id int64) (models.MLModel, error)
}

type MLScorer interface {
	Score(ctx context.Context, modelName string, modelData []byte, featureVector []float32) (float32, error)
}

type RequestLogRepository interface {
	Create(ctx context.Context, logRecord models.RequestLog) (models.RequestLog, error)
}

type EventLogRepository interface {
	Create(ctx context.Context, logRecord models.EventLog) (models.EventLog, error)
}

type Usecase interface {
	Check(ctx context.Context, req *authpb.CheckRequest) (*Decision, error)
}
