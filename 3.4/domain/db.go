package domain

import "context"

type DB interface {
	Create(ctx context.Context, original string) (string, error)
	Update(ctx context.Context, id string) error
	UpdatePathsAndStatus(ctx context.Context, id, thumb, watermark, resized string) error
}
