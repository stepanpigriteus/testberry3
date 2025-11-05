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

func (s *Serv) Gets(ctx context.Context, eventId string) (domain.Event, error) {

	event, err := s.db.GetEvent(ctx, eventId)
	if err != nil {
		return domain.Event{}, err
	}

	return event, nil
}

func (s *Serv) Confirm(ctx context.Context, eventId string, bookId string) error {
	err := s.db.Update(ctx, eventId, bookId)
	if err != nil {
		return err
	}

	return nil

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

func (s *Serv) Book(ctx context.Context, eventId string, userId string) (string, error) {

	bookId, err := s.db.Book(ctx, eventId, userId)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateKey) {
			return "", domain.ErrAlreadyExists
		}
		return "", fmt.Errorf("create event: %w", err)
	}

	return bookId, nil
}

func (s *Serv) CreateUser(ctx context.Context, user domain.User) (string, error) {

	if user.Email == "" {
		return "", fmt.Errorf("email is required")
	}
	if user.Name == "" {
		return "", fmt.Errorf("name is required")
	}

	if user.Role == "" {
		user.Role = "user"
	}

	userId, err := s.db.CreateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("create user: %w", err)
	}

	return userId, nil
}
