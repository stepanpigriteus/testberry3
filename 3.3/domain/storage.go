package domain

import "context"

type Storage interface {
	CreateComments(ctx context.Context, comment Comment) error
	GetComments(ctx context.Context, id int) (Comment, error)
}
