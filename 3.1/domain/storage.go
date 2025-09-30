package domain

import "context"

type Storage interface {
	CreateNotify(ctx context.Context, notify Notify) error
	GetNotify(ctx context.Context, id string) (Notify, error)
	DeleteNotify(ctx context.Context, id string) error
}
