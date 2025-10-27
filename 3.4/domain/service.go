package domain

import (
	"context"
)

type Service interface {
	Upload(ctx context.Context, image ImageData) (string, error)
	Get(ctx context.Context, id string) ([][]byte, error)
	Delete(ctx context.Context, id string) error
}
