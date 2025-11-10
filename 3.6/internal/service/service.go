package service

import (
	"context"
	"errors"
	"fmt"
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

func (s *Serv) CreateItem(ctx context.Context, input domain.Item) (domain.Item, error) {
	item, err := s.db.Create(ctx, &input)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			s.logger.Warn().
				Str("type", input.Type).
				Str("category", input.Category).
				Msg("Attempted to create duplicate item")
			return domain.Item{}, fmt.Errorf("item already exists: %w", err)
		}

		s.logger.Error().Err(err).
			Str("type", input.Type).
			Str("category", input.Category).
			Msg("Failed to create item")
		return domain.Item{}, fmt.Errorf("create item: %w", err)
	}

	s.logger.Info().Str("id", item.ID).Str("type", item.Type).Float64("amount", item.Amount).Msg("Item created successfully")
	return item, nil
}

func (s *Serv) GetItems(ctx context.Context, filter domain.Filter) ([]domain.Item, error) {
	items, err := s.db.GetAll(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).Interface("filter", filter).Msg("Failed to get items")
		return nil, fmt.Errorf("get items: %w", err)
	}

	s.logger.Debug().Int("count", len(items)).Msg("Items retrieved successfully")

	return items, nil
}

func (s *Serv) GetItem(ctx context.Context, id string) (domain.Item, error) {
	if id == "" {
		return domain.Item{}, fmt.Errorf("id cannot be empty")
	}
	item, err := s.db.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			s.logger.Warn().Str("id", id).Msg("Item not found")
			return domain.Item{}, domain.ErrNotFound
		}
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to get item")
		return domain.Item{}, fmt.Errorf("get item: %w", err)
	}

	return *item, nil
}

func (s *Serv) UpdateItem(ctx context.Context, id string, input domain.Item) (*domain.Item, error) {
	if id == "" {
		return nil, fmt.Errorf("id cannot be empty")
	}

	input.ID = id
	updated, err := s.db.Update(ctx, &input)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			s.logger.Warn().Str("id", id).Msg("Item not found for update")
			return nil, domain.ErrNotFound
		}
		s.logger.Error().Err(err).
			Str("id", id).
			Msg("Failed to update item")
		return nil, fmt.Errorf("update item: %w", err)
	}

	s.logger.Info().Str("id", updated.ID).Msg("Item updated successfully")

	return &updated, nil
}

func (s *Serv) DeleteItem(ctx context.Context, id string) error {
	if id == "" {
		return fmt.Errorf("id cannot be empty")
	}

	err := s.db.Delete(ctx, id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			s.logger.Warn().Str("id", id).Msg("Item not found for deletion")
			return domain.ErrNotFound
		}
		s.logger.Error().Err(err).Str("id", id).Msg("Failed to delete item")
		return fmt.Errorf("delete item: %w", err)
	}
	s.logger.Info().Str("id", id).Msg("Item deleted successfully")
	return nil
}

func (s *Serv) GetAnalytics(ctx context.Context, filter domain.AnalyticsFilter) (*domain.AnalyticsResult, error) {
	analit, err := s.db.GetAnalytics(ctx, filter)
	if err != nil {
		s.logger.Error().Err(err).
			Interface("filter", filter).
			Msg("Failed to get analytics")
		return nil, fmt.Errorf("get analytics: %w", err)
	}

	s.logger.Debug().Int("grouped_count", len(analit.Grouped)).Msg("Analytics retrieved successfully")

	return analit, nil
}
