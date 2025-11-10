package domain

import "context"

type DB interface {
	Create(ctx context.Context, item *Item) (Item, error)
	GetAll(ctx context.Context, filter Filter) ([]Item, error)
	GetByID(ctx context.Context, id string) (Item, error)
	Update(ctx context.Context, item *Item) (Item, error)
	Delete(ctx context.Context, id string) error

	GetAnalytics(ctx context.Context, filter AnalyticsFilter) (*AnalyticsResult, error)
}

type Filter struct {
	From     string
	To       string
	Type     string
	Category string
	SortBy   string
	Order    string
}

type AnalyticsFilter struct {
	From    string
	To      string
	Type    string
	GroupBy string
}
