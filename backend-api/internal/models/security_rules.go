package models

import (
	"encoding/json"
	"time"
)

type SecurityRuleType string

const (
	SecurityRuleTypeDeterministic SecurityRuleType = "deterministic"
	SecurityRuleTypeML            SecurityRuleType = "ml"
)

type SecurityRuleAction string

const (
	SecurityRuleActionAllow SecurityRuleAction = "allow"
	SecurityRuleActionBlock SecurityRuleAction = "block"
)

type SecurityRule struct {
	ID          int64
	ProfileID   int64
	Name        string
	Description *string
	Priority    int32
	RuleType    SecurityRuleType
	Action      SecurityRuleAction
	Conditions  json.RawMessage
	MLModelID   *int64
	MLThreshold *int16
	DryRun      bool
	IsEnabled   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SecurityRuleToCreate struct {
	ProfileID   int64
	Name        string
	Description *string
	Priority    int32
	RuleType    SecurityRuleType
	Action      SecurityRuleAction
	Conditions  json.RawMessage
	MLModelID   *int64
	MLThreshold *int16
	DryRun      bool
	IsEnabled   bool
}

type SecurityRuleToUpdate struct {
	ID          int64
	ProfileID   int64
	Name        string
	Description *string
	Priority    int32
	RuleType    SecurityRuleType
	Action      SecurityRuleAction
	Conditions  json.RawMessage
	MLModelID   *int64
	MLThreshold *int16
	DryRun      bool
	IsEnabled   bool
}
