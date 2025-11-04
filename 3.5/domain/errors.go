package domain

import "errors"

var (
	ErrInvalidSeats   = errors.New("invalid number of seats")
	ErrAlreadyExists  = errors.New("event already exists")
	ErrNotFound       = errors.New("not found")
	ErrInvalidBooking = errors.New("invalid booking data")
	ErrDuplicateKey   = errors.New("duplicate key")
)
