package service

import (
	"WarehouseControl/internal/db"
	"context"

	"github.com/rs/zerolog"
)

type Serv struct {
	logger zerolog.Logger
	db     db.DB
}

func NewService(ctx context.Context, logger zerolog.Logger, db db.DB) *Serv {
	return &Serv{
		logger: logger,
		db:     db,
	}
}
