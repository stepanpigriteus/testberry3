package domain

import "context"

type Service interface {
	CreateComments(ctx context.Context, comment Comment) error
	GetComments(ctx context.Context, id int) (Comment, error)
	DeleteComments(ctx context.Context, id int) error
}
