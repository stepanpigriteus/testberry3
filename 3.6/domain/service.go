package domain

import "context"

type Service interface {
	CreateItem(ctx context.Context, input Item) (*Item, error)
	GetItems(ctx context.Context, filter Filter) ([]Item, error)
	GetItem(ctx context.Context, id string) (*Item, error)
	UpdateItem(ctx context.Context, id string, input Item) (*Item, error)
	DeleteItem(ctx context.Context, id string) error
	GetAnalytics(ctx context.Context, filter AnalyticsFilter) (*AnalyticsResult, error)
}
