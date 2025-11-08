package db

import (
	"context"
	"database/sql"
	"threeSixth/domain"

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

func (d *DB) Create(ctx context.Context, item *domain.Item) error {
	return nil
}

func (d *DB) GetAll(ctx context.Context, filter domain.Filter) ([]domain.Item, error) {
	return nil, nil
}

func (d *DB) GetByID(ctx context.Context, id string) (*domain.Item, error) {
	return nil, nil
}

func (d *DB) Update(ctx context.Context, item *domain.Item) error {
	return nil
}

func (d *DB) Delete(ctx context.Context, id string) error {
	return nil
}

func (d *DB) GetAnalytics(ctx context.Context, filter domain.AnalyticsFilter) (*domain.AnalyticsResult, error) {
	return nil, nil
}

func (d *DB) Close() error {
	if d.db == nil {
		return nil
	}

	d.logger.Info().Msg("Closing database connection...")
	if err := d.db.Master.Close(); err != nil {
		d.logger.Error().Err(err).Msg("Failed to close database connection")
		return err
	}

	d.logger.Info().Msg("Database connection closed successfully")
	return nil
}

func (d *DB) GetMaster() *sql.DB {
	return d.db.Master
}
