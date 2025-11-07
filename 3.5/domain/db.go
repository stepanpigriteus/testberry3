package domain

import (
	"context"
	"database/sql"
)

type DB interface {
	Create(context.Context, Event) (string, error)
	Book(context.Context, string) (string, error)
	GetMaster() *sql.DB
	Update(context.Context, string, string) error
	CreateUser(context.Context, User) error
	GetAll(context.Context) ([]Event, error)
}
