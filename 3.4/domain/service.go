package domain

import (
	"context"
)

type Service interface {
	Upload(ctx context.Context, image ImageData) error
	Get(ctx context.Context, id int) (ImageData, error)
	Delete(ctx context.Context, id int) error
}
