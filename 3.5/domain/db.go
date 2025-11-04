package domain

import "context"

type DB interface {
	Create(context.Context, Event) (string, error)
	Book(context.Context, string) (string, error)
}
