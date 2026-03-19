package envoy_filter

import (
	"encoding/json"
	"time"
)

type Decision struct {
	Allowed         bool
	Action          string
	StatusCode      int32
	RuleID          *int64
	ProfileID       *int64
	MatchedRule     *string
	MatchedRuleType *string
	DryRun          bool
	Reason          *string
	DecisionSource  *string
	MLModelID       *int64
	MLModelName     *string
	MLThreshold     *int16
	MLScore         *float32
	MLScorePercent  *float32
	MLError         *string
	DecidedAt       time.Time
}

type RequestContext struct {
	OccurredAt         time.Time
	Method             string
	Path               string
	Query              string
	Host               string
	Scheme             string
	Protocol           string
	Authority          string
	Headers            map[string]string
	ClientIP           string
	SourcePort         int32
	DestinationIP      string
	DestinationPort    int32
	UserAgent          string
	RequestID          string
	RequestHTTPID      string
	Fragment           string
	SourcePrincipal    string
	SourceService      string
	DestinationService string
	RequestBodySize    int32
	RequestBody        string
	SourceLabels       map[string]string
	DestinationLabels  map[string]string
	ContextExtensions  map[string]string
	MetadataContext    json.RawMessage
	RouteMetadataCtx   json.RawMessage
	MLFeatureVector    []float32
}
