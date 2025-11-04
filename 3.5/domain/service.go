package domain

import "context"

type Service interface {
	Gets()
	Book(context.Context, string) (string, error)
	Confirm()
	Create(context.Context, Event) (string, error)
}
