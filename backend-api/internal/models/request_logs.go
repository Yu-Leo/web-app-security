package models

import (
	"encoding/json"
	"time"
)

type RequestLog struct {
	ID                 int64
	ResourceID         *int64
	OccurredAt         time.Time
	ClientIP           string
	Method             string
	Path               string
	StatusCode         int32
	Action             string
	RuleID             *int64
	ProfileID          *int64
	UserAgent          *string
	Country            *string
	LatencyMs          *int32
	RequestID          *string
	Metadata           json.RawMessage
	Host               *string
	Scheme             *string
	Protocol           *string
	Authority          *string
	Query              *string
	SourcePort         *int32
	DestinationIP      *string
	DestinationPort    *int32
	SourcePrincipal    *string
	SourceService      *string
	SourceLabels       json.RawMessage
	DestinationService *string
	DestinationLabels  json.RawMessage
	RequestHTTPID      *string
	Fragment           *string
	RequestHeaders     json.RawMessage
	RequestBodySize    *int32
	RequestBody        *string
	ContextExtensions  json.RawMessage
	MetadataContext    json.RawMessage
	RouteMetadataCtx   json.RawMessage
}
