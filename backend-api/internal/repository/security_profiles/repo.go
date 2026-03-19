package security_profiles

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/storages/db"
)

type Repository struct {
	database database
}

func NewRepository(dbProvider database) *Repository {
	return &Repository{database: dbProvider}
}

func (r *Repository) Create(ctx context.Context, profile models.SecurityProfileToCreate) (models.SecurityProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateSecurityProfile(ctx, db.CreateSecurityProfileParams{
		Name:        profile.Name,
		Description: toNullString(profile.Description),
		BaseAction:  string(profile.BaseAction),
		LogEnabled:  profile.LogEnabled,
		IsEnabled:   profile.IsEnabled,
	})
	if err != nil {
		return models.SecurityProfile{}, err
	}

	return securityProfileFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.SecurityProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	profile, err := querier.GetSecurityProfile(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SecurityProfile{}, models.ErrNotFound
		}
		return models.SecurityProfile{}, err
	}

	return securityProfileFromDB(profile), nil
}

func (r *Repository) List(ctx context.Context) ([]models.SecurityProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	profiles, err := querier.ListSecurityProfiles(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.SecurityProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, securityProfileFromDB(profile))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, update models.SecurityProfileToUpdate) (models.SecurityProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	profile, err := querier.UpdateSecurityProfile(ctx, db.UpdateSecurityProfileParams{
		ID:          update.ID,
		Name:        update.Name,
		Description: toNullString(update.Description),
		BaseAction:  string(update.BaseAction),
		LogEnabled:  update.LogEnabled,
		IsEnabled:   update.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SecurityProfile{}, models.ErrNotFound
		}
		return models.SecurityProfile{}, err
	}

	return securityProfileFromDB(profile), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	_, err := querier.DeleteSecurityProfile(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

func securityProfileFromDB(profile db.SecurityProfile) models.SecurityProfile {
	return models.SecurityProfile{
		ID:          profile.ID,
		Name:        profile.Name,
		Description: fromNullString(profile.Description),
		BaseAction:  models.SecurityProfileBaseAction(profile.BaseAction),
		LogEnabled:  profile.LogEnabled,
		IsEnabled:   profile.IsEnabled,
		CreatedAt:   profile.CreatedAt,
		UpdatedAt:   profile.UpdatedAt,
	}
}

func toNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *value,
		Valid:  true,
	}
}

func fromNullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}
