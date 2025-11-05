package domain

import "context"

type Service interface {
	Gets(context.Context, string) (Event, error)
	Book(context.Context, string) (string, error)
	Confirm()
	Create(context.Context, Event) (string, error)
}
