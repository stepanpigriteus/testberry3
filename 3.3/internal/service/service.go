package service

import (
	"context"
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
