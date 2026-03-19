package models

import (
	"encoding/json"
	"time"
)

type TrafficRule struct {
	ID            int64
	ProfileID     int64
	Name          string
	Description   *string
	Priority      int32
	DryRun        bool
	MatchAll      bool
	RequestsLimit int32
	PeriodSeconds int32
	Conditions    json.RawMessage
	IsEnabled     bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type TrafficRuleToCreate struct {
	ProfileID     int64
	Name          string
	Description   *string
	Priority      int32
	DryRun        bool
	MatchAll      bool
	RequestsLimit int32
	PeriodSeconds int32
	Conditions    json.RawMessage
	IsEnabled     bool
}

type TrafficRuleToUpdate struct {
	ID            int64
	ProfileID     int64
	Name          string
	Description   *string
	Priority      int32
	DryRun        bool
	MatchAll      bool
	RequestsLimit int32
	PeriodSeconds int32
	Conditions    json.RawMessage
	IsEnabled     bool
}
