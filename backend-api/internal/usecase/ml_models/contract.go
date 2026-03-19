package ml_models

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, model models.MLModelToCreate) (models.MLModel, error)
	Get(ctx context.Context, id int64) (models.MLModel, error)
	List(ctx context.Context) ([]models.MLModel, error)
	Update(ctx context.Context, update models.MLModelToUpdate) (models.MLModel, error)
	Delete(ctx context.Context, id int64) error
}
