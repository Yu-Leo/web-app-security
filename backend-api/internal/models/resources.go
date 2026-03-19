package models

import "time"

type Resource struct {
	ID                int64
	Name              string
	URLPattern        string
	SecurityProfileID *int64
	TrafficProfileID  *int64
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type ResourceToCreate struct {
	Name              string
	URLPattern        string
	SecurityProfileID *int64
	TrafficProfileID  *int64
}

type ResourceToUpdate struct {
	ID                int64
	Name              string
	URLPattern        string
	SecurityProfileID *int64
	TrafficProfileID  *int64
}
