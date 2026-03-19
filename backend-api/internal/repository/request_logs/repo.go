package request_logs

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

func (r *Repository) Create(ctx context.Context, logRecord models.RequestLog) (models.RequestLog, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	created, err := querier.CreateRequestLog(ctx, db.CreateRequestLogParams{
		ResourceID:           toNullInt64(logRecord.ResourceID),
		OccurredAt:           logRecord.OccurredAt,
		ClientIp:             logRecord.ClientIP,
		Method:               logRecord.Method,
		Path:                 logRecord.Path,
		StatusCode:           logRecord.StatusCode,
		Action:               logRecord.Action,
		RuleID:               toNullInt64(logRecord.RuleID),
		ProfileID:            toNullInt64(logRecord.ProfileID),
		UserAgent:            toNullString(logRecord.UserAgent),
		Country:              toNullString(logRecord.Country),
		LatencyMs:            toNullInt32(logRecord.LatencyMs),
		RequestID:            toNullString(logRecord.RequestID),
		Metadata:             toNullRawMessage(logRecord.Metadata),
		Host:                 toNullString(logRecord.Host),
		Scheme:               toNullString(logRecord.Scheme),
		Protocol:             toNullString(logRecord.Protocol),
		Authority:            toNullString(logRecord.Authority),
		Query:                toNullString(logRecord.Query),
		SourcePort:           toNullInt32(logRecord.SourcePort),
		DestinationIp:        toNullString(logRecord.DestinationIP),
		DestinationPort:      toNullInt32(logRecord.DestinationPort),
		SourcePrincipal:      toNullString(logRecord.SourcePrincipal),
		SourceService:        toNullString(logRecord.SourceService),
		SourceLabels:         toNullRawMessage(logRecord.SourceLabels),
		DestinationService:   toNullString(logRecord.DestinationService),
		DestinationLabels:    toNullRawMessage(logRecord.DestinationLabels),
		RequestHttpID:        toNullString(logRecord.RequestHTTPID),
		Fragment:             toNullString(logRecord.Fragment),
		RequestHeaders:       toNullRawMessage(logRecord.RequestHeaders),
		RequestBodySize:      toNullInt32(logRecord.RequestBodySize),
		RequestBody:          toNullString(logRecord.RequestBody),
		ContextExtensions:    toNullRawMessage(logRecord.ContextExtensions),
		MetadataContext:      toNullRawMessage(logRecord.MetadataContext),
		RouteMetadataContext: toNullRawMessage(logRecord.RouteMetadataCtx),
	})
	if err != nil {
		return models.RequestLog{}, err
	}

	return requestLogFromCreateDB(created), nil
}

func (r *Repository) Get(ctx context.Context, id int64) (models.RequestLog, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	logRecord, err := querier.GetRequestLog(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.RequestLog{}, models.ErrNotFound
		}
		return models.RequestLog{}, err
	}

	return requestLogFromGetDB(logRecord), nil
}

func (r *Repository) List(ctx context.Context) ([]models.RequestLog, error) {
	tx := r.database.ProvideTransaction(ctx)
	querier := db.New(tx)

	logs, err := querier.ListRequestLogs(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]models.RequestLog, 0, len(logs))
	for _, logRecord := range logs {
		result = append(result, requestLogFromListDB(logRecord))
	}

	return result, nil
}

