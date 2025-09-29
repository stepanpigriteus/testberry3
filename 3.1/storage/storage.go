package storage

import (
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
	"github.com/wb-go/wbf/zlog"
)

type StorageImpl struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewStorage(masterDSN string, slaveDSNs []string, logger zerolog.Logger) *StorageImpl {
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

func (st *StorageImpl) CreateNotify(notify domain.Notify) error {
	return nil
}

func (st *StorageImpl) GetNotify(id string) (error, domain.Notify) {
	var n domain.Notify
	return nil, n
}

func (st *StorageImpl) DeleteNotify(id string) error {
	return nil
}
