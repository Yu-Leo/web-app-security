package traffic_rules

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, rule models.TrafficRuleToCreate) (models.TrafficRule, error)
	Get(ctx context.Context, id int64) (models.TrafficRule, error)
	List(ctx context.Context) ([]models.TrafficRule, error)
	Update(ctx context.Context, update models.TrafficRuleToUpdate) (models.TrafficRule, error)
	Delete(ctx context.Context, id int64) error
}
