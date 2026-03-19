package event_logs

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Repository interface {
	Create(ctx context.Context, logRecord models.EventLog) (models.EventLog, error)
	Get(ctx context.Context, id int64) (models.EventLog, error)
	List(ctx context.Context) ([]models.EventLog, error)
}