func requestLogFromCreateDB(logRecord db.CreateRequestLogRow) models.RequestLog {
	return models.RequestLog{
		ID:                 logRecord.ID,
		ResourceID:         fromNullInt64(logRecord.ResourceID),
		OccurredAt:         logRecord.OccurredAt,
		ClientIP:           logRecord.ClientIp,
		Method:             logRecord.Method,
		Path:               logRecord.Path,
		StatusCode:         logRecord.StatusCode,
		Action:             logRecord.Action,
		RuleID:             fromNullInt64(logRecord.RuleID),
		ProfileID:          fromNullInt64(logRecord.ProfileID),
		UserAgent:          fromNullString(logRecord.UserAgent),
		Country:            fromNullString(logRecord.Country),
		LatencyMs:          fromNullInt32(logRecord.LatencyMs),
		RequestID:          fromNullString(logRecord.RequestID),
		Metadata:           fromNullRawMessage(logRecord.Metadata),
		Host:               fromNullString(logRecord.Host),
		Scheme:             fromNullString(logRecord.Scheme),
		Protocol:           fromNullString(logRecord.Protocol),
		Authority:          fromNullString(logRecord.Authority),
		Query:              fromNullString(logRecord.Query),
		SourcePort:         fromNullInt32(logRecord.SourcePort),
		DestinationIP:      fromNullString(logRecord.DestinationIp),
		DestinationPort:    fromNullInt32(logRecord.DestinationPort),
		SourcePrincipal:    fromNullString(logRecord.SourcePrincipal),
		SourceService:      fromNullString(logRecord.SourceService),
		SourceLabels:       fromNullRawMessage(logRecord.SourceLabels),
		DestinationService: fromNullString(logRecord.DestinationService),
		DestinationLabels:  fromNullRawMessage(logRecord.DestinationLabels),
		RequestHTTPID:      fromNullString(logRecord.RequestHttpID),
		Fragment:           fromNullString(logRecord.Fragment),
		RequestHeaders:     fromNullRawMessage(logRecord.RequestHeaders),
		RequestBodySize:    fromNullInt32(logRecord.RequestBodySize),
		RequestBody:        fromNullString(logRecord.RequestBody),
		ContextExtensions:  fromNullRawMessage(logRecord.ContextExtensions),
		MetadataContext:    fromNullRawMessage(logRecord.MetadataContext),
		RouteMetadataCtx:   fromNullRawMessage(logRecord.RouteMetadataContext),
	}
}

func requestLogFromGetDB(logRecord db.GetRequestLogRow) models.RequestLog {
	return models.RequestLog{
		ID:                 logRecord.ID,
		ResourceID:         fromNullInt64(logRecord.ResourceID),
		OccurredAt:         logRecord.OccurredAt,
		ClientIP:           logRecord.ClientIp,
		Method:             logRecord.Method,
		Path:               logRecord.Path,
		StatusCode:         logRecord.StatusCode,
		Action:             logRecord.Action,
		RuleID:             fromNullInt64(logRecord.RuleID),
		ProfileID:          fromNullInt64(logRecord.ProfileID),
		UserAgent:          fromNullString(logRecord.UserAgent),
		Country:            fromNullString(logRecord.Country),
		LatencyMs:          fromNullInt32(logRecord.LatencyMs),
		RequestID:          fromNullString(logRecord.RequestID),
		Metadata:           fromNullRawMessage(logRecord.Metadata),
		Host:               fromNullString(logRecord.Host),
		Scheme:             fromNullString(logRecord.Scheme),
		Protocol:           fromNullString(logRecord.Protocol),
		Authority:          fromNullString(logRecord.Authority),
		Query:              fromNullString(logRecord.Query),
		SourcePort:         fromNullInt32(logRecord.SourcePort),
		DestinationIP:      fromNullString(logRecord.DestinationIp),
		DestinationPort:    fromNullInt32(logRecord.DestinationPort),
		SourcePrincipal:    fromNullString(logRecord.SourcePrincipal),
		SourceService:      fromNullString(logRecord.SourceService),
		SourceLabels:       fromNullRawMessage(logRecord.SourceLabels),
		DestinationService: fromNullString(logRecord.DestinationService),
		DestinationLabels:  fromNullRawMessage(logRecord.DestinationLabels),
		RequestHTTPID:      fromNullString(logRecord.RequestHttpID),
		Fragment:           fromNullString(logRecord.Fragment),
		RequestHeaders:     fromNullRawMessage(logRecord.RequestHeaders),
		RequestBodySize:    fromNullInt32(logRecord.RequestBodySize),
		RequestBody:        fromNullString(logRecord.RequestBody),
		ContextExtensions:  fromNullRawMessage(logRecord.ContextExtensions),
		MetadataContext:    fromNullRawMessage(logRecord.MetadataContext),
		RouteMetadataCtx:   fromNullRawMessage(logRecord.RouteMetadataContext),
	}
}

