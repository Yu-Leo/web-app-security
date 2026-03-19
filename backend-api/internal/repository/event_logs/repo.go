package event_logs

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

func (r *Repository) Create(ctx context.Context, logRecord models.EventLog) (models.EventLog, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateEventLog(ctx, db.CreateEventLogParams{
		ResourceID: logRecord.ResourceID,
		OccurredAt: logRecord.OccurredAt,
		EventType:  logRecord.EventType,
		Severity:   logRecord.Severity,
		Message:    logRecord.Message,
		RuleID:     toNullInt64(logRecord.RuleID),
		ProfileID:  toNullInt64(logRecord.ProfileID),
		Metadata:   toNullRawMessage(logRecord.Metadata),
		RequestID:  toNullString(logRecord.RequestID),
		ClientIp:   toNullString(logRecord.ClientIP),
		Method:     toNullString(logRecord.Method),
		Path:       toNullString(logRecord.Path),
	})
	if err != nil {
		return models.EventLog{}, err
	}

	return eventLogFromDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.EventLog, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	logRecord, err := querier.GetEventLog(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.EventLog{}, models.ErrNotFound
		}
		return models.EventLog{}, err
	}

	return eventLogFromDB(logRecord), nil
}

func (r *Repository) List(ctx context.Context) ([]models.EventLog, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	logs, err := querier.ListEventLogs(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.EventLog, 0, len(logs))
	for _, logRecord := range logs {
		result = append(result, eventLogFromDB(logRecord))
	}

	return result, nil
}

func eventLogFromDB(logRecord db.EventLog) models.EventLog {
	return models.EventLog{
		ID:         logRecord.ID,
		ResourceID: logRecord.ResourceID,
		OccurredAt: logRecord.OccurredAt,
		EventType:  logRecord.EventType,
		Severity:   logRecord.Severity,
		Message:    logRecord.Message,
		RuleID:     fromNullInt64(logRecord.RuleID),
		ProfileID:  fromNullInt64(logRecord.ProfileID),
		Metadata:   fromNullRawMessage(logRecord.Metadata),
		RequestID:  fromNullString(logRecord.RequestID),
		ClientIP:   fromNullString(logRecord.ClientIp),
		Method:     fromNullString(logRecord.Method),
		Path:       fromNullString(logRecord.Path),
	}
}

func fromNullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

func fromNullRawMessage(value pqtype.NullRawMessage) json.RawMessage {
	if !value.Valid {
		return nil
	}
	return value.RawMessage
}

func fromNullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}

func toNullString(value *string) sql.NullString {
	if value == nil {
		return sql.NullString{}
	}
	return sql.NullString{String: *value, Valid: true}
}

func toNullInt64(value *int64) sql.NullInt64 {
	if value == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{Int64: *value, Valid: true}
}

func toNullRawMessage(value json.RawMessage) pqtype.NullRawMessage {
	if value == nil {
		return pqtype.NullRawMessage{}
	}
	return pqtype.NullRawMessage{RawMessage: value, Valid: true}
}
