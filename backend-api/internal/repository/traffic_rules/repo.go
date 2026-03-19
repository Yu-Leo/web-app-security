package traffic_rules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/Yu-Leo/web-app-security/backend-api/internal/models"
	"github.com/Yu-Leo/web-app-security/backend-api/internal/storages/db"
	"github.com/sqlc-dev/pqtype"
)

type Repository struct {
	database database
}

func NewRepository(dbProvider database) *Repository {
	return &Repository{database: dbProvider}
}

func (r *Repository) Create(ctx context.Context, rule models.TrafficRuleToCreate) (models.TrafficRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateTrafficRule(ctx, db.CreateTrafficRuleParams{
		ProfileID:     rule.ProfileID,
		Name:          rule.Name,
		Description:   toNullString(rule.Description),
		Priority:      rule.Priority,
		DryRun:        rule.DryRun,
		MatchAll:      rule.MatchAll,
		RequestsLimit: rule.RequestsLimit,
		PeriodSeconds: rule.PeriodSeconds,
		Conditions:    toNullRawMessage(rule.Conditions),
		IsEnabled:     rule.IsEnabled,
	})
	if err != nil {
		return models.TrafficRule{}, err
	}

	return trafficRuleFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.TrafficRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	rule, err := querier.GetTrafficRule(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TrafficRule{}, models.ErrNotFound
		}
		return models.TrafficRule{}, err
	}

	return trafficRuleFromDB(rule), nil
}

func (r *Repository) List(ctx context.Context) ([]models.TrafficRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	rules, err := querier.ListTrafficRules(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.TrafficRule, 0, len(rules))
	for _, rule := range rules {
		result = append(result, trafficRuleFromDB(rule))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, update models.TrafficRuleToUpdate) (models.TrafficRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	rule, err := querier.UpdateTrafficRule(ctx, db.UpdateTrafficRuleParams{
		ID:            update.ID,
		ProfileID:     update.ProfileID,
		Name:          update.Name,
		Description:   toNullString(update.Description),
		Priority:      update.Priority,
		DryRun:        update.DryRun,
		MatchAll:      update.MatchAll,
		RequestsLimit: update.RequestsLimit,
		PeriodSeconds: update.PeriodSeconds,
		Conditions:    toNullRawMessage(update.Conditions),
		IsEnabled:     update.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.TrafficRule{}, models.ErrNotFound
		}
		return models.TrafficRule{}, err
	}

	return trafficRuleFromDB(rule), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	_, err := querier.DeleteTrafficRule(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

func trafficRuleFromDB(rule db.TrafficRule) models.TrafficRule {
	return models.TrafficRule{
		ID:            rule.ID,
		ProfileID:     rule.ProfileID,
		Name:          rule.Name,
		Description:   fromNullString(rule.Description),
		Priority:      rule.Priority,
		DryRun:        rule.DryRun,
		MatchAll:      rule.MatchAll,
		RequestsLimit: rule.RequestsLimit,
		PeriodSeconds: rule.PeriodSeconds,
		Conditions:    fromNullRawMessage(rule.Conditions),
		IsEnabled:     rule.IsEnabled,
		CreatedAt:     rule.CreatedAt,
		UpdatedAt:     rule.UpdatedAt,
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

func toNullRawMessage(value json.RawMessage) pqtype.NullRawMessage {
	if len(value) == 0 {
		return pqtype.NullRawMessage{}
	}
	return pqtype.NullRawMessage{
		RawMessage: value,
		Valid:      true,
	}
}

func fromNullRawMessage(value pqtype.NullRawMessage) json.RawMessage {
	if !value.Valid {
		return nil
	}
	return value.RawMessage
}