func requestLogFromListDB(logRecord db.ListRequestLogsRow) models.RequestLog {
	return models.RequestLog{
		ID:                 logRecord.ID,
		ResourceID:         fromNullInt64(logRecord.ResourceID),
		OccurredAt:         logRecord.OccurredAt,
		ClientIP:           logRecord.ClientIp,
		Method:             logRecord.Method,
		Path:               logRecord.Path,
		StatusCode:         logRecord.StatusCode,
		Action:             logRecord.Action,
		RuleID:             fromNullInt64(logRecord.RuleID),
		ProfileID:          fromNullInt64(logRecord.ProfileID),
		UserAgent:          fromNullString(logRecord.UserAgent),
		Country:            fromNullString(logRecord.Country),
		LatencyMs:          fromNullInt32(logRecord.LatencyMs),
		RequestID:          fromNullString(logRecord.RequestID),
		Metadata:           fromNullRawMessage(logRecord.Metadata),
		Host:               fromNullString(logRecord.Host),
		Scheme:             fromNullString(logRecord.Scheme),
		Protocol:           fromNullString(logRecord.Protocol),
		Authority:          fromNullString(logRecord.Authority),
		Query:              fromNullString(logRecord.Query),
		SourcePort:         fromNullInt32(logRecord.SourcePort),
		DestinationIP:      fromNullString(logRecord.DestinationIp),
		DestinationPort:    fromNullInt32(logRecord.DestinationPort),
		SourcePrincipal:    fromNullString(logRecord.SourcePrincipal),
		SourceService:      fromNullString(logRecord.SourceService),
		SourceLabels:       fromNullRawMessage(logRecord.SourceLabels),
		DestinationService: fromNullString(logRecord.DestinationService),
		DestinationLabels:  fromNullRawMessage(logRecord.DestinationLabels),
		RequestHTTPID:      fromNullString(logRecord.RequestHttpID),
		Fragment:           fromNullString(logRecord.Fragment),
		RequestHeaders:     fromNullRawMessage(logRecord.RequestHeaders),
		RequestBodySize:    fromNullInt32(logRecord.RequestBodySize),
		RequestBody:        fromNullString(logRecord.RequestBody),
		ContextExtensions:  fromNullRawMessage(logRecord.ContextExtensions),
		MetadataContext:    fromNullRawMessage(logRecord.MetadataContext),
		RouteMetadataCtx:   fromNullRawMessage(logRecord.RouteMetadataContext),
	}
}

func fromNullString(value sql.NullString) *string {
	if !value.Valid {
		return nil
	}
	result := value.String
	return &result
}

func fromNullInt64(value sql.NullInt64) *int64 {
	if !value.Valid {
		return nil
	}
	result := value.Int64
	return &result
}

func fromNullInt32(value sql.NullInt32) *int32 {
	if !value.Valid {
		return nil
	}
	result := value.Int32
	return &result
}

func fromNullRawMessage(value pqtype.NullRawMessage) json.RawMessage {
	if !value.Valid {
		return nil
	}
	return value.RawMessage
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

func toNullInt32(value *int32) sql.NullInt32 {
	if value == nil {
		return sql.NullInt32{}
	}
	return sql.NullInt32{Int32: *value, Valid: true}
}

func toNullRawMessage(value json.RawMessage) pqtype.NullRawMessage {
	if value == nil {
		return pqtype.NullRawMessage{}
	}
	return pqtype.NullRawMessage{RawMessage: value, Valid: true}
}
