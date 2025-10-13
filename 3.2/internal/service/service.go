package service

import (
	"context"
	"encoding/json"
	"treeTwo/domain"
	"treeTwo/pkg"

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

func (s *Serv) CreateShorten(ctx context.Context, link string) error {
	short := pkg.Shortener(link)
	s.logger.Log().Msg(short)
	linkSt := domain.ShortLink{
		ShortCode:   short,
		OriginalURL: link,
	}

	err := s.db.CreateShorten(ctx, linkSt)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to save shorten in db")
		return err
	}

	return nil
}

func (s *Serv) GetShorten(ctx context.Context, link string) (domain.ShortLink, error) {
	var shortLink domain.ShortLink
	val, err := s.redis.Get(ctx, link)
	if err == nil {

		if e := json.Unmarshal([]byte(val), &shortLink); e == nil {
			shortLink.ClickCount++
			_ = s.updateCounters(ctx, shortLink)
			s.logger.Info().Msg("shortlink get from redis")
			return shortLink, nil
		}
		s.logger.Warn().Msg("failed to unmarshal shortLink from redis")
	} else if err.Error() != "redis: nil" {
		s.logger.Warn().Err(err).Msg("redis get failed")
	}

	shortLink, err = s.db.GetShorten(ctx, link)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to Get shorten in db")
		return shortLink, err
	}
	data, _ := json.Marshal(shortLink)
	err = s.redis.Set(ctx, shortLink.ShortCode, data)
	if err != nil {
		s.logger.Warn().Err(err).Msg("failed to cache short link in Redis")
	}

	shortLink.ClickCount++
	_ = s.updateCounters(ctx, shortLink)
	return shortLink, nil

}

func (s *Serv) GetAnalytics(ctx context.Context, shortCode string, group string) (domain.VisitStats, error) {
	
	stats, err := s.db.GetAnalytics(ctx, shortCode, group)
	if err != nil {
		s.logger.Error().Err(err).Msg("failed to get analytics")
	}
	return stats, err
}

func (s *Serv) updateCounters(ctx context.Context, shortLink domain.ShortLink) error {
	if err := s.db.UpdateClickCount(ctx, shortLink.ID, shortLink.ClickCount); err != nil {
		s.logger.Warn().Err(err).Msg("failed to update click_count in db")
	}

	data, _ := json.Marshal(shortLink)
	if err := s.redis.Set(ctx, shortLink.ShortCode, data); err != nil {
		s.logger.Warn().Err(err).Msg("failed to update redis cache")
	}

	return nil
}

func (s *Serv) SaveVisit(ctx context.Context, visit domain.Visit) error {
	if err := s.db.SaveVisit(ctx, visit); err != nil {
		s.logger.Warn().Err(err).Msg("failed to save visit in db")
	}
	return nil
}
