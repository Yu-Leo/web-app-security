package ml_models

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

func (r *Repository) Create(ctx context.Context, model models.MLModelToCreate) (models.MLModel, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateMLModel(ctx, db.CreateMLModelParams{
		Name:      model.Name,
		ModelData: model.ModelData,
	})
	if err != nil {
		return models.MLModel{}, err
	}

	return mlModelFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.MLModel, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	model, err := querier.GetMLModel(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.MLModel{}, models.ErrNotFound
		}
		return models.MLModel{}, err
	}

	return mlModelFromDB(model), nil
}

func (r *Repository) List(ctx context.Context) ([]models.MLModel, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	modelsDB, err := querier.ListMLModels(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.MLModel, 0, len(modelsDB))
	for _, model := range modelsDB {
		result = append(result, mlModelFromDB(model))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, update models.MLModelToUpdate) (models.MLModel, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	model, err := querier.UpdateMLModel(ctx, db.UpdateMLModelParams{
		ID:        update.ID,
		Name:      update.Name,
		ModelData: update.ModelData,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.MLModel{}, models.ErrNotFound
		}
		return models.MLModel{}, err
	}

	return mlModelFromDB(model), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	_, err := querier.DeleteMLModel(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

func mlModelFromDB(model db.MlModel) models.MLModel {
	return models.MLModel{
		ID:        model.ID,
		Name:      model.Name,
		ModelData: model.ModelData,
	}
}
