package resources

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

func (r *Repository) Create(ctx context.Context, resource models.ResourceToCreate) (models.Resource, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateResource(ctx, db.CreateResourceParams{
		Name:              resource.Name,
		UrlPattern:        resource.URLPattern,
		SecurityProfileID: toNullInt64(resource.SecurityProfileID),
		TrafficProfileID:  toNullInt64(resource.TrafficProfileID),
	})
	if err != nil {
		return models.Resource{}, err
	}

	return resourceFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.Resource, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	resource, err := querier.GetResource(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Resource{}, models.ErrNotFound
		}
		return models.Resource{}, err
	}

	return resourceFromDB(resource), nil
}

func (r *Repository) List(ctx context.Context) ([]models.Resource, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	resources, err := querier.ListResources(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.Resource, 0, len(resources))
	for _, resource := range resources {
		result = append(result, resourceFromDB(resource))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, update models.ResourceToUpdate) (models.Resource, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	resource, err := querier.UpdateResource(ctx, db.UpdateResourceParams{
		ID:                update.ID,
		Name:              update.Name,
		UrlPattern:        update.URLPattern,
		SecurityProfileID: toNullInt64(update.SecurityProfileID),
		TrafficProfileID:  toNullInt64(update.TrafficProfileID),
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Resource{}, models.ErrNotFound
		}
		return models.Resource{}, err
	}

	return resourceFromDB(resource), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	_, err := querier.DeleteResource(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

func resourceFromDB(resource db.Resource) models.Resource {
	return models.Resource{
		ID:                resource.ID,
		Name:              resource.Name,
		URLPattern:        resource.UrlPattern,
		SecurityProfileID: fromNullInt64(resource.SecurityProfileID),
		TrafficProfileID:  fromNullInt64(resource.TrafficProfileID),
		CreatedAt:         resource.CreatedAt,
		UpdatedAt:         resource.UpdatedAt,
	}
}

func toNullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}

func fromNullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	v := value.Int64
	return &v
}
