package security_rules

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

func (u *Usecase) Create(ctx context.Context, rule models.SecurityRuleToCreate) (models.SecurityRule, error) {
	return u.repository.Create(ctx, rule)
}

func (u *Usecase) Get(ctx context.Context, id int64) (models.SecurityRule, error) {
	return u.repository.Get(ctx, id)
}

func (u *Usecase) List(ctx context.Context) ([]models.SecurityRule, error) {
	return u.repository.List(ctx)
}

func (u *Usecase) Update(ctx context.Context, rule models.SecurityRuleToUpdate) (models.SecurityRule, error) {
	return u.repository.Update(ctx, rule)
}

func (u *Usecase) Delete(ctx context.Context, id int64) error {
	return u.repository.Delete(ctx, id)
}
