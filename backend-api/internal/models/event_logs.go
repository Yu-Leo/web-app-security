package models

import (
	"encoding/json"
	"time"
)

type EventLog struct {
	ID         int64
	ResourceID int64
	OccurredAt time.Time
	EventType  string
	Severity   string
	Message    string
	RuleID     *int64
	ProfileID  *int64
	Metadata   json.RawMessage
	RequestID  *string
	ClientIP   *string
	Method     *string
	Path       *string
}
