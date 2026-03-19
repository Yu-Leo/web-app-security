package traffic_profiles

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

func (r *Repository) Create(ctx context.Context, profile models.TrafficProfileToCreate) (models.TrafficProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateTrafficProfile(ctx, db.CreateTrafficProfileParams{
		Name:        profile.Name,
		Description: toNullString(profile.Description),
		IsEnabled:   profile.IsEnabled,
	})
	if err != nil {
		return models.TrafficProfile{}, err
	}

	return trafficProfileFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.TrafficProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	profile, err := querier.GetTrafficProfile(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TrafficProfile{}, models.ErrNotFound
		}
		return models.TrafficProfile{}, err
	}

	return trafficProfileFromDB(profile), nil
}

func (r *Repository) List(ctx context.Context) ([]models.TrafficProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	profiles, err := querier.ListTrafficProfiles(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.TrafficProfile, 0, len(profiles))
	for _, profile := range profiles {
		result = append(result, trafficProfileFromDB(profile))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, update models.TrafficProfileToUpdate) (models.TrafficProfile, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	profile, err := querier.UpdateTrafficProfile(ctx, db.UpdateTrafficProfileParams{
		ID:          update.ID,
		Name:        update.Name,
		Description: toNullString(update.Description),
		IsEnabled:   update.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TrafficProfile{}, models.ErrNotFound
		}
		return models.TrafficProfile{}, err
	}

	return trafficProfileFromDB(profile), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	_, err := querier.DeleteTrafficProfile(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

func trafficProfileFromDB(profile db.TrafficProfile) models.TrafficProfile {
	return models.TrafficProfile{
		ID:          profile.ID,
		Name:        profile.Name,
		Description: fromNullString(profile.Description),
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
