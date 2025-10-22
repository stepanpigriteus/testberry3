package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"treethree/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/redis"
)

type Serv struct {
	db     domain.Storage
	logger zerolog.Logger
	redis  redis.Client
}

func NewService(ctx context.Context, storage domain.Storage, logger zerolog.Logger, redis redis.Client) *Serv {
	return &Serv{
		db:     storage,
		logger: logger,
		redis:  redis,
	}
}

func (s *Serv) CreateComments(ctx context.Context, comment domain.Comment) error {

	err := s.db.CreateComments(ctx, comment)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to save comment in db")
		return err
	}
	return nil
}

func (s *Serv) GetComments(ctx context.Context, id int) (domain.Comment, error) {
	comments, err := s.db.GetComments(ctx, id)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to Get comments in db")
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Comment{}, fmt.Errorf("comment not found: %w", err)
		}
		return domain.Comment{}, fmt.Errorf("db error: %w", err)
	}
	return comments, nil
}
func (s *Serv) DeleteComments(ctx context.Context, id int) error {
	err := s.db.DeleteComments(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.logger.Warn().Int("id", id).Msg("comment not found")
			return fmt.Errorf("comment not found: %w", err)
		}
		s.logger.Error().Err(err).Msg("failed to delete comments in db")
		return err
	}
	return nil
}
