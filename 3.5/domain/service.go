package domain

import "context"

type Service interface {
	Gets(context.Context, string) (Event, error)
	Book(context.Context, string, string) (string, error)
	Confirm(context.Context, string, string) error
	Create(context.Context, Event) (string, error)
	CreateUser(context.Context, User) (string, error)
}
