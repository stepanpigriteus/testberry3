package service

import (
	"context"
	"threeSixth/domain"
	"threeSixth/internal/db"

	"github.com/rs/zerolog"
)

type Serv struct {
	logger zerolog.Logger
	db     db.DB
}

func NewService(ctx context.Context, logger zerolog.Logger, db db.DB) *Serv {
	return &Serv{
		logger: logger,
		db:     db,
	}
}

func (s *Serv) CreateItem(ctx context.Context, input domain.Item) (*domain.Item, error) {
	return nil, nil
}

func (s *Serv) GetItems(ctx context.Context, filter domain.Filter) ([]domain.Item, error) {
	return nil, nil
}

func (s *Serv) GetItem(ctx context.Context, id string) (*domain.Item, error) {
	return nil, nil
}

func (s *Serv) UpdateItem(ctx context.Context, id string, input domain.Item) (*domain.Item, error) {
	return nil, nil
}

func (s *Serv) DeleteItem(ctx context.Context, id string) error {
	return nil
}

func (s *Serv) GetAnalytics(ctx context.Context, filter domain.AnalyticsFilter) (*domain.AnalyticsResult, error) {
	return nil, nil
}
