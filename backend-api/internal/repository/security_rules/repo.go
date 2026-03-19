package security_rules

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

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

func (r *Repository) Create(ctx context.Context, rule models.SecurityRuleToCreate) (models.SecurityRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateSecurityRule(ctx, db.CreateSecurityRuleParams{
		ProfileID:   rule.ProfileID,
		Name:        rule.Name,
		Description: toNullString(rule.Description),
		Priority:    int32(rule.Priority),
		RuleType:    string(rule.RuleType),
		Action:      string(rule.Action),
		Conditions:  toNullRawMessage(rule.Conditions),
		MlModelID:   toNullInt64(rule.MLModelID),
		MlThreshold: toNullInt16(rule.MLThreshold),
		DryRun:      rule.DryRun,
		IsEnabled:   rule.IsEnabled,
	})
	if err != nil {
		return models.SecurityRule{}, err
	}

	return securityRuleFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.SecurityRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	rule, err := querier.GetSecurityRule(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SecurityRule{}, models.ErrNotFound
		}
		return models.SecurityRule{}, err
	}

	return securityRuleFromDB(rule), nil
}

func (r *Repository) List(ctx context.Context) ([]models.SecurityRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	rules, err := querier.ListSecurityRules(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.SecurityRule, 0, len(rules))
	for _, rule := range rules {
		result = append(result, securityRuleFromDB(rule))
	}

	return result, nil
}

func (r *Repository) Update(ctx context.Context, update models.SecurityRuleToUpdate) (models.SecurityRule, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	rule, err := querier.UpdateSecurityRule(ctx, db.UpdateSecurityRuleParams{
		ID:          update.ID,
		ProfileID:   update.ProfileID,
		Name:        update.Name,
		Description: toNullString(update.Description),
		Priority:    int32(update.Priority),
		RuleType:    string(update.RuleType),
		Action:      string(update.Action),
		Conditions:  toNullRawMessage(update.Conditions),
		MlModelID:   toNullInt64(update.MLModelID),
		MlThreshold: toNullInt16(update.MLThreshold),
		DryRun:      update.DryRun,
		IsEnabled:   update.IsEnabled,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.SecurityRule{}, models.ErrNotFound
		}
		return models.SecurityRule{}, err
	}

	return securityRuleFromDB(rule), nil
}

func (r *Repository) Delete(ctx context.Context, id int64) error {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	_, err := querier.DeleteSecurityRule(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrNotFound
		}
		return err
	}

	return nil
}

func securityRuleFromDB(rule any) models.SecurityRule {
	switch row := rule.(type) {
	case db.CreateSecurityRuleRow:
		return securityRuleFromFields(row.ID, row.ProfileID, row.Name, row.Description, row.Priority, row.RuleType, row.Action, row.Conditions, row.MlModelID, row.MlThreshold, row.DryRun, row.IsEnabled, row.CreatedAt, row.UpdatedAt)
	case db.GetSecurityRuleRow:
		return securityRuleFromFields(row.ID, row.ProfileID, row.Name, row.Description, row.Priority, row.RuleType, row.Action, row.Conditions, row.MlModelID, row.MlThreshold, row.DryRun, row.IsEnabled, row.CreatedAt, row.UpdatedAt)
	case db.ListSecurityRulesRow:
		return securityRuleFromFields(row.ID, row.ProfileID, row.Name, row.Description, row.Priority, row.RuleType, row.Action, row.Conditions, row.MlModelID, row.MlThreshold, row.DryRun, row.IsEnabled, row.CreatedAt, row.UpdatedAt)
	case db.UpdateSecurityRuleRow:
		return securityRuleFromFields(row.ID, row.ProfileID, row.Name, row.Description, row.Priority, row.RuleType, row.Action, row.Conditions, row.MlModelID, row.MlThreshold, row.DryRun, row.IsEnabled, row.CreatedAt, row.UpdatedAt)
	default:
		panic("unsupported security rule row type")
	}
}

func securityRuleFromFields(
	id int64,
	profileID int64,
	name string,
	description sql.NullString,
	priority int32,
	ruleType string,
	action string,
	conditions pqtype.NullRawMessage,
	mlModelID sql.NullInt64,
	mlThreshold sql.NullInt16,
	dryRun bool,
	isEnabled bool,
	createdAt time.Time,
	updatedAt time.Time,
) models.SecurityRule {
	return models.SecurityRule{
		ID:          id,
		ProfileID:   profileID,
		Name:        name,
		Description: fromNullString(description),
		Priority:    priority,
		RuleType:    models.SecurityRuleType(ruleType),
		Action:      models.SecurityRuleAction(action),
		Conditions:  fromNullRawMessage(conditions),
		MLModelID:   fromNullInt64(mlModelID),
		MLThreshold: fromNullInt16(mlThreshold),
		DryRun:      dryRun,
		IsEnabled:   isEnabled,
		CreatedAt:   createdAt,
		UpdatedAt:   updatedAt,
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

func toNullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: *value,
		Valid: true,
	}
}

func fromNullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

func toNullInt16(value *int16) sql.NullInt16 {
	if value == nil {
		return sql.NullInt16{}
	}
	return sql.NullInt16{
		Int16: *value,
		Valid: true,
	}
}

func fromNullInt16(value sql.NullInt16) *int16 {
	if !value.Valid {
		return nil
	}
	result := value.Int16
	return &result
}
