package service

import (
	"context"

	"threeFive/internal/db"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/kafka"
)

type Serv struct {
	logger   zerolog.Logger
	producer *kafka.Producer
	consumer *kafka.Consumer
	db       db.DB
}

func NewService(ctx context.Context, producer *kafka.Producer, consumer *kafka.Consumer, logger zerolog.Logger, db db.DB) *Serv {
	return &Serv{
		logger:   logger,
		producer: producer,
		consumer: consumer,
		db:       db,
	}
}
