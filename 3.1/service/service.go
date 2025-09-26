package service

import (
	"treeOne/domain"

	"github.com/rs/zerolog"
	"github.com/wb-go/wbf/dbpg"
)

type Service struct {
	db     *dbpg.DB
	logger zerolog.Logger
}

func NewService(db *dbpg.DB, logger zerolog.Logger) *Service {
	return &Service{
		db:     db,
		logger: logger,
	}
}

func (s *Service) CreateNotify(notify domain.Notify) error {
	return nil
}

func (s *Service) GetNotify(id string) (error, domain.Notify) {
	var not domain.Notify
	return nil, not
}

func (s *Service) DeleteNotify(id string) error {
	return nil
}
