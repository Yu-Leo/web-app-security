package models

import "time"

type TrafficProfile struct {
	ID          int64
	Name        string
	Description *string
	IsEnabled   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type TrafficProfileToCreate struct {
	Name        string
	Description *string
	IsEnabled   bool
}

type TrafficProfileToUpdate struct {
	ID          int64
	Name        string
	Description *string
	IsEnabled   bool
}
