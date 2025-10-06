package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/redis"
	"github.com/wb-go/wbf/retry"
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

	delay := time.Until(notify.Timing)
	if delay < 0 {
		delay = 0
	}

	strategy := retry.Strategy{
		Attempts: 5,
		Delay:    500 * time.Millisecond,
		Backoff:  2,
	}

	err = retry.Do(func() error {
		return s.rabbit.Publish(
			data,
			"delay_key",
			"application/json",
			rabbitmq.PublishingOptions{
				Expiration: delay,
			},
		)
	}, strategy)

	if err != nil {
		s.logger.Error().Err(err).Msg("Ошибка при публикации")
		return err
	}

	return nil
}

func (s *Service) GetNotify(ctx context.Context, id string) (domain.Notify, error) {
	var notify domain.Notify
	val, err := s.redis.Get(ctx, id)
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

	err = s.db.DeleteNotify(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) CheckStatus(ctx context.Context, msg []byte) (error, bool) {
	var notify domain.Notify
	if err := json.Unmarshal(msg, &notify); err != nil {
		s.logger.Warn().Err(err).Msg("failed to unmarshal notify from request")
		return err, false
	}
	idStr := strconv.Itoa(notify.Id)
	val, err := s.redis.Get(ctx, idStr)
	if err != nil {
		if err.Error() == "redis: nil" {
			s.logger.Info().Msgf("ключ %d не найден в redis, проверяем БД", notify.Id)
			dbNotify, dbErr := s.db.GetNotify(ctx, idStr)
			if dbErr != nil {
				if errors.Is(dbErr, sql.ErrNoRows) {
					s.logger.Warn().Msgf("уведомление id=%d не найдено в БД", notify.Id)
					return fmt.Errorf("уведомление id=%d не найдено", notify.Id), false
				}
				return dbErr, false
			}
			data, e := json.Marshal(dbNotify)
			if e == nil {
				if setErr := s.redis.Set(ctx, idStr, string(data)); setErr != nil {
					s.logger.Warn().Err(setErr).Msgf("не удалось сохранить уведомление id=%d в Redis", notify.Id)
				}
			}

			if dbNotify.Status == "cancelled" {
				s.logger.Warn().Msgf("уведомление id=%d отменено", notify.Id)
				return nil, false
			}
			return nil, true
		}
		s.logger.Error().Err(err).Msg("ошибка при чтении из redis")
		return err, false
	}

	if err := json.Unmarshal([]byte(val), &notify); err != nil {
		s.logger.Warn().Err(err).Msg("failed to unmarshal notify from redis")
		return err, false
	}

	if notify.Status == "cancelled" {
		s.logger.Warn().Msgf("уведомление id=%d отменено", notify.Id)
		return nil, false
	}

	return nil, true
}
