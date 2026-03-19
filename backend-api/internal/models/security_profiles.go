package models

import "time"

type SecurityProfileBaseAction string

const (
	SecurityProfileBaseActionAllow SecurityProfileBaseAction = "allow"
	SecurityProfileBaseActionBlock SecurityProfileBaseAction = "block"
)

type SecurityProfile struct {
	ID          int64
	Name        string
	Description *string
	BaseAction  SecurityProfileBaseAction
	LogEnabled  bool
	IsEnabled   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type SecurityProfileToCreate struct {
	Name        string
	Description *string
	BaseAction  SecurityProfileBaseAction
	LogEnabled  bool
	IsEnabled   bool
}

type SecurityProfileToUpdate struct {
	ID          int64
	Name        string
	Description *string
	BaseAction  SecurityProfileBaseAction
	LogEnabled  bool
	IsEnabled   bool
}
