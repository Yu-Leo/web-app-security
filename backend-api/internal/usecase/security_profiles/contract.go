package security_profiles

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, profile models.SecurityProfileToCreate) (models.SecurityProfile, error)
	Get(ctx context.Context, id int64) (models.SecurityProfile, error)
	List(ctx context.Context) ([]models.SecurityProfile, error)
	Update(ctx context.Context, update models.SecurityProfileToUpdate) (models.SecurityProfile, error)
	Delete(ctx context.Context, id int64) error
}
