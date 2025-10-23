package domain

import "context"

type DB interface {
	Create(ctx context.Context, original string) (string, error)
}
