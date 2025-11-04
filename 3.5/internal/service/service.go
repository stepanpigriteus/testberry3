package service

import (
	"context"
	"errors"
	"fmt"

	"threeFive/domain"
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

func (s *Serv) Gets() {

}

func (s *Serv) Confirm() {

}

func (s *Serv) Create(ctx context.Context, event domain.Event) (string, error) {

	id, err := s.db.Create(ctx, event)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			return "", domain.ErrAlreadyExists
		}
		return "", fmt.Errorf("create event: %w", err)
	}

	return id, nil
}

func (s *Serv) Book(ctx context.Context, eventId string) (string, error) {

	bookId, err := s.db.Book(ctx, eventId)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			return "", domain.ErrAlreadyExists
		}
		return "", fmt.Errorf("create event: %w", err)
	}

	return bookId, nil
}
