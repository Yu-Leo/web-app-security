package resources

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

func (u *Usecase) Create(ctx context.Context, resource models.ResourceToCreate) (models.Resource, error) {
	return u.repository.Create(ctx, resource)
}

func (u *Usecase) Get(ctx context.Context, id int64) (models.Resource, error) {
	return u.repository.Get(ctx, id)
}

func (u *Usecase) List(ctx context.Context) ([]models.Resource, error) {
	return u.repository.List(ctx)
}

func (u *Usecase) Update(ctx context.Context, resource models.ResourceToUpdate) (models.Resource, error) {
	return u.repository.Update(ctx, resource)
}

func (u *Usecase) Delete(ctx context.Context, id int64) error {
	return u.repository.Delete(ctx, id)
}
