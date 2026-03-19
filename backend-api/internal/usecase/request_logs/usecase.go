package request_logs

import (
	"context"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
)

type Usecase struct {
	repository Repository
}

func NewUsecase(repository Repository) *Usecase {
	return &Usecase{repository: repository}
}

func (u *Usecase) Create(ctx context.Context, logRecord models.RequestLog) (models.RequestLog, error) {
	return u.repository.Create(ctx, logRecord)
}

func (u *Usecase) Get(ctx context.Context, id int64) (models.RequestLog, error) {
	return u.repository.Get(ctx, id)
}

func (u *Usecase) List(ctx context.Context) ([]models.RequestLog, error) {
	return u.repository.List(ctx)
}
