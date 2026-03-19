package traffic_profiles

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, profile models.TrafficProfileToCreate) (models.TrafficProfile, error)
	Get(ctx context.Context, id int64) (models.TrafficProfile, error)
	List(ctx context.Context) ([]models.TrafficProfile, error)
	Update(ctx context.Context, update models.TrafficProfileToUpdate) (models.TrafficProfile, error)
	Delete(ctx context.Context, id int64) error
}
