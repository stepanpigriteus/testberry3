package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
)

type Service struct {
	db     domain.Storage
	logger zerolog.Logger
	redis  redis.Client
	rabbit *rabbitmq.Publisher
}

func NewService(db domain.Storage, logger zerolog.Logger, redis redis.Client, rabbit *rabbitmq.Publisher) *Service {
	return &Service{
		db:     db,
		logger: logger,
		redis:  redis,
		rabbit: rabbit,
	}
}

func (s *Service) CreateNotify(ctx context.Context, notify domain.Notify) error {
	data, err := json.Marshal(notify)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to marshal notify")
		return err
	}

	if err := s.redis.Set(ctx, strconv.Itoa(notify.Id), data); err != nil {
		s.logger.Error().Err(err).Msg("failed to save notify in redis")
		return err
	}

	if err := s.db.CreateNotify(ctx, notify); err != nil {
		s.logger.Error().Err(err).Msg("failed to save notify in db")
		return err
	}

	jnotify, _ := json.Marshal(notify)

	err = s.rabbit.Publish([]byte(string(jnotify)), "test_key", "text/plain")
	if err != nil {
		log.Fatalf("Ошибка при публикации: %v", err)
	}

	return nil
}

func (s *Service) GetNotify(ctx context.Context, id string) (domain.Notify, error) {
	var notify domain.Notify
	val, err := s.redis.Get(ctx, id)
	fmt.Println(val)
	if err == nil {
		if e := json.Unmarshal([]byte(val), &notify); e == nil {
			return notify, nil
		}
		s.logger.Warn().Err(err).Msg("failed to unmarshal notify from redis")
	} else if err.Error() != "redis: nil" {
		return notify, err
	}

	notify, err = s.db.GetNotify(ctx, id)
	if err != nil {
		return notify, err
	}

	data, e := json.Marshal(notify)
	if e == nil {
		if err := s.redis.Set(ctx, id, string(data)); err != nil {
			s.logger.Warn().Err(err).Msg("failed to cache notify in redis")
		}
	}

	return notify, nil
}

func (s *Service) DeleteNotify(ctx context.Context, id string) error {
	res, err := s.redis.Del(ctx, id).Result()
	if err != nil {
		return fmt.Errorf("failed to delete key %s: %w", id, err)
	}
	if res > 0 {
		s.logger.Info().Msgf("%d keys from redis deleted", res)
	}

	return nil
}
