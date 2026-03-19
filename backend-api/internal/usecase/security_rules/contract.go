package security_rules

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, rule models.SecurityRuleToCreate) (models.SecurityRule, error)
	Get(ctx context.Context, id int64) (models.SecurityRule, error)
	List(ctx context.Context) ([]models.SecurityRule, error)
	Update(ctx context.Context, update models.SecurityRuleToUpdate) (models.SecurityRule, error)
	Delete(ctx context.Context, id int64) error
}
