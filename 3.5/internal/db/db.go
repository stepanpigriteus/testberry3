package db

import (
	"context"

	"threeFive/domain"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type DB struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewDb(ctx context.Context, masterDSN string, slaveDSNs []string, logger zerolog.Logger) *DB {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("init database error %s", err)
	}
	return &DB{
		db:     db,
		logger: logger,
	}
}

func (d *DB) Create(ctx context.Context, original string) (string, error) {
	id := uuid.New().String()

	return id, nil
}

func (d *DB) GetEvent(ctx context.Context, id string) (domain.Event, error) {
	var event domain.Event
	return event, nil
}

func (d *DB) Update(ctx context.Context, id string) error {
	return nil
}
