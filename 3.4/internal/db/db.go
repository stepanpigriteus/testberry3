package db

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type DB struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewStorage(ctx context.Context, masterDSN string, slaveDSNs []string, logger zerolog.Logger) *DB {
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


