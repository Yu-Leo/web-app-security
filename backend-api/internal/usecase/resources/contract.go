package resources

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, resource models.ResourceToCreate) (models.Resource, error)
	Get(ctx context.Context, id int64) (models.Resource, error)
	List(ctx context.Context) ([]models.Resource, error)
	Update(ctx context.Context, update models.ResourceToUpdate) (models.Resource, error)
	Delete(ctx context.Context, id int64) error
}
