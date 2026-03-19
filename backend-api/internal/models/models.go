package models

import "errors"

var (
	ErrNilRequest = errors.New("nil request")
	ErrNotFound   = errors.New("not found")
)
