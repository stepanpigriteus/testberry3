package storage

import (
	"context"
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type StorageImpl struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewStorage(ctx context.Context, masterDSN string, slaveDSNs []string, logger zerolog.Logger) *StorageImpl {
	opts := &dbpg.Options{MaxOpenConns: 10, MaxIdleConns: 5}
	db, err := dbpg.New(masterDSN, slaveDSNs, opts)
	if err != nil {
		zlog.Logger.Error().Msgf("init database error %s", err)
	}
	return &StorageImpl{
		db:     db,
		logger: logger,
	}
}

func (st *StorageImpl) CreateNotify(ctx context.Context, notify domain.Notify) error {
	return nil
}

func (st *StorageImpl) GetNotify(ctx context.Context, id string) (domain.Notify, error) {
	var n domain.Notify
	return n, nil
}

func (st *StorageImpl) DeleteNotify(ctx context.Context, id string) error {
	return nil
}
