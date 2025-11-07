package domain

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrAlreadyExists  = errors.New("already exists")
	ErrInvalidSeats   = errors.New("invalid seats")
	ErrInvalidBooking = errors.New("invalid booking")
	ErrDuplicateKey   = errors.New("duplicate key")
	ErrUserNotFound   = errors.New("user not found")
	ErrInvalidStatus  = errors.New("invalid status")
)
