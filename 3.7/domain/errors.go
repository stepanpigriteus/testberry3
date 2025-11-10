package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrAlreadyExists      = errors.New("already exists")
	ErrInvalidAggregation = errors.New("invalid aggregation")
	ErrDuplicateKey       = errors.New("duplicate key")
	ErrItemNotFound       = errors.New("item not found")
	ErrInvalidStatus      = errors.New("invalid status")
)
